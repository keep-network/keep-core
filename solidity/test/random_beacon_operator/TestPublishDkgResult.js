import { sign } from '../helpers/signature';
import mineBlocks from '../helpers/mineBlocks';
import increaseTime from '../helpers/increaseTime';
import packTicket from '../helpers/packTicket';
import generateTickets from '../helpers/generateTickets';
import stakeDelegate from '../helpers/stakeDelegate';
import expectThrow from '../helpers/expectThrow';
import shuffleArray from '../helpers/shuffle';
import {initContracts} from '../helpers/initContracts';
import {createSnapshot, restoreSnapshot} from '../helpers/snapshot';
import {bls} from '../helpers/data';
const { expectRevert } = require("@openzeppelin/test-helpers")

contract('KeepRandomBeaconOperator/PublishDkgResult', function(accounts) {

  let resultPublicationTime, token, stakingContract, operatorContract,
  owner = accounts[0], magpie = accounts[4], ticket,
  operator1 = accounts[0],
  operator2 = accounts[1],
  operator3 = accounts[2],
  operator4 = accounts[3],
  selectedParticipants, signatures, signingMemberIndices = [],
  misbehaved = '0x',
  groupPubKey = bls.groupPubKey,
  resultHash = web3.utils.soliditySha3(groupPubKey, misbehaved);

  const groupSize = 20;
  const groupThreshold = 15;
  const resultPublicationBlockStep = 3;

  before(async () => {

    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./stubs/KeepRandomBeaconOperatorStub.sol')
    );

    token = contracts.token;
    stakingContract = contracts.stakingContract;
    operatorContract = contracts.operatorContract;

    operatorContract.setGroupSize(groupSize);
    operatorContract.setGroupThreshold(groupThreshold);

    const operator1StakingWeight = 100;
    const operator2StakingWeight = 200;
    const operator3StakingWeight = 300;
    let minimumStake = await stakingContract.minimumStake()

    await stakeDelegate(stakingContract, token, owner, operator1, magpie, owner, minimumStake.muln(operator1StakingWeight))
    await stakeDelegate(stakingContract, token, owner, operator2, magpie, owner, minimumStake.muln(operator2StakingWeight))
    await stakeDelegate(stakingContract, token, owner, operator3, magpie, owner, minimumStake.muln(operator3StakingWeight))

    await stakingContract.authorizeOperatorContract(operator1, operatorContract.address, {from: owner})
    await stakingContract.authorizeOperatorContract(operator2, operatorContract.address, {from: owner})
    await stakingContract.authorizeOperatorContract(operator3, operatorContract.address, {from: owner})

    increaseTime((await stakingContract.initializationPeriod()).toNumber() + 1);

    let tickets1 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator1, operator1StakingWeight);
    let tickets2 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator2, operator2StakingWeight);
    let tickets3 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator3, operator3StakingWeight);

    for(let i = 0; i < groupSize; i++) {
      ticket = packTicket(tickets1[i].valueHex, tickets1[i].virtualStakerIndex, operator1);
      await operatorContract.submitTicket(ticket, {from: operator1});
    }

    for(let i = 0; i < groupSize; i++) {
      ticket = packTicket(tickets2[i].valueHex, tickets2[i].virtualStakerIndex, operator2);
      await operatorContract.submitTicket(ticket, {from: operator2});
    }

    for(let i = 0; i < groupSize; i++) {
      ticket = packTicket(tickets3[i].valueHex, tickets3[i].virtualStakerIndex, operator3);
      await operatorContract.submitTicket(ticket, {from: operator3});
    }

    let ticketSubmissionStartBlock = (await operatorContract.getTicketSubmissionStartBlock()).toNumber();
    let submissionTimeout = (await operatorContract.ticketSubmissionTimeout()).toNumber();
    let timeDKG = (await operatorContract.timeDKG()).toNumber();
    resultPublicationTime = ticketSubmissionStartBlock + submissionTimeout + timeDKG;

    let currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(ticketSubmissionStartBlock + submissionTimeout - currentBlock);

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

    await operatorContract.submitDkgResult(1, groupPubKey, misbehaved, signatures, signingMemberIndices, {from: selectedParticipants[0]})
    assert.isTrue(await operatorContract.isGroupRegistered(groupPubKey), "group should be registered");
    assert.equal(await operatorContract.numberOfGroups(), 1, "expected 1 group to be registered")
  });

  it("should send reward to the DKG submitter", async function() {
    // Jump in time to when submitter becomes eligible to submit
    let currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(resultPublicationTime - currentBlock);

    let magpieBalance = web3.utils.toBN(await web3.eth.getBalance(magpie));
    let dkgGasEstimate = await operatorContract.dkgGasEstimate();
    let submitterCustomGasPrice = web3.utils.toWei(web3.utils.toBN(35), 'gwei');
    let expectedSubmitterReward = dkgGasEstimate.mul(await operatorContract.gasPriceCeiling());

    await operatorContract.submitDkgResult(
      1, groupPubKey, misbehaved, signatures, signingMemberIndices,
      {from: selectedParticipants[0], gasPrice: submitterCustomGasPrice}
    )

    let updatedMagpieBalance = web3.utils.toBN(await web3.eth.getBalance(magpie));
    assert.isTrue(updatedMagpieBalance.eq(magpieBalance.add(expectedSubmitterReward)), "Submitter should receive expected reward.");
  });

  it("should send max dkgSubmitterReimbursementFee to the submitter in case of a much higher price than gas price ceiling", async function() {
    // Jump in time to when submitter becomes eligible to submit
    let currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(resultPublicationTime - currentBlock);

    let dkgSubmitterReimbursementFee = web3.utils.toBN(await web3.eth.getBalance(operatorContract.address));
    let magpieBalance = web3.utils.toBN(await web3.eth.getBalance(magpie));

    await operatorContract.setGasPriceCeiling(web3.utils.toWei(web3.utils.toBN(100), 'gwei'));

    await operatorContract.submitDkgResult(
      1, groupPubKey, misbehaved, signatures, signingMemberIndices,
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

    await operatorContract.submitDkgResult(1, groupPubKey, misbehaved, unorderedSignatures, unorderedSigningMembersIndexes, {from: selectedParticipants[0]})
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
      submitter2MemberIndex, groupPubKey, misbehaved, signatures, signingMemberIndices,
      {from: submitter2})
    );

    // Jump in time to when submitter 2 becomes eligible to submit
    currentBlock = await web3.eth.getBlockNumber();
    mineBlocks(eligibleBlockForSubmitter2 - currentBlock);

    await operatorContract.submitDkgResult(submitter2MemberIndex, groupPubKey, misbehaved, signatures, signingMemberIndices, {from: submitter2})
    assert.isTrue(await operatorContract.isGroupRegistered(groupPubKey), "group should be registered");
    assert.equal(await operatorContract.numberOfGroups(), 1, "expected 1 group to be registered")
  });

  it("should not be able to submit if submitter was not selected to be part of the group.", async function() {
    await expectThrow(operatorContract.submitDkgResult(
      1, groupPubKey, misbehaved, signatures, signingMemberIndices, 
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
      1, groupPubKey, misbehaved, signatures, signingMemberIndices,
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
      1, groupPubKey, misbehaved, signatures, signingMemberIndices,
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
      1, groupPubKey, misbehaved, signatures, signingMemberIndices,
      {from: selectedParticipants[0]})
    );

    assert.isFalse(await operatorContract.isGroupRegistered(groupPubKey), "group should not be registered");
  });

  it("should fail to submit with a public key having less than 128 bytes", async () => {
      // Jump in time to when submitter becomes eligible to submit
      let currentBlock = await web3.eth.getBlockNumber();
      mineBlocks(resultPublicationTime - currentBlock);
    
      let invalidGroupPubKey = groupPubKey.slice(0, -2)

      await expectRevert(
        operatorContract.submitDkgResult(
          1, invalidGroupPubKey, misbehaved, signatures, 
          signingMemberIndices, {from: selectedParticipants[0]}
        ),
        "Malformed group public key"
      )
  })

  it("should fail to submit with a public key having more than 128 bytes", async () => {
      // Jump in time to when submitter becomes eligible to submit
      let currentBlock = await web3.eth.getBlockNumber();
      mineBlocks(resultPublicationTime - currentBlock);
    
      let invalidGroupPubKey = groupPubKey + 'ff';

      await expectRevert(
        operatorContract.submitDkgResult(
          1, invalidGroupPubKey, misbehaved, signatures, 
          signingMemberIndices, {from: selectedParticipants[0]}
        ),
        "Malformed group public key"
      ) 
  })
})
