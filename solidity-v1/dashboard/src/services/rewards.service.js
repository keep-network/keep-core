import web3Utils from "web3-utils"
import { add, gt } from "../utils/arithmetics.utils"
import { getContractDeploymentBlockNumber, ContractsLoaded } from "../contracts"
import {
  OPERATOR_CONTRACT_NAME,
  REWARD_STATUS,
  SIGNING_GROUP_STATUS,
} from "../constants/constants"
import { getOperatorsOfBeneficiary } from "./token-staking.service"

const fetchAvailableRewards = async (address) => {
  const {
    keepRandomBeaconOperatorContract,
    keepRandomBeaconOperatorStatistics,
  } = await ContractsLoaded

  try {
    let totalRewardsBalance = web3Utils.toBN(0)
    const operatorEventsSearchFilters = {
      fromBlock: await getContractDeploymentBlockNumber(OPERATOR_CONTRACT_NAME),
    }

    // get all created groups
    const groupPubKeys = (
      await keepRandomBeaconOperatorContract.getPastEvents(
        "DkgResultSubmittedEvent",
        operatorEventsSearchFilters
      )
    ).map((event) => event.returnValues.groupPubKey)

    const operatorsOfBeneficiary = await getOperatorsOfBeneficiary(address)
    const rewards = []
    const groups = {}

    for (let groupIndex = 0; groupIndex < groupPubKeys.length; groupIndex++) {
      const groupPubKey = groupPubKeys[groupIndex]
      for (const memberAddress of operatorsOfBeneficiary) {
        const awaitingRewards = await keepRandomBeaconOperatorStatistics.methods
          .awaitingRewards(memberAddress, groupIndex)
          .call()

        if (!gt(awaitingRewards, 0)) {
          continue
        }

        let groupInfo = {}
        if (groups.hasOwnProperty(groupIndex)) {
          groupInfo = { ...groups[groupIndex] }
        } else {
          const isStale = await keepRandomBeaconOperatorContract.methods
            .isStaleGroup(groupPubKey)
            .call()

          const isTerminated =
            !isStale &&
            (await keepRandomBeaconOperatorContract.methods
              .isGroupTerminated(groupIndex)
              .call())

          let status = REWARD_STATUS.ACCUMULATING
          let groupStatus = SIGNING_GROUP_STATUS.ACTIVE
          if (isTerminated) {
            status = REWARD_STATUS.ACCUMULATING
            groupStatus = SIGNING_GROUP_STATUS.TERMINATED
          } else if (isStale) {
            status = REWARD_STATUS.AVAILABLE
            groupStatus = SIGNING_GROUP_STATUS.COMPLETED
          }

          groupInfo = {
            groupPublicKey: groupPubKey,
            isStale,
            status,
            groupStatus,
          }

          groups[groupIndex] = groupInfo
        }

        totalRewardsBalance = add(totalRewardsBalance, awaitingRewards)
        rewards.push({
          groupIndex: groupIndex.toString(),
          ...groupInfo,
          operatorAddress: memberAddress,
          reward: awaitingRewards,
        })
      }
    }
    return [rewards, totalRewardsBalance]
  } catch (error) {
    throw error
  }
}

const fetchWithdrawalHistory = async (address) => {
  const { keepRandomBeaconOperatorContract } = await ContractsLoaded
  const searchFilters = {
    fromBlock: await getContractDeploymentBlockNumber(OPERATOR_CONTRACT_NAME),
    filter: { beneficiary: address },
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
            reward: amount,
            transactionHash,
            operatorAddress: operator,
            status: REWARD_STATUS.WITHDRAWN,
            groupStatus: SIGNING_GROUP_STATUS.COMPLETED,
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
  fetchWithdrawalHistory,
}

export default rewardsService
