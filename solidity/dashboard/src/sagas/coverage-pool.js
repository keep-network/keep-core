import {
  put,
  call,
  takeLatest,
  select,
  take,
  actionChannel,
} from "redux-saga/effects"
import { takeOnlyOnce } from "./effects"
import {
  COVERAGE_POOL_FETCH_TVL_REQUEST,
  COVERAGE_POOL_FETCH_TVL_ERROR,
  COVERAGE_POOL_FETCH_COV_POOL_DATA_REQUEST,
  COVERAGE_POOL_COV_TOKEN_TRANSFER_EVENT_EMITTED,
  fetchTvlStart,
  fetchTvlSuccess,
  fetchCovPoolDataStart,
  fetchCovPoolDataSuccess,
  covTokenUpdated,
} from "../actions/coverage-pool"
import { identifyTaskByAddress, logErrorAndThrow, logError } from "./utils"
import { Keep } from "../contracts"
import { add, sub } from "../utils/arithmetics.utils"
import { isSameEthAddress } from "../utils/general.utils"
import { ZERO_ADDRESS } from "../utils/ethereum.utils"

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

function* fetchCovPoolData(action) {
  const { address } = action.payload
  try {
    yield put(fetchCovPoolDataStart())

    const balanceOf = yield call(Keep.coveragePoolV1.covBalanceOf, address)
    const totalSupply = yield call(Keep.coveragePoolV1.covTotalSupply)
    const shareOfPool = Keep.coveragePoolV1.shareOfPool(totalSupply, balanceOf)
    const estimatedRewards = yield call(
      Keep.coveragePoolV1.estimatedRewards,
      shareOfPool
    )

    yield put(
      fetchCovPoolDataSuccess({
        shareOfPool,
        covBalance: balanceOf,
        covTotalSupply: totalSupply,
        estimatedRewards,
      })
    )
  } catch (error) {
    yield* logErrorAndThrow(COVERAGE_POOL_FETCH_TVL_ERROR, error)
  }
}

export function* watchFetchCovPoolData() {
  yield takeOnlyOnce(
    COVERAGE_POOL_FETCH_COV_POOL_DATA_REQUEST,
    identifyTaskByAddress,
    fetchCovPoolData
  )
}

export function* subscribeToCovTokenTransferEvent() {
  const requestChan = yield actionChannel(
    COVERAGE_POOL_COV_TOKEN_TRANSFER_EVENT_EMITTED
  )

  while (true) {
    const {
      payload: { event },
    } = yield take(requestChan)
    const {
      returnValues: { from, to, value },
    } = event
    const { covTotalSupply, covBalance } = yield select(
      (state) => state.coveragePool
    )

    const address = yield select((state) => state.app.address)
    let updatedCovTotalSupply = 0
    if (isSameEthAddress(from, ZERO_ADDRESS)) {
      updatedCovTotalSupply = add(covTotalSupply, value)
    } else if (isSameEthAddress(to, ZERO_ADDRESS)) {
      updatedCovTotalSupply = add(covTotalSupply, value)
    }

    let arithmeticOpration = null
    if (isSameEthAddress(address, from)) {
      arithmeticOpration = sub
    } else if (isSameEthAddress(address, to)) {
      arithmeticOpration = add
    }

    const updatedCovBalance = arithmeticOpration
      ? arithmeticOpration(covBalance, value)
      : covBalance

    const shareOfPool = Keep.coveragePoolV1.shareOfPool(
      updatedCovTotalSupply,
      updatedCovBalance
    )

    // TODO update estimated rewards.
    const estimatedRewards = 0

    yield put(
      covTokenUpdated({
        covBalance: updatedCovBalance,
        covTotalSupply: updatedCovTotalSupply,
        shareOfPool,
      })
    )
  }
}
