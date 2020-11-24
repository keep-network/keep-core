import { all, fork } from "redux-saga/effects"
import * as messagesSaga from "./messages"
import * as delegateStakeSaga from "./staking"
import * as tokenGrantSaga from "./token-grant"
import { watchSendTransactionRequest } from "./web3"
import * as copyStakeSaga from "./copy-stake"
import * as subscriptions from "./subscriptions"
import * as keepTokenBalance from "./keep-balance"

export default function* rootSaga() {
  yield all(
    [
      ...Object.values(messagesSaga),
      ...Object.values(delegateStakeSaga),
      watchSendTransactionRequest,
      ...Object.values(tokenGrantSaga),
      ...Object.values(copyStakeSaga),
      ...Object.values(subscriptions),
      ...Object.values(keepTokenBalance),
    ].map(fork)
  )
}
