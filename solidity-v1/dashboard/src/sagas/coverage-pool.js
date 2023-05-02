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
  fetchTvlStart,
  fetchTvlSuccess,
  fetchCovPoolDataStart,
  fetchCovPoolDataSuccess,
  covTokenUpdated,
  COVERAGE_POOL_FETCH_COV_POOL_DATA_ERROR,
  fetchAPYStart,
  COVERAGE_POOL_FETCH_APY_ERROR,
  fetchAPYSuccess,
  COVERAGE_POOL_FETCH_APY_REQUEST,
  COVERAGE_POOL_WITHDRAW_ASSET_POOL,
  COVERAGE_POOL_CLAIM_TOKENS_FROM_WITHDRAWAL,
  COVERAGE_POOL_WITHDRAWAL_COMPLETED_EVENT_EMITTED,
  COVERAGE_POOL_WITHDRAWAL_INITIATED_EVENT_EMITTED,
  RISK_MANAGER_AUCTION_CREATED_EVENT_EMITTED,
  RISK_MANAGER_AUCTION_CLOSED_EVENT_EMITTED,
  COVERAGE_POOL_REINITAITE_WITHDRAW,
  COVERAGE_POOL_INCREASE_WITHDRAWAL,
} from "../actions/coverage-pool"
import {
  identifyTaskByAddress,
  logErrorAndThrow,
  logError,
  submitButtonHelper,
  confirmModalSaga,
} from "./utils"
import { Keep } from "../contracts"
import { add, gt, sub } from "../utils/arithmetics.utils"
import { isSameEthAddress } from "../utils/general.utils"
import { sendTransaction } from "./web3"
import { KEEP } from "../utils/token.utils"
import selectors from "./selectors"
import { showModal } from "../actions/modal"
import { MODAL_TYPES } from "../constants/constants"
import { getPendingWithdrawalStatus } from "../utils/coverage-pools.utils"

