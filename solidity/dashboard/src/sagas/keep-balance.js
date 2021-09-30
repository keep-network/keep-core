import {
  takeLatest,
  call,
  put,
  retry,
  actionChannel,
  take,
} from "redux-saga/effects"
import { logError } from "./utils"
import { getWeb3Context } from "./utils"
import keepToken from "../services/keepToken"
import { keepBalanceActions } from "../actions"

export function* watchKeepTokenBalanceRequest() {
  yield takeLatest("keep-token/balance_request", fetchKeepTokenBalance)
}

export function* fetchKeepTokenBalanceWithRetry() {
  yield retry(3, 5000, fetchKeepTokenBalance)
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

export function* subscribeToKeepTokenTransferFromEvent() {
  const requestChan = yield actionChannel(
    keepBalanceActions.KEEP_TOKEN_TRANSFER_FROM_EVENT_EMITTED
  )

  // Observe and dispatch an action that updates keep token balance.
  while (true) {
    const {
      payload: { event },
    } = yield take(requestChan)
    const {
      returnValues: { value },
    } = event

    yield put({
      type: "keep-token/transferred_from",
      payload: { value },
    })
  }
}

export function* subscribeToKeepTokenTransferToEvent() {
  const requestChan = yield actionChannel(
    keepBalanceActions.KEEP_TOKEN_TRANSFER_TO_EVENT_EMITTED
  )

  // Observe and dispatch an action that updates keep token balance.
  while (true) {
    const {
      payload: { event },
    } = yield take(requestChan)
    const {
      returnValues: { value },
    } = event

    yield put({
      type: "keep-token/transferred_to",
      payload: { value },
    })
  }
}
