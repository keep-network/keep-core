const crypto = require("crypto")
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot"
import stakeDelegate from './helpers/stakeDelegate'
import {initContracts} from './helpers/initContracts'
import expectThrowWithMessage from './helpers/expectThrowWithMessage';
import mineBlocks from './helpers/mineBlocks'

contract('KeepRandomBeaconOperator', function(accounts) {

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
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./stubs/KeepRandomBeaconOperatorPricingRewardsWithdrawStub.sol')
    )

    token = contracts.token
    stakingContract = contracts.stakingContract
    operatorContract = contracts.operatorContract
    serviceContract = contracts.serviceContract

    groupSize = web3.utils.toBN(3)
    await operatorContract.setGroupSize(groupSize)

    await stakeDelegate(stakingContract, token, owner, operator1, beneficiary1, 0)
    await stakeDelegate(stakingContract, token, owner, operator2, beneficiary2, 0)
    await stakeDelegate(stakingContract, token, owner, operator3, beneficiary3, 0)

    group1 = crypto.randomBytes(128)
    group2 = crypto.randomBytes(128)
    group3 = crypto.randomBytes(128)

    await operatorContract.registerNewGroup(group1)
    await operatorContract.addGroupMember(group1, operator1)
    await operatorContract.addGroupMember(group1, operator2)
    await operatorContract.addGroupMember(group1, operator2)

    entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: requestor})

    // Stub relay entry for the current signing request and accumulate rewards
    await operatorContract.relayEntry()

    // Register second group with all members as operator1
    await operatorContract.registerNewGroup(group2)
    await operatorContract.addGroupMember(group2, operator1)
    await operatorContract.addGroupMember(group2, operator1)
    await operatorContract.addGroupMember(group2, operator1)

    // New request will expire the first group
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: requestor})
    await operatorContract.relayEntry()

    let entryFee = await serviceContract.entryFeeBreakdown()
    memberBaseReward = entryFee.groupProfitFee.div(groupSize)
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
    let memberIndices = await operatorContract.getGroupMemberIndices(group1, operator1)
    await operatorContract.withdrawGroupMemberRewards(operator1, 0, memberIndices)
    memberIndices = await operatorContract.getGroupMemberIndices(group2, operator1)
    await operatorContract.withdrawGroupMemberRewards(operator1, 1, memberIndices)

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
    let memberIndices = await operatorContract.getGroupMemberIndices(group1, operator2)
    await operatorContract.withdrawGroupMemberRewards(operator2, 0, memberIndices)
    assert.isTrue((web3.utils.toBN(await web3.eth.getBalance(beneficiary2))).eq(beneficiary2balance.add(expectedReward)), "Unexpected beneficiary balance")
  })

  it("should not be able to withdraw group rewards without correct data", async () => {
    // Register new group and request new entry so we can expire the previous two groups
    await operatorContract.registerNewGroup(group3)
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: requestor})

    mineBlocks(10)
    assert.isTrue(await operatorContract.isStaleGroup('0x' + group1.toString('hex')), "Group should be stale")

    // get indices for operator2 to be used by operator3 to withdraw
    let memberIndices = await operatorContract.getGroupMemberIndices(group1, operator2)

    let beneficiary3balance = web3.utils.toBN(await web3.eth.getBalance(beneficiary3))

    // operator3 doesn't have any group members, nothing can be withdrawn even using valid indices from other members
    await operatorContract.withdrawGroupMemberRewards(operator3, 0, memberIndices)
    assert.isTrue((web3.utils.toBN(await web3.eth.getBalance(beneficiary3))).eq(beneficiary3balance), "Unexpected beneficiary balance")
  })

  it("should not be able to withdraw group rewards if group is active", async () => {
    // operator1 has 3 members in group2
    let memberIndices = await operatorContract.getGroupMemberIndices(group2, operator1)

    assert.isFalse(await operatorContract.isExpiredGroup('0x' + group2.toString('hex')), "Group should not be expired")
    let beneficiary1balance = web3.utils.toBN(await web3.eth.getBalance(beneficiary1))

    // Nothing can be withdrawn
    await expectThrowWithMessage(
      operatorContract.withdrawGroupMemberRewards(operator1, 1, memberIndices),
      "Group must be expired and stale"
    )
    assert.isTrue((web3.utils.toBN(await web3.eth.getBalance(beneficiary1))).eq(beneficiary1balance), "Unexpected beneficiary balance")
  })

  it("should not be able to withdraw group rewards if group is expired but not stale", async () => {
    // operator2 has 2 members in group1
    let memberIndices = await operatorContract.getGroupMemberIndices(group1, operator2)

    assert.isTrue(await operatorContract.isExpiredGroup('0x' + group1.toString('hex')), "Group should be expired")
    assert.isFalse(await operatorContract.isStaleGroup('0x' + group1.toString('hex')), "Group should not be stale")
    let beneficiary2balance = web3.utils.toBN(await web3.eth.getBalance(beneficiary2))

    // Nothing can be withdrawn
    await expectThrowWithMessage(
      operatorContract.withdrawGroupMemberRewards(operator2, 0, memberIndices),
      "Group must be expired and stale"
    )
    assert.isTrue((web3.utils.toBN(await web3.eth.getBalance(beneficiary2))).eq(beneficiary2balance), "Unexpected beneficiary balance")
  })
})
