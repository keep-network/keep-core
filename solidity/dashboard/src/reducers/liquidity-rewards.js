const wrappedTokenInitialData = {
  shareOfPoolInPercent: 0,
  reward: 0,
  wrappedTokenBalance: 0,
  lpBalance: 0,
  isFetching: false,
  error: null,
}

const initialState = {
  KEEP_ETH: { ...wrappedTokenInitialData },
  TBTC_ETH: { ...wrappedTokenInitialData },
  KEEP_TBTC: { ...wrappedTokenInitialData },
}

const liquidityRewardsReducer = (state = initialState, action) => {
  if (!action.payload) {
    return state
  }

  const { liquidityRewardPair, ...restPayload } = action.payload

  switch (action.type) {
    case `liquidity_rewards/${liquidityRewardPair}_fetch_data_start`:
      return {
        ...state,
        [liquidityRewardPair]: {
          ...state[liquidityRewardPair],
          isFetching: true,
        },
      }

    case `liquidity_rewards/${liquidityRewardPair}_fetch_data_success`:
      return {
        ...state,
        [liquidityRewardPair]: {
          ...state[liquidityRewardPair],
          ...restPayload,
          isFetching: false,
          error: null,
        },
      }
    case `liquidity_rewards/${liquidityRewardPair}_fetch_data_failure`:
      return {
        ...state,
        [liquidityRewardPair]: {
          ...state[liquidityRewardPair],
          isFetching: false,
          error: action.payload.error,
        },
      }
    default:
      return state
  }
}

export default liquidityRewardsReducer
