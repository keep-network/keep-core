import web3Utils from "web3-utils"
import { add, gt } from "../utils/arithmetics.utils"
import { CONTRACT_DEPLOY_BLOCK_NUMBER } from "../contracts"
import {
  OPERATOR_CONTRACT_NAME,
  REWARD_STATUS,
  SIGNING_GROUP_STATUS,
} from "../constants/constants"
import { contractService } from "./contracts.service"
import { getOperatorsOfBeneficiary } from "./token-staking.service"

const fetchAvailableRewards = async (web3Context) => {
  const {
    keepRandomBeaconOperatorContract,
    keepRandomBeaconOperatorStatistics,
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

    const operatorsOfBeneficiary = await getOperatorsOfBeneficiary(
      web3Context,
      yourAddress
    )
    const rewards = []
    // { groupIndex: { isStale, isTerminated, groupPubKey } }
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
          reward: web3Utils.fromWei(awaitingRewards, "ether"),
        })
      }
    }
    return [rewards, web3Utils.fromWei(totalRewardsBalance, "ether")]
  } catch (error) {
    throw error
  }
}

const withdrawRewardFromGroup = async (
  web3Context,
  data,
  onTransactionHash
) => {
  const { keepRandomBeaconOperatorContract, yourAddress } = web3Context
  const { operatorAddress, groupIndex } = data

  await keepRandomBeaconOperatorContract.methods
    .withdrawGroupMemberRewards(operatorAddress, groupIndex)
    .send({ from: yourAddress })
    .on("transactionHash", onTransactionHash)
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
  withdrawRewardFromGroup,
  fetchWithdrawalHistory,
}

export default rewardsService
