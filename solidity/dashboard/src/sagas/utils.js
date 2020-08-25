import { call } from "redux-saga/effects"
import { Web3Loaded, ContractsLoaded } from "../contracts"

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
