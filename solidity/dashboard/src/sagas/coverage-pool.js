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
  COVERAGE_POOL_FETCH_APY_REQUEST,
  COVERAGE_POOL_WITHDRAW_ASSET_POOL,
  COVERAGE_POOL_CLAIM_TOKENS_FROM_WITHDRAWAL,
  COVERAGE_POOL_WITHDRAWAL_COMPLETED_EVENT_EMITTED,
  COVERAGE_POOL_WITHDRAWAL_INITIATED_EVENT_EMITTED,
} from "../actions/coverage-pool"
import {
  identifyTaskByAddress,
  logErrorAndThrow,
  logError,
  submitButtonHelper,
} from "./utils"
import { Keep } from "../contracts"
import { add, eq, gt, sub } from "../utils/arithmetics.utils"
import { isSameEthAddress } from "../utils/general.utils"
import { ZERO_ADDRESS } from "../utils/ethereum.utils"
import { sendTransaction } from "./web3"
import { KEEP } from "../utils/token.utils"
import selectors from "./selectors"
import { showModal } from "../actions/modal"
import { modalComponentType } from "../components/Modal"

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
      selectors.getCoveragePool
    )

    const address = yield select(selectors.getUserAddress)
    let updatedCovTotalSupply = covTotalSupply
    if (isSameEthAddress(from, ZERO_ADDRESS)) {
      updatedCovTotalSupply = add(covTotalSupply, value).toString()
    } else if (isSameEthAddress(to, ZERO_ADDRESS)) {
      updatedCovTotalSupply = sub(covTotalSupply, value).toString()
    }

    let arithmeticOpration = null
    if (isSameEthAddress(address, from)) {
      arithmeticOpration = sub
    } else if (isSameEthAddress(address, to)) {
      arithmeticOpration = add
    }

    const updatedCovBalance = arithmeticOpration
      ? arithmeticOpration(covBalance, value).toString()
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

    if (isSameEthAddress(to, address)) {
      yield put(
        showModal({
          modalComponentType:
            modalComponentType.COV_POOLS.KEEP_DEPOSITED_SUCCESS,
          componentProps: {
            transactionFinished: true,
            transactionHash: event.transactionHash,
            amount: value,
            balanceAmount: updatedCovBalance,
            estimatedBalanceAmountInKeep: estimatedKeepBalance,
          },
          modalProps: {
            title: "Claim tokens",
            classes: {
              modalWrapperClassName: "modal-wrapper__claim-tokens",
            },
          },
        })
      )
    }

    const estimatedRewards = yield call(
      Keep.coveragePoolV1.estimatedRewards,
      address,
      shareOfPool
    )

    const tvl = yield call(Keep.coveragePoolV1.totalValueLocked)
    const keepInUSD = yield call(Keep.exchangeService.getKeepTokenPriceInUSD)
    const tvlInUSD = keepInUSD.multipliedBy(KEEP.toTokenUnit(tvl)).toFormat(2)
    const apy = yield call(Keep.coveragePoolV1.apy)

    const pendingWithdrawal = yield call(
      Keep.coveragePoolV1.pendingWithdrawal,
      address
    )

    const withdrawalInitiatedTimestamp = yield call(
      Keep.coveragePoolV1.withdrawalInitiatedTimestamp,
      address
    )

    const covTokensAvailableToWithdraw = sub(
      updatedCovBalance,
      pendingWithdrawal
    )

    yield put(
      covTokenUpdated({
        covBalance: updatedCovBalance,
        covTokensAvailableToWithdraw,
        covTotalSupply: updatedCovTotalSupply,
        shareOfPool,
        estimatedKeepBalance,
        estimatedRewards,
        totalValueLocked: tvl,
        totalValueLockedInUSD: tvlInUSD,
        apy,
        pendingWithdrawal,
        withdrawalInitiatedTimestamp,
      })
    )
  }
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
      returnValues: { underwriter, covAmount },
    } = event

    const address = yield select(selectors.getUserAddress)
    const { componentProps } = yield select(selectors.getModalData)

    if (isSameEthAddress(address, underwriter)) {
      let modalType = modalComponentType.COV_POOLS.INITIATE_WITHDRAWAL
      let title = "Withdraw"
      let amount = covAmount
      if (
        componentProps?.pendingWithdrawalBalance &&
        componentProps?.amount &&
        gt(componentProps?.pendingWithdrawalBalance, 0) &&
        eq(componentProps?.amount, 0)
      ) {
        modalType = modalComponentType.COV_POOLS.RE_INITIATE_WITHDRAWAL
        title = "Re-initiate withdrawal"
      } else if (
        componentProps?.pendingWithdrawalBalance &&
        componentProps?.amount &&
        gt(componentProps?.pendingWithdrawalBalance, 0) &&
        gt(componentProps?.amount, 0)
      ) {
        modalType = modalComponentType.COV_POOLS.INCREASE_WITHDRAWAL
        title = "Re-initiate withdrawal"
        amount = componentProps.amount
      }

      yield put(
        showModal({
          modalComponentType: modalType,
          componentProps: {
            transactionFinished: true,
            transactionHash: event.transactionHash,
            pendingWithdrawalBalance: componentProps?.pendingWithdrawalBalance,
            amount: amount,
          },
          modalProps: {
            title,
            classes: {
              modalWrapperClassName: "modal-wrapper__initiate-withdrawal",
            },
          },
        })
      )
    }
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
      returnValues: { underwriter, amount },
    } = event
    const address = yield select(selectors.getUserAddress)

    if (isSameEthAddress(address, underwriter)) {
      yield put(
        showModal({
          modalComponentType: modalComponentType.COV_POOLS.WITHDRAWAL_COMPLETED,
          componentProps: {
            transactionFinished: true,
            amount: amount,
          },
          modalProps: {
            title: "Claim tokens",
            classes: {
              modalWrapperClassName: "modal-wrapper__claim-tokens",
            },
          },
        })
      )
    }
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
