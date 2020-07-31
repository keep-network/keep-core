import { takeEvery } from "redux-saga/effects"

function* showMessage(action) {}

export function* showMessageSaga() {
  yield takeEvery("SHOW_MESSAGE", showMessage)
}
