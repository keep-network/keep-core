import { call, put, take } from "redux-saga/effects"
import { Web3Loaded, ContractsLoaded } from "../contracts"
import { LiquidityRewardsFactory } from "../services/liquidity-rewards"
import { createSubcribeToContractEventChannel } from "./web3"
import { isString } from "../utils/general.utils"

/** @typedef { import("../services/liquidity-rewards").LiquidityRewards} LiquidityRewards */
/** @typedef { import("web3-eth-contract").Contract} Web3jsContract */

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

export const identifyTaskByAddress = (action) =>
  identifyTaskBy("address")(action)

export const identifyTaskBy = (indentificationField) => (action) =>
  action.payload[indentificationField]

/**
 * A helper saga that subscribes to the contract event and emits an action with
 * the event data.
 * @param {Web3jsContract} contractInstance A web3 js contract instance.
 * @param {string} eventName The contract event name.
 * @param {string | Function} action The action to be emitted to the redux.
 * @param {string} debugName A debug event name.
 */
export function* subscribeToEventAndEmitData(
  contractInstance,
  eventName,
  action,
  debugName
) {
  const contractEventCahnnel = yield call(
    createSubcribeToContractEventChannel,
    contractInstance,
    eventName
  )
  const _debugName = debugName || eventName
  const _isString = isString(action)

  while (true) {
    try {
      const event = yield take(contractEventCahnnel)
      const _action = _isString
        ? { type: action, payload: { event } }
        : action(event)
      yield put(_action)
    } catch (error) {
      console.error(`Failed subscribing to ${_debugName} event`, error)
      contractEventCahnnel.close()
    }
  }
}
