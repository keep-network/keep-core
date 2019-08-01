import expectThrow from './helpers/expectThrow';
import mineBlocks from './helpers/mineBlocks';
import {bls} from './helpers/data';
import generateTickets from './helpers/generateTickets';
import stakeDelegate from './helpers/stakeDelegate';
import {initContracts} from './helpers/initContracts';


contract('TestKeepRandomBeaconOperatorGroupSelection', function(accounts) {

  let token, stakingContract, serviceContract, operatorContract,
  owner = accounts[0], magpie = accounts[1],
  operator1 = accounts[2], tickets1,
  operator2 = accounts[3], tickets2,
  operator3 = accounts[4], tickets3;

  const minimumStake = web3.utils.toBN(200000);
  const ticketInitialSubmissionTimeout = 20;
  const ticketReactiveSubmissionTimeout = 100;

  before(async () => {

    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./KeepRandomBeaconOperatorStub.sol')
    );
    
    token = contracts.token;
    serviceContract = contracts.serviceContract;
    operatorContract = contracts.operatorContract;
    stakingContract = contracts.stakingContract;

    operatorContract.setMinimumStake(minimumStake)
    operatorContract.setTicketInitialSubmissionTimeout(ticketInitialSubmissionTimeout);
    operatorContract.setTicketReactiveSubmissionTimeout(ticketReactiveSubmissionTimeout);

    await stakeDelegate(stakingContract, token, owner, operator1, magpie, minimumStake.mul(web3.utils.toBN(2000)))
    await stakeDelegate(stakingContract, token, owner, operator2, magpie, minimumStake.mul(web3.utils.toBN(2000)))
    await stakeDelegate(stakingContract, token, owner, operator3, magpie, minimumStake.mul(web3.utils.toBN(3000)))

    tickets1 = generateTickets(await operatorContract.groupSelectionRelayEntry(), operator1, 2000);
    tickets2 = generateTickets(await operatorContract.groupSelectionRelayEntry(), operator2, 2000);
    tickets3 = generateTickets(await operatorContract.groupSelectionRelayEntry(), operator3, 3000);

    // Using stub method to add first group to help testing.
    await operatorContract.registerNewGroup(bls.groupPubKey);

  });

  it("should be able to get staking weight", async function() {
    assert.isTrue(web3.utils.toBN(2000).eq(await operatorContract.stakingWeight(operator1)), "Should have expected staking weight.");
    assert.isTrue(web3.utils.toBN(3000).eq(await operatorContract.stakingWeight(operator3)), "Should have expected staking weight.");
  });

  it("should fail to get selected tickets before challenge period is over", async function() {
    await expectThrow(operatorContract.selectedTickets());
  });

  it("should fail to get selected participants before challenge period is over", async function() {
    await expectThrow(operatorContract.selectedParticipants());
  });

  it("should be able to output submited tickets in ascending ordered", async function() {

    let tickets = [];

    await operatorContract.submitTicket(tickets1[0].value, operator1, tickets1[0].virtualStakerIndex, {from: operator1});
    tickets.push(tickets1[0].value);

    await operatorContract.submitTicket(tickets2[0].value, operator2, tickets2[0].virtualStakerIndex, {from: operator2});
    tickets.push(tickets2[0].value);

    await operatorContract.submitTicket(tickets3[0].value, operator3, tickets3[0].virtualStakerIndex, {from: operator3});
    tickets.push(tickets3[0].value);

    tickets = tickets.sort(function(a, b){return a-b}); // Sort numbers in ascending order

    // Test tickets ordering
    let orderedTickets = await operatorContract.orderedTickets();
    assert.isTrue(orderedTickets[0].eq(tickets[0]), "Tickets should be in ascending order.");
    assert.isTrue(orderedTickets[1].eq(tickets[1]), "Tickets should be in ascending order.");
    assert.isTrue(orderedTickets[2].eq(tickets[2]), "Tickets should be in ascending order.");

  });

  it("should be able to submit a ticket during ticket submission period", async function() {
    await operatorContract.submitTicket(tickets1[0].value, operator1, tickets1[0].virtualStakerIndex, {from: operator1});
    let proof = await operatorContract.getTicketProof(tickets1[0].value);
    assert.isTrue(proof[1].eq(web3.utils.toBN(operator1)), "Should be able to get submitted ticket proof.");
    assert.equal(proof[2], tickets1[0].virtualStakerIndex, "Should be able to get submitted ticket proof.");
  });

  it("should be able to verify a ticket", async function() {

    await operatorContract.submitTicket(tickets1[0].value, operator1, 1, {from: operator1});

    assert.isTrue(await operatorContract.cheapCheck(
      operator1, operator1, 1
    ), "Should be able to verify a valid ticket.");
    
    assert.isTrue(await operatorContract.costlyCheck(
      operator1, tickets1[0].value, operator1, tickets1[0].virtualStakerIndex
    ), "Should be able to verify a valid ticket.");
  
    assert.isFalse(await operatorContract.costlyCheck(
      operator1, 0, operator1, tickets1[0].virtualStakerIndex
    ), "Should fail verifying invalid ticket.");

  });

  it("should not trigger group selection while one is in progress", async function() {
    let groupSelectionStartBlock = await operatorContract.ticketSubmissionStartBlock();
    await serviceContract.requestRelayEntry(bls.seed, {value: 10});
    await operatorContract.relayEntry(bls.nextGroupSignature);

    assert.isTrue((await operatorContract.ticketSubmissionStartBlock()).eq(groupSelectionStartBlock), "Group selection start block should not be updated.");
    assert.isTrue((await operatorContract.groupSelectionRelayEntry()).eq(bls.groupSignature), "Random beacon value for the current group selection should not change.");
  });

  it("should be able to get selected tickets and participants after challenge period is over", async function() {

    let groupSize = await operatorContract.groupSize();

    for (let i = 0; i < groupSize*2; i++) {
      await operatorContract.submitTicket(tickets1[i].value, operator1, tickets1[i].virtualStakerIndex, {from: operator1});
    }

    mineBlocks(await operatorContract.ticketChallengeTimeout());
    let selectedTickets = await operatorContract.selectedTickets();
    assert.equal(selectedTickets.length, groupSize, "Should be trimmed to groupSize length.");

    let selectedParticipants = await operatorContract.selectedParticipants();
    assert.equal(selectedParticipants.length, groupSize, "Should be trimmed to groupSize length.");
  });

  it("should trigger new group selection when the last one is over", async function() {
    let groupSelectionStartBlock = await operatorContract.ticketSubmissionStartBlock();

    // Calculate the block time when the group selection should be finished
    let timeoutChallenge = (await operatorContract.ticketChallengeTimeout()).toNumber();
    let timeDKG = (await operatorContract.timeDKG()).toNumber();
    let groupSize = (await operatorContract.groupSize()).toNumber();
    let resultPublicationBlockStep = (await operatorContract.resultPublicationBlockStep()).toNumber();
    mineBlocks(timeoutChallenge + timeDKG + groupSize * resultPublicationBlockStep);

    await serviceContract.requestRelayEntry(bls.seed, {value: 10});
    await operatorContract.relayEntry(bls.nextNextGroupSignature);

    assert.isFalse((await operatorContract.ticketSubmissionStartBlock()).eq(groupSelectionStartBlock), "Group selection start block should be updated.");
    assert.isTrue((await operatorContract.groupSelectionRelayEntry()).eq(bls.nextNextGroupSignature), "Random beacon value for the current group selection should be updated.");
  });

});
