import { sign } from './helpers/signature';
import mineBlocks from './helpers/mineBlocks';
import generateTickets from './helpers/generateTickets';
import stakeDelegate from './helpers/stakeDelegate';
import expectThrow from './helpers/expectThrow';
import shuffleArray from './helpers/shuffle';
import {initContracts} from './helpers/initContracts';
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot";


contract('TestKeepRandomBeaconOperatorPublishDkgResult', function(accounts) {

  let resultPublicationTime, token, stakingContract, operatorContract,
  owner = accounts[0], magpie = accounts[4],
  operator1 = accounts[0],
  operator2 = accounts[1],
  operator3 = accounts[2],
  operator4 = accounts[3],  
  selectedParticipants, signatures, signingMemberIndices = [],
  disqualified = '0x0000000000000000000000000000000000000000',
  inactive = '0x0000000000000000000000000000000000000000',
  groupPubKey = "0x1000000000000000000000000000000000000000000000000000000000000000",
  resultHash = web3.utils.soliditySha3(groupPubKey, disqualified, inactive);

  const groupSize = 20;
  const groupThreshold = 15;  
  const minimumStake = web3.utils.toBN(200000);
  const resultPublicationBlockStep = 3;

  before(async () => {

    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./stubs/KeepRandomBeaconOperatorStub.sol'),
      artifacts.require('./KeepRandomBeaconOperatorGroups.sol')
    );

    token = contracts.token;
    stakingContract = contracts.stakingContract;
    operatorContract = contracts.operatorContract;

    operatorContract.setGroupSize(groupSize);
    operatorContract.setMinimumStake(minimumStake);

    await stakeDelegate(stakingContract, token, owner, operator1, magpie, minimumStake.mul(web3.utils.toBN(2000)))
    await stakeDelegate(stakingContract, token, owner, operator2, magpie, minimumStake.mul(web3.utils.toBN(2000)))
    await stakeDelegate(stakingContract, token, owner, operator3, magpie, minimumStake.mul(web3.utils.toBN(3000)))

    let tickets1 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator1, 2000);
    let tickets2 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator2, 2000);
    let tickets3 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator3, 3000);

    for(let i = 0; i < groupSize; i++) {
      await operatorContract.submitTicket(tickets1[i].value, operator1, tickets1[i].virtualStakerIndex, {from: operator1});
    }

    for(let i = 0; i < groupSize; i++) {
      await operatorContract.submitTicket(tickets2[i].value, operator2, tickets2[i].virtualStakerIndex, {from: operator2});
    }

    for(let i = 0; i < groupSize; i++) {
      await operatorContract.submitTicket(tickets3[i].value, operator3, tickets3[i].virtualStakerIndex, {from: operator3});
    }

    let ticketSubmissionStartBlock = (await operatorContract.getTicketSubmissionStartBlock()).toNumber();
    let timeoutChallenge = (await operatorContract.ticketSubmissionTimeout()).toNumber();
    let timeDKG = (await operatorContract.timeDKG()).toNumber();
    resultPublicationTime = ticketSubmissionStartBlock + timeoutChallenge + timeDKG;

    selectedParticipants = await operatorContract.selectedParticipants();

    signingMemberIndices = [];
    signatures = undefined;

    for(let i = 0; i < selectedParticipants.length; i++) {
      let signature = await sign(resultHash, selectedParticipants[i]);
      signingMemberIndices.push(i+1);
      if (signatures == undefined) signatures = signature
      else signatures += signature.slice(2, signature.length);
    }
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

  it("should be able to submit correct result as first member after DKG finished.", async function() {
    // Jump in time to when submitter becomes eligible to submit
    let currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(resultPublicationTime - currentBlock);

    await operatorContract.submitDkgResult(1, groupPubKey, disqualified, inactive, signatures, signingMemberIndices, {from: selectedParticipants[0]})
    assert.isTrue(await operatorContract.isGroupRegistered(groupPubKey), "group should be registered");
    assert.equal(await operatorContract.numberOfGroups(), 1, "expected 1 group to be registered")
  });

  it("should send reward to the DKG submitter.", async function() {
    // Jump in time to when submitter becomes eligible to submit
    let currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(resultPublicationTime - currentBlock);

    let magpieBalance = web3.utils.toBN(await web3.eth.getBalance(magpie));
    let dkgGasEstimate = await operatorContract.dkgGasEstimate();
    let submitterCustomGasPrice = web3.utils.toWei(web3.utils.toBN(25), 'gwei');
    let expectedSubmitterReward = dkgGasEstimate.mul(await operatorContract.priceFeedEstimate());

    await operatorContract.submitDkgResult(
      1, groupPubKey, disqualified, inactive, signatures, signingMemberIndices,
      {from: selectedParticipants[0], gasPrice: submitterCustomGasPrice}
    )

    let updatedMagpieBalance = web3.utils.toBN(await web3.eth.getBalance(magpie));
    assert.isTrue(updatedMagpieBalance.eq(magpieBalance.add(expectedSubmitterReward)), "Submitter should receive expected reward.");
  });

  it("should send max dkgSubmitterReimbursementFee to the submitter in case of a much higher price than priceFeedEstimate.", async function() {
    // Jump in time to when submitter becomes eligible to submit
    let currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(resultPublicationTime - currentBlock);

    let dkgSubmitterReimbursementFee = web3.utils.toBN(await web3.eth.getBalance(operatorContract.address));
    let magpieBalance = web3.utils.toBN(await web3.eth.getBalance(magpie));

    await operatorContract.setPriceFeedEstimate(web3.utils.toWei(web3.utils.toBN(100), 'gwei'));

    await operatorContract.submitDkgResult(
      1, groupPubKey, disqualified, inactive, signatures, signingMemberIndices,
      {from: selectedParticipants[0], gasPrice: web3.utils.toWei(web3.utils.toBN(100), 'gwei')}
    )
    let updatedMagpieBalance = web3.utils.toBN(await web3.eth.getBalance(magpie));
    assert.isTrue(updatedMagpieBalance.eq(magpieBalance.add(dkgSubmitterReimbursementFee)), "Submitter should receive dkgSubmitterReimbursementFee");
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
    mineBlocks(resultPublicationTime - currentBlock);

    await operatorContract.submitDkgResult(1, groupPubKey, disqualified, inactive, unorderedSignatures, unorderedSigningMembersIndexes, {from: selectedParticipants[0]})
    assert.isTrue(await operatorContract.isGroupRegistered(groupPubKey), "group should be registered");
    assert.equal(await operatorContract.numberOfGroups(), 1, "expected 1 group to be registered")
  });

  it("should only be able to submit result at eligible block time based on member index.", async function() {
    let submitter1MemberIndex = 4;
    let submitter2MemberIndex = 5;
    let submitter2 = selectedParticipants[submitter2MemberIndex - 1];
    let eligibleBlockForSubmitter1 = resultPublicationTime + (submitter1MemberIndex-1)*resultPublicationBlockStep;
    let eligibleBlockForSubmitter2 = resultPublicationTime + (submitter2MemberIndex-1)*resultPublicationBlockStep;

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
    let lastParticipantIdx = groupThreshold - 1;

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
    mineBlocks(resultPublicationTime - currentBlock);

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
    for(let i = 0; i < groupThreshold; i++) {
      let signature = await sign(resultHash, selectedParticipants[i]);
      signingMemberIndices.push(i+1);
      if (signatures == undefined) signatures = signature
      else signatures += signature.slice(2, signature.length);
    }

    // Jump in time to when first member is eligible to submit
    let currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(resultPublicationTime - currentBlock);

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
    for(let i = 0; i < groupThreshold - 1; i++) {
      let signature = await sign(resultHash, selectedParticipants[i]);
      signingMemberIndices.push(i+1);
      if (signatures == undefined) signatures = signature
      else signatures += signature.slice(2, signature.length);
    }

    // Jump in time to when first member is eligible to submit
    let currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(resultPublicationTime - currentBlock);

    await expectThrow(operatorContract.submitDkgResult(
      1, groupPubKey, disqualified, inactive, signatures, signingMemberIndices,
      {from: selectedParticipants[0]})
    );

    assert.isFalse(await operatorContract.isGroupRegistered(groupPubKey), "group should not be registered");
  });
})
