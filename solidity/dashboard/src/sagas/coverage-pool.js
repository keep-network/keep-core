import { put, call, takeLatest } from "redux-saga/effects"
import { takeOnlyOnce } from "./effects"
import {
  COVERAGE_POOL_FETCH_TVL_REQUEST,
  COVERAGE_POOL_FETCH_TVL_ERROR,
  fetchTvlStart,
  fetchTvlSuccess,
  fetchShareOfPoolSuccess,
  COVERAGE_POOL_FETCH_SHARE_OF_POOL_REQUEST,
} from "../actions/coverage-pool"
import { identifyTaskByAddress, logErrorAndThrow, logError } from "./utils"
import { Keep } from "../contracts"

function* fetchTvl() {
  try {
    const assetPoolAddress = Keep.coveragePoolV1.assetPoolContract.address

    yield put(fetchTvlStart())
    const tvl = yield call(
      Keep.coveragePoolV1.corateralTokenContract.makeCall,
      "balanceOf",
      assetPoolAddress
    )
    yield put(fetchTvlSuccess(tvl))
  } catch (error) {
    yield* logError(COVERAGE_POOL_FETCH_TVL_ERROR, error)
  }
}

export function* watchFetchTvl() {
  yield takeLatest(COVERAGE_POOL_FETCH_TVL_REQUEST, fetchTvl)
}

function* fetchShareOfPool(action) {
  const { address } = action.payload
  try {
    const balanceOf = yield call(Keep.coveragePoolV1.covBalanceOf, address)
    const totalSupply = yield call(Keep.coveragePoolV1.covTotalSupply)
    const shareOfPool = Keep.coveragePoolV1.shareOfPool(totalSupply, balanceOf)

    yield put(
      fetchShareOfPoolSuccess({
        shareOfPool,
        covBalance: balanceOf,
        covTotalSupply: totalSupply,
      })
    )
  } catch (error) {
    yield* logErrorAndThrow(COVERAGE_POOL_FETCH_TVL_ERROR, error)
  }
}

export function* watchFetchShareOfPool() {
  yield takeOnlyOnce(
    COVERAGE_POOL_FETCH_SHARE_OF_POOL_REQUEST,
    identifyTaskByAddress,
    fetchShareOfPool
  )
}
