import { takeLatest, call, all } from "redux-saga/effects"
import { sendTransaction } from "./web3"

import { createBondERC20Contract } from "../contracts"
import { getContractsContext, getWeb3Context, submitButtonHelper } from "./utils"

function* mintBondTokens({ payload: { amount, address } }) {
  const web3 = yield getWeb3Context()
  const { keepBondingContract } = yield getContractsContext()

  const bondTokenAddress = yield call(keepBondingContract.methods.bondTokenAddress().call)
  const bondToken = yield call(createBondERC20Contract, web3, bondTokenAddress)

  yield call(sendTransaction, {
    payload: {
      contract: bondToken,
      methodName: "mint",
      args: [address, amount],
    }
  })
}

function* mintBondTokensWrapper(action) {
  yield call(submitButtonHelper, mintBondTokens, action)
}

function* transferKeepTokens({ payload: { amount, address } }) {
  const { token } = yield getContractsContext()

  yield call(sendTransaction, {
    payload: {
      contract: token,
      methodName: "transfer",
      args: [address, amount],
    }
  })
}

function* transferKeepTokensWrapper(action) {
  yield call(submitButtonHelper, transferKeepTokens, action)
}

export function* debugMintingSaga() {
  yield all([
    takeLatest("debug-minting/mint-bondTokens", mintBondTokensWrapper),
    takeLatest("debug-minting/transfer-keep-tokens", transferKeepTokensWrapper)
  ])
}