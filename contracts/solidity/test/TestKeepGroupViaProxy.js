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

  let token, stakingProxy, stakingContract, minimumStake,
    keepRandomBeaconImplV1, keepRandomBeaconProxy, keepRandomBeaconImplViaProxy,
    keepGroupImplV1, keepGroupProxy, keepGroupImplViaProxy, groupOnePubKey, groupTwoPubKey,
    account_one = accounts[0],
    account_two = accounts[1],
    account_three = accounts[2],
    account_four = accounts[3];

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
    minimumStake = 200;
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

    // Send tokens to account_two and stake
    await token.transfer(account_two, minimumStake, {from: account_one});
    await token.approveAndCall(stakingContract.address, minimumStake, "", {from: account_two});

    // Send tokens to account_three and stake
    await token.transfer(account_three, minimumStake*3, {from: account_one});
    await token.approveAndCall(stakingContract.address, minimumStake*3, "", {from: account_three});

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
    assert.equal(await keepGroupImplViaProxy.stakingWeight(account_three), 3, "Should have staking weight of 3.");
  });

  it("should be able to submit a ticket during initial ticket submission", async function() {

    let randomBeaconValue = 123456789;
    await keepGroupImplViaProxy.runGroupSelection(randomBeaconValue);

    let stakerValue = account_one;
    let virtualStakerIndex = 1;
    let tickets = [];

    let ticketValue = new BigNumber('0x' + abi.soliditySHA3(
      ["uint", "uint", "uint"],
      [randomBeaconValue, stakerValue, virtualStakerIndex]
    ).toString('hex'));

    await keepGroupImplViaProxy.submitTicket(ticketValue, stakerValue, virtualStakerIndex);
    tickets.push(ticketValue);

    stakerValue = account_two;
    ticketValue = new BigNumber('0x' + abi.soliditySHA3(
      ["uint", "uint", "uint"],
      [randomBeaconValue, stakerValue, virtualStakerIndex]
    ).toString('hex'));

    await keepGroupImplViaProxy.submitTicket(ticketValue, stakerValue, virtualStakerIndex, {from: account_two});
    tickets.push(ticketValue);

    stakerValue = account_three;
    ticketValue = new BigNumber('0x' + abi.soliditySHA3(
      ["uint", "uint", "uint"],
      [randomBeaconValue, stakerValue, virtualStakerIndex]
    ).toString('hex'));

    await keepGroupImplViaProxy.submitTicket(ticketValue, stakerValue, virtualStakerIndex, {from: account_three});
    tickets.push(ticketValue);

    tickets = tickets.sort(function(a, b){return a-b}); // Sort numbers in ascending order

    // Test tickets ordering
    let orderedTickets = await keepGroupImplViaProxy.orderedTickets();
    assert.equal(orderedTickets[0].equals(tickets[0]), true, "Tickets should be in ascending order.");
    assert.equal(orderedTickets[1].equals(tickets[1]), true, "Tickets should be in ascending order.");
    assert.equal(orderedTickets[2].equals(tickets[2]), true, "Tickets should be in ascending order.");

    // Test can't submit group pubkey if haven't submitted a ticket
    await exceptThrow(keepGroupImplViaProxy.submitGroupPublicKey(groupOnePubKey, {from: account_four}));

    // Test submit group pubkey
    await keepGroupImplViaProxy.submitGroupPublicKey(groupOnePubKey, {from: account_one});

    // Test vote for submission of the group key
    await keepGroupImplViaProxy.voteForSubmission(groupOnePubKey, {from: account_two});

    // Test group is selected
    await keepGroupImplViaProxy.getFinalResult();

    let proof = await keepGroupImplViaProxy.getTicketProof(ticketValue);
    assert.equal(proof[0].equals(new BigNumber(stakerValue)), true , "Should be able to get submitted ticket proof.");
    assert.equal(proof[1], virtualStakerIndex, "Should be able to get submitted ticket proof.");
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

    let stakerValue = account_one;
    let virtualStakerIndex = 1;
    let ticketValue = new BigNumber('0x' + abi.soliditySHA3(
      ["uint", "uint", "uint"],
      [randomBeaconValue, stakerValue, virtualStakerIndex]
    ).toString('hex'));
    await keepGroupImplViaProxy.submitTicket(ticketValue, stakerValue, virtualStakerIndex);

    stakerValue = account_two;
    ticketValue = new BigNumber('0x' + abi.soliditySHA3(
      ["uint", "uint", "uint"],
      [randomBeaconValue, stakerValue, virtualStakerIndex]
    ).toString('hex'));
    await keepGroupImplViaProxy.submitTicket(ticketValue, stakerValue, virtualStakerIndex, {from: account_two});

    let proof = await keepGroupImplViaProxy.getTicketProof(ticketValue);
    assert.equal(proof[0].equals(new BigNumber(stakerValue)), true , "Should be able to get submitted ticket proof.");
    assert.equal(proof[1], virtualStakerIndex, "Should be able to get submitted ticket proof.");

  });

  it("should not be able to submit a ticket during reactive ticket submission after enough tickets received", async function() {

    let randomBeaconValue = 123456789;
    await keepGroupImplViaProxy.runGroupSelection(randomBeaconValue);

    let stakerValue = account_one;
    let virtualStakerIndex = 1;
    let ticketValue = new BigNumber('0x' + abi.soliditySHA3(
      ["uint", "uint", "uint"],
      [randomBeaconValue, stakerValue, virtualStakerIndex]
    ).toString('hex'));
    await keepGroupImplViaProxy.submitTicket(ticketValue, stakerValue, virtualStakerIndex);

    stakerValue = account_two;
    ticketValue = new BigNumber('0x' + abi.soliditySHA3(
      ["uint", "uint", "uint"],
      [randomBeaconValue, stakerValue, virtualStakerIndex]
    ).toString('hex'));
    await keepGroupImplViaProxy.submitTicket(ticketValue, stakerValue, virtualStakerIndex, {from: account_two});

    // Mine one block
    web3.currentProvider.sendAsync({
      jsonrpc: "2.0",
      method: "evm_mine",
      id: 12345
    }, function(err, _) {
      if (err) console.log("Error mining a block.")
    });

    stakerValue = account_one;
    virtualStakerIndex = 2;
    ticketValue = new BigNumber('0x' + abi.soliditySHA3(
      ["uint", "uint", "uint"],
      [randomBeaconValue, stakerValue, virtualStakerIndex]
    ).toString('hex'));
    await exceptThrow(keepGroupImplViaProxy.submitTicket(ticketValue, stakerValue, virtualStakerIndex));
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

  it("should be able to challenge a ticket", async function() {

    // TODO: replace with a secure authorization protocol (addressed in RFC 4).
    await keepGroupImplViaProxy.authorizeStakingContract(stakingContract.address);

    let randomBeaconValue = 123456789;
    await keepGroupImplViaProxy.runGroupSelection(randomBeaconValue);

    // Submit tickets as account_one
    let stakerValue = account_one;
    let virtualStakerIndex = 1;
    let ticketValue = new BigNumber('0x' + abi.soliditySHA3(
      ["uint", "uint", "uint"],
      [randomBeaconValue, stakerValue, virtualStakerIndex]
    ).toString('hex'));

    await keepGroupImplViaProxy.submitTicket(ticketValue, stakerValue, virtualStakerIndex);
    await keepGroupImplViaProxy.submitTicket(1, stakerValue, virtualStakerIndex); // invalid ticket

    // Challenging valid ticket
    let previousBalance = await stakingContract.stakeBalanceOf(account_two);
    await keepGroupImplViaProxy.challenge(ticketValue, {from: account_two});
    assert.equal(await stakingContract.stakeBalanceOf(account_two), previousBalance.toNumber() - minimumStake, "Should result slashing challenger's balance");

    // Challenging invalid ticket
    previousBalance = await stakingContract.stakeBalanceOf(account_two);
    await keepGroupImplViaProxy.challenge(1, {from: account_two});
    assert.equal(await stakingContract.stakeBalanceOf(account_two), previousBalance.toNumber() + minimumStake, "Should result rewarding challenger's balance");

  });
});
