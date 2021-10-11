import { put, call } from "redux-saga/effects"
import { takeOnlyOnce } from "./effects"
import {
  FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_REQUEST,
  FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_START,
  FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_FAILURE,
  FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_SUCCESS,
} from "../actions"
import { identifyTaskByAddress, logErrorAndThrow } from "./utils"
import { beaconAuthorizationService } from "../services/beacon-authorization.service"

function* fetchKeepRandomBeaconAuthData(action) {
  try {
    const {
      payload: { address },
    } = action
    yield put({ type: FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_START })
    const data = yield call(
      beaconAuthorizationService.fetchRandomBeaconAuthorizationData,
      address
    )
    yield put({
      type: FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_SUCCESS,
      payload: data,
    })
  } catch (error) {
    yield* logErrorAndThrow(FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_FAILURE, error)
  }
}

export function* watchFetchKeepRandomBeaconAuthData() {
  yield takeOnlyOnce(
    FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_REQUEST,
    identifyTaskByAddress,
    fetchKeepRandomBeaconAuthData
  )
}
