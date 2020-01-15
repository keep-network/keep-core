const fetchAvailableRewards = async (web3Context) => {
  const { keepRandomBeaconOperatorContract, stakingContract, yourAddress, utils } = web3Context
  try {
    const expiredGroupsCount = await keepRandomBeaconOperatorContract.methods.getFirstActiveGroupIndex().call()
    const groups = []
    const groupMemberIndices = {}
    // TODO iterate trough expired groups
    for (let groupIndex = 0; groupIndex < 2; groupIndex++) {
      const groupPubKey = await keepRandomBeaconOperatorContract.methods.getGroupPublicKey(groupIndex).call()
      const isStale = await keepRandomBeaconOperatorContract.methods.isStaleGroup(groupPubKey).call()
      if (isStale) {
        continue
      }

      const groupMembers = new Set(await keepRandomBeaconOperatorContract.methods.getGroupMembers(groupPubKey).call())
      groupMemberIndices[groupPubKey] = {}
      for (const memberAddress of groupMembers) {
        const beneficiaryAddressForMember = await stakingContract.methods.magpieOf(memberAddress).call()
        if (utils.toChecksumAddress(yourAddress) !== utils.toChecksumAddress(beneficiaryAddressForMember)) {
          continue
        }
        groupMemberIndices[groupPubKey][memberAddress] = await keepRandomBeaconOperatorContract.methods.getGroupMemberIndices(groupPubKey, memberAddress).call()
      }
      if (Object.keys(groupMemberIndices[groupPubKey]).length === 0) {
        continue
      }
      const reward = await getAvailableRewardFromGroupInEther(Object.keys(groupMemberIndices[groupPubKey]), groupPubKey, groupMemberIndices, web3Context)
      groups.push({ groupIndex, groupPubKey, membersIndeces: groupMemberIndices[groupPubKey], reward })
    }
    return Promise.resolve(groups)
  } catch (error) {
    return Promise.reject(error)
  }
}

const getAvailableRewardFromGroupInEther = async (membersInGroup, groupPubKey, groupMemberIndices, web3Context) => {
  const { utils, keepRandomBeaconOperatorContract } = web3Context
  const rewardsMultiplier = membersInGroup.length === 1 ?
    groupMemberIndices[groupPubKey][membersInGroup[0]].length :
    membersInGroup.reduce((prev, current) => groupMemberIndices[groupPubKey][prev].length + groupMemberIndices[groupPubKey][current].length)
  const groupMemberReward = await keepRandomBeaconOperatorContract.methods.getGroupMemberRewards(groupPubKey).call()
  const wholeReward = utils.toBN(groupMemberReward).mul(utils.toBN(rewardsMultiplier))

  return utils.fromWei(wholeReward, 'ether')
}

const rewardsService = {
  fetchAvailableRewards,
}

export default rewardsService
