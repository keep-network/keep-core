import { take, call, put, fork } from "redux-saga/effects"
import { logError } from "./utils"
import rewardsService from "../services/rewards.service"

function* fetchBeaconDistributedRewards(address) {
  try {
    yield put({ type: "rewards/beacon_fetch_distributed_rewards_start" })
    const balance = yield call(
      rewardsService.fetchtDistributedBeaconRewards,
      address
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
