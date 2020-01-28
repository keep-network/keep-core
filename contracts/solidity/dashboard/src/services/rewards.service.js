import { formatDate, wait } from '../utils'

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
      const { reward, rewardPerMemberInWei } = await getAvailableRewardFromGroupInEther(groupPublicKey, groupMemberIndices, web3Context)
      totalRewardsBalance = totalRewardsBalance.add(utils.toBN(utils.toWei(reward, 'ether')))
      groups.push({ groupIndex, groupPublicKey, membersIndeces: groupMemberIndices[groupPublicKey], reward, rewardPerMemberInWei })
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
    membersInGroup.reduce((prev, current, index) => {
      const prevValue = index === 1 ? groupMemberIndices[groupPublicKey][prev].length : prev
      return prevValue + groupMemberIndices[groupPublicKey][current].length
    })
  const groupMemberReward = await keepRandomBeaconOperatorContract.methods.getGroupMemberRewards(groupPublicKey).call()
  const wholeReward = utils.toBN(groupMemberReward).mul(utils.toBN(rewardsMultiplier))

  return { reward: utils.fromWei(wholeReward, 'ether'), rewardPerMemberInWei: groupMemberReward }
}

const withdrawRewardFromGroup = async (groupIndex, groupMembersIndices, web3Context) => {
  const { web3, keepRandomBeaconOperatorContract, yourAddress } = web3Context

  try {
    const batchRequest = new web3.BatchRequest()
    const groupMembers = Object.keys(groupMembersIndices)

    const promises = groupMembers.map((memberAddress) => {
      return new Promise((resolve, reject) => {
        const request = keepRandomBeaconOperatorContract
          .methods
          .withdrawGroupMemberRewards(memberAddress, groupIndex, groupMembersIndices[memberAddress])
          .send.request({ from: yourAddress }, (error, transactionHash) => {
            if (error) {
              resolve({ memberAddress, memberIndices: groupMembersIndices[memberAddress], isError: true, error })
            } else {
              resolve({ transactionHash })
            }
          })
        batchRequest.add(request)
      })
    })

    batchRequest.execute()
    const transactions = await Promise.all(promises)
    const pendingTransactions = transactions.filter((t) => !t.isError)
    let allTransactionsCompleted = !(pendingTransactions.length > 0)

    while (!allTransactionsCompleted) {
      for (let i = 0; i < pendingTransactions.length; i++) {
        const recipt = await web3.eth.getTransactionReceipt(pendingTransactions[i].transactionHash)
        if (!recipt) {
          continue
        }
        const isLastIdex = i === pendingTransactions.length -1
        if (isLastIdex) {
          allTransactionsCompleted = true
        }
      }
      await wait(2000)
    }

    return transactions
  } catch (error) {
    throw error
  }
}

const fetchWithdrawalHistory = async (web3Context) => {
  const { keepRandomBeaconOperatorContract, yourAddress, utils, eth } = web3Context
  const searchFilters = { fromBlock: 0, filter: { beneficiary: yourAddress } }

  try {
    const events = await keepRandomBeaconOperatorContract.getPastEvents('GroupMemberRewardsWithdrawn', searchFilters)
    return Promise.all(
      events.map(async (event) => {
        const { blockNumber, returnValues: { groupIndex, amount } } = event
        const withdrawnAt = (await eth.getBlock(blockNumber)).timestamp
        const groupPublicKey = await keepRandomBeaconOperatorContract.methods.getGroupPublicKey(groupIndex).call()
        return { blockNumber, groupPublicKey, date: formatDate(withdrawnAt * 1000), amount: utils.fromWei(amount, 'ether') }
      }).reverse()
    )
  } catch (error) {
    throw error
  }
}

const rewardsService = {
  fetchAvailableRewards,
  withdrawRewardFromGroup,
  fetchWithdrawalHistory,
}

export default rewardsService
