import { takeLatest, fork, call, put } from "redux-saga/effects"
import BigNumber from "bignumber.js"
import { getContractsContext, logError } from "./utils"
import {
  fetchStakedBalance,
  fetchWrappedTokenBalance,
  fetchLPRewardsTotalSupply,
  fetchRewardBalance,
} from "../services/liquidity-rewards"
import { gt } from "../utils/arithmetics.utils"
import { LIQUIDITY_REWARD_PAIRS } from "../constants/constants"

function* fetchAllLiquidtyRewardsData(action) {
  const { address } = action.payload

  for (const [pairName, value] of Object.entries(LIQUIDITY_REWARD_PAIRS)) {
    yield fork(fetchLiquidityRewardsData, { name: pairName, ...value }, address)
  }
}

function* fetchLiquidityRewardsData(liquidityRewardPair, address) {
  const contracts = yield getContractsContext()

  const LPRewardsContract = contracts[liquidityRewardPair.contractName]
  try {
    yield put({
      type: `liquidity_rewards/${liquidityRewardPair.name}_fetch_data_start`,
      payload: { liquidityRewardPairName: liquidityRewardPair.name },
    })
    // Fetching balance of liquidity token for a given uniswap pair deposited in
    // the `LPRewards` contract.
    const lpBalance = yield call(fetchStakedBalance, address, LPRewardsContract)
    // Fetching balance of liquidity token for a given uniswap pair.
    const wrappedTokenBalance = yield call(
      fetchWrappedTokenBalance,
      address,
      LPRewardsContract
    )
    let reward = 0
    let shareOfPoolInPercent = 0
    if (gt(lpBalance, 0)) {
      // Fetching available reward balance from `LPRewards` contract.
      reward = yield call(fetchRewardBalance, address, LPRewardsContract)
      // Fetching total deposited liqidity tokens in the `LPRewards` contract.
      const totalSupply = yield call(
        fetchLPRewardsTotalSupply,
        LPRewardsContract
      )
      // % of total pool in the `LPRewards` contract.
      shareOfPoolInPercent = new BigNumber(lpBalance)
        .div(totalSupply)
        .multipliedBy(100)
        .toString()
    }
    yield put({
      type: `liquidity_rewards/${liquidityRewardPair.name}_fetch_data_success`,
      payload: {
        liquidityRewardPairName: liquidityRewardPair.name,
        lpBalance,
        wrappedTokenBalance,
        reward,
        shareOfPoolInPercent,
      },
    })
  } catch (error) {
    yield* logError(
      `liquidity_rewards/${liquidityRewardPair.name}_fetch_data_failure`,
      error,
      { liquidityRewardPairName: liquidityRewardPair.name }
    )
  }
}

export function* watchFetchLiquidityRewardsData() {
  yield takeLatest(
    "liquidity_rewards/fetch_data_request",
    fetchAllLiquidtyRewardsData
  )
}
