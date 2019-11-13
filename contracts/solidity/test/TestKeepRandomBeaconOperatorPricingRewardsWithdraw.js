const crypto = require("crypto")
import expectThrow from './helpers/expectThrow'
import {bls} from './helpers/data'
import stakeDelegate from './helpers/stakeDelegate'
import {initContracts} from './helpers/initContracts'

contract('KeepRandomBeaconOperator', function(accounts) {

  let token, stakingContract, operatorContract, groupContract, serviceContract,
    groupSize, memberBaseReward, entryFeeEstimate,
    owner = accounts[0],
    requestor = accounts[1],
    operator1 = accounts[2],
    operator2 = accounts[3],
    operator3 = accounts[4],
    beneficiary1 = accounts[5],
    beneficiary2 = accounts[6],
    beneficiary3 = accounts[7]

  beforeEach(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./stubs/KeepRandomBeaconOperatorStub.sol'),
      artifacts.require('./KeepRandomBeaconOperatorGroupTerminationStub.sol')
    )

    token = contracts.token
    stakingContract = contracts.stakingContract
    operatorContract = contracts.operatorContract
    groupContract = contracts.groupContract
    serviceContract = contracts.serviceContract

    groupSize = web3.utils.toBN(3)
    await operatorContract.setGroupSize(groupSize)

    await stakeDelegate(stakingContract, token, owner, operator1, beneficiary1, 0)
    await stakeDelegate(stakingContract, token, owner, operator2, beneficiary2, 0)
    await stakeDelegate(stakingContract, token, owner, operator3, beneficiary3, 0)

    let group = crypto.randomBytes(128)
    await operatorContract.registerNewGroup(group)
    await operatorContract.addGroupMember(group, operator1)
    await operatorContract.addGroupMember(group, operator2)
    await operatorContract.addGroupMember(group, operator2)

    entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
    await serviceContract.methods['requestRelayEntry(uint256)'](bls.seed, {value: entryFeeEstimate, from: requestor})

    // Stub relay entry for the current signing request and accumulate rewards
    await operatorContract.stubRelayEntry()

    // Register second group with all members as operator1
    group = crypto.randomBytes(128)
    await operatorContract.registerNewGroup(group)
    await operatorContract.addGroupMember(group, operator1)
    await operatorContract.addGroupMember(group, operator1)
    await operatorContract.addGroupMember(group, operator1)

    // New request will expire the first group
    await serviceContract.methods['requestRelayEntry(uint256)'](bls.seed, {value: entryFeeEstimate, from: requestor})
    await operatorContract.stubRelayEntry()

    let entryFee = await serviceContract.entryFeeBreakdown()
    memberBaseReward = entryFee.groupProfitFee.div(groupSize)
  })

  it("should be able to withdraw group rewards from multiple staled groups", async function() {
    // Register new group and request new entry so we can expire the previous two groups
    await operatorContract.registerNewGroup(crypto.randomBytes(128))
    await serviceContract.methods['requestRelayEntry(uint256)'](bls.seed, {value: entryFeeEstimate, from: requestor})

    let beneficiary1balance = web3.utils.toBN(await web3.eth.getBalance(beneficiary1))
    let groupIndices = []
    let groupMemberIndices = []
    let staleGroupIndices = await groupContract.getStaleGroupsIndices()
    for (let i = 0; i < staleGroupIndices.length; i++) {
      let groupPubKey = await groupContract.getGroupPublicKey(staleGroupIndices[i])
      let memberIndices = await groupContract.getGroupMemberIndices(groupPubKey, operator1)

      for (let j = 0; j < memberIndices.length; j++) {
        groupIndices.push(staleGroupIndices[i])
        groupMemberIndices.push(memberIndices[j])
      }
    }

    // operator1 has 1 member in group1 and 3 members in group2
    let expectedReward = memberBaseReward.muln(groupMemberIndices.length)

    await operatorContract.withdrawGroupMemberRewards(operator1, groupMemberIndices, groupIndices)
    assert.isTrue((web3.utils.toBN(await web3.eth.getBalance(beneficiary1))).eq(beneficiary1balance.add(expectedReward)), "Unexpected beneficiary balance")
  })

  it("should be able to withdraw group rewards from a staled group", async function() {
    // Register new group and request new entry so we can expire the previous two groups
    await operatorContract.registerNewGroup(crypto.randomBytes(128))
    await serviceContract.methods['requestRelayEntry(uint256)'](bls.seed, {value: entryFeeEstimate, from: requestor})

    let beneficiary2balance = web3.utils.toBN(await web3.eth.getBalance(beneficiary2))
    let groupIndices = []
    let groupMemberIndices = []
    let staleGroupIndices = await groupContract.getStaleGroupsIndices()
    for (let i = 0; i < staleGroupIndices.length; i++) {
      let groupPubKey = await groupContract.getGroupPublicKey(staleGroupIndices[i])
      let memberIndices = await groupContract.getGroupMemberIndices(groupPubKey, operator2)

      for (let j = 0; j < memberIndices.length; j++) {
        groupIndices.push(staleGroupIndices[i])
        groupMemberIndices.push(memberIndices[j])
      }
    }

    // operator2 has 2 members in group1 only
    let expectedReward = memberBaseReward.muln(groupMemberIndices.length)

    await operatorContract.withdrawGroupMemberRewards(operator2, groupMemberIndices, groupIndices)
    assert.isTrue((web3.utils.toBN(await web3.eth.getBalance(beneficiary2))).eq(beneficiary2balance.add(expectedReward)), "Unexpected beneficiary balance")
  })

  it("should not be able to withdraw group rewards without correct data", async function() {
    let staleGroupIndices = await groupContract.getStaleGroupsIndices()
    assert.isTrue(staleGroupIndices.length == 1, "Unexpected amount of staled groups")

    // should revert as indices arrays length can't be larger than number of stale groups
    await expectThrow(operatorContract.withdrawGroupMemberRewards(operator1, [0, 1, 2, 3, 4], [0, 1, 2, 3, 4]))

    // Register new group and request new entry so we can expire the previous two groups
    await operatorContract.registerNewGroup(crypto.randomBytes(128))
    await serviceContract.methods['requestRelayEntry(uint256)'](bls.seed, {value: entryFeeEstimate, from: requestor})
    assert.isTrue((await groupContract.getStaleGroupsIndices()).length == 2, "Unexpected amount of staled groups")

    // get indices for operator2 to be used by operator3 to withdraw
    let groupIndices = []
    let groupMemberIndices = []
    staleGroupIndices = await groupContract.getStaleGroupsIndices()
    for (let i = 0; i < staleGroupIndices.length; i++) {
      let groupPubKey = await groupContract.getGroupPublicKey(staleGroupIndices[i])
      let memberIndices = await groupContract.getGroupMemberIndices(groupPubKey, operator2)

      for (let j = 0; j < memberIndices.length; j++) {
        groupIndices.push(staleGroupIndices[i])
        groupMemberIndices.push(memberIndices[j])
      }
    }

    let beneficiary3balance = web3.utils.toBN(await web3.eth.getBalance(beneficiary3))

    // operator3 doesn't have any group members, nothing can be withdrawn even using valid indices from other members
    await operatorContract.withdrawGroupMemberRewards(operator3, groupMemberIndices, groupIndices)
    assert.isTrue((web3.utils.toBN(await web3.eth.getBalance(beneficiary3))).eq(beneficiary3balance), "Unexpected beneficiary balance")
  })
})
