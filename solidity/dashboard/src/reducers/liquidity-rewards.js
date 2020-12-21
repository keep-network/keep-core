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
  const { wrappedToken, ...restPayload } = action.payload
  switch (action.type) {
    case `liquidity_rewards/${wrappedToken}_fetch_data_start`:
      return {
        ...state,
        [wrappedToken]: { ...state[wrappedToken], isFetching: true },
      }

    case `liquidity_rewards/${wrappedToken}_fetch_data_success`:
      return {
        ...state,
        [wrappedToken]: {
          ...state[wrappedToken],
          ...restPayload,
          isFetching: true,
          error: null,
        },
      }
    case `liquidity_rewards/${wrappedToken}_fetch_data_failure`:
      return {
        ...state,
        [wrappedToken]: {
          ...state[wrappedToken],
          isFetching: false,
          error: action.payload.error,
        },
      }
    default:
      return state
  }
}

export default liquidityRewardsReducer
