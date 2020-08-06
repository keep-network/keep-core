import { all, fork } from "redux-saga/effects"
import * as messagesSaga from "./messages"
import * as delegateStakeSaga from "./delegate-stake"

export default function* rootSaga() {
  yield all(
    [...Object.values(messagesSaga), ...Object.values(delegateStakeSaga)].map(
      fork
    )
  )
}
