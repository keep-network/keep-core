import { all, fork } from "redux-saga/effects"
import * as messagesSaga from "./messages"
import * as delegateStakeSaga from "./staking"
import * as tokenGrantSaga from "./token-grant"
import { watchSendTransactionRequest } from "./web3"

export default function* rootSaga() {
  yield all(
    [
      ...Object.values(messagesSaga),
      ...Object.values(delegateStakeSaga),
      watchSendTransactionRequest,
      ...Object.values(tokenGrantSaga),
    ].map(fork)
  )
}
