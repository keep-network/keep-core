import web3Utils from "web3-utils"
import { formatDate, wait, isSameEthAddress } from "../utils/general.utils"
import { add, mul } from "../utils/arithmetics.utils"
import { CONTRACT_DEPLOY_BLOCK_NUMBER } from "../contracts"
import { OPERATOR_CONTRACT_NAME } from "../constants/constants"

const fetchAvailableRewards = async (web3Context) => {
  const {
    keepRandomBeaconOperatorContract,
    stakingContract,
    yourAddress,
  } = web3Context
  try {
    let totalRewardsBalance = web3Utils.toBN(0)
    const expiredGroupsCount = await keepRandomBeaconOperatorContract.methods
      .getFirstActiveGroupIndex()
      .call()
    const acitveGroups = await keepRandomBeaconOperatorContract.methods
      .numberOfGroups()
      .call()
    const allGroups = add(expiredGroupsCount, acitveGroups).toNumber()
    const groups = []
    const groupMemberIndices = {}

    for (let groupIndex = 0; groupIndex < allGroups; groupIndex++) {
      const groupPublicKey = await keepRandomBeaconOperatorContract.methods
        .getGroupPublicKey(groupIndex)
        .call()
      const groupMembers = new Set(
        await keepRandomBeaconOperatorContract.methods
          .getGroupMembers(groupPublicKey)
          .call()
      )
      groupMemberIndices[groupPublicKey] = {}
      for (const memberAddress of groupMembers) {
        const beneficiaryAddressForMember = await stakingContract.methods
          .magpieOf(memberAddress)
          .call()
        if (!isSameEthAddress(yourAddress, beneficiaryAddressForMember)) {
          continue
        }
        groupMemberIndices[groupPublicKey][
          memberAddress
        ] = await keepRandomBeaconOperatorContract.methods
          .getGroupMemberIndices(groupPublicKey, memberAddress)
          .call()
      }
      if (Object.keys(groupMemberIndices[groupPublicKey]).length === 0) {
        continue
      }
      const {
        reward,
        rewardPerMemberInWei,
      } = await getAvailableRewardFromGroupInEther(
        groupPublicKey,
        groupMemberIndices,
        web3Context
      )
      const isStale = await keepRandomBeaconOperatorContract.methods
        .isStaleGroup(groupPublicKey)
        .call()

      totalRewardsBalance = add(
        totalRewardsBalance,
        web3Utils.toWei(reward, "ether")
      )
      groups.push({
        groupIndex: groupIndex.toString(),
        groupPublicKey,
        membersIndeces: groupMemberIndices[groupPublicKey],
        reward,
        rewardPerMemberInWei,
        isStale,
      })
    }
    return [groups, web3Utils.fromWei(totalRewardsBalance.toString(), "ether")]
  } catch (error) {
    throw error
  }
}

const getAvailableRewardFromGroupInEther = async (
  groupPublicKey,
  groupMemberIndices,
  web3Context
) => {
  const { keepRandomBeaconOperatorContract } = web3Context
  const membersInGroup = Object.keys(groupMemberIndices[groupPublicKey])
  const rewardsMultiplier =
    membersInGroup.length === 1
      ? groupMemberIndices[groupPublicKey][membersInGroup[0]].length
      : membersInGroup.reduce((prev, current, index) => {
          const prevValue =
            index === 1 ? groupMemberIndices[groupPublicKey][prev].length : prev
          return prevValue + groupMemberIndices[groupPublicKey][current].length
        })
  const groupMemberReward = await keepRandomBeaconOperatorContract.methods
    .getGroupMemberRewards(groupPublicKey)
    .call()
  const wholeReward = mul(groupMemberReward, rewardsMultiplier)

  return {
    reward: web3Utils.fromWei(wholeReward, "ether"),
    rewardPerMemberInWei: groupMemberReward,
  }
}

const withdrawRewardFromGroup = async (
  groupIndex,
  groupMembersIndices,
  web3Context
) => {
  const { web3, keepRandomBeaconOperatorContract, yourAddress } = web3Context

  try {
    const batchRequest = new web3.BatchRequest()
    const groupMembers = Object.keys(groupMembersIndices)

    const promises = groupMembers.map((memberAddress) => {
      return new Promise((resolve, reject) => {
        const request = keepRandomBeaconOperatorContract.methods
          .withdrawGroupMemberRewards(
            memberAddress,
            groupIndex,
            groupMembersIndices[memberAddress]
          )
          .send.request({ from: yourAddress }, (error, transactionHash) => {
            if (error) {
              resolve({
                memberAddress,
                memberIndices: groupMembersIndices[memberAddress],
                isError: true,
                error,
              })
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
        const recipt = await web3.eth.getTransactionReceipt(
          pendingTransactions[i].transactionHash
        )
        if (!recipt) {
          continue
        }
        const isLastIdex = i === pendingTransactions.length - 1
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
  const {
    keepRandomBeaconOperatorContract,
    yourAddress,
    utils,
    eth,
  } = web3Context
  const searchFilters = {
    fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[OPERATOR_CONTRACT_NAME],
    filter: { beneficiary: yourAddress },
  }

  try {
    const events = await keepRandomBeaconOperatorContract.getPastEvents(
      "GroupMemberRewardsWithdrawn",
      searchFilters
    )
    return Promise.all(
      events
        .map(async (event) => {
          const {
            transactionHash,
            blockNumber,
            returnValues: { groupIndex, amount },
          } = event
          const withdrawnAt = (await eth.getBlock(blockNumber)).timestamp
          const groupPublicKey = await keepRandomBeaconOperatorContract.methods
            .getGroupPublicKey(groupIndex)
            .call()
          return {
            blockNumber,
            groupPublicKey,
            date: formatDate(withdrawnAt * 1000),
            amount: utils.fromWei(amount, "ether"),
            transactionHash,
          }
        })
        .reverse()
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
