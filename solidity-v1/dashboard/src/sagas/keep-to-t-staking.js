import { getThresholdTokenStakingAddress, Keep } from "../contracts"
import { actionChannel, call, put, take, takeEvery } from "redux-saga/effects"
import { sendTransaction } from "./web3"
import {
  getContractsContext,
  identifyTaskByAddress,
  logErrorAndThrow,
  submitButtonHelper,
} from "./utils"
import {
  STAKE_KEEP_TO_T,
  stakedToT,
  THRESHOLD_STAKE_KEEP_EVENT_EMITTED,
  thresholdContractAuthorized,
} from "../actions/keep-to-t-staking"
import { hideModal, showModal } from "../actions/modal"
import { MODAL_TYPES } from "../constants/constants"
import {
  FETCH_THRESHOLD_AUTH_DATA_FAILURE,
  FETCH_THRESHOLD_AUTH_DATA_REQUEST,
  FETCH_THRESHOLD_AUTH_DATA_START,
  FETCH_THRESHOLD_AUTH_DATA_SUCCESS,
} from "../actions"
import { thresholdAuthorizationService } from "../services/threshold-authorization.service"
import { takeOnlyOnce } from "./effects"
import { fromThresholdTokenAmount } from "../utils/stake-to-t.utils"

export function* subscribeToStakeKeepEvent() {
  const requestChan = yield actionChannel(THRESHOLD_STAKE_KEEP_EVENT_EMITTED)

  while (true) {
    const {
      payload: { event },
    } = yield take(requestChan)

    const {
      returnValues: {
        authorizer,
        beneficiary,
        stakingProvider: operator,
        amount: tAmount,
      },
    } = event

    console.log("tAmount", tAmount)

    yield put(
      showModal({
        modalType: MODAL_TYPES.StakeOnThresholdConfirmed,
        modalProps: {
          transactionHash: event.transactionHash,
          authorizer,
          beneficiary,
          operator,
          keepAmount: fromThresholdTokenAmount(tAmount),
        },
      })
    )
  }
}

function* authorizeAndStakeKeepToT(action) {
  const { payload } = action
  const { operator, isAuthorized } = payload
  const { stakingContract } = yield getContractsContext()
  const operatorContractAddress = getThresholdTokenStakingAddress()

  if (!isAuthorized) {
    yield put(
      showModal({
        modalType: MODAL_TYPES.ThresholdAuthorizationLoadingModal,
        modalProps: {
          text: "Please, authorize in your wallet",
        },
      })
    )

    try {
      yield call(sendTransaction, {
        payload: {
          contract: stakingContract,
          methodName: "authorizeOperatorContract",
          args: [operator, operatorContractAddress],
        },
      })
    } catch (err) {
      yield put(hideModal())
      return
    }

    yield put(thresholdContractAuthorized(operator))
  }

  yield put(
    showModal({
      modalType: MODAL_TYPES.ThresholdStakeConfirmationLoadingModal,
      modalProps: {
        text: "Please, confirm in your wallet",
      },
    })
  )

  try {
    yield call(sendTransaction, {
      payload: {
        contract: Keep.keepToTStaking.thresholdStakingContract.instance,
        methodName: "stakeKeep",
        args: [operator],
      },
    })
  } catch (err) {
    yield put(hideModal())
    return
  }

  yield put(stakedToT(operator))
}

function* authorizeAndStakeKeepToTWorker(action) {
  yield call(submitButtonHelper, authorizeAndStakeKeepToT, action)
}

export function* watchAuthorizeAndStakeKeepToT() {
  yield takeEvery(STAKE_KEEP_TO_T, authorizeAndStakeKeepToTWorker)
}

function* fetchThresholdAuthData(action) {
  try {
    const {
      payload: { address },
    } = action
    yield put({ type: FETCH_THRESHOLD_AUTH_DATA_START })
    const data = yield call(
      thresholdAuthorizationService.fetchThresholdAuthorizationData,
      address
    )
    yield put({
      type: FETCH_THRESHOLD_AUTH_DATA_SUCCESS,
      payload: data,
    })
  } catch (error) {
    yield* logErrorAndThrow(FETCH_THRESHOLD_AUTH_DATA_FAILURE, error)
  }
}

export function* watchFetchThresholdAuthData() {
  yield takeOnlyOnce(
    FETCH_THRESHOLD_AUTH_DATA_REQUEST,
    identifyTaskByAddress,
    fetchThresholdAuthData
  )
}
