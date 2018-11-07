import increaseTime, { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import exceptThrow from './helpers/expectThrow';
import encodeCall from './helpers/encodeCall';
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
    await keepGroupImplViaProxy.initialize(stakingProxy.address, minimumStake, 6, 10, 1, 1, 1);

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

  it("should be able to submit a ticket within initial timeout", async function() {

    await keepGroupImplViaProxy.runGroupSelection();

    assert.equal(await keepGroupImplViaProxy.submitTicket(), true, "Should be able to submit ticket.");

    // Mine one block
    web3.currentProvider.sendAsync({
      jsonrpc: "2.0",
      method: "evm_mine",
      id: 12345
    }, function(err, _) {
      if (err) console.log("Error mining a block.")
    });

    assert.equal(await keepGroupImplViaProxy.submitTicket(), false, "Should not be able to submit ticket after initial timeout is reached.");

  });
});
