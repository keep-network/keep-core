import { call, put } from "redux-saga/effects"
import { Web3Loaded, ContractsLoaded } from "../contracts"
import { LiquidityRewardsFactory } from "../services/liquidity-rewards"
/** @typedef { import("../services/liquidity-rewards").LiquidityRewards} LiquidityRewards */

export function* getWeb3Context() {
  return yield Web3Loaded
}

export function* getContractsContext() {
  return yield ContractsLoaded
}

export function* submitButtonHelper(saga, action) {
  const { resolve, reject } = action.meta

  try {
    yield call(saga, action)
    yield call(resolve, "success")
  } catch (error) {
    console.error(error)
    yield call(reject, error)
  }
}

export function* logError(errorActionType, error, payload = {}) {
  const { message, reason, stack } = error
  yield put({
    type: errorActionType,
    payload: {
      ...payload,
      error: reason ? `Error: ${reason}` : message,
    },
  })
  console.error({
    reason,
    message,
    originalStack: stack.split("\n").map((s) => s.trim()),
    error,
  })
}

export function* logErrorAndThrow(errorActionType, error, payload = {}) {
  yield* logError(errorActionType, error, payload)
  throw error
}

/**
 *
 * @param {Object} liquidityRewardPair - Liquidity reward data.
 * @param {string} liquidityRewardPair.pool - The type of pool.
 * @param {string} liquidityRewardPair.contractName - The LPRewards contract
 * name for a given liquidity pair.
 * @return {LiquidityRewards} Liquidity rewards wrapper.
 */
export function* getLPRewardsWrapper(liquidityRewardPair) {
  const contracts = yield getContractsContext()
  const web3 = yield getWeb3Context()

  const LPRewardsContract = contracts[liquidityRewardPair.contractName]
  const LiquidityRewards = yield call(
    [LiquidityRewardsFactory, LiquidityRewardsFactory.initialize],
    liquidityRewardPair.pool,
    LPRewardsContract,
    web3
  )

  return LiquidityRewards
}
