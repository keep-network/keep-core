import { put, call } from "redux-saga/effects"
import { takeOnlyOnce } from "./effects"
import {
  COVERAGE_POOL_FETCH_TVL_REQUEST,
  COVERAGE_POOL_FETCH_TVL_ERROR,
  fetchTvlStart,
  fetchTvlSuccess,
} from "../actions/coverage-pool"
import { identifyTaskByAddress, logErrorAndThrow } from "./utils"

function* fetchTvl() {
  try {
    yield put(fetchTvlStart())
    const tvl = yield call()
    yield put(fetchTvlSuccess(tvl))
  } catch (error) {
    yield* logErrorAndThrow(COVERAGE_POOL_FETCH_TVL_ERROR, error)
  }
}

export function* watchFetchTvl() {
  yield takeOnlyOnce(
    COVERAGE_POOL_FETCH_TVL_REQUEST,
    identifyTaskByAddress,
    fetchTvl
  )
}
