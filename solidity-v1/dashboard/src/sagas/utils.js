import { call, put, take, race } from "redux-saga/effects"
import { Web3Loaded, ContractsLoaded } from "../contracts"
import { createSubcribeToContractEventChannel } from "./web3"
import { isString } from "../utils/general.utils"
import { modalActions } from "../actions"

/** @typedef { import("web3-eth-contract").Contract} Web3jsContract */
/** @typedef {import("../constants/constants").MODAL_TYPES} ModalTypes */

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
 * @param {Object} options (optional) Options used when calling an event
 */
export function* subscribeToEventAndEmitData(
  contractInstance,
  eventName,
  action,
  debugName,
  options = null
) {
  const contractEventCahnnel = yield call(
    createSubcribeToContractEventChannel,
    contractInstance,
    eventName,
    options
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

export function* confirmModalSaga(modalType, modalProps) {
  yield put(modalActions.openConfirmationModal(modalType, modalProps))
  const { yes } = yield race({
    yes: take(modalActions.CONFIRM),
    no: take(modalActions.CANCEL),
  })
  yield put(modalActions.hideModal())
  const isConfirmed = Boolean(yes)
  return { isConfirmed, task: yes }
}
