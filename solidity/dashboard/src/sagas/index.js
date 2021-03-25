import { all, fork, take, cancel, put, delay } from "redux-saga/effects"
import * as messagesSaga from "./messages"
import * as delegateStakeSaga from "./staking"
import * as tokenGrantSaga from "./token-grant"
import { watchSendTransactionRequest } from "./web3"
import * as copyStakeSaga from "./copy-stake"
import * as subscriptions from "./subscriptions"
import * as keepTokenBalance from "./keep-balance"
import * as rewards from "./rewards"
import * as liquidityRewards from "./liquidity-rewards"

export default function* rootSaga() {
  yield take("app/set_account")
  while (true) {
    const tasks = yield all(
      [
        ...Object.values(messagesSaga),
        ...Object.values(delegateStakeSaga),
        watchSendTransactionRequest,
        ...Object.values(tokenGrantSaga),
        ...Object.values(copyStakeSaga),
        ...Object.values(subscriptions),
        ...Object.values(keepTokenBalance),
        ...Object.values(rewards),
        ...Object.values(liquidityRewards),
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
