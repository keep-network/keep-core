import { Keep } from "../contracts"
import { call, takeEvery } from "redux-saga/effects"
import { sendTransaction } from "./web3"
import { submitButtonHelper } from "./utils"
import { STAKE_KEEP_TO_T } from "../actions/keep-to-t-staking"

function* stakeKeepToT(action) {
  const { payload } = action
  const { operator } = payload

  yield call(sendTransaction, {
    payload: {
      contract: Keep.keepToTStaking.thresholdStakingContract.instance,
      methodName: "stakeKeep",
      args: [operator],
    },
  })
}

function* stakeKeepToTWorker(action) {
  yield call(submitButtonHelper, stakeKeepToT, action)
}

export function* watchStakeKeepToT() {
  yield takeEvery(STAKE_KEEP_TO_T, stakeKeepToTWorker)
}
