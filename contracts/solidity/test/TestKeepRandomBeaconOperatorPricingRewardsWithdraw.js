const crypto = require("crypto")
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot"
import {bls} from './helpers/data'
import stakeDelegate from './helpers/stakeDelegate'
import {initContracts} from './helpers/initContracts'
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
    await serviceContract.methods['requestRelayEntry(uint256)'](bls.seed, {value: entryFeeEstimate, from: requestor})

    // Stub relay entry for the current signing request and accumulate rewards
    await operatorContract.relayEntry()

    // Register second group with all members as operator1
    await operatorContract.registerNewGroup(group2)
    await operatorContract.addGroupMember(group2, operator1)
    await operatorContract.addGroupMember(group2, operator1)
    await operatorContract.addGroupMember(group2, operator1)

    // New request will expire the first group
    await serviceContract.methods['requestRelayEntry(uint256)'](bls.seed, {value: entryFeeEstimate, from: requestor})
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

  it("should be able to withdraw group rewards from multiple staled groups", async () => {
    // Register new group and request new entry so we can expire the previous two groups
    await operatorContract.registerNewGroup(group3)
    await serviceContract.methods['requestRelayEntry(uint256)'](bls.seed, {value: entryFeeEstimate, from: requestor})
    let beneficiary1balance = web3.utils.toBN(await web3.eth.getBalance(beneficiary1))

    // Make sure expired groups become stale
    mineBlocks(10)

    // operator1 has 1 member in group1 and 3 members in group2
    let expectedReward = memberBaseReward.muln(4)
    let memberIndices = await operatorContract.getGroupMemberIndices(group1, operator1)
    let groupIndex = await operatorContract.getGroupIndex('0x' + group1.toString('hex'))
    await operatorContract.withdrawGroupMemberRewards(groupIndex, memberIndices, {from: operator1})
    memberIndices = await operatorContract.getGroupMemberIndices(group2, operator1)
    groupIndex = await operatorContract.getGroupIndex('0x' + group2.toString('hex'))
    await operatorContract.withdrawGroupMemberRewards(groupIndex, memberIndices, {from: operator1})

    assert.isTrue((web3.utils.toBN(await web3.eth.getBalance(beneficiary1))).eq(beneficiary1balance.add(expectedReward)), "Unexpected beneficiary balance")
  })

  it("should be able to withdraw group rewards from a staled group", async () => {
    // Register new group and request new entry so we can expire the previous two groups
    await operatorContract.registerNewGroup(group3)
    await serviceContract.methods['requestRelayEntry(uint256)'](bls.seed, {value: entryFeeEstimate, from: requestor})
    let beneficiary2balance = web3.utils.toBN(await web3.eth.getBalance(beneficiary2))

    // Make sure expired groups become stale
    mineBlocks(10)

    // operator2 has 2 members in group1 only
    let expectedReward = memberBaseReward.muln(2)
    let memberIndices = await operatorContract.getGroupMemberIndices(group1, operator2)
    let groupIndex = await operatorContract.getGroupIndex('0x' + group1.toString('hex'))
    await operatorContract.withdrawGroupMemberRewards(groupIndex, memberIndices, {from: operator2})
    assert.isTrue((web3.utils.toBN(await web3.eth.getBalance(beneficiary2))).eq(beneficiary2balance.add(expectedReward)), "Unexpected beneficiary balance")
  })

  it("should not be able to withdraw group rewards without correct data", async () => {
    // Register new group and request new entry so we can expire the previous two groups
    await operatorContract.registerNewGroup(group3)
    await serviceContract.methods['requestRelayEntry(uint256)'](bls.seed, {value: entryFeeEstimate, from: requestor})

    // Make sure expired groups become stale
    mineBlocks(10)

    // get indices for operator2 to be used by operator3 to withdraw
    let memberIndices = await operatorContract.getGroupMemberIndices(group1, operator2)
    let groupIndex = await operatorContract.getGroupIndex('0x' + group1.toString('hex'))

    let beneficiary3balance = web3.utils.toBN(await web3.eth.getBalance(beneficiary3))

    // operator3 doesn't have any group members, nothing can be withdrawn even using valid indices from other members
    await operatorContract.withdrawGroupMemberRewards(groupIndex, memberIndices, {from: operator3})
    assert.isTrue((web3.utils.toBN(await web3.eth.getBalance(beneficiary3))).eq(beneficiary3balance), "Unexpected beneficiary balance")
  })

  it("should not be able to withdraw group rewards if group is active", async () => {
    // operator1 has 3 members in group2
    let memberIndices = await operatorContract.getGroupMemberIndices(group2, operator1)
    let groupIndex = await operatorContract.getGroupIndex('0x' + group2.toString('hex'))

    assert.isFalse(await operatorContract.isExpiredGroup('0x' + group2.toString('hex')), "Group should not be expired")
    assert.isFalse(await operatorContract.isStaleGroup('0x' + group2.toString('hex')), "Group should not be stale")
    let beneficiary1balance = web3.utils.toBN(await web3.eth.getBalance(beneficiary1))

    // Nothing can be withdrawn
    await operatorContract.withdrawGroupMemberRewards(groupIndex, memberIndices, {from: operator1})
    assert.isTrue((web3.utils.toBN(await web3.eth.getBalance(beneficiary1))).eq(beneficiary1balance), "Unexpected beneficiary balance")
  })

  it("should not be able to withdraw group rewards if group is expired but not stale", async () => {
    // operator2 has 2 members in group1
    let memberIndices = await operatorContract.getGroupMemberIndices(group1, operator2)
    let groupIndex = await operatorContract.getGroupIndex('0x' + group1.toString('hex'))

    assert.isTrue(await operatorContract.isExpiredGroup('0x' + group1.toString('hex')), "Group should be expired")
    assert.isFalse(await operatorContract.isStaleGroup('0x' + group1.toString('hex')), "Group should not be stale")
    let beneficiary2balance = web3.utils.toBN(await web3.eth.getBalance(beneficiary2))

    // Nothing can be withdrawn
    await operatorContract.withdrawGroupMemberRewards(groupIndex, memberIndices, {from: operator2})
    assert.isTrue((web3.utils.toBN(await web3.eth.getBalance(beneficiary2))).eq(beneficiary2balance), "Unexpected beneficiary balance")
  })
})
