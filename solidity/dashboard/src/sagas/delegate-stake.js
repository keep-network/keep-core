import { takeLatest, call } from "redux-saga/effects"
import web3Utils from "web3-utils"
import { getContractsContext } from "./utils"
import { sendTransaction } from "./web3"

function* delegateStake(action) {
  const { token, stakingContract } = yield getContractsContext()
  const {
    stakeTokens,
    beneficiaryAddress,
    operatorAddress,
    authorizerAddress,
  } = action.payload
  const { setSubmitting, resetForm } = action.meta
  const amount = web3Utils
    .toBN(stakeTokens)
    .mul(web3Utils.toBN(10).pow(web3Utils.toBN(18)))
    .toString()

  const stakingContractAddress = stakingContract.options.address

  const extraData =
    "0x" +
    Buffer.concat([
      Buffer.from(beneficiaryAddress.substr(2), "hex"),
      Buffer.from(operatorAddress.substr(2), "hex"),
      Buffer.from(authorizerAddress.substr(2), "hex"),
    ]).toString("hex")

  try {
    yield call(sendTransaction, {
      type: "web3/send_transaction",
      payload: {
        contract: token,
        methodName: "approveAndCall",
        args: [stakingContractAddress, amount, extraData],
      },
    })
    yield call(resetForm)
  } catch (error) {}
  yield call(setSubmitting, false)
}

export function* watchDelegateStakeRequest() {
  yield takeLatest("DELEGATE_STAKE_REQUEST", delegateStake)
}
