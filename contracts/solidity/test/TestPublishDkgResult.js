import { duration } from './helpers/increaseTime';
import {bls} from './helpers/data';
import mineBlocks from './helpers/mineBlocks';
import generateTickets from './helpers/generateTickets';
import stakeDelegate from './helpers/stakeDelegate';
import expectThrow from './helpers/expectThrow';
import shuffleArray from './helpers/shuffle';
const KeepToken = artifacts.require('./KeepToken.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');
const Staking = artifacts.require('./Staking.sol');
const KeepRandomBeaconProxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconImplV1 = artifacts.require('./KeepRandomBeaconImplV1.sol');
const KeepGroupProxy = artifacts.require('./KeepGroup.sol');
const KeepGroupImplV1 = artifacts.require('./KeepGroupImplV1.sol');


contract('TestPublishDkgResult', function(accounts) {

  const minimumStake = 200000;
  const groupThreshold = 15;
  const groupSize = 20;
  const timeoutInitial = 20;
  const timeoutSubmission = 100;
  const timeoutChallenge = 60;
  const timeDKG = 20;
  const resultPublicationBlockStep = 3;
  const relayRequestTimeout = 10;

  let disqualified, inactive, resultHash,
  token, stakingProxy, stakingContract, randomBeaconValue, requestId,
  keepRandomBeaconImplV1, keepRandomBeaconProxy, keepRandomBeaconImplViaProxy,
  keepGroupImplV1, keepGroupProxy, keepGroupImplViaProxy, groupPubKey,
  ticketSubmissionStartBlock, selectedParticipants, signatures, signingMemberIndices = [],
  owner = accounts[0], magpie = accounts[0],
  operator1 = accounts[0], tickets1,
  operator2 = accounts[1], tickets2,
  operator3 = accounts[2], tickets3,
  operator4 = accounts[3];
  requestId = 0;
  disqualified = '0x0000000000000000000000000000000000000000'
  inactive = '0x0000000000000000000000000000000000000000'
  groupPubKey = "0x1000000000000000000000000000000000000000000000000000000000000000"

  resultHash = web3.utils.soliditySha3(groupPubKey, disqualified, inactive);

  beforeEach(async () => {
    token = await KeepToken.new();

    // Initialize staking contract under proxy
    stakingProxy = await StakingProxy.new();
    stakingContract = await Staking.new(token.address, stakingProxy.address, duration.days(30));
    await stakingProxy.authorizeContract(stakingContract.address, {from: owner})

    // Initialize Keep Random Beacon contract
    keepRandomBeaconImplV1 = await KeepRandomBeaconImplV1.new();
    keepRandomBeaconProxy = await KeepRandomBeaconProxy.new(keepRandomBeaconImplV1.address);
    keepRandomBeaconImplViaProxy = await KeepRandomBeaconImplV1.at(keepRandomBeaconProxy.address);

    // Initialize Keep Group contract
    keepGroupImplV1 = await KeepGroupImplV1.new();
    keepGroupProxy = await KeepGroupProxy.new(keepGroupImplV1.address);
    keepGroupImplViaProxy = await KeepGroupImplV1.at(keepGroupProxy.address);
    await keepGroupImplViaProxy.initialize(
      stakingProxy.address, keepRandomBeaconProxy.address, minimumStake, groupThreshold,
      groupSize, timeoutInitial, timeoutSubmission, timeoutChallenge, timeDKG, resultPublicationBlockStep
    );

    randomBeaconValue = bls.groupSignature;

    await keepRandomBeaconImplViaProxy.initialize(1,1, randomBeaconValue, bls.groupPubKey, keepGroupProxy.address,
      relayRequestTimeout);
    await keepRandomBeaconImplViaProxy.relayEntry(1, bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

    await stakeDelegate(stakingContract, token, owner, operator1, magpie, minimumStake*2000)
    await stakeDelegate(stakingContract, token, owner, operator2, magpie, minimumStake*2000)
    await stakeDelegate(stakingContract, token, owner, operator3, magpie, minimumStake*3000)

    tickets1 = generateTickets(randomBeaconValue, operator1, 2000);
    tickets2 = generateTickets(randomBeaconValue, operator2, 2000);
    tickets3 = generateTickets(randomBeaconValue, operator3, 3000);

    for(let i = 0; i < groupSize; i++) {
      await keepGroupImplViaProxy.submitTicket(tickets1[i].value, operator1, tickets1[i].virtualStakerIndex, {from: operator1});
    }

    for(let i = 0; i < groupSize; i++) {
      await keepGroupImplViaProxy.submitTicket(tickets2[i].value, operator2, tickets2[i].virtualStakerIndex, {from: operator2});
    }

    for(let i = 0; i < groupSize; i++) {
      await keepGroupImplViaProxy.submitTicket(tickets3[i].value, operator3, tickets3[i].virtualStakerIndex, {from: operator3});
    }

    ticketSubmissionStartBlock = await keepGroupImplViaProxy.ticketSubmissionStartBlock();
    selectedParticipants = await keepGroupImplViaProxy.selectedParticipants();

    for(let i = 0; i < selectedParticipants.length; i++) {
      let signature = await web3.eth.sign(resultHash, selectedParticipants[i]);
      signingMemberIndices.push(i+1);
      if (signatures == undefined) signatures = signature
      else signatures += signature.slice(2, signature.length);
    }
  });

  it("should be able to submit correct result as first member after DKG finished.", async function() {

    // Jump in time to when submitter becomes eligible to submit
    let currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(ticketSubmissionStartBlock.toNumber() + timeoutChallenge + timeDKG - currentBlock);

    await keepGroupImplViaProxy.submitDkgResult(requestId, 1, groupPubKey, disqualified, inactive, signatures, signingMemberIndices, {from: selectedParticipants[0]})
    let submitted = await keepGroupImplViaProxy.isDkgResultSubmitted.call(requestId);
    assert.equal(submitted, true, "DkgResult should should be submitted");
  });

  it("should be able to submit correct result with unordered signatures and indexes.", async function() {

    let unorderedSigningMembersIndexes = [];
    for (let i = 0; i < selectedParticipants.length; i++) {
      unorderedSigningMembersIndexes[i] = i + 1;
    }

    unorderedSigningMembersIndexes = shuffleArray(unorderedSigningMembersIndexes);
    let unorderedSignatures;

    for(let i = 0; i < selectedParticipants.length; i++) {
      let signature = await web3.eth.sign(resultHash, selectedParticipants[unorderedSigningMembersIndexes[i] - 1]);
      if (unorderedSignatures == undefined) unorderedSignatures = signature
      else unorderedSignatures += signature.slice(2, signature.length);
    }

    // Jump in time to when submitter becomes eligible to submit
    let currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(ticketSubmissionStartBlock.toNumber() + timeoutChallenge + timeDKG - currentBlock);

    await keepGroupImplViaProxy.submitDkgResult(requestId, 1, groupPubKey, disqualified, inactive, unorderedSignatures, unorderedSigningMembersIndexes, {from: selectedParticipants[0]})
    let submitted = await keepGroupImplViaProxy.isDkgResultSubmitted.call(requestId);
    assert.equal(submitted, true, "DkgResult should should be submitted");
  });

  it("should only be able to submit result at eligible block time based on member index.", async function() {

    let submitter1MemberIndex = 4;
    let submitter2MemberIndex = 5;
    let submitter2 = selectedParticipants[submitter2MemberIndex - 1];
    let eligibleBlockForSubmitter1 = ticketSubmissionStartBlock.toNumber() + timeoutChallenge + timeDKG + (submitter1MemberIndex-1)*resultPublicationBlockStep;
    let eligibleBlockForSubmitter2 = ticketSubmissionStartBlock.toNumber() + timeoutChallenge + timeDKG + (submitter2MemberIndex-1)*resultPublicationBlockStep;

    // Jump in time to when submitter 1 becomes eligible to submit
    let currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(eligibleBlockForSubmitter1 - currentBlock);

    // Should throw if non eligible submitter 2 tries to submit
    await expectThrow(keepGroupImplViaProxy.submitDkgResult(
      requestId, submitter2MemberIndex, groupPubKey, disqualified, inactive, signatures, signingMemberIndices,
      {from: submitter2})
    );

    // Jump in time to when submitter 2 becomes eligible to submit
    currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(eligibleBlockForSubmitter2 - currentBlock);

    await keepGroupImplViaProxy.submitDkgResult(requestId, submitter2MemberIndex, groupPubKey, disqualified, inactive, signatures, signingMemberIndices, {from: submitter2})
    let submitted = await keepGroupImplViaProxy.isDkgResultSubmitted.call(requestId);
    assert.equal(submitted, true, "DkgResult should be submitted");
  });

  it("should not be able to submit if submitter was not selected to be part of the group.", async function() {
    await expectThrow(keepGroupImplViaProxy.submitDkgResult(
      requestId, 1, groupPubKey, disqualified, inactive, signatures, signingMemberIndices, 
      {from: operator4})
    );
  });

  it("should reject the result with invalid signatures.", async function() {

    signingMemberIndices = [];
    signatures = undefined;
    let lastParticipantIdx = groupThreshold - 1;

    // Create less than minimum amount of valid signatures
    for(let i = 0; i < lastParticipantIdx; i++) {
      let signature = await web3.eth.sign(resultHash, selectedParticipants[i]);
      signingMemberIndices.push(i+1);
      if (signatures == undefined) signatures = signature
      else signatures += signature.slice(2, signature.length);
    }

    // Add invalid signature as the last one
    let nonsenseHash = web3.utils.soliditySha3("ducky duck");
    let invalidSignature = await web3.eth.sign(nonsenseHash, selectedParticipants[lastParticipantIdx]);
    signatures += invalidSignature.slice(2, invalidSignature.length);
    signingMemberIndices.push(lastParticipantIdx);

    // Jump in time to when first member is eligible to submit
    let currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(ticketSubmissionStartBlock.toNumber() + timeoutChallenge + timeDKG - currentBlock);

    await expectThrow(keepGroupImplViaProxy.submitDkgResult(
      requestId, 1, groupPubKey, disqualified, inactive, signatures, signingMemberIndices,
      {from: selectedParticipants[0]})
    );
  });

  it("should be able to submit the result with minimum number of valid signatures", async function() {

    signingMemberIndices = [];
    signatures = undefined;

    // Create minimum amount of valid signatures
    for(let i = 0; i < groupThreshold; i++) {
      let signature = await web3.eth.sign(resultHash, selectedParticipants[i]);
      signingMemberIndices.push(i+1);
      if (signatures == undefined) signatures = signature
      else signatures += signature.slice(2, signature.length);
    }

    // Jump in time to when first member is eligible to submit
    let currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(ticketSubmissionStartBlock.toNumber() + timeoutChallenge + timeDKG - currentBlock);

    await keepGroupImplViaProxy.submitDkgResult(
      requestId, 1, groupPubKey, disqualified, inactive, signatures, signingMemberIndices,
      {from: selectedParticipants[0]})
    let submitted = await keepGroupImplViaProxy.isDkgResultSubmitted.call(requestId);
    assert.equal(submitted, true, "DkgResult should should be submitted");

  });

  it("should not be able to submit without minimum number of signatures", async function() {

    signingMemberIndices = [];
    signatures = undefined;

    // Create less than minimum amount of valid signatures
    for(let i = 0; i < groupThreshold - 1; i++) {
      let signature = await web3.eth.sign(resultHash, selectedParticipants[i]);
      signingMemberIndices.push(i+1);
      if (signatures == undefined) signatures = signature
      else signatures += signature.slice(2, signature.length);
    }

    // Jump in time to when first member is eligible to submit
    let currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(ticketSubmissionStartBlock.toNumber() + timeoutChallenge + timeDKG - currentBlock);

    await expectThrow(keepGroupImplViaProxy.submitDkgResult(
      requestId, 1, groupPubKey, disqualified, inactive, signatures, signingMemberIndices,
      {from: selectedParticipants[0]})
    );

  });
})
