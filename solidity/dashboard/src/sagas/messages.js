import { takeEvery, put } from "redux-saga/effects"
import {
  SHOW_MESSAGE,
  CLOSE_MEESSAGE,
  REMOVE_MESSAGE,
  ADD_MESSAGE,
} from "../actions/messages"
import { messageType } from "../components/Message"

function* addMessage(action) {
  yield put({ type: ADD_MESSAGE, payload: action.payload })
}

function* removeMessage(action) {
  console.log("action", action)
  yield put({ type: REMOVE_MESSAGE, payload: action.payload.id })
  if (action.messageType) {
    // switch for future notification data clearance for other message types
    switch (action.payload.messageType) {
      case messageType.LIQUIDITY_REWARDS_EARNED:
        yield put({
          type:
            "notifications_data/liquidityRewardNotification/pairs_displayed_updated",
          payload: [],
        })
        break
      default:
        break
    }
  }
  yield put({
    type:
      "notifications_data/liquidityRewardNotification/pairs_displayed_updated",
    payload: [],
  })
}

export function* showMessageSaga() {
  yield takeEvery(SHOW_MESSAGE, addMessage)
}

export function* removeMessageSaga() {
  yield takeEvery(CLOSE_MEESSAGE, removeMessage)
}
