import { takeLatest, call } from "redux-saga/effects"
import { getContractsContext, getWeb3Context, submitButtonHelper } from "./utils"
import { createBondERC20Contract } from "../contracts"
import { sendTransaction } from "./web3"

function* depositForOperator({ payload: { amount, operatorAddress }}) {
  const { keepBondingContract } = yield getContractsContext()

  const web3 = yield getWeb3Context()

  const bondTokenAddress = yield call(keepBondingContract.methods.bondTokenAddress().call)
  const bondToken = yield call(createBondERC20Contract, web3, bondTokenAddress)

  yield call(sendTransaction, {
    payload: {
      contract: bondToken,
      methodName: "approve",
      args: [keepBondingContract.options.address, amount],
    },
  })

  yield call(sendTransaction, {
    payload: {
      contract: keepBondingContract,
      methodName: "deposit",
      args: [operatorAddress, amount],
    },
  })
}

function * deposit (action) {
  yield call(submitButtonHelper, depositForOperator, action)
}

export function* watchDepositForOperator() {
  yield takeLatest("bonding/deposit_start", deposit)
}