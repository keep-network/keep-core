import { all, fork } from "redux-saga/effects"
import * as messagesSaga from "./messages"

export default function* rootSaga() {
  yield all([...Object.values(messagesSaga)].map(fork))
}
