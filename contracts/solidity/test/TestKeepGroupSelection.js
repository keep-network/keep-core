import { duration } from './helpers/increaseTime';
import exceptThrow from './helpers/expectThrow';
import BigNumber from 'bignumber.js';
import abi from 'ethereumjs-abi';
const KeepToken = artifacts.require('./KeepToken.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const KeepRandomBeaconProxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepGroupProxy = artifacts.require('./KeepGroup.sol');
const KeepGroupImplV1 = artifacts.require('./KeepGroupImplV1.sol');


function generateTickets(randomBeaconValue, stakerValue, stakerWeight) {
  let tickets = [];
  for (let i = 1; i <= stakerWeight; i++) {
    let ticketValue = new BigNumber('0x' + abi.soliditySHA3(
      ["uint", "uint", "uint"],
      [randomBeaconValue, stakerValue, i]
    ).toString('hex'));
    tickets.push(ticketValue);
  }
  return tickets
}

contract('TestKeepGroupSelection', function(accounts) {

  let token, stakingProxy, stakingContract, minimumStake, groupSize,
    randomBeaconValue = 123456789,
    keepGroupImplV1, keepGroupProxy, keepGroupImplViaProxy, groupPubKey,
    staker1 = accounts[0], tickets1,
    staker2 = accounts[1], tickets2,
    staker3 = accounts[2], tickets3,
    staker4 = accounts[3], tickets4;

  beforeEach(async () => {
    token = await KeepToken.new();
    
    // Initialize staking contract under proxy
    stakingProxy = await StakingProxy.new();
    stakingContract = await TokenStaking.new(token.address, stakingProxy.address, duration.days(30));
    await stakingProxy.authorizeContract(stakingContract.address, {from: staker1})

    // Initialize Keep Group contract
    minimumStake = 200;
    groupSize = 200;

    keepGroupImplV1 = await KeepGroupImplV1.new();
    keepGroupProxy = await KeepGroupProxy.new(keepGroupImplV1.address);
    keepGroupImplViaProxy = await KeepGroupImplV1.at(keepGroupProxy.address);
    await keepGroupImplViaProxy.initialize(stakingProxy.address, minimumStake, groupSize, 1, 3, 4);

    groupPubKey = "0x1000000000000000000000000000000000000000000000000000000000000000";

    // Stake tokens as account one so it has minimum stake to be able to get into a group.
    await token.approveAndCall(stakingContract.address, minimumStake*2, "", {from: staker1});
    tickets1 = generateTickets(randomBeaconValue, staker1, 2);

    // Send tokens to staker2 and stake
    await token.transfer(staker2, minimumStake, {from: staker1});
    await token.approveAndCall(stakingContract.address, minimumStake, "", {from: staker2});
    tickets2 = generateTickets(randomBeaconValue, staker2, 1);

    // Send tokens to staker3 and stake
    await token.transfer(staker3, minimumStake*3, {from: staker1});
    await token.approveAndCall(stakingContract.address, minimumStake*3, "", {from: staker3});
    tickets3 = generateTickets(randomBeaconValue, staker3, 3);

    await keepGroupImplViaProxy.runGroupSelection(randomBeaconValue);
  });

  it("should be able to get staking weight", async function() {
    assert.equal(await keepGroupImplViaProxy.stakingWeight(staker1), 2, "Should have the staking weight of 1.");
    assert.equal(await keepGroupImplViaProxy.stakingWeight(staker3), 3, "Should have staking weight of 3.");
  });

  it("should be able to submit a ticket during initial ticket submission", async function() {

    let tickets = [];

    await keepGroupImplViaProxy.submitTicket(tickets1[0], staker1, 1);
    tickets.push(tickets1[0]);

    await keepGroupImplViaProxy.submitTicket(tickets2[0], staker2, 1, {from: staker2});
    tickets.push(tickets2[0]);

    await keepGroupImplViaProxy.submitTicket(tickets3[0], staker3, 1, {from: staker3});
    tickets.push(tickets3[0]);

    tickets = tickets.sort(function(a, b){return a-b}); // Sort numbers in ascending order

    // Test getting ticket proof
    let proof = await keepGroupImplViaProxy.getTicketProof(tickets1[0]);
    assert.equal(proof[0].equals(new BigNumber(staker1)), true , "Should be able to get submitted ticket proof.");
    assert.equal(proof[1], 1, "Should be able to get submitted ticket proof.");

    // Test tickets ordering
    let orderedTickets = await keepGroupImplViaProxy.orderedTickets();
    assert.equal(orderedTickets[0].equals(tickets[0]), true, "Tickets should be in ascending order.");
    assert.equal(orderedTickets[1].equals(tickets[1]), true, "Tickets should be in ascending order.");
    assert.equal(orderedTickets[2].equals(tickets[2]), true, "Tickets should be in ascending order.");

    // Test can't submit group pubkey if haven't submitted a ticket
    await exceptThrow(keepGroupImplViaProxy.submitGroupPublicKey(groupPubKey, {from: staker4}));

    // Test submit group pubkey
    await keepGroupImplViaProxy.submitGroupPublicKey(groupPubKey, {from: staker1});

    // Test vote for submission of the group key
    await keepGroupImplViaProxy.voteForSubmission(groupPubKey, {from: staker2});

    // Test group is selected
    await keepGroupImplViaProxy.getFinalResult();

  });

  it("should be able to submit a ticket during reactive ticket submission", async function() {

    // Mine one block
    web3.currentProvider.sendAsync({
      jsonrpc: "2.0",
      method: "evm_mine",
      id: 12345
    }, function(err, _) {
      if (err) console.log("Error mining a block.")
    });

    await keepGroupImplViaProxy.submitTicket(tickets1[0], staker1, 1);
    await keepGroupImplViaProxy.submitTicket(tickets2[0], staker2, 1, {from: staker2});

    let proof = await keepGroupImplViaProxy.getTicketProof(tickets2[0]);
    assert.equal(proof[0].equals(new BigNumber(staker2)), true , "Should be able to get submitted ticket proof.");
    assert.equal(proof[1], 1, "Should be able to get submitted ticket proof.");

  });

  it("should not be able to submit a ticket during reactive ticket submission after enough tickets received", async function() {

    await keepGroupImplViaProxy.submitTicket(tickets1[0], staker1, 1);
    await keepGroupImplViaProxy.submitTicket(tickets2[0], staker2, 1, {from: staker2});

    // Mine one block
    web3.currentProvider.sendAsync({
      jsonrpc: "2.0",
      method: "evm_mine",
      id: 12345
    }, function(err, _) {
      if (err) console.log("Error mining a block.")
    });

    await exceptThrow(keepGroupImplViaProxy.submitTicket(tickets1[1], staker1, 2));
  });

  it("should be able to verify a ticket", async function() {

    await keepGroupImplViaProxy.submitTicket(tickets1[0], staker1, 1);

    assert.equal(await keepGroupImplViaProxy.cheapCheck(
      staker1, staker1, 1
    ), true, "Should be able to verify a valid ticket.");
    
    assert.equal(await keepGroupImplViaProxy.costlyCheck(
      staker1, tickets1[0], staker1, 1
    ), true, "Should be able to verify a valid ticket.");
  
    assert.equal(await keepGroupImplViaProxy.costlyCheck(
      staker1, 0, staker1, 1
    ), false, "Should fail verifying invalid ticket.");

  });

  it("should be able to challenge a ticket", async function() {

    // TODO: replace with a secure authorization protocol (addressed in RFC 4).
    await keepGroupImplViaProxy.authorizeStakingContract(stakingContract.address);

    await keepGroupImplViaProxy.submitTicket(tickets1[0], staker1, 1);
    await keepGroupImplViaProxy.submitTicket(1, staker1, 2); // invalid ticket

    // Challenging valid ticket
    let previousBalance = await stakingContract.stakeBalanceOf(staker2);
    await keepGroupImplViaProxy.challenge(tickets1[0], {from: staker2});
    //assert.equal(await stakingContract.stakeBalanceOf(staker2), previousBalance.toNumber() - minimumStake, "Should result slashing challenger's balance");

    // Challenging invalid ticket
    previousBalance = await stakingContract.stakeBalanceOf(staker2);
    await keepGroupImplViaProxy.challenge(1, {from: staker2});
    //assert.equal(await stakingContract.stakeBalanceOf(staker2), previousBalance.toNumber() + minimumStake, "Should result rewarding challenger's balance");

  });
});
