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
const BN = require('bn.js');


function generateTickets(randomBeaconValue, stakerValue, stakerWeight) {
  let tickets = [];
  for (let i = 1; i <= stakerWeight; i++) {
    let ticketValue = new BigNumber('0x' + abi.soliditySHA3(
      ["uint", "uint", "uint"],
      [new BN(randomBeaconValue), stakerValue, i]
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
    previousEntry, randomBeaconValue, naturalThreshold,
    timeoutInitial, timeoutSubmission, timeoutChallenge,
    keepRandomBeaconImplV1, keepRandomBeaconProxy, keepRandomBeaconImplViaProxy,
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

    // Initialize Keep Random Beacon contract
    keepRandomBeaconImplV1 = await KeepRandomBeaconImplV1.new();
    keepRandomBeaconProxy = await KeepRandomBeaconProxy.new(keepRandomBeaconImplV1.address);
    keepRandomBeaconImplViaProxy = await KeepRandomBeaconImplV1.at(keepRandomBeaconProxy.address);
    await keepRandomBeaconImplViaProxy.initialize(1,1);

    // Initialize Keep Group contract
    minimumStake = 200000;
    groupThreshold = 15;
    groupSize = 20;
    timeoutInitial = 20;
    timeoutSubmission = 40;
    timeoutChallenge = 60;

    previousEntry = '61647192470559497961835987041772849466909869278998228659325527399967350471301';

    keepGroupImplV1 = await KeepGroupImplV1.new();
    keepGroupProxy = await KeepGroupProxy.new(keepGroupImplV1.address);
    keepGroupImplViaProxy = await KeepGroupImplV1.at(keepGroupProxy.address);
    await keepGroupImplViaProxy.initialize(
      stakingProxy.address, keepRandomBeaconProxy.address, minimumStake, groupThreshold, groupSize, timeoutInitial, timeoutSubmission, timeoutChallenge
    );

    naturalThreshold = await keepGroupImplViaProxy.naturalThreshold();

    // Data generated using client Go code with master secret key 123
    groupPubKey = "0x1f1954b33144db2b5c90da089e8bde287ec7089d5d6433f3b6becaefdb678b1b2a9de38d14bef2cf9afc3c698a4211fa7ada7b4f036a2dfef0dc122b423259d0";
    randomBeaconValue = '14338479826386391162410575915611605755963697664966940200639296307898208140504';

    // Stake tokens as account one so it has minimum stake to be able to get into a group.
    await token.approveAndCall(stakingContract.address, minimumStake*1000, "", {from: staker1});
    tickets1 = generateTickets(randomBeaconValue, staker1, 1000);

    // Send tokens to staker2 and stake
    await token.transfer(staker2, minimumStake*2000, {from: staker1});
    await token.approveAndCall(stakingContract.address, minimumStake*2000, "", {from: staker2});
    tickets2 = generateTickets(randomBeaconValue, staker2, 2000);

    // Send tokens to staker3 and stake
    await token.transfer(staker3, minimumStake*3000, {from: staker1});
    await token.approveAndCall(stakingContract.address, minimumStake*3000, "", {from: staker3});
    tickets3 = generateTickets(randomBeaconValue, staker3, 3000);

    await keepRandomBeaconImplViaProxy.setGroupContract(keepGroupProxy.address);
    await keepRandomBeaconImplViaProxy.relayEntry(1, new BigNumber(randomBeaconValue), groupPubKey, new BigNumber(previousEntry));
  });

  it("should be able to get staking weight", async function() {
    assert.equal(await keepGroupImplViaProxy.stakingWeight(staker1), 1000, "Should have expected staking weight.");
    assert.equal(await keepGroupImplViaProxy.stakingWeight(staker3), 3000, "Should have expected staking weight.");
  });

  it("should fail to get selected tickets before challenge period is over", async function() {
    await exceptThrow(keepGroupImplViaProxy.selectedTickets());
  });

  it("should be able to get selected tickets after challenge period is over", async function() {

    for (let i = 0; i < groupSize*2; i++) {
      await keepGroupImplViaProxy.submitTicket(tickets1[i].value, staker1, tickets1[i].virtualStakerIndex, {from: staker1});
    }

    mineBlocks(timeoutChallenge);
    let selectedTickets = await keepGroupImplViaProxy.selectedTickets();
    assert.equal(selectedTickets.length, groupSize, "Should be trimmed to groupSize length.");
  });

  it("should be able to output submited tickets in ascending ordered", async function() {

    let tickets = [];

    await keepGroupImplViaProxy.submitTicket(tickets1[0].value, staker1, tickets1[0].virtualStakerIndex);
    tickets.push(tickets1[0].value);

    await keepGroupImplViaProxy.submitTicket(tickets2[0].value, staker2, tickets2[0].virtualStakerIndex, {from: staker2});
    tickets.push(tickets2[0].value);

    await keepGroupImplViaProxy.submitTicket(tickets3[0].value, staker3, tickets3[0].virtualStakerIndex, {from: staker3});
    tickets.push(tickets3[0].value);

    tickets = tickets.sort(function(a, b){return a-b}); // Sort numbers in ascending order

    // Test tickets ordering
    let orderedTickets = await keepGroupImplViaProxy.orderedTickets();
    assert.equal(orderedTickets[0].equals(tickets[0]), true, "Tickets should be in ascending order.");
    assert.equal(orderedTickets[1].equals(tickets[1]), true, "Tickets should be in ascending order.");
    assert.equal(orderedTickets[2].equals(tickets[2]), true, "Tickets should be in ascending order.");

  });

  it("should be able to submit a ticket during ticket submission period", async function() {
    await keepGroupImplViaProxy.submitTicket(tickets1[0].value, staker1, tickets1[0].virtualStakerIndex);
    let proof = await keepGroupImplViaProxy.getTicketProof(tickets1[0].value);
    assert.equal(proof[1].equals(new BigNumber(staker1)), true , "Should be able to get submitted ticket proof.");
    assert.equal(proof[2], tickets1[0].virtualStakerIndex, "Should be able to get submitted ticket proof.");
  });

  it("should be able to verify a ticket", async function() {

    await keepGroupImplViaProxy.submitTicket(tickets1[0].value, staker1, 1);

    assert.equal(await keepGroupImplViaProxy.cheapCheck(
      staker1, staker1, 1
    ), true, "Should be able to verify a valid ticket.");
    
    assert.equal(await keepGroupImplViaProxy.costlyCheck(
      staker1, tickets1[0].value, staker1, tickets1[0].virtualStakerIndex
    ), true, "Should be able to verify a valid ticket.");
  
    assert.equal(await keepGroupImplViaProxy.costlyCheck(
      staker1, 0, staker1, tickets1[0].virtualStakerIndex
    ), false, "Should fail verifying invalid ticket.");

  });
});
