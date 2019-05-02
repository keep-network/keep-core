import { duration } from './helpers/increaseTime';
import exceptThrow from './helpers/expectThrow';
import mineBlocks from './helpers/mineBlocks';
import generateTickets from './helpers/generateTickets';
import {bls} from './helpers/data';
const KeepToken = artifacts.require('./KeepToken.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const KeepRandomBeaconFrontendProxy = artifacts.require('./KeepRandomBeaconFrontendProxy.sol');
const KeepRandomBeaconFrontendImplV1 = artifacts.require('./KeepRandomBeaconFrontendImplV1.sol');
const KeepRandomBeaconBackend = artifacts.require('./KeepRandomBeaconBackend.sol');


contract('TestKeepGroupSelection', function(accounts) {

  let token, stakingProxy, stakingContract, minimumStake, groupThreshold, groupSize,
    randomBeaconValue,
    timeoutInitial, timeoutSubmission, timeoutChallenge, timeDKG, resultPublicationBlockStep,
    frontendImplV1, frontendProxy, frontend,
    backend,
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
    frontendImplV1 = await KeepRandomBeaconFrontendImplV1.new();
    frontendProxy = await KeepRandomBeaconFrontendProxy.new(frontendImplV1.address);
    frontend = await KeepRandomBeaconFrontendImplV1.at(frontendProxy.address);

    // Initialize Keep Random Beacon backend contract
    minimumStake = 200000;
    groupThreshold = 15;
    groupSize = 20;
    timeoutInitial = 20;
    timeoutSubmission = 50;
    timeoutChallenge = 60;
    timeDKG = 20;
    resultPublicationBlockStep = 3;

    randomBeaconValue = bls.groupSignature;

    backend = await KeepRandomBeaconBackend.new();
    await backend.initialize(
      stakingProxy.address, frontendProxy.address, minimumStake, groupThreshold,
      groupSize, timeoutInitial, timeoutSubmission, timeoutChallenge, timeDKG, resultPublicationBlockStep,
      randomBeaconValue, bls.groupPubKey
    );

    await frontend.initialize(1, 1, backend.address);
    await backend.relayEntry(1, bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

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
    assert.isTrue(web3.utils.toBN(2000).eq(await backend.stakingWeight(operator1)), "Should have expected staking weight.");
    assert.isTrue(web3.utils.toBN(3000).eq(await backend.stakingWeight(operator3)), "Should have expected staking weight.");
  });

  it("should fail to get selected tickets before challenge period is over", async function() {
    await exceptThrow(backend.selectedTickets());
  });

  it("should fail to get selected participants before challenge period is over", async function() {
    await exceptThrow(backend.selectedParticipants());
  });

  it("should be able to get selected tickets and participants after challenge period is over", async function() {

    for (let i = 0; i < groupSize*2; i++) {
      await backend.submitTicket(tickets1[i].value, operator1, tickets1[i].virtualStakerIndex, {from: operator1});
    }

    mineBlocks(timeoutChallenge);
    let selectedTickets = await backend.selectedTickets();
    assert.equal(selectedTickets.length, groupSize, "Should be trimmed to groupSize length.");

    let selectedParticipants = await backend.selectedParticipants();
    assert.equal(selectedParticipants.length, groupSize, "Should be trimmed to groupSize length.");
  });

  it("should be able to output submited tickets in ascending ordered", async function() {

    let tickets = [];

    await backend.submitTicket(tickets1[0].value, operator1, tickets1[0].virtualStakerIndex, {from: operator1});
    tickets.push(tickets1[0].value);

    await backend.submitTicket(tickets2[0].value, operator2, tickets2[0].virtualStakerIndex, {from: operator2});
    tickets.push(tickets2[0].value);

    await backend.submitTicket(tickets3[0].value, operator3, tickets3[0].virtualStakerIndex, {from: operator3});
    tickets.push(tickets3[0].value);

    tickets = tickets.sort(function(a, b){return a-b}); // Sort numbers in ascending order

    // Test tickets ordering
    let orderedTickets = await backend.orderedTickets();
    assert.isTrue(orderedTickets[0].eq(tickets[0]), "Tickets should be in ascending order.");
    assert.isTrue(orderedTickets[1].eq(tickets[1]), "Tickets should be in ascending order.");
    assert.isTrue(orderedTickets[2].eq(tickets[2]), "Tickets should be in ascending order.");

  });

  it("should be able to submit a ticket during ticket submission period", async function() {
    await backend.submitTicket(tickets1[0].value, operator1, tickets1[0].virtualStakerIndex, {from: operator1});
    let proof = await backend.getTicketProof(tickets1[0].value);
    assert.isTrue(proof[1].eq(web3.utils.toBN(operator1)), "Should be able to get submitted ticket proof.");
    assert.equal(proof[2], tickets1[0].virtualStakerIndex, "Should be able to get submitted ticket proof.");
  });

  it("should be able to verify a ticket", async function() {

    await backend.submitTicket(tickets1[0].value, operator1, 1, {from: operator1});

    assert.isTrue(await backend.cheapCheck(
      operator1, operator1, 1
    ), "Should be able to verify a valid ticket.");
    
    assert.isTrue(await backend.costlyCheck(
      operator1, tickets1[0].value, operator1, tickets1[0].virtualStakerIndex
    ), "Should be able to verify a valid ticket.");
  
    assert.isFalse(await backend.costlyCheck(
      operator1, 0, operator1, tickets1[0].virtualStakerIndex
    ), "Should fail verifying invalid ticket.");

  });
});
