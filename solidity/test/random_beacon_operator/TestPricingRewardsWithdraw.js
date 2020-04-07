const initContracts = require('../helpers/initContracts')
const assert = require('chai').assert
const mineBlocks = require("../helpers/mineBlocks")
const {expectRevert} = require("@openzeppelin/test-helpers")
const {createSnapshot, restoreSnapshot} = require("../helpers/snapshot.js")
const {contract, accounts, web3} = require("@openzeppelin/test-environment")
const crypto = require("crypto")
const stakeDelegate = require('../helpers/stakeDelegate')

describe('KeepRandomBeaconOperator/PricingRewardsWithdraw', function() {

  let token, stakingContract, operatorContract, serviceContract,
    groupSize, memberBaseReward, entryFeeEstimate,
    group1, group2, group3,
    owner = accounts[0],
    requestor = accounts[1],
    operator1 = accounts[2],
    operator2 = accounts[3],
    operator3 = accounts[4],
    beneficiary1 = accounts[5],
    beneficiary2 = accounts[6],
    beneficiary3 = accounts[7]

  before(async () => {
    let contracts = await initContracts(
      contract.fromArtifact('KeepToken'),
      contract.fromArtifact('TokenStaking'),
      contract.fromArtifact('KeepRandomBeaconService'),
      contract.fromArtifact('KeepRandomBeaconServiceImplV1'),
      contract.fromArtifact('KeepRandomBeaconOperatorPricingRewardsWithdrawStub')
    )

    token = contracts.token
    stakingContract = contracts.stakingContract
    operatorContract = contracts.operatorContract
    serviceContract = contracts.serviceContract

    groupSize = web3.utils.toBN(3)
    await operatorContract.setGroupSize(groupSize)
    let minimumStake = await stakingContract.minimumStake()

    await stakeDelegate(stakingContract, token, owner, operator1, beneficiary1, operator1, minimumStake)
    await stakeDelegate(stakingContract, token, owner, operator2, beneficiary2, operator2, minimumStake)
    await stakeDelegate(stakingContract, token, owner, operator3, beneficiary3, operator3, minimumStake)

    group1 = crypto.randomBytes(128)
    group2 = crypto.randomBytes(128)
    group3 = crypto.randomBytes(128)

    await operatorContract.registerNewGroup(group1)
    await operatorContract.setGroupMembers(group1, [operator1, operator2, operator2])

    entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: requestor})

    // Stub relay entry for the current signing request and accumulate rewards
    await operatorContract.relayEntry()

    // Register second group with all members as operator1
    await operatorContract.registerNewGroup(group2)
    await operatorContract.setGroupMembers(group2, [operator1, operator1, operator1])

    // New request will expire the first group
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: requestor})
    await operatorContract.relayEntry()

    let entryFee = await serviceContract.entryFeeBreakdown()
    memberBaseReward = entryFee.groupProfitFee.div(groupSize)

    mineBlocks(4) // make sure groups become stale in tests
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  it("should allow fetching public key of active group", async() => {
    let groupPublicKey = await operatorContract.getGroupPublicKey(1)
    assert.equal(groupPublicKey, '0x' + group2.toString('hex'))
  })

  it("should allow fetching public key of stale group", async() => {
    mineBlocks(10)
    assert.isTrue(await operatorContract.isStaleGroup('0x' + group1.toString('hex')), "Group should be stale")

    let groupPublicKey = await operatorContract.getGroupPublicKey(0)
    assert.equal(groupPublicKey, '0x' + group1.toString('hex'))
  })

  it("should be able to withdraw group rewards from multiple staled groups", async () => {
    // Register new group and request new entry so we can expire the previous two groups
    await operatorContract.registerNewGroup(group3)
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: requestor})
    let beneficiary1balance = web3.utils.toBN(await web3.eth.getBalance(beneficiary1))

    mineBlocks(10)
    assert.isTrue(await operatorContract.isStaleGroup('0x' + group1.toString('hex')), "Group should be stale")
    assert.isTrue(await operatorContract.isStaleGroup('0x' + group2.toString('hex')), "Group should be stale")

    // operator1 has 1 member in group1 and 3 members in group2
    let expectedReward = memberBaseReward.muln(4)
    await operatorContract.withdrawGroupMemberRewards(operator1, 0)
    await operatorContract.withdrawGroupMemberRewards(operator1, 1)

    assert.isTrue((web3.utils.toBN(await web3.eth.getBalance(beneficiary1))).eq(beneficiary1balance.add(expectedReward)), "Unexpected beneficiary balance")
  })

  it("should be able to withdraw group rewards from a staled group", async () => {
    // Register new group and request new entry so we can expire the previous two groups
    await operatorContract.registerNewGroup(group3)
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: requestor})
    let beneficiary2balance = web3.utils.toBN(await web3.eth.getBalance(beneficiary2))

    mineBlocks(10)
    assert.isTrue(await operatorContract.isStaleGroup('0x' + group1.toString('hex')), "Group should be stale")

    // operator2 has 2 members in group1 only
    let expectedReward = memberBaseReward.muln(2)
    await operatorContract.withdrawGroupMemberRewards(operator2, 0)
    assert.isTrue((web3.utils.toBN(await web3.eth.getBalance(beneficiary2))).eq(beneficiary2balance.add(expectedReward)), "Unexpected beneficiary balance")
  })

  it("should not be able to withdraw group rewards without correct data", async () => {
    // Register new group and request new entry so we can expire the previous two groups
    await operatorContract.registerNewGroup(group3)
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: requestor})

    mineBlocks(10)
    assert.isTrue(await operatorContract.isStaleGroup('0x' + group1.toString('hex')), "Group should be stale")

    let beneficiary3balance = web3.utils.toBN(await web3.eth.getBalance(beneficiary3))

    // operator3 doesn't have any group members, nothing can be withdrawn
    await operatorContract.withdrawGroupMemberRewards(operator3, 0)
    assert.isTrue((web3.utils.toBN(await web3.eth.getBalance(beneficiary3))).eq(beneficiary3balance), "Unexpected beneficiary balance")
  })

  it("should not be able to withdraw group rewards if group is active", async () => {
    // operator1 has 3 members in group2

    assert.isFalse(await operatorContract.isExpiredGroup('0x' + group2.toString('hex')), "Group should not be expired")
    let beneficiary1balance = web3.utils.toBN(await web3.eth.getBalance(beneficiary1))

    // Nothing can be withdrawn
    await expectRevert(
      operatorContract.withdrawGroupMemberRewards(operator1, 1),
      "Group must be expired and stale"
    )
    assert.isTrue((web3.utils.toBN(await web3.eth.getBalance(beneficiary1))).eq(beneficiary1balance), "Unexpected beneficiary balance")
  })

  it("should not be able to withdraw group rewards if group is expired but not stale", async () => {
    // operator2 has 2 members in group1

    assert.isTrue(await operatorContract.isExpiredGroup('0x' + group1.toString('hex')), "Group should be expired")
    assert.isFalse(await operatorContract.isStaleGroup('0x' + group1.toString('hex')), "Group should not be stale")
    let beneficiary2balance = web3.utils.toBN(await web3.eth.getBalance(beneficiary2))

    // Nothing can be withdrawn
    await expectRevert(
      operatorContract.withdrawGroupMemberRewards(operator2, 0),
      "Group must be expired and stale"
    )
    assert.isTrue((web3.utils.toBN(await web3.eth.getBalance(beneficiary2))).eq(beneficiary2balance), "Unexpected beneficiary balance")
  })

  it("should not be able to withdraw group rewards multiple times", async () => {
    // Register new group and request new entry so we can expire the previous two groups
    await operatorContract.registerNewGroup(group3)
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: requestor})
    let beneficiary2balance = web3.utils.toBN(await web3.eth.getBalance(beneficiary2))

    mineBlocks(10)
    assert.isTrue(await operatorContract.isStaleGroup('0x' + group1.toString('hex')), "Group should be stale")

    // operator2 has 2 members in group1 only
    let expectedReward = memberBaseReward.muln(2)
    await operatorContract.withdrawGroupMemberRewards(operator2, 0)
    assert.isTrue((web3.utils.toBN(await web3.eth.getBalance(beneficiary2))).eq(beneficiary2balance.add(expectedReward)), "Unexpected beneficiary balance")

    await expectRevert(
      operatorContract.withdrawGroupMemberRewards(operator2, 0),
      "Rewards already withdrawn"
    );
  })
})
