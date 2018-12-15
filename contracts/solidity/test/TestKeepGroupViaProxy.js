import increaseTime, { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import exceptThrow from './helpers/expectThrow';
import encodeCall from './helpers/encodeCall';
import BigNumber from 'bignumber.js';
import abi from 'ethereumjs-abi';
const KeepToken = artifacts.require('./KeepToken.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const KeepRandomBeaconProxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconImplV1 = artifacts.require('./KeepRandomBeaconImplV1.sol');
const KeepGroupProxy = artifacts.require('./KeepGroup.sol');
const KeepGroupImplV1 = artifacts.require('./KeepGroupImplV1.sol');

contract('TestKeepGroupViaProxy', function(accounts) {

  let token, stakingProxy, stakingContract, 
    keepRandomBeaconImplV1, keepRandomBeaconProxy, keepRandomBeaconImplViaProxy,
    keepGroupImplV1, keepGroupProxy, keepGroupImplViaProxy, groupOnePubKey, groupTwoPubKey,
    account_one = accounts[0],
    account_two = accounts[1];

  beforeEach(async () => {
    token = await KeepToken.new();
    
    // Initialize staking contract under proxy
    stakingProxy = await StakingProxy.new();
    stakingContract = await TokenStaking.new(token.address, stakingProxy.address, duration.days(30));
    await stakingProxy.authorizeContract(stakingContract.address, {from: account_one})

    // Initialize Keep Random Beacon
    keepRandomBeaconImplV1 = await KeepRandomBeaconImplV1.new();
    keepRandomBeaconProxy = await KeepRandomBeaconProxy.new(keepRandomBeaconImplV1.address);
    keepRandomBeaconImplViaProxy = await KeepRandomBeaconImplV1.at(keepRandomBeaconProxy.address);
    await keepRandomBeaconImplViaProxy.initialize(100, duration.days(30));

    // Initialize Keep Group contract
    let minimumStake = 200;
    keepGroupImplV1 = await KeepGroupImplV1.new();
    keepGroupProxy = await KeepGroupProxy.new(keepGroupImplV1.address);
    keepGroupImplViaProxy = await KeepGroupImplV1.at(keepGroupProxy.address);
    await keepGroupImplViaProxy.initialize(stakingProxy.address, minimumStake, 1, 2, 1, 3, 4);

    // Create test groups.
    groupOnePubKey = "0x1000000000000000000000000000000000000000000000000000000000000000";
    await keepGroupImplViaProxy.createGroup(groupOnePubKey);
    groupTwoPubKey = "0x2000000000000000000000000000000000000000000000000000000000000000";
    await keepGroupImplViaProxy.createGroup(groupTwoPubKey);

    // Stake tokens as account one so it has minimum stake to be able to get into a group.
    await token.approveAndCall(stakingContract.address, minimumStake, "", {from: account_one});

  });

  it("should fail to update minimum stake by non owner", async function() {
    await exceptThrow(keepGroupImplViaProxy.setMinimumStake(123, {from: account_two}));
  });

  it("should be able to update minimum stake by the owner", async function() {
    await keepGroupImplViaProxy.setMinimumStake(123);
    let newMinStake = await keepGroupImplViaProxy.minimumStake();
    assert.equal(newMinStake, 123, "Should be able to get updated minimum stake.");
  });

  it("should be able to check if the implementation contract was initialized", async function() {
    assert.equal(await keepGroupImplViaProxy.initialized(), true, "Implementation contract should be initialized.");
  });

  it("should be able to return a total number of created group", async function() {
    assert.equal(await keepGroupImplViaProxy.numberOfGroups(), 2, "Should get correct total group count.");
  });

  it("should be able to get a group index number with provided group public key", async function() {
    assert.equal(await keepGroupImplViaProxy.getGroupIndex(groupOnePubKey), 0, "Should get correct group index number for group one.");
    assert.equal(await keepGroupImplViaProxy.getGroupIndex(groupTwoPubKey), 1, "Should get correct group index number for group two.");
  });

  it("should be able to get group public key by group index number", async function() {
    assert.equal(await keepGroupImplViaProxy.getGroupPubKey(0), groupOnePubKey, "Should get group public key.");
  });

  it("should be able to get staking weight", async function() {
    assert.equal(await keepGroupImplViaProxy.stakingWeight(account_one), 1, "Should have the staking weight of 1.");
    assert.equal(await keepGroupImplViaProxy.stakingWeight(account_two), 0, "Should have staking weight of 0.");
  });

  it("should be able to submit a ticket during initial ticket submission", async function() {

    let randomBeaconValue = 123456789;
    await keepGroupImplViaProxy.runGroupSelection(randomBeaconValue);

    await keepGroupImplViaProxy.submitTicket(1, 2, 3);

    let proof = await keepGroupImplViaProxy.getTicketProof(1);
    assert.equal(proof[0], 2, "Should be able to get submitted ticket proof.");
    assert.equal(proof[1], 3, "Should be able to get submitted ticket proof.");
  });

  it("should be able to submit a ticket during reactive ticket submission", async function() {

    let randomBeaconValue = 123456789;
    await keepGroupImplViaProxy.runGroupSelection(randomBeaconValue);

    // Mine one block
    web3.currentProvider.sendAsync({
      jsonrpc: "2.0",
      method: "evm_mine",
      id: 12345
    }, function(err, _) {
      if (err) console.log("Error mining a block.")
    });

    await keepGroupImplViaProxy.submitTicket(1, 2, 3);
    await keepGroupImplViaProxy.submitTicket(2, 2, 3);

    let proof = await keepGroupImplViaProxy.getTicketProof(2);
    assert.equal(proof[0], 2, "Should be able to get submitted ticket proof.");
    assert.equal(proof[1], 3, "Should be able to get submitted ticket proof.");

  });

  it("should not be able to submit a ticket during reactive ticket submission after enough tickets received", async function() {

    let randomBeaconValue = 123456789;
    await keepGroupImplViaProxy.runGroupSelection(randomBeaconValue);

    await keepGroupImplViaProxy.submitTicket(1, 2, 3);
    await keepGroupImplViaProxy.submitTicket(2, 2, 3);

    // Mine one block
    web3.currentProvider.sendAsync({
      jsonrpc: "2.0",
      method: "evm_mine",
      id: 12345
    }, function(err, _) {
      if (err) console.log("Error mining a block.")
    });

    await exceptThrow(keepGroupImplViaProxy.submitTicket(3, 2, 3));
  });

  it("should be able to verify a ticket", async function() {

    let randomBeaconValue = 123456789;
    await keepGroupImplViaProxy.runGroupSelection(randomBeaconValue);

    let stakerValue = account_one;
    let virtualStakerIndex = 1;

    let ticketValue = new BigNumber('0x' + abi.soliditySHA3(
      ["uint", "uint", "uint"],
      [randomBeaconValue, stakerValue, virtualStakerIndex]
    ).toString('hex'));

    assert.equal(await keepGroupImplViaProxy.cheapCheck(
      account_one, stakerValue, virtualStakerIndex
    ), true, "Should be able to verify a valid ticket.");
    
    assert.equal(await keepGroupImplViaProxy.costlyCheck(
      account_one, ticketValue, stakerValue, virtualStakerIndex
    ), true, "Should be able to verify a valid ticket.");
  
    assert.equal(await keepGroupImplViaProxy.costlyCheck(
      account_one, 1, stakerValue, virtualStakerIndex
    ), false, "Should fail verifying invalid ticket.");

  });
});
