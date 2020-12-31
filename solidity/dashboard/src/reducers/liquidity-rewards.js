import { sub, add, percentageOf } from "../utils/arithmetics.utils"

const liquidityPairInitialData = {
  apy: 0,
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
    case `liquidity_rewards/${liquidityRewardPairName}_staked`: {
      const lpBalance = add(
        state[liquidityRewardPairName].lpBalance,
        restPayload.amount
      ).toString()

      return {
        ...state,
        [liquidityRewardPairName]: {
          ...state[liquidityRewardPairName],
          wrappedTokenBalance: sub(
            state[liquidityRewardPairName].wrappedTokenBalance,
            restPayload.amount
          ).toString(),
          lpBalance,
          shareOfPoolInPercent: percentageOf(
            lpBalance,
            restPayload.totalSupply
          ).toString(),
          reward: restPayload.reward,
          apy: restPayload.apy,
        },
      }
    }
    case `liquidity_rewards/${liquidityRewardPairName}_withdrawn`: {
      const lpBalance = sub(
        state[liquidityRewardPairName].lpBalance,
        restPayload.amount
      ).toString()

      return {
        ...state,
        [liquidityRewardPairName]: {
          ...state[liquidityRewardPairName],
          wrappedTokenBalance: add(
            state[liquidityRewardPairName].wrappedTokenBalance,
            restPayload.amount
          ).toString(),
          lpBalance,
          shareOfPoolInPercent: percentageOf(
            lpBalance,
            restPayload.reward
          ).toString(),
          reward: restPayload.reward,
          apy: restPayload.apy,
        },
      }
    }
    case `liquidity_rewards/${liquidityRewardPairName}_reward_paid`:
      return {
        ...state,
        [liquidityRewardPairName]: {
          ...state[liquidityRewardPairName],
          reward: "0",
        },
      }
    case `liquidity_rewards/${liquidityRewardPairName}_apy_updated`:
      return {
        ...state,
        [liquidityRewardPairName]: {
          ...state[liquidityRewardPairName],
          apy: restPayload.apy,
        },
      }
    default:
      return state
  }
}

export default liquidityRewardsReducer