function* fetchTvl() {
  try {
    yield put(fetchTvlStart())
    const tvl = yield call(Keep.coveragePoolV1.totalValueLocked)
    const keepInUSD = yield call(Keep.exchangeService.getKeepTokenPriceInUSD)
    const tvlInUSD = keepInUSD.multipliedBy(KEEP.toTokenUnit(tvl)).toFormat(2)
    const totalAllocatedRewards = yield call(
      Keep.coveragePoolV1.totalAllocatedRewards
    )
    const totalCoverageClaimed = yield call(
      Keep.coveragePoolV1.totalCoverageClaimed
    )
    yield put(
      fetchTvlSuccess({
        tvl,
        tvlInUSD,
        totalAllocatedRewards,
        totalCoverageClaimed,
      })
    )
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
  yield takeLatest(COVERAGE_POOL_FETCH_APY_REQUEST, fetchAPY)
}

function* fetchCovPoolData(action) {
  const { address } = action.payload
  try {
    yield put(fetchCovPoolDataStart())

    const balanceOf = yield call(Keep.coveragePoolV1.covBalanceOf, address)
    const totalSupply = yield call(Keep.coveragePoolV1.covTotalSupply)

    const withdrawalDelays = yield call(Keep.coveragePoolV1.withdrawalDelays)

    const pendingWithdrawal = yield call(
      Keep.coveragePoolV1.pendingWithdrawal,
      address
    )

    const withdrawalInitiatedTimestamp = yield call(
      Keep.coveragePoolV1.withdrawalInitiatedTimestamp,
      address
    )

    const covBalance = add(balanceOf, pendingWithdrawal).toString()

    const shareOfPool = yield call(
      Keep.coveragePoolV1.shareOfPool,
      totalSupply,
      covBalance
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

    const hasRiskManagerOpenAuctions = yield call(
      Keep.coveragePoolV1.hasRiskManagerOpenAuctions
    )

    yield put(
      fetchCovPoolDataSuccess({
        shareOfPool,
        covBalance,
        covTokensAvailableToWithdraw: balanceOf,
        covTotalSupply: totalSupply,
        estimatedRewards,
        estimatedKeepBalance,
        withdrawalDelay: withdrawalDelays.withdrawalDelay,
        withdrawalTimeout: withdrawalDelays.withdrawalTimeout,
        pendingWithdrawal,
        withdrawalInitiatedTimestamp,
        hasRiskManagerOpenAuctions,
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

export function* subscribeToWithdrawalInitiatedEvent() {
  const requestChan = yield actionChannel(
    COVERAGE_POOL_WITHDRAWAL_INITIATED_EVENT_EMITTED
  )

  while (true) {
    const {
      payload: { event },
    } = yield take(requestChan)
    const {
      returnValues: { underwriter, covAmount, timestamp },
    } = event

    const address = yield select(selectors.getUserAddress)
    const {
      covBalance,
      totalValueLocked,
      covTotalSupply,
      withdrawalDelay,
      withdrawalTimeout,
    } = yield select(selectors.getCoveragePool)

    if (!isSameEthAddress(address, underwriter)) {
      continue
    }
    // TODO: display modal with `WithdrawalOverview` component if a user
    // increased existing withdrawal.
    yield put(
      showModal({
        modalType: MODAL_TYPES.CovPoolWithdrawInitialized,
        modalProps: {
          amount: covAmount,
          transactionHash: event.transactionHash,
          totalValueLocked,
          covTotalSupply,
          covBalanceOf: covBalance,
          timestamp,
          withdrawalDelay,
          withdrawalTimeout,
        },
      })
    )

    yield put(
      covTokenUpdated({
        pendingWithdrawal: covAmount,
        withdrawalInitiatedTimestamp: timestamp,
        covTokensAvailableToWithdraw: sub(covBalance, covAmount).toString(),
      })
    )
  }
}

export function* subscribeToWithdrawalCompletedEvent() {
  const requestChan = yield actionChannel(
    COVERAGE_POOL_WITHDRAWAL_COMPLETED_EVENT_EMITTED
  )

  while (true) {
    const {
      payload: { event },
    } = yield take(requestChan)
    const {
      returnValues: { underwriter, amount, covAmount },
    } = event
    const address = yield select(selectors.getUserAddress)
    const isAddressedToCurrentAddress = isSameEthAddress(address, underwriter)

    const { covTotalSupply, covBalance } = yield select(
      selectors.getCoveragePool
    )

    if (isAddressedToCurrentAddress) {
      yield put(
        showModal({
          modalType: MODAL_TYPES.CovPoolClaimTokens,
          modalProps: {
            transactionHash: event.transactionHash,
            collateralTokenAmount: amount,
            covAmount,
            address: underwriter,
          },
        })
      )
    }

    const updatedCovTotalSupply = sub(covTotalSupply, covAmount).toString()
    const totalValueLocked = yield call(Keep.coveragePoolV1.totalValueLocked)
    const keepInUSD = yield call(Keep.exchangeService.getKeepTokenPriceInUSD)
    const totalValueLockedInUSD = keepInUSD
      .multipliedBy(KEEP.toTokenUnit(totalValueLocked))
      .toFormat(2)
    const apy = yield call(Keep.coveragePoolV1.apy)

    const updatedCovBalance = isAddressedToCurrentAddress
      ? sub(covBalance, covAmount).toString()
      : covBalance

    const shareOfPool = yield call(
      Keep.coveragePoolV1.shareOfPool,
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

    const covTokenUpdatedData = {
      shareOfPool,
      covBalance: updatedCovBalance,
      covTotalSupply: updatedCovTotalSupply,
      estimatedRewards,
      estimatedKeepBalance,
      totalValueLockedInUSD,
      totalValueLocked,
      apy,
    }

    if (isAddressedToCurrentAddress) {
      covTokenUpdatedData.pendingWithdrawal = "0"
      covTokenUpdatedData.withdrawalInitiatedTimestamp = "0"
    }

    yield put(covTokenUpdated(covTokenUpdatedData))
  }
}

export function* subscribeToAuctionCreatedEvent() {
  const requestChan = yield actionChannel(
    RISK_MANAGER_AUCTION_CREATED_EVENT_EMITTED
  )

  while (true) {
    yield take(requestChan)

    const hasRiskManagerOpenAuctions = true

    yield put(
      covTokenUpdated({
        hasRiskManagerOpenAuctions,
      })
    )
  }
}

export function* subscribeToAuctionClosedEvent() {
  const requestChan = yield actionChannel(
    RISK_MANAGER_AUCTION_CLOSED_EVENT_EMITTED
  )

  while (true) {
    yield take(requestChan)

    const hasRiskManagerOpenAuctions = yield call(
      Keep.coveragePoolV1.hasRiskManagerOpenAuctions
    )

    yield put(
      covTokenUpdated({
        hasRiskManagerOpenAuctions,
      })
    )
  }
}

function* withdrawAssetPool(action) {
  const { payload } = action
  const { amount } = payload

  const address = yield select(selectors.getUserAddress)
  const assetPoolAddress = Keep.coveragePoolV1.assetPoolContract.address

  const covTokensAllowed = yield call(
    Keep.coveragePoolV1.covTokensAllowed,
    address,
    assetPoolAddress
  )

  if (gt(amount, covTokensAllowed)) {
    yield call(sendTransaction, {
      payload: {
        contract: Keep.coveragePoolV1.covTokenContract.instance,
        methodName: "approve",
        args: [assetPoolAddress, amount],
      },
    })
  }

  yield call(sendTransaction, {
    payload: {
      contract: Keep.coveragePoolV1.assetPoolContract.instance,
      methodName: "initiateWithdrawal",
      args: [amount],
    },
  })
}

function* withdrawAssetPoolWorker(action) {
  yield call(submitButtonHelper, withdrawAssetPool, action)
}

export function* watchWithdrawAssetPool() {
  yield takeEvery(COVERAGE_POOL_WITHDRAW_ASSET_POOL, withdrawAssetPoolWorker)
}

function* claimTokensFromWithdrawal() {
  const address = yield select(selectors.getUserAddress)

  yield call(sendTransaction, {
    payload: {
      contract: Keep.coveragePoolV1.assetPoolContract.instance,
      methodName: "completeWithdrawal",
      args: [address],
    },
  })
}

function* claimTokensFromWithdrawalWorker(action) {
  yield call(submitButtonHelper, claimTokensFromWithdrawal, action)
}

export function* watchClaimTokensFromWithdrawal() {
  yield takeEvery(
    COVERAGE_POOL_CLAIM_TOKENS_FROM_WITHDRAWAL,
    claimTokensFromWithdrawalWorker
  )
}

function* reinitiateWithdraw(action) {
  const { pendingWithdrawal, totalValueLocked, covTotalSupply, covBalance } =
    yield select(selectors.getCoveragePool)

  yield put(
    showModal({
      modalType: MODAL_TYPES.InitiateCovPoolWithdraw,
      modalProps: {
        amount: pendingWithdrawal,
        covBalanceOf: covBalance,
        estimatedBalanceAmountInKeep: Keep.coveragePoolV1.estimatedBalanceFor(
          covBalance,
          covTotalSupply,
          totalValueLocked
        ),
        totalValueLocked,
        covTotalSupply,
        isReinitialization: true,
        bodyTitle: "You are about to re-withdraw",
      },
    })
  )
}

export function* watchReInitiateWithdraw() {
  yield takeEvery(COVERAGE_POOL_REINITAITE_WITHDRAW, reinitiateWithdraw)
}

function* increaseWithdrawal(amount) {
  const {
    pendingWithdrawal,
    totalValueLocked,
    covTotalSupply,
    covBalance,
    withdrawalDelay,
    withdrawalTimeout,
    withdrawalInitiatedTimestamp,
  } = yield select(selectors.getCoveragePool)
  const address = yield select(selectors.getUserAddress)

  const withdrawalStatus = getPendingWithdrawalStatus(
    withdrawalDelay,
    withdrawalTimeout,
    withdrawalInitiatedTimestamp
  )
  const { isConfirmed } = yield call(
    confirmModalSaga,
    MODAL_TYPES.ConfirmCovPoolIncreaseWithdrawal,
    {
      covAmountToAdd: amount,
      existingWithdrawalCovAmount: pendingWithdrawal,
      withdrawalDelay,
      withdrawalStatus,
      withdrawalInitiatedTimestamp,
    }
  )
  if (!isConfirmed) {
    return
  }

  yield put(
    showModal({
      modalType: MODAL_TYPES.IncreaseCovPoolWithdrawal,
      modalProps: {
        existingWithdrawalCovAmount: pendingWithdrawal,
        covAmountToAdd: amount,
        address,
        covBalanceOf: covBalance,
        totalValueLocked,
        covTotalSupply,
        withdrawalDelay,
        withdrawalStatus,
        withdrawalInitiatedTimestamp,
      },
    })
  )
}

export function* watchIncreaseWithdrawal() {
  yield takeEvery(COVERAGE_POOL_INCREASE_WITHDRAWAL, function* (action) {
    const {
      payload: { amount },
    } = action
    yield* increaseWithdrawal(amount)
  })
}
