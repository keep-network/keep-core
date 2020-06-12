import web3Utils from "web3-utils"
import { wait, isSameEthAddress } from "../utils/general.utils"
import { add, gt } from "../utils/arithmetics.utils"
import { CONTRACT_DEPLOY_BLOCK_NUMBER } from "../contracts"
import { OPERATOR_CONTRACT_NAME } from "../constants/constants"
import { contractService } from "./contracts.service"

const fetchAvailableRewards = async (web3Context) => {
  const {
    keepRandomBeaconOperatorContract,
    keepRandomBeaconOperatorStatistics,
    stakingContract,
    yourAddress,
  } = web3Context
  try {
    let totalRewardsBalance = web3Utils.toBN(0)
    const operatorEventsSearchFilters = {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[OPERATOR_CONTRACT_NAME],
    }

    // get all created groups
    const groupPubKeys = (
      await contractService.getPastEvents(
        web3Context,
        OPERATOR_CONTRACT_NAME,
        "DkgResultSubmittedEvent",
        operatorEventsSearchFilters
      )
    ).map((event) => event.returnValues.groupPubKey)
    const rewards = []
    const groupMemberIndices = {}

    for (let groupIndex = 0; groupIndex < groupPubKeys.length; groupIndex++) {
      const groupPubKey = groupPubKeys[groupIndex]
      const groupMembers = new Set(
        await keepRandomBeaconOperatorContract.methods
          .getGroupMembers(groupPubKey)
          .call()
      )
      groupMemberIndices[groupPubKey] = {}
      for (const memberAddress of groupMembers) {
        const beneficiaryAddressForMember = await stakingContract.methods
          .beneficiaryOf(memberAddress)
          .call()

        if (!isSameEthAddress(yourAddress, beneficiaryAddressForMember)) {
          continue
        }

        const awaitingRewards = await keepRandomBeaconOperatorStatistics.methods
          .awaitingRewards(memberAddress, groupIndex)
          .call()

        if (gt(awaitingRewards, 0)) {
          groupMemberIndices[groupPubKey][memberAddress] = awaitingRewards
        }
      }

      if (Object.keys(groupMemberIndices[groupPubKey]).length !== 0) {
        const reward = getAvailableRewardForGroup(
          groupMemberIndices[groupPubKey]
        )
        const isStale = await keepRandomBeaconOperatorContract.methods
          .isStaleGroup(groupPubKey)
          .call()

        totalRewardsBalance = add(totalRewardsBalance, reward)
        rewards.push({
          groupIndex: groupIndex.toString(),
          groupPublicKey: groupPubKey,
          membersIndeces: groupMemberIndices[groupPubKey],
          reward: web3Utils.fromWei(reward, "ether"),
          isStale,
        })
      }
    }
    return [rewards, web3Utils.fromWei(totalRewardsBalance, "ether")]
  } catch (error) {
    throw error
  }
}

const getAvailableRewardForGroup = (operatorsAmount) => {
  let wholeReward = 0
  for (const operator in operatorsAmount) {
    if (operatorsAmount.hasOwnProperty(operator)) {
      wholeReward = add(wholeReward, operatorsAmount[operator])
    }
  }
  return wholeReward
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
          .withdrawGroupMemberRewards(memberAddress, groupIndex)
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
  const { keepRandomBeaconOperatorContract, yourAddress, utils } = web3Context
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
            returnValues: { groupIndex, amount, operator },
          } = event
          const groupPublicKey = await keepRandomBeaconOperatorContract.methods
            .getGroupPublicKey(groupIndex)
            .call()
          return {
            blockNumber,
            groupPublicKey,
            reward: utils.fromWei(amount, "ether"),
            transactionHash,
            operator,
            status: "WITHDRAWN",
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
