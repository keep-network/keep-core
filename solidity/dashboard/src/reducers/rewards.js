import { sub } from "../utils/arithmetics.utils"
import { isSameEthAddress } from "../utils/general.utils"

const initialState = {
  // Beacon distributted rewards
  beaconRewardsFetching: false,
  becaonRewardsBalance: 0,
  error: null,

  // ECDSA available rewards
  ecdsaAvailableRewardsFetching: false,
  ecdsaAvailableRewardsBalance: 0,
  ecdsaAvailableRewards: [],
  ecdsaAvailableRewardsError: null,
  ecdsaRewardsHistory: [],
}

const rewardsReducer = (state = initialState, action) => {
  switch (action.type) {
    case "rewards/beacon_fetch_distributed_rewards_start":
      return { ...state, beaconRewardsFetching: true, error: null }
    case "rewards/beacon_fetch_distributed_rewards_success":
      return {
        ...state,
        beaconRewardsFetching: false,
        becaonRewardsBalance: action.payload,
      }
    case "rewards/beacon_fetch_distributed_rewards_failure":
      return {
        ...state,
        beaconRewardsFetching: false,
        error: action.payload.error,
      }
    case "rewards/ecdsa_fetch_rewards_data_start":
      return {
        ...state,
        ecdsaAvailableRewardsFetching: true,
        ecdsaAvailableRewardsError: null,
      }
    case "rewards/ecdsa_fetch_rewards_data_success":
      return {
        ...state,
        ecdsaAvailableRewardsFetching: false,
        ecdsaAvailableRewardsBalance: action.payload.totalAvailableAmount,
        ecdsaAvailableRewards: action.payload.availableRewards,
        ecdsaRewardsHistory: action.payload.rewardsHistory,
      }
    case "rewards/ecdsa_fetch_rewards_data_failure":
      return {
        ...state,
        ecdsaAvailableRewardsFetching: false,
        ecdsaAvailableRewardsError: action.payload.error,
      }
    case "rewards/ecdsa_withdrawn":
      return ecdsaRewardsWithdrawn({ ...state }, action.payload)
    case "rewards/ecdsa_update_available_rewards":
      return {
        ...state,
        ecdsaAvailableRewards: action.payload,
      }
    default:
      return state
  }
}

const ecdsaRewardsWithdrawn = (state, { merkleRoot, operator, amount }) => {
  const isSameRewardRecord = (reward) =>
    reward.merkleRoot === merkleRoot &&
    isSameEthAddress(reward.operator, operator)

  const reward = state.ecdsaAvailableRewards.find(isSameRewardRecord)

  if (!reward) {
    return state
  }

  const ecdsaAvailableRewardsBalance = sub(
    state.ecdsaAvailableRewardsBalance,
    amount
  )

  const ecdsaAvailableRewards = state.ecdsaAvailableRewards.filter(
    (reward) => !isSameRewardRecord(reward)
  )

  const ecdsaRewardsHistory = state.ecdsaRewardsHistory.map((reward) => {
    if (isSameRewardRecord(reward)) {
      return { ...reward, status: "WITHDRAWN" }
    }

    return reward
  })

  return {
    ...state,
    ecdsaAvailableRewardsBalance,
    ecdsaAvailableRewards,
    ecdsaRewardsHistory,
  }
}

export default rewardsReducer
