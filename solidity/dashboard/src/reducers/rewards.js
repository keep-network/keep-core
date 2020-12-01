import { sub, add } from "../utils/arithmetics.utils"

const initialState = {
  // Beacon distributted rewards
  beaconRewardsFetching: false,
  becaonRewardsBalance: 0,
  error: null,

  // ECDSA distributed rewards
  ecdsaDistributedRewardsFetching: false,
  ecdsaDistributedBalance: 0,
  ecdsaDistributedBalanceError: null,

  // ECDSA available rewards
  ecdsaAvailableRewardsFetching: false,
  ecdsaAvailableRewardsBalance: 0,
  ecdsaAvailableRewards: [],
  ecdsaAvailableRewardsError: null,
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
    case "rewards/ecdsa_fetch_distributed_rewards_start":
      return {
        ...state,
        ecdsaDistributedRewardsFetching: true,
        ecdsaDistributedBalanceError: null,
      }
    case "rewards/ecdsa_fetch_distributed_rewards_success":
      return {
        ...state,
        ecdsaDistributedRewardsFetching: false,
        ecdsaDistributedBalance: action.payload,
      }
    case "rewards/ecdsa_fetch_distributed_rewards_failure":
      return {
        ...state,
        ecdsaDistributedRewardsFetching: false,
        ecdsaDistributedBalanceError: action.payload.error,
      }
    case "rewards/ecdsa_fetch_available_rewards_start":
      return {
        ...state,
        ecdsaAvailableRewardsFetching: true,
        ecdsaAvailableRewardsError: null,
      }
    case "rewards/ecdsa_fetch_available_rewards_success":
      return {
        ...state,
        ecdsaAvailableRewardsFetching: false,
        ecdsaAvailableRewardsBalance: action.payload.totalAvailableRewards,
        ecdsaAvailableRewards: action.payload.toWithdrawn,
      }
    case "rewards/ecdsa_fetch_available_rewards_failure":
      return {
        ...state,
        ecdsaAvailableRewardsFetching: false,
        ecdsaAvailableRewardsError: action.payload.error,
      }
    case "rewards/ecdsa_withdrawn":
      return {
        ...state,
        ecdsaAvailableRewardsBalance: sub(
          state.ecdsaAvailableRewardsBalance,
          action.payload
        ),
        ecdsaDistributedBalance: add(
          state.ecdsaDistributedBalance,
          action.payload
        ),
      }
    case "rewards/ecdsa_update_available_rewards":
      return {
        ...state,
        ecdsaAvailableRewards: action.payload,
      }
    default:
      return state
  }
}

export default rewardsReducer
