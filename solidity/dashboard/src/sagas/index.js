import { all, fork, take, cancel, put } from "redux-saga/effects"
import * as messagesSaga from "./messages"
import * as delegateStakeSaga from "./staking"
import * as tokenGrantSaga from "./token-grant"
import {
  watchSendRawTransactionsInSequenceRequest,
  watchSendTransactionRequest,
} from "./web3"
import * as copyStakeSaga from "./copy-stake"
import * as subscriptions from "./subscriptions"
import * as keepTokenBalance from "./keep-balance"
import * as rewards from "./rewards"
import * as liquidityRewards from "./liquidity-rewards"
import * as operator from "./operartor"
import * as authrization from "./authorization"

export default function* rootSaga() {
  while (true) {
    const {
      payload: { address },
    } = yield take("app/login")
    yield put({ type: "app/set_account", payload: { address } })
    const task = yield fork(runTasks)
    yield take("app/logout")
    yield cancel(task)
    yield put({ type: "app/reset_store" })
  }
}

export function* runTasks() {
  while (true) {
    const tasks = yield all(
      [
        ...Object.values(messagesSaga),
        ...Object.values(delegateStakeSaga),
        watchSendTransactionRequest,
        watchSendRawTransactionsInSequenceRequest,
        ...Object.values(tokenGrantSaga),
        ...Object.values(copyStakeSaga),
        ...Object.values(subscriptions),
        ...Object.values(keepTokenBalance),
        ...Object.values(rewards),
        ...Object.values(liquidityRewards),
        ...Object.values(operator),
        ...Object.values(authrization),
      ].map(fork)
    )

    const {
      payload: { address },
    } = yield take("app/account_changed")
    yield cancel(tasks)
    yield put({ type: "app/reset_store" })
    yield put({ type: "app/set_account", payload: { address } })
  }
}
