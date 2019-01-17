import { duration } from './helpers/increaseTime';
import exceptThrow from './helpers/expectThrow';
import BigNumber from 'bignumber.js';
import abi from 'ethereumjs-abi';
const KeepToken = artifacts.require('./KeepToken.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const KeepRandomBeaconProxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconImplV1 = artifacts.require('./KeepRandomBeaconImplV1.sol');
const KeepGroupProxy = artifacts.require('./KeepGroup.sol');
const KeepGroupImplV1 = artifacts.require('./KeepGroupImplV1.sol');


function generateTickets(randomBeaconValue, stakerValue, stakerWeight) {
  let tickets = [];
  for (let i = 1; i <= stakerWeight; i++) {
    let ticketValue = new BigNumber('0x' + abi.soliditySHA3(
      ["uint", "uint", "uint"],
      [randomBeaconValue, stakerValue, i]
    ).toString('hex'));
    let ticket = {
      value: ticketValue,
      virtualStakerIndex: i
    }
    tickets.push(ticket);
  }
  return tickets
}

function mineBlocks(blocks) {
  for (let i = 0; i <= blocks; i++) {
    web3.currentProvider.sendAsync({
      jsonrpc: "2.0",
      method: "evm_mine",
      id: 12345
    }, function(err, _) {
      if (err) console.log("Error mining a block.")
    });
  }
}

contract('TestKeepGroupSelection', function(accounts) {

  let token, stakingProxy, stakingContract, minimumStake, groupThreshold, groupSize,
    randomBeaconValue, naturalThreshold,
    timeoutInitial, timeoutSubmission, timeoutChallenge,
    keepRandomBeaconImplV1, keepRandomBeaconProxy, keepRandomBeaconImplViaProxy,
    keepGroupImplV1, keepGroupProxy, keepGroupImplViaProxy, groupPubKey,
    staker1 = accounts[0], tickets1, tickets1BelowNatT, tickets1AboveNatT,
    staker2 = accounts[1], tickets2, tickets2BelowNatT, tickets2AboveNatT,
    staker3 = accounts[2], tickets3, tickets3BelowNatT, tickets3AboveNatT,
    staker4 = accounts[3], tickets4, tickets4BelowNatT, tickets4AboveNatT;

  beforeEach(async () => {
    token = await KeepToken.new();
    
    // Initialize staking contract under proxy
    stakingProxy = await StakingProxy.new();
    stakingContract = await TokenStaking.new(token.address, stakingProxy.address, duration.days(30));
    await stakingProxy.authorizeContract(stakingContract.address, {from: staker1})

    // Initialize Keep Random Beacon contract
    keepRandomBeaconImplV1 = await KeepRandomBeaconImplV1.new();
    keepRandomBeaconProxy = await KeepRandomBeaconProxy.new(keepRandomBeaconImplV1.address);
    keepRandomBeaconImplViaProxy = await KeepRandomBeaconImplV1.at(keepRandomBeaconProxy.address);
    await keepRandomBeaconImplViaProxy.initialize(1,1);

    // Initialize Keep Group contract
    minimumStake = 20000000;
    groupThreshold = 15;
    groupSize = 20;
    timeoutInitial = 20;
    timeoutSubmission = 40;
    timeoutChallenge = 60;

    randomBeaconValue = 123456789;

    keepGroupImplV1 = await KeepGroupImplV1.new();
    keepGroupProxy = await KeepGroupProxy.new(keepGroupImplV1.address);
    keepGroupImplViaProxy = await KeepGroupImplV1.at(keepGroupProxy.address);
    await keepGroupImplViaProxy.initialize(
      stakingProxy.address, keepRandomBeaconProxy.address, minimumStake, groupThreshold, groupSize, timeoutInitial, timeoutSubmission, timeoutChallenge
    );

    naturalThreshold = await keepGroupImplViaProxy.naturalThreshold();

    groupPubKey = "0x1000000000000000000000000000000000000000000000000000000000000000";

    // Stake tokens as account one so it has minimum stake to be able to get into a group.
    await token.approveAndCall(stakingContract.address, minimumStake*100, "", {from: staker1});
    tickets1 = generateTickets(randomBeaconValue, staker1, 100);
    tickets1BelowNatT = tickets1.filter(function(ticket) {
      return ticket.value.lessThan(naturalThreshold);
    });
    tickets1AboveNatT = tickets1.filter(function(ticket) {
      return ticket.value.greaterThan(naturalThreshold);
    });

    // Send tokens to staker2 and stake
    await token.transfer(staker2, minimumStake*200, {from: staker1});
    await token.approveAndCall(stakingContract.address, minimumStake*200, "", {from: staker2});
    tickets2 = generateTickets(randomBeaconValue, staker2, 200);
    tickets2BelowNatT = tickets2.filter(function(ticket) {
      return ticket.value.lessThan(naturalThreshold);
    });
    tickets2AboveNatT = tickets2.filter(function(ticket) {
      return ticket.value.greaterThan(naturalThreshold);
    });

    // Send tokens to staker3 and stake
    await token.transfer(staker3, minimumStake*3000, {from: staker1});
    await token.approveAndCall(stakingContract.address, minimumStake*3000, "", {from: staker3});
    tickets3 = generateTickets(randomBeaconValue, staker3, 3000);
    tickets3BelowNatT = tickets3.filter(function(ticket) {
      return ticket.value.lessThan(naturalThreshold);
    });
    tickets3AboveNatT = tickets3.filter(function(ticket) {
      return ticket.value.greaterThan(naturalThreshold);
    });

    await keepRandomBeaconImplViaProxy.setGroupContract(keepGroupProxy.address);
    await keepRandomBeaconImplViaProxy.relayEntry(1, randomBeaconValue, 1, 1);
  });

  it("should be able to get staking weight", async function() {
    assert.equal(await keepGroupImplViaProxy.stakingWeight(staker1), 100, "Should have expected staking weight.");
    assert.equal(await keepGroupImplViaProxy.stakingWeight(staker3), 3000, "Should have expected staking weight.");
  });

  it("should be able to output submited tickets in ascending ordered", async function() {

    let tickets = [];

    await keepGroupImplViaProxy.submitTicket(tickets1BelowNatT[0].value, staker1, tickets1BelowNatT[0].virtualStakerIndex);
    tickets.push(tickets1BelowNatT[0].value);

    await keepGroupImplViaProxy.submitTicket(tickets2BelowNatT[0].value, staker2, tickets2BelowNatT[0].virtualStakerIndex, {from: staker2});
    tickets.push(tickets2BelowNatT[0].value);

    await keepGroupImplViaProxy.submitTicket(tickets3BelowNatT[0].value, staker3, tickets3BelowNatT[0].virtualStakerIndex, {from: staker3});
    tickets.push(tickets3BelowNatT[0].value);

    tickets = tickets.sort(function(a, b){return a-b}); // Sort numbers in ascending order

    // Test tickets ordering
    let orderedTickets = await keepGroupImplViaProxy.orderedTickets();
    assert.equal(orderedTickets[0].equals(tickets[0]), true, "Tickets should be in ascending order.");
    assert.equal(orderedTickets[1].equals(tickets[1]), true, "Tickets should be in ascending order.");
    assert.equal(orderedTickets[2].equals(tickets[2]), true, "Tickets should be in ascending order.");

  });

  it("should be able to submit a ticket during initial ticket submission", async function() {
    await keepGroupImplViaProxy.submitTicket(tickets1BelowNatT[0].value, staker1, tickets1BelowNatT[0].virtualStakerIndex);
    let proof = await keepGroupImplViaProxy.getTicketProof(tickets1BelowNatT[0].value);
    assert.equal(proof[1].equals(new BigNumber(staker1)), true , "Should be able to get submitted ticket proof.");
    assert.equal(proof[2], tickets1BelowNatT[0].virtualStakerIndex, "Should be able to get submitted ticket proof.");
  });

  it("should be able to submit a high value ticket during reactive ticket submission", async function() {
    mineBlocks(timeoutInitial);
    await keepGroupImplViaProxy.submitTicket(tickets1AboveNatT[0].value, staker1, tickets1AboveNatT[0].virtualStakerIndex);
    let proof = await keepGroupImplViaProxy.getTicketProof(tickets1AboveNatT[0].value);
    assert.equal(proof[1].equals(new BigNumber(staker1)), true , "Should be able to get submitted ticket proof.");
    assert.equal(proof[2], tickets1AboveNatT[0].virtualStakerIndex, "Should be able to get submitted ticket proof.");
  });

  it("should be able to verify a ticket", async function() {

    await keepGroupImplViaProxy.submitTicket(tickets1BelowNatT[0].value, staker1, 1);

    assert.equal(await keepGroupImplViaProxy.cheapCheck(
      staker1, staker1, 1
    ), true, "Should be able to verify a valid ticket.");
    
    assert.equal(await keepGroupImplViaProxy.costlyCheck(
      staker1, tickets1BelowNatT[0].value, staker1, tickets1BelowNatT[0].virtualStakerIndex
    ), true, "Should be able to verify a valid ticket.");
  
    assert.equal(await keepGroupImplViaProxy.costlyCheck(
      staker1, 0, staker1, tickets1BelowNatT[0].virtualStakerIndex
    ), false, "Should fail verifying invalid ticket.");

  });

  it("should be able to challenge a ticket", async function() {

    // TODO: replace with a secure authorization protocol (addressed in RFC 4).
    await keepGroupImplViaProxy.authorizeStakingContract(stakingContract.address);

    await keepGroupImplViaProxy.submitTicket(tickets1BelowNatT[0].value, staker1, tickets1BelowNatT[0].virtualStakerIndex);
    await keepGroupImplViaProxy.submitTicket(1, staker1, tickets1BelowNatT[1].virtualStakerIndex); // invalid ticket

    // Challenging valid ticket
    let previousBalance = await stakingContract.stakeBalanceOf(staker2);
    await keepGroupImplViaProxy.challenge(tickets1BelowNatT[0].value, {from: staker2});
    //assert.equal(await stakingContract.stakeBalanceOf(staker2), previousBalance.toNumber() - minimumStake, "Should result slashing challenger's balance");

    // Challenging invalid ticket
    previousBalance = await stakingContract.stakeBalanceOf(staker2);
    await keepGroupImplViaProxy.challenge(1, {from: staker2});
    //assert.equal(await stakingContract.stakeBalanceOf(staker2), previousBalance.toNumber() + minimumStake, "Should result rewarding challenger's balance");

  });

});
