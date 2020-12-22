const liquidityPairInitialData = {
  shareOfPoolInPercent: 0,
  reward: 0,
  wrappedTokenBalance: 0,
  lpBalance: 0,
  isFetching: false,
  error: null,
}

const initialState = {
  KEEP_ETH: { ...liquidityPairInitialData },
  TBTC_ETH: { ...liquidityPairInitialData },
  KEEP_TBTC: { ...liquidityPairInitialData },
}

const liquidityRewardsReducer = (state = initialState, action) => {
  if (!action.payload) {
    return state
  }

  const { liquidityRewardPairName, ...restPayload } = action.payload

  switch (action.type) {
    case `liquidity_rewards/${liquidityRewardPairName}_fetch_data_start`:
      return {
        ...state,
        [liquidityRewardPairName]: {
          ...state[liquidityRewardPairName],
          isFetching: true,
        },
      }

    case `liquidity_rewards/${liquidityRewardPairName}_fetch_data_success`:
      return {
        ...state,
        [liquidityRewardPairName]: {
          ...state[liquidityRewardPairName],
          ...restPayload,
          isFetching: false,
          error: null,
        },
      }
    case `liquidity_rewards/${liquidityRewardPairName}_fetch_data_failure`:
      return {
        ...state,
        [liquidityRewardPairName]: {
          ...state[liquidityRewardPairName],
          isFetching: false,
          error: action.payload.error,
        },
      }
    default:
      return state
  }
}

export default liquidityRewardsReducer
