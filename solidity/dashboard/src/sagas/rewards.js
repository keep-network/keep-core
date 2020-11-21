import { take, takeLatest, call, put, fork } from "redux-saga/effects"
import { logError, submitButtonHelper, getContractsContext } from "./utils"
import {
  fetchtTotalDistributedRewards,
  fetchECDSAAvailableRewards,
} from "../services/rewards"
import { sendTransaction } from "./web3"

function* fetchBeaconDistributedRewards(address) {
  try {
    yield put({ type: "rewards/beacon_fetch_distributed_rewards_start" })
    const balance = yield call(
      fetchtTotalDistributedRewards,
      address,
      "beaconRewardsContract"
    )
    yield put({
      type: "rewards/beacon_fetch_distributed_rewards_success",
      payload: balance,
    })
  } catch (error) {
    yield* logError("rewards/beacon_fetch_distributed_rewards_failure", error)
  }
}

export function* watchFetchBeaconDistributedRewards() {
  const { payload } = yield take(
    "rewards/beacon_fetch_distributed_rewards_request"
  )
  yield fork(fetchBeaconDistributedRewards, payload)
}

function* fetchECDSADistributedRewards(address) {
  try {
    yield put({ type: "rewards/ecdsa_fetch_distributed_rewards_start" })
    const balance = yield call(
      fetchtTotalDistributedRewards,
      address,
      "ECDSARewardsContract"
    )
    yield put({
      type: "rewards/ecdsa_fetch_distributed_rewards_success",
      payload: balance,
    })
  } catch (error) {
    yield* logError("rewards/ecdsa_fetch_distributed_rewards_failure", error)
  }
}

export function* watchFetchECDSADistributedRewards() {
  const { payload } = yield take(
    "rewards/ecdsa_fetch_distributed_rewards_request"
  )
  yield fork(fetchECDSADistributedRewards, payload)
}

function* fetchECDSAAvailabledRewards(address) {
  try {
    yield put({ type: "rewards/ecdsa_fetch_available_rewards_start" })
    const { totalAvailableRewards, toWithdrawn } = yield call(
      fetchECDSAAvailableRewards,
      address
    )
    yield put({
      type: "rewards/ecdsa_fetch_available_rewards_success",
      payload: { totalAvailableRewards, toWithdrawn },
    })
  } catch (error) {
    yield* logError("rewards/ecdsa_fetch_available_rewards_failure", error)
  }
}

export function* watchFetchECDSAAvailableRewards() {
  const { payload } = yield take(
    "rewards/ecdsa_fetch_available_rewards_request"
  )
  yield fork(fetchECDSAAvailabledRewards, payload)
}

function* withdrawECDSARewards(action) {
  const {
    payload: { availableRewards },
  } = action
  const { ECDSARewardsContract } = yield getContractsContext()

  for (const { operator, interval } of availableRewards) {
    try {
      yield call(sendTransaction, {
        payload: {
          contract: ECDSARewardsContract,
          methodName: "withdrawRewards",
          args: [interval, operator],
        },
      })
    } catch (error) {
      continue
    }
  }
}

export function* watchWithdrawECDSARewards() {
  yield takeLatest("rewards/ecdsa_withdraw", function* (action) {
    yield call(submitButtonHelper, withdrawECDSARewards, action)
  })
}
