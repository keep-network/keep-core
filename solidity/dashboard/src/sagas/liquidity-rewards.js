import { takeLatest, call, put } from "redux-saga/effects"
import BigNumber from "bignumber.js"
import { logError } from "./utils"
import {
  fetchLPRewardsBalance,
  fetchUniTokenBalance,
} from "../services/liquidity_rewards"
import { gt } from "../utils/arithmetics.utils"

function* fetchLiquidityRewardsData(payload) {
  const { wrappedToken, address } = payload

  try {
    yield put({ type: `liquidity_rewards/${wrappedToken}_fetch_data_start` })
    // Fetching balance of liquidity token for a given uniswap pair deposited in
    // the `LPRewards` contract.
    const lpBalance = yield call(fetchLPRewardsBalance, address, wrappedToken)
    // Fetching balance of liquidity token for a given uniswap pair.
    const uniswapTokenBalance = yield call(
      fetchUniTokenBalance,
      address,
      wrappedToken
    )
    let reward = 0
    let shareOfPoolInPercent = 0
    if (gt(lpBalance, 0)) {
      // Fetching available reward balance from `LPRewards` contract.
      reward = yield call(fetchLPRewardsBalance, address, wrappedToken)
      // Fetching total deposited liqidity tokens in the `LPRewards` contract.
      const totalSupply = yield call(
        fetchLPRewardsTotalSupply,
        address,
        wrappedToken
      )
      // % of total pool in the `LPRewards` contract.
      shareOfPoolInPercent = new BigNumber(lpBalance)
        .div(totalSupply)
        .multipliedBy(100)
        .toFixed(2, BigNumber.ROUND_DOWN)
    }
    yield put({
      type: `liquidity_rewards/${wrappedToken}_fetch_data_success`,
      payload: {
        lpBalance,
        uniswapTokenBalance,
        reward,
        shareOfPoolInPercent,
      },
    })
  } catch (error) {
    yield* logError(
      `liquidity_rewards/${wrappedToken}_fetch_data_failure`,
      error
    )
  }
}

export function* watchFetchLiquidityRewardsData() {
  yield takeLatest(
    "liquidity_rewards/fetch_data_request",
    fetchLiquidityRewardsData
  )
}
