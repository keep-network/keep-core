import {
  put,
  call,
  takeLatest,
  select,
  take,
  actionChannel,
  takeEvery,
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
  COVERAGE_POOL_FETCH_COV_POOL_DATA_ERROR,
  COVERAGE_POOL_DEPOSIT_ASSET_POOL,
  fetchAPYStart,
  COVERAGE_POOL_FETCH_APY_ERROR,
  fetchAPYSuccess,
} from "../actions/coverage-pool"
import {
  identifyTaskByAddress,
  logErrorAndThrow,
  logError,
  submitButtonHelper,
} from "./utils"
import { Keep } from "../contracts"
import { add, sub } from "../utils/arithmetics.utils"
import { isSameEthAddress } from "../utils/general.utils"
import { ZERO_ADDRESS } from "../utils/ethereum.utils"
import { sendTransaction } from "./web3"
import { KEEP } from "../utils/token.utils"

function* fetchTvl() {
  try {
    yield put(fetchTvlStart())
    const tvl = yield call(Keep.coveragePoolV1.totalValueLocked)
    const keepInUSD = yield call(Keep.exchangeService.getKeepTokenPriceInUSD)
    const tvlInUSD = keepInUSD.multipliedBy(KEEP.toTokenUnit(tvl)).toFormat(2)
    const totalAllocatedRewards = yield call(
      Keep.coveragePoolV1.totalAllocatedRewards
    )
    yield put(fetchTvlSuccess({ tvl, tvlInUSD, totalAllocatedRewards }))
  } catch (error) {
    yield* logError(COVERAGE_POOL_FETCH_TVL_ERROR, error)
  }
}

export function* watchFetchTvl() {
  yield takeLatest(COVERAGE_POOL_FETCH_TVL_REQUEST, fetchTvl)
}

function* fetchAPY() {
  try {
    yield put(fetchAPYStart())
    const apy = yield call(Keep.coveragePoolV1.apy)
    yield put(fetchAPYSuccess(apy))
  } catch (error) {
    yield* logError(COVERAGE_POOL_FETCH_APY_ERROR, error)
  }
}

export function* watchFetchAPY() {
  yield takeLatest(COVERAGE_POOL_FETCH_TVL_REQUEST, fetchAPY)
}

function* fetchCovPoolData(action) {
  const { address } = action.payload
  try {
    yield put(fetchCovPoolDataStart())

    const balanceOf = yield call(Keep.coveragePoolV1.covBalanceOf, address)
    const totalSupply = yield call(Keep.coveragePoolV1.covTotalSupply)
    const shareOfPool = yield call(
      Keep.coveragePoolV1.shareOfPool,
      totalSupply,
      balanceOf
    )
    const estimatedKeepBalance = yield call(
      Keep.coveragePoolV1.estimatedCollateralTokenBalance,
      shareOfPool
    )

    const estimatedRewards = yield call(
      Keep.coveragePoolV1.estimatedRewards,
      address,
      shareOfPool
    )

    yield put(
      fetchCovPoolDataSuccess({
        shareOfPool,
        covBalance: balanceOf,
        covTotalSupply: totalSupply,
        estimatedRewards,
        estimatedKeepBalance,
      })
    )
  } catch (error) {
    yield* logErrorAndThrow(COVERAGE_POOL_FETCH_COV_POOL_DATA_ERROR, error)
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

    const estimatedKeepBalance = yield call(
      Keep.coveragePoolV1.estimatedCollateralTokenBalance,
      shareOfPool
    )

    const estimatedRewards = yield call(
      Keep.coveragePoolV1.estimatedRewards,
      address,
      shareOfPool
    )

    const tvl = yield call(Keep.coveragePoolV1.totalValueLocked)
    const keepInUSD = yield call(Keep.exchangeService.getKeepTokenPriceInUSD)
    const tvlInUSD = keepInUSD.multipliedBy(KEEP.toTokenUnit(tvl)).toFormat(2)
    const apy = yield call(Keep.coveragePoolV1.apy)

    yield put(
      covTokenUpdated({
        covBalance: updatedCovBalance,
        covTotalSupply: updatedCovTotalSupply,
        shareOfPool,
        estimatedKeepBalance,
        estimatedRewards,
        totalValueLocked: tvl,
        totalValueLockedInUSD: tvlInUSD,
        apy,
      })
    )
  }
}

function* depositAssetPool(action) {
  const { payload } = action
  const { amount } = payload

  const assetPoolAddress = Keep.coveragePoolV1.assetPoolContract.address

  yield call(sendTransaction, {
    payload: {
      contract: Keep.coveragePoolV1.collateralToken.instance,
      methodName: "approveAndCall",
      args: [assetPoolAddress, amount, []],
    },
  })
}

function* depositAssetPoolWorker(action) {
  yield call(submitButtonHelper, depositAssetPool, action)
}

export function* watchDepositAssetPool() {
  yield takeEvery(COVERAGE_POOL_DEPOSIT_ASSET_POOL, depositAssetPoolWorker)
}
