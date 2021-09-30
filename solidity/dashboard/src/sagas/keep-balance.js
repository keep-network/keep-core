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
  yield takeLatest(
    keepBalanceActions.KEEP_TOKEN_BALANCE_REQUEST,
    fetchKeepTokenBalance
  )
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
      type: keepBalanceActions.KEEP_TOKEN_BALANCE_REQUEST_SUCCESS,
      payload: keepTokenBalance,
    })
  } catch (error) {
    yield* logError(
      keepBalanceActions.KEEP_TOKEN_BALANCE_REQUEST_FAILURE,
      error
    )
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
      type: keepBalanceActions.KEEP_TOKEN_TRANSFERRED_FROM,
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
      type: keepBalanceActions.KEEP_TOKEN_TRANSFERRED_TO,
      payload: { value },
    })
  }
}
