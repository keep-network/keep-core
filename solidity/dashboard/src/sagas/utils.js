import { getContext } from "redux-saga/effects"

export function* getWeb3Context() {
  const web3Context = yield getContext("web3")

  return yield web3Context
}

export function* getContractsContext() {
  const contractsContext = yield getContext("contracts")

  return yield contractsContext
}
