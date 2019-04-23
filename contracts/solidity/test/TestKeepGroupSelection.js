import { duration } from './helpers/increaseTime';
import exceptThrow from './helpers/expectThrow';
import mineBlocks from './helpers/mineBlocks';
import generateTickets from './helpers/generateTickets';
import {bls} from './helpers/data';
const KeepToken = artifacts.require('./KeepToken.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const KeepRandomBeaconProxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconImplV1 = artifacts.require('./KeepRandomBeaconImplV1.sol');
const KeepGroupProxy = artifacts.require('./KeepGroup.sol');
const KeepGroupImplV1 = artifacts.require('./KeepGroupImplV1.sol');


contract('TestKeepGroupSelection', function(accounts) {

  let token, stakingProxy, stakingContract, minimumStake, groupThreshold, groupSize,
    randomBeaconValue,
<<<<<<< HEAD
    timeoutInitial, timeoutSubmission, timeoutChallenge,
    groupExpirationTimeout, numberOfActiveGroups,
=======
    timeoutInitial, timeoutSubmission, timeoutChallenge, timeDKG, resultPublicationBlockStep,
>>>>>>> origin/master
    keepRandomBeaconImplV1, keepRandomBeaconProxy, keepRandomBeaconImplViaProxy,
    keepGroupImplV1, keepGroupProxy, keepGroupImplViaProxy,
    owner = accounts[0], magpie = accounts[1], signature, delegation,
    operator1 = accounts[2], tickets1,
    operator2 = accounts[3], tickets2,
    operator3 = accounts[4], tickets3,
    operator4 = accounts[5], tickets4;

  beforeEach(async () => {
    token = await KeepToken.new();
    
    // Initialize staking contract under proxy
    stakingProxy = await StakingProxy.new();
    stakingContract = await TokenStaking.new(token.address, stakingProxy.address, duration.days(30));
    await stakingProxy.authorizeContract(stakingContract.address, {from: owner})

    // Initialize Keep Random Beacon contract
    keepRandomBeaconImplV1 = await KeepRandomBeaconImplV1.new();
    keepRandomBeaconProxy = await KeepRandomBeaconProxy.new(keepRandomBeaconImplV1.address);
    keepRandomBeaconImplViaProxy = await KeepRandomBeaconImplV1.at(keepRandomBeaconProxy.address);

    // Initialize Keep Group contract
    minimumStake = 200000;
    groupThreshold = 15;
    groupSize = 20;
    timeoutInitial = 20;
    timeoutSubmission = 50;
    timeoutChallenge = 60;
    timeDKG = 20;
    resultPublicationBlockStep = 3;
    numberOfActiveGroups = 1;
    groupExpirationTimeout = 1;

    randomBeaconValue = bls.groupSignature;

    keepGroupImplV1 = await KeepGroupImplV1.new();
    keepGroupProxy = await KeepGroupProxy.new(keepGroupImplV1.address);
    keepGroupImplViaProxy = await KeepGroupImplV1.at(keepGroupProxy.address);
    await keepGroupImplViaProxy.initialize(
      stakingProxy.address, keepRandomBeaconProxy.address, minimumStake, groupThreshold,
      groupSize, timeoutInitial, timeoutSubmission, timeoutChallenge, timeDKG, resultPublicationBlockStep
    );

    await keepRandomBeaconImplViaProxy.initialize(1,1, randomBeaconValue, bls.groupPubKey, keepGroupProxy.address);
    await keepRandomBeaconImplViaProxy.relayEntry(1, bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

    // Stake delegate tokens to operator1
    signature = Buffer.from((await web3.eth.sign(web3.utils.soliditySha3(owner), operator1)).substr(2), 'hex');
    delegation = '0x' + Buffer.concat([Buffer.from(magpie.substr(2), 'hex'), signature]).toString('hex');
    await token.approveAndCall(stakingContract.address, minimumStake*2000, delegation, {from: owner});
    tickets1 = generateTickets(randomBeaconValue, operator1, 2000);

    // Stake delegate tokens to operator2
    signature = Buffer.from((await web3.eth.sign(web3.utils.soliditySha3(owner), operator2)).substr(2), 'hex');
    delegation = '0x' + Buffer.concat([Buffer.from(magpie.substr(2), 'hex'), signature]).toString('hex');
    await token.approveAndCall(stakingContract.address, minimumStake*2000, delegation, {from: owner});
    tickets2 = generateTickets(randomBeaconValue, operator2, 2000);

    // Stake delegate tokens to operator3
    signature = Buffer.from((await web3.eth.sign(web3.utils.soliditySha3(owner), operator3)).substr(2), 'hex');
    delegation = '0x' + Buffer.concat([Buffer.from(magpie.substr(2), 'hex'), signature]).toString('hex');
    await token.approveAndCall(stakingContract.address, minimumStake*3000, delegation, {from: owner});
    tickets3 = generateTickets(randomBeaconValue, operator3, 3000);

  });

  it("should be able to get staking weight", async function() {
    assert.equal(web3.utils.toBN(2000).eq(await keepGroupImplViaProxy.stakingWeight(operator1)), true, "Should have expected staking weight.");
    assert.equal(web3.utils.toBN(3000).eq(await keepGroupImplViaProxy.stakingWeight(operator3)), true, "Should have expected staking weight.");
  });

  it("should fail to get selected tickets before challenge period is over", async function() {
    await exceptThrow(keepGroupImplViaProxy.selectedTickets());
  });

  it("should fail to get selected participants before challenge period is over", async function() {
    await exceptThrow(keepGroupImplViaProxy.selectedParticipants());
  });

  it("should be able to get selected tickets and participants after challenge period is over", async function() {

    for (let i = 0; i < groupSize*2; i++) {
      await keepGroupImplViaProxy.submitTicket(tickets1[i].value, operator1, tickets1[i].virtualStakerIndex, {from: operator1});
    }

    mineBlocks(timeoutChallenge);
    let selectedTickets = await keepGroupImplViaProxy.selectedTickets();
    assert.equal(selectedTickets.length, groupSize, "Should be trimmed to groupSize length.");

    let selectedParticipants = await keepGroupImplViaProxy.selectedParticipants();
    assert.equal(selectedParticipants.length, groupSize, "Should be trimmed to groupSize length.");
  });

  it("should be able to output submited tickets in ascending ordered", async function() {

    let tickets = [];

    await keepGroupImplViaProxy.submitTicket(tickets1[0].value, operator1, tickets1[0].virtualStakerIndex, {from: operator1});
    tickets.push(tickets1[0].value);

    await keepGroupImplViaProxy.submitTicket(tickets2[0].value, operator2, tickets2[0].virtualStakerIndex, {from: operator2});
    tickets.push(tickets2[0].value);

    await keepGroupImplViaProxy.submitTicket(tickets3[0].value, operator3, tickets3[0].virtualStakerIndex, {from: operator3});
    tickets.push(tickets3[0].value);

    tickets = tickets.sort(function(a, b){return a-b}); // Sort numbers in ascending order

    // Test tickets ordering
    let orderedTickets = await keepGroupImplViaProxy.orderedTickets();
    assert.equal(orderedTickets[0].eq(tickets[0]), true, "Tickets should be in ascending order.");
    assert.equal(orderedTickets[1].eq(tickets[1]), true, "Tickets should be in ascending order.");
    assert.equal(orderedTickets[2].eq(tickets[2]), true, "Tickets should be in ascending order.");

  });

  it("should be able to submit a ticket during ticket submission period", async function() {
    await keepGroupImplViaProxy.submitTicket(tickets1[0].value, operator1, tickets1[0].virtualStakerIndex, {from: operator1});
    let proof = await keepGroupImplViaProxy.getTicketProof(tickets1[0].value);
    assert.equal(proof[1].eq(web3.utils.toBN(operator1)), true , "Should be able to get submitted ticket proof.");
    assert.equal(proof[2], tickets1[0].virtualStakerIndex, "Should be able to get submitted ticket proof.");
  });

  it("should be able to verify a ticket", async function() {

    await keepGroupImplViaProxy.submitTicket(tickets1[0].value, operator1, 1, {from: operator1});

    assert.equal(await keepGroupImplViaProxy.cheapCheck(
      operator1, operator1, 1
    ), true, "Should be able to verify a valid ticket.");
    
    assert.equal(await keepGroupImplViaProxy.costlyCheck(
      operator1, tickets1[0].value, operator1, tickets1[0].virtualStakerIndex
    ), true, "Should be able to verify a valid ticket.");
  
    assert.equal(await keepGroupImplViaProxy.costlyCheck(
      operator1, 0, operator1, tickets1[0].virtualStakerIndex
    ), false, "Should fail verifying invalid ticket.");

  });
});
