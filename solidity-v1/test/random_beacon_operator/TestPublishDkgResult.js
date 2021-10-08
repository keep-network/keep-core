const blsData = require("../helpers/data")
const sign = require("../helpers/signature")
const packTicket = require("../helpers/packTicket")
const generateTickets = require("../helpers/generateTickets")
const shuffleArray = require("../helpers/shuffle")
const { initContracts } = require("../helpers/initContracts")
const assert = require("chai").assert
const { createSnapshot, restoreSnapshot } = require("../helpers/snapshot.js")
const { contract, accounts, web3 } = require("@openzeppelin/test-environment")
const { expectRevert, time } = require("@openzeppelin/test-helpers")
const stakeDelegate = require("../helpers/stakeDelegate")

describe("KeepRandomBeaconOperator/PublishDkgResult", function () {
  const groupSize = 20
  const groupThreshold = 11
  const dkgResultSignatureThreshold = 15
  const resultPublicationBlockStep = 6

  let resultPublicationTime
  let token
  let stakingContract
  let operatorContract
  const owner = accounts[0]
  const beneficiary = accounts[4]
  let ticket
  const operator1 = accounts[0]
  const operator2 = accounts[1]
  const operator3 = accounts[2]
  const operator4 = accounts[3]
  let selectedParticipants
  let signatures
  let signingMemberIndices = []
  const noMisbehaved = "0x"
  const maxMisbehaved = "0x0102030405" // 20 - 15 = 5 max could misbehave
  const groupPubKey = blsData.groupPubKey
  const resultHash = web3.utils.soliditySha3(groupPubKey, noMisbehaved)

  before(async () => {
    const contracts = await initContracts(
      contract.fromArtifact("TokenStaking"),
      contract.fromArtifact("KeepRandomBeaconService"),
      contract.fromArtifact("KeepRandomBeaconServiceImplV1"),
      contract.fromArtifact("KeepRandomBeaconOperatorDKGResultStub")
    )

    token = contracts.token
    stakingContract = contracts.stakingContract
    operatorContract = contracts.operatorContract

    await operatorContract.setGroupSize(groupSize)
    await operatorContract.setGroupThreshold(groupThreshold)
    await operatorContract.setDKGResultSignatureThreshold(
      dkgResultSignatureThreshold
    )

    const operator1StakingWeight = 100
    const operator2StakingWeight = 200
    const operator3StakingWeight = 300
    const minimumStake = await stakingContract.minimumStake()

    await stakeDelegate(
      stakingContract,
      token,
      owner,
      operator1,
      beneficiary,
      owner,
      minimumStake.muln(operator1StakingWeight)
    )
    await stakeDelegate(
      stakingContract,
      token,
      owner,
      operator2,
      beneficiary,
      owner,
      minimumStake.muln(operator2StakingWeight)
    )
    await stakeDelegate(
      stakingContract,
      token,
      owner,
      operator3,
      beneficiary,
      owner,
      minimumStake.muln(operator3StakingWeight)
    )

    await stakingContract.authorizeOperatorContract(
      operator1,
      operatorContract.address,
      { from: owner }
    )
    await stakingContract.authorizeOperatorContract(
      operator2,
      operatorContract.address,
      { from: owner }
    )
    await stakingContract.authorizeOperatorContract(
      operator3,
      operatorContract.address,
      { from: owner }
    )

    time.increase((await stakingContract.initializationPeriod()).addn(1))

    const groupSelectionRelayEntry = await operatorContract.getGroupSelectionRelayEntry()
    const tickets1 = generateTickets(
      groupSelectionRelayEntry,
      operator1,
      operator1StakingWeight
    )
    const tickets2 = generateTickets(
      groupSelectionRelayEntry,
      operator2,
      operator2StakingWeight
    )
    const tickets3 = generateTickets(
      groupSelectionRelayEntry,
      operator3,
      operator3StakingWeight
    )

    for (let i = 0; i < groupSize; i++) {
      ticket = packTicket(
        tickets1[i].valueHex,
        tickets1[i].virtualStakerIndex,
        operator1
      )
      await operatorContract.submitTicket(ticket, { from: operator1 })
    }

    for (let i = 0; i < groupSize; i++) {
      ticket = packTicket(
        tickets2[i].valueHex,
        tickets2[i].virtualStakerIndex,
        operator2
      )
      await operatorContract.submitTicket(ticket, { from: operator2 })
    }

    for (let i = 0; i < groupSize; i++) {
      ticket = packTicket(
        tickets3[i].valueHex,
        tickets3[i].virtualStakerIndex,
        operator3
      )
      await operatorContract.submitTicket(ticket, { from: operator3 })
    }

    const ticketSubmissionStartBlock = await operatorContract.getTicketSubmissionStartBlock()
    const submissionTimeout = await operatorContract.ticketSubmissionTimeout()
    const timeDKG = await operatorContract.timeDKG()
    resultPublicationTime = ticketSubmissionStartBlock
      .add(submissionTimeout)
      .add(timeDKG)

    await time.advanceBlockTo(ticketSubmissionStartBlock.add(submissionTimeout))

    selectedParticipants = await operatorContract.selectedParticipants()

    signingMemberIndices = []
    signatures = undefined

    for (let i = 0; i < selectedParticipants.length; i++) {
      const signature = await sign(resultHash, selectedParticipants[i])
      signingMemberIndices.push(i + 1)
      if (signatures == undefined) signatures = signature
      else signatures += signature.slice(2, signature.length)
    }
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  it("allows to submit correct result as the first member after DKG finished", async () => {
    // Jump in time to when submitter becomes eligible to submit
    await time.advanceBlockTo(resultPublicationTime)

    await operatorContract.submitDkgResult(
      1,
      groupPubKey,
      noMisbehaved,
      signatures,
      signingMemberIndices,
      { from: selectedParticipants[0] }
    )
    assert.isTrue(
      await operatorContract.isGroupRegistered(groupPubKey),
      "group should be registered"
    )
    assert.equal(
      await operatorContract.numberOfGroups(),
      1,
      "expected 1 group to be registered"
    )
  })

  it("sends reward to the DKG submitter", async () => {
    // Jump in time to when submitter becomes eligible to submit
    await time.advanceBlockTo(resultPublicationTime)

    const beneficiaryBalance = web3.utils.toBN(
      await web3.eth.getBalance(beneficiary)
    )
    const dkgGasEstimate = await operatorContract.dkgGasEstimate()
    const submitterCustomGasPrice = web3.utils.toWei(
      web3.utils.toBN(65),
      "gwei"
    )
    const expectedSubmitterReward = dkgGasEstimate.mul(
      await operatorContract.gasPriceCeiling()
    )

    await operatorContract.submitDkgResult(
      1,
      groupPubKey,
      noMisbehaved,
      signatures,
      signingMemberIndices,
      { from: selectedParticipants[0], gasPrice: submitterCustomGasPrice }
    )

    const updatedBeneficiaryBalance = web3.utils.toBN(
      await web3.eth.getBalance(beneficiary)
    )
    assert.isTrue(
      updatedBeneficiaryBalance.eq(
        beneficiaryBalance.add(expectedSubmitterReward)
      ),
      "Submitter should receive expected reward."
    )
  })

  it("sends max dkgSubmitterReimbursementFee to the submitter in case of a much higher price than gas price ceiling", async () => {
    // Jump in time to when submitter becomes eligible to submit
    await time.advanceBlockTo(resultPublicationTime)

    const dkgSubmitterReimbursementFee = web3.utils.toBN(
      await web3.eth.getBalance(operatorContract.address)
    )
    const beneficiaryBalance = web3.utils.toBN(
      await web3.eth.getBalance(beneficiary)
    )

    await operatorContract.setGasPriceCeiling(
      web3.utils.toWei(web3.utils.toBN(100), "gwei")
    )

    await operatorContract.submitDkgResult(
      1,
      groupPubKey,
      noMisbehaved,
      signatures,
      signingMemberIndices,
      {
        from: selectedParticipants[0],
        gasPrice: web3.utils.toWei(web3.utils.toBN(100), "gwei"),
      }
    )
    const updatedBeneficiaryBalance = web3.utils.toBN(
      await web3.eth.getBalance(beneficiary)
    )
    assert.isTrue(
      updatedBeneficiaryBalance.eq(
        beneficiaryBalance.add(dkgSubmitterReimbursementFee)
      ),
      "Submitter should receive dkgSubmitterReimbursementFee"
    )
  })

  it("allows to submit correct result with unordered signatures and indexes", async () => {
    let unorderedSigningMembersIndexes = []
    for (let i = 0; i < selectedParticipants.length; i++) {
      unorderedSigningMembersIndexes[i] = i + 1
    }

    unorderedSigningMembersIndexes = shuffleArray(
      unorderedSigningMembersIndexes
    )
    let unorderedSignatures

    for (let i = 0; i < selectedParticipants.length; i++) {
      const signature = await sign(
        resultHash,
        selectedParticipants[unorderedSigningMembersIndexes[i] - 1]
      )
      if (unorderedSignatures == undefined) unorderedSignatures = signature
      else unorderedSignatures += signature.slice(2, signature.length)
    }

    // Jump in time to when submitter becomes eligible to submit
    await time.advanceBlockTo(resultPublicationTime)

    await operatorContract.submitDkgResult(
      1,
      groupPubKey,
      noMisbehaved,
      unorderedSignatures,
      unorderedSigningMembersIndexes,
      { from: selectedParticipants[0] }
    )
    assert.isTrue(
      await operatorContract.isGroupRegistered(groupPubKey),
      "group should be registered"
    )
    assert.equal(
      await operatorContract.numberOfGroups(),
      1,
      "expected 1 group to be registered"
    )
  })

  it("allows to submit result at eligible block time based on member index", async () => {
    const submitter1MemberIndex = 4
    const submitter2MemberIndex = 5
    const submitter2 = selectedParticipants[submitter2MemberIndex - 1]
    const eligibleBlockForSubmitter1 = resultPublicationTime.addn(
      (submitter1MemberIndex - 1) * resultPublicationBlockStep
    )
    const eligibleBlockForSubmitter2 = resultPublicationTime.addn(
      (submitter2MemberIndex - 1) * resultPublicationBlockStep
    )

    // Jump in time to when submitter 1 becomes eligible to submit
    await time.advanceBlockTo(eligibleBlockForSubmitter1)

    // Should throw if non eligible submitter 2 tries to submit
    await expectRevert(
      operatorContract.submitDkgResult(
        submitter2MemberIndex,
        groupPubKey,
        noMisbehaved,
        signatures,
        signingMemberIndices,
        { from: submitter2 }
      ),
      "Submitter not eligible"
    )

    // Jump in time to when submitter 2 becomes eligible to submit
    await time.advanceBlockTo(eligibleBlockForSubmitter2)

    await operatorContract.submitDkgResult(
      submitter2MemberIndex,
      groupPubKey,
      noMisbehaved,
      signatures,
      signingMemberIndices,
      { from: submitter2 }
    )
    assert.isTrue(
      await operatorContract.isGroupRegistered(groupPubKey),
      "group should be registered"
    )
    assert.equal(
      await operatorContract.numberOfGroups(),
      1,
      "expected 1 group to be registered"
    )
  })

  it("reverts if submitter was not selected to the group.", async () => {
    await expectRevert(
      operatorContract.submitDkgResult(
        1,
        groupPubKey,
        noMisbehaved,
        signatures,
        signingMemberIndices,
        { from: operator4 }
      ),
      "Unexpected submitter index"
    )

    assert.isFalse(
      await operatorContract.isGroupRegistered(groupPubKey),
      "group should not be registered"
    )
  })

  it("reverts for invalid signatures", async () => {
    signingMemberIndices = []
    signatures = undefined
    const lastParticipantIdx = dkgResultSignatureThreshold - 1

    for (let i = 0; i < lastParticipantIdx; i++) {
      const signature = await sign(resultHash, selectedParticipants[i])
      signingMemberIndices.push(i + 1)
      if (signatures == undefined) signatures = signature
      else signatures += signature.slice(2, signature.length)
    }

    // Add invalid signature as the last one
    const nonsenseHash = web3.utils.soliditySha3("ducky duck")
    const invalidSignature = await sign(
      nonsenseHash,
      selectedParticipants[lastParticipantIdx]
    )
    signatures += invalidSignature.slice(2, invalidSignature.length)
    signingMemberIndices.push(lastParticipantIdx + 1)

    // Jump in time to when first member is eligible to submit
    await time.advanceBlockTo(resultPublicationTime)

    await expectRevert(
      operatorContract.submitDkgResult(
        1,
        groupPubKey,
        noMisbehaved,
        signatures,
        signingMemberIndices,
        { from: selectedParticipants[0] }
      ),
      "Invalid signature"
    )

    assert.isFalse(
      await operatorContract.isGroupRegistered(groupPubKey),
      "group should not be registered"
    )
  })

  it("reverts for duplicate member indices", async () => {
    signingMemberIndices = []
    signatures = undefined
    const lastParticipantIdx = dkgResultSignatureThreshold - 1

    for (let i = 0; i < lastParticipantIdx; i++) {
      const signature = await sign(resultHash, selectedParticipants[i])
      signingMemberIndices.push(i + 1)
      if (signatures == undefined) signatures = signature
      else signatures += signature.slice(2, signature.length)
    }

    // Duplicate member and signature
    const signature = await sign(resultHash, selectedParticipants[0])
    signatures += signature.slice(2, signature.length)
    signingMemberIndices.push(1)

    // Jump in time to when first member is eligible to submit
    await time.advanceBlockTo(resultPublicationTime)

    await expectRevert(
      operatorContract.submitDkgResult(
        1,
        groupPubKey,
        noMisbehaved,
        signatures,
        signingMemberIndices,
        { from: selectedParticipants[0] }
      ),
      "Duplicate member index"
    )

    assert.isFalse(
      await operatorContract.isGroupRegistered(groupPubKey),
      "group should not be registered"
    )
  })

  it("allows to submit the result with minimum number of valid signatures", async () => {
    signingMemberIndices = []
    signatures = undefined

    // Create minimum amount of valid signatures
    for (let i = 0; i < dkgResultSignatureThreshold; i++) {
      const signature = await sign(resultHash, selectedParticipants[i])
      signingMemberIndices.push(i + 1)
      if (signatures == undefined) signatures = signature
      else signatures += signature.slice(2, signature.length)
    }

    // Jump in time to when first member is eligible to submit
    await time.advanceBlockTo(resultPublicationTime)

    await operatorContract.submitDkgResult(
      1,
      groupPubKey,
      noMisbehaved,
      signatures,
      signingMemberIndices,
      { from: selectedParticipants[0] }
    )

    assert.isTrue(
      await operatorContract.isGroupRegistered(groupPubKey),
      "group should be registered"
    )
    assert.equal(
      await operatorContract.numberOfGroups(),
      1,
      "expected 1 group to be registered"
    )
  })

  it("reverts without minimum number of signatures", async () => {
    signingMemberIndices = []
    signatures = undefined

    // Create less than minimum amount of valid signatures
    for (let i = 0; i < dkgResultSignatureThreshold - 1; i++) {
      const signature = await sign(resultHash, selectedParticipants[i])
      signingMemberIndices.push(i + 1)
      if (signatures == undefined) signatures = signature
      else signatures += signature.slice(2, signature.length)
    }

    // Jump in time to when first member is eligible to submit
    await time.advanceBlockTo(resultPublicationTime)

    await expectRevert(
      operatorContract.submitDkgResult(
        1,
        groupPubKey,
        noMisbehaved,
        signatures,
        signingMemberIndices,
        { from: selectedParticipants[0] }
      ),
      "Too few signatures"
    )

    assert.isFalse(
      await operatorContract.isGroupRegistered(groupPubKey),
      "group should not be registered"
    )
  })

  it("reverts for a public key having less than 128 bytes", async () => {
    // Jump in time to when submitter becomes eligible to submit
    await time.advanceBlockTo(resultPublicationTime)

    const invalidGroupPubKey = groupPubKey.slice(0, -2)

    const s = await signResult(invalidGroupPubKey, noMisbehaved)
    await expectRevert(
      operatorContract.submitDkgResult(
        1,
        invalidGroupPubKey,
        noMisbehaved,
        s.signatures,
        s.signingMemberIndices,
        { from: selectedParticipants[0] }
      ),
      "Malformed group public key"
    )
  })

  it("reverts for a public key having more than 128 bytes", async () => {
    // Jump in time to when submitter becomes eligible to submit
    await time.advanceBlockTo(resultPublicationTime)

    const invalidGroupPubKey = groupPubKey + "ff"

    const s = await signResult(invalidGroupPubKey, noMisbehaved)
    await expectRevert(
      operatorContract.submitDkgResult(
        1,
        invalidGroupPubKey,
        noMisbehaved,
        s.signatures,
        s.signingMemberIndices,
        { from: selectedParticipants[0] }
      ),
      "Malformed group public key"
    )
  })

  it("reverts for too many misbehaved", async () => {
    // Jump in time to when submitter becomes eligible to submit
    await time.advanceBlockTo(resultPublicationTime)

    const invalidMisbehaved = maxMisbehaved + "ff"

    const s = await signResult(groupPubKey, invalidMisbehaved)
    await expectRevert(
      operatorContract.submitDkgResult(
        1,
        groupPubKey,
        invalidMisbehaved,
        s.signatures,
        s.signingMemberIndices,
        { from: selectedParticipants[0] }
      ),
      "Malformed misbehaved"
    )
  })

  it("allows to submit with maximum possible misbehaved", async () => {
    // Jump in time to when submitter becomes eligible to submit
    await time.advanceBlockTo(resultPublicationTime)

    const s = await signResult(groupPubKey, maxMisbehaved)

    await operatorContract.submitDkgResult(
      1,
      groupPubKey,
      maxMisbehaved,
      s.signatures,
      s.signingMemberIndices,
      { from: selectedParticipants[0] }
    )
    // ok, no exceptions
  })

  it("allows to submit with more signatures than the group size", async () => {
    // Jump in time to when submitter becomes eligible to submit
    await time.advanceBlockTo(resultPublicationTime)

    const s = await signResult(groupPubKey, noMisbehaved)

    const anotherSignature = await sign(resultHash, selectedParticipants[0])
    s.signatures += anotherSignature.slice(2, anotherSignature.length)
    s.signingMemberIndices.push(s.signingMemberIndices.length + 1)

    await expectRevert(
      operatorContract.submitDkgResult(
        1,
        groupPubKey,
        noMisbehaved,
        s.signatures,
        s.signingMemberIndices,
        { from: selectedParticipants[0] }
      ),
      "Too many signatures"
    )
  })

  async function signResult(groupPublicKey, misbehaved) {
    const resultHash = web3.utils.soliditySha3(groupPublicKey, misbehaved)

    signingMemberIndices = []
    signatures = undefined

    for (let i = 0; i < selectedParticipants.length; i++) {
      const signature = await sign(resultHash, selectedParticipants[i])
      signingMemberIndices.push(i + 1)
      if (signatures == undefined) signatures = signature
      else signatures += signature.slice(2, signature.length)
    }

    return {
      signingMemberIndices: signingMemberIndices,
      signatures: signatures,
    }
  }
})
