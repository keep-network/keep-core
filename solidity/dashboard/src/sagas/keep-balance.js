import { takeLatest, call, put } from "redux-saga/effects"
import { logError } from "./utils"
import { getWeb3Context } from "./utils"
import keepToken from "../services/keepToken"

export function* watchKeepTokenBalanceRequest() {
  yield takeLatest("keep-token/balance_request", fetchKeepTokenBalance)
}

function* fetchKeepTokenBalance() {
  const {
    eth: { defaultAccount },
  } = yield getWeb3Context()
  try {
    const keepTokenBalance = yield call(
      [keepToken, keepToken.balanceOf],
      defaultAccount
    )
    yield put({
      type: "keep-token/balance_request_success",
      payload: keepTokenBalance,
    })
  } catch (error) {
    yield* logError("keep-token/balance_request_failure", error)
  }
}
