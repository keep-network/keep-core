const initialState = {
  beaconRewardsFetching: false,
  becaonRewardsBalance: 0,
  error: null,
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
        eror: action.payload.error,
      }
    default:
      return state
  }
}

export default rewardsReducer
