import { takeLatest, takeEvery, fork, call, put } from "redux-saga/effects"
import { getContractsContext, submitButtonHelper, logError } from "./utils"
import { sendTransaction } from "./web3"
import {
  fetchStakedBalance,
  fetchWrappedTokenBalance,
  fetchLPRewardsTotalSupply,
  fetchRewardBalance,
  getWrappedTokenConctract, calculateAPY, fetchTotalLPTokensCreatedInUniswap,
} from "../services/liquidity-rewards"
import { gt, percentageOf, eq } from "../utils/arithmetics.utils"
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
    const apy = 0
    // Fetching total deposited liqidity tokens in the `LPRewards` contract.
    const totalSupply = yield call(fetchLPRewardsTotalSupply, LPRewardsContract)

    // if (gt(totalSupply, 0)) {
    //   const totalLPTokensCreatedInUniswap = yield call(
    //     fetchTotalLPTokensCreatedInUniswap,
    //     LPRewardsContract
    //   )
    //
    //   apy = calculateAPY(
    //     totalSupply,
    //     totalLPTokensCreatedInUniswap,
    //     liquidityRewardPair.name
    //   )
    // }

    let reward = 0
    let shareOfPoolInPercent = 0

    if (gt(lpBalance, 0)) {
      // Fetching available reward balance from `LPRewards` contract.
      reward = yield call(fetchRewardBalance, address, LPRewardsContract)

      // % of total pool in the `LPRewards` contract.
      shareOfPoolInPercent = percentageOf(lpBalance, totalSupply).toString()
    }

    yield put({
      type: `liquidity_rewards/${liquidityRewardPair.name}_fetch_data_success`,
      payload: {
        liquidityRewardPairName: liquidityRewardPair.name,
        apy,
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

function* stakeTokens(action) {
  const { contractName, address, amount } = action.payload

  const contracts = yield getContractsContext()
  const LPRewardsContract = contracts[contractName]
  const lpRewardsContractAddress = LPRewardsContract.options.address

  const WrappedTokenContract = yield call(
    getWrappedTokenConctract,
    LPRewardsContract
  )

  const approvedAmount = yield call(
    WrappedTokenContract.methods.allowance(address, lpRewardsContractAddress)
      .call
  )

  if (!eq(amount, approvedAmount)) {
    yield call(sendTransaction, {
      payload: {
        contract: WrappedTokenContract,
        methodName: "approve",
        args: [lpRewardsContractAddress, amount],
      },
    })
  }

  yield call(sendTransaction, {
    payload: {
      contract: LPRewardsContract,
      methodName: "stake",
      args: [amount],
    },
  })
}

function* stakeTokensWorker(action) {
  yield call(submitButtonHelper, stakeTokens, action)
}

export function* watchStakeTokens() {
  yield takeEvery("liquidity_rewards/stake_tokens", stakeTokensWorker)
}
