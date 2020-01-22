import { formatDate } from '../utils'

const fetchAvailableRewards = async (web3Context) => {
  const { keepRandomBeaconOperatorContract, stakingContract, yourAddress, utils } = web3Context
  try {
    let totalRewardsBalance = utils.toBN(0)
    const expiredGroupsCount = await keepRandomBeaconOperatorContract.methods.getFirstActiveGroupIndex().call()
    const groups = []
    const groupMemberIndices = {}
    for (let groupIndex = 0; groupIndex < expiredGroupsCount; groupIndex++) {
      const groupPublicKey = await keepRandomBeaconOperatorContract.methods.getGroupPublicKey(groupIndex).call()
      const isStale = await keepRandomBeaconOperatorContract.methods.isStaleGroup(groupPublicKey).call()
      if (!isStale) {
        continue
      }

      const groupMembers = new Set(await keepRandomBeaconOperatorContract.methods.getGroupMembers(groupPublicKey).call())
      groupMemberIndices[groupPublicKey] = {}
      for (const memberAddress of groupMembers) {
        const beneficiaryAddressForMember = await stakingContract.methods.magpieOf(memberAddress).call()
        if (utils.toChecksumAddress(yourAddress) !== utils.toChecksumAddress(beneficiaryAddressForMember)) {
          continue
        }
        groupMemberIndices[groupPublicKey][memberAddress] = await keepRandomBeaconOperatorContract.methods.getGroupMemberIndices(groupPublicKey, memberAddress).call()
      }
      if (Object.keys(groupMemberIndices[groupPublicKey]).length === 0) {
        continue
      }
      const reward = await getAvailableRewardFromGroupInEther(groupPublicKey, groupMemberIndices, web3Context)
      totalRewardsBalance = totalRewardsBalance.add(utils.toBN(utils.toWei(reward, 'ether')))
      groups.push({ groupIndex, groupPublicKey, membersIndeces: groupMemberIndices[groupPublicKey], reward })
    }
    return [groups, utils.fromWei(totalRewardsBalance.toString(), 'ether')]
  } catch (error) {
    throw error
  }
}

const getAvailableRewardFromGroupInEther = async (groupPublicKey, groupMemberIndices, web3Context) => {
  const { utils, keepRandomBeaconOperatorContract } = web3Context
  const membersInGroup = Object.keys(groupMemberIndices[groupPublicKey])
  const rewardsMultiplier = membersInGroup.length === 1 ?
    groupMemberIndices[groupPublicKey][membersInGroup[0]].length :
    membersInGroup.reduce((prev, current) => groupMemberIndices[groupPublicKey][prev].length + groupMemberIndices[groupPublicKey][current].length)
  const groupMemberReward = await keepRandomBeaconOperatorContract.methods.getGroupMemberRewards(groupPublicKey).call()
  const wholeReward = utils.toBN(groupMemberReward).mul(utils.toBN(rewardsMultiplier))

  return utils.fromWei(wholeReward, 'ether')
}

const fetchWithdrawalHistory = async (web3Context) => {
  const { keepRandomBeaconOperatorContract, yourAddress, utils } = web3Context
  const searchFilters = { fromBlock: 0, filter: { beneficiary: yourAddress } }

  try {
    const events = await keepRandomBeaconOperatorContract.getPastEvents('GroupMemberRewardsWithdrawn', searchFilters)
    return events.map((event) => {
      const { returnValues: { groupPublicKey, withdrawnAt, amount } } = event
      return { groupPublicKey, date: formatDate(withdrawnAt * 1000), amount: utils.fromWei(amount, 'ether') }
    })
  } catch (error) {
    throw error
  }
}

const rewardsService = {
  fetchAvailableRewards,
  fetchWithdrawalHistory,
}

export default rewardsService
