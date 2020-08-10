import { getContext, call } from "redux-saga/effects"

export function* getWeb3Context() {
  const web3Context = yield getContext("web3")

  return yield web3Context
}

export function* getContractsContext() {
  const contractsContext = yield getContext("contracts")

  return yield contractsContext
}

export function* submitButtonHelper(saga, action) {
  const { resolve, reject } = action.meta

  try {
    yield call(saga, action)
    yield call(resolve, "success")
  } catch (error) {
    yield call(reject, error)
  }
}
