import { sign } from './helpers/signature';
import mineBlocks from './helpers/mineBlocks';
import generateTickets from './helpers/generateTickets';
import stakeDelegate from './helpers/stakeDelegate';
import expectThrow from './helpers/expectThrow';
import shuffleArray from './helpers/shuffle';
import {initContracts} from './helpers/initContracts';


contract('TestKeepRandomBeaconOperatorPublishDkgResult', function(accounts) {

  let config, token, stakingContract, operatorContract,
  owner = accounts[0], magpie = accounts[0],
  operator1 = accounts[0],
  operator2 = accounts[1],
  operator3 = accounts[2],
  operator4 = accounts[3],
  selectedParticipants, signatures, signingMemberIndices = [],
  disqualified = '0x0000000000000000000000000000000000000000',
  inactive = '0x0000000000000000000000000000000000000000',
  groupPubKey = "0x1000000000000000000000000000000000000000000000000000000000000000",
  resultHash = web3.utils.soliditySha3(groupPubKey, disqualified, inactive);

  beforeEach(async () => {

    let contracts = await initContracts(
      accounts,
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./StakingProxy.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./KeepRandomBeaconOperator.sol')
    );
    config = contracts.config;
    token = contracts.token;
    stakingContract = contracts.stakingContract;
    operatorContract = contracts.operatorContract;

    await stakeDelegate(stakingContract, token, owner, operator1, magpie, config.minimumStake.mul(web3.utils.toBN(2000)))
    await stakeDelegate(stakingContract, token, owner, operator2, magpie, config.minimumStake.mul(web3.utils.toBN(2000)))
    await stakeDelegate(stakingContract, token, owner, operator3, magpie, config.minimumStake.mul(web3.utils.toBN(3000)))

    let tickets1 = generateTickets(await operatorContract.groupSelectionSeed(), operator1, 2000);
    let tickets2 = generateTickets(await operatorContract.groupSelectionSeed(), operator2, 2000);
    let tickets3 = generateTickets(await operatorContract.groupSelectionSeed(), operator3, 3000);

    for(let i = 0; i < config.groupSize; i++) {
      await operatorContract.submitTicket(tickets1[i].value, operator1, tickets1[i].virtualStakerIndex, {from: operator1});
    }

    for(let i = 0; i < config.groupSize; i++) {
      await operatorContract.submitTicket(tickets2[i].value, operator2, tickets2[i].virtualStakerIndex, {from: operator2});
    }

    for(let i = 0; i < config.groupSize; i++) {
      await operatorContract.submitTicket(tickets3[i].value, operator3, tickets3[i].virtualStakerIndex, {from: operator3});
    }

    let ticketSubmissionStartBlock = (await operatorContract.ticketSubmissionStartBlock()).toNumber();
    let timeoutChallenge = (await operatorContract.ticketChallengeTimeout()).toNumber();
    let timeDKG = (await operatorContract.timeDKG()).toNumber();
    config.resultPublicationTime = ticketSubmissionStartBlock + timeoutChallenge + timeDKG;

    selectedParticipants = await operatorContract.selectedParticipants();

    for(let i = 0; i < selectedParticipants.length; i++) {
      let signature = await sign(resultHash, selectedParticipants[i]);
      signingMemberIndices.push(i+1);
      if (signatures == undefined) signatures = signature
      else signatures += signature.slice(2, signature.length);
    }
  });

  it("should be able to submit correct result as first member after DKG finished.", async function() {
    // Jump in time to when submitter becomes eligible to submit
    let currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(config.resultPublicationTime - currentBlock);

    await operatorContract.submitDkgResult(1, groupPubKey, disqualified, inactive, signatures, signingMemberIndices, {from: selectedParticipants[0]})
    assert.isTrue(await operatorContract.isGroupRegistered(groupPubKey), "group should be registered");
    assert.equal(await operatorContract.numberOfGroups(), 1, "expected 1 group to be registered")
  });

  it("should be able to submit correct result with unordered signatures and indexes.", async function() {
    let unorderedSigningMembersIndexes = [];
    for (let i = 0; i < selectedParticipants.length; i++) {
      unorderedSigningMembersIndexes[i] = i + 1;
    }

    unorderedSigningMembersIndexes = shuffleArray(unorderedSigningMembersIndexes);
    let unorderedSignatures;

    for(let i = 0; i < selectedParticipants.length; i++) {
      let signature = await sign(resultHash, selectedParticipants[unorderedSigningMembersIndexes[i] - 1]);
      if (unorderedSignatures == undefined) unorderedSignatures = signature
      else unorderedSignatures += signature.slice(2, signature.length);
    }

    // Jump in time to when submitter becomes eligible to submit
    let currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(config.resultPublicationTime - currentBlock);

    await operatorContract.submitDkgResult(1, groupPubKey, disqualified, inactive, unorderedSignatures, unorderedSigningMembersIndexes, {from: selectedParticipants[0]})
    assert.isTrue(await operatorContract.isGroupRegistered(groupPubKey), "group should be registered");
    assert.equal(await operatorContract.numberOfGroups(), 1, "expected 1 group to be registered")
  });

  it("should only be able to submit result at eligible block time based on member index.", async function() {
    let resultPublicationBlockStep = (await operatorContract.resultPublicationBlockStep()).toNumber();
    let submitter1MemberIndex = 4;
    let submitter2MemberIndex = 5;
    let submitter2 = selectedParticipants[submitter2MemberIndex - 1];
    let eligibleBlockForSubmitter1 = config.resultPublicationTime + (submitter1MemberIndex-1)*resultPublicationBlockStep;
    let eligibleBlockForSubmitter2 = config.resultPublicationTime + (submitter2MemberIndex-1)*resultPublicationBlockStep;

    // Jump in time to when submitter 1 becomes eligible to submit
    let currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(eligibleBlockForSubmitter1 - currentBlock);

    // Should throw if non eligible submitter 2 tries to submit
    await expectThrow(operatorContract.submitDkgResult(
      submitter2MemberIndex, groupPubKey, disqualified, inactive, signatures, signingMemberIndices,
      {from: submitter2})
    );

    // Jump in time to when submitter 2 becomes eligible to submit
    currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(eligibleBlockForSubmitter2 - currentBlock);

    await operatorContract.submitDkgResult(submitter2MemberIndex, groupPubKey, disqualified, inactive, signatures, signingMemberIndices, {from: submitter2})
    assert.isTrue(await operatorContract.isGroupRegistered(groupPubKey), "group should be registered");
    assert.equal(await operatorContract.numberOfGroups(), 1, "expected 1 group to be registered")
  });

  it("should not be able to submit if submitter was not selected to be part of the group.", async function() {
    await expectThrow(operatorContract.submitDkgResult(
      1, groupPubKey, disqualified, inactive, signatures, signingMemberIndices, 
      {from: operator4})
    );

    assert.isFalse(await operatorContract.isGroupRegistered(groupPubKey), "group should not be registered");
  });

  it("should reject the result with invalid signatures.", async function() {
    signingMemberIndices = [];
    signatures = undefined;
    let lastParticipantIdx = config.groupThreshold - 1;

    // Create less than minimum amount of valid signatures
    for(let i = 0; i < lastParticipantIdx; i++) {
      let signature = await sign(resultHash, selectedParticipants[i]);
      signingMemberIndices.push(i+1);
      if (signatures == undefined) signatures = signature
      else signatures += signature.slice(2, signature.length);
    }

    // Add invalid signature as the last one
    let nonsenseHash = web3.utils.soliditySha3("ducky duck");
    let invalidSignature = await sign(nonsenseHash, selectedParticipants[lastParticipantIdx]);
    signatures += invalidSignature.slice(2, invalidSignature.length);
    signingMemberIndices.push(lastParticipantIdx);

    // Jump in time to when first member is eligible to submit
    let currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(config.resultPublicationTime - currentBlock);

    await expectThrow(operatorContract.submitDkgResult(
      1, groupPubKey, disqualified, inactive, signatures, signingMemberIndices,
      {from: selectedParticipants[0]})
    );

    assert.isFalse(await operatorContract.isGroupRegistered(groupPubKey), "group should not be registered");
  });

  it("should be able to submit the result with minimum number of valid signatures", async function() {
    signingMemberIndices = [];
    signatures = undefined;

    // Create minimum amount of valid signatures
    for(let i = 0; i < config.groupThreshold; i++) {
      let signature = await sign(resultHash, selectedParticipants[i]);
      signingMemberIndices.push(i+1);
      if (signatures == undefined) signatures = signature
      else signatures += signature.slice(2, signature.length);
    }

    // Jump in time to when first member is eligible to submit
    let currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(config.resultPublicationTime - currentBlock);

    await operatorContract.submitDkgResult(
      1, groupPubKey, disqualified, inactive, signatures, signingMemberIndices,
      {from: selectedParticipants[0]})

      assert.isTrue(await operatorContract.isGroupRegistered(groupPubKey), "group should be registered");
      assert.equal(await operatorContract.numberOfGroups(), 1, "expected 1 group to be registered")
  });

  it("should not be able to submit without minimum number of signatures", async function() {
    signingMemberIndices = [];
    signatures = undefined;

    // Create less than minimum amount of valid signatures
    for(let i = 0; i < config.groupThreshold - 1; i++) {
      let signature = await sign(resultHash, selectedParticipants[i]);
      signingMemberIndices.push(i+1);
      if (signatures == undefined) signatures = signature
      else signatures += signature.slice(2, signature.length);
    }

    // Jump in time to when first member is eligible to submit
    let currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(config.resultPublicationTime - currentBlock);

    await expectThrow(operatorContract.submitDkgResult(
      1, groupPubKey, disqualified, inactive, signatures, signingMemberIndices,
      {from: selectedParticipants[0]})
    );

    assert.isFalse(await operatorContract.isGroupRegistered(groupPubKey), "group should not be registered");
  });
})
