import { getThresholdTokenStakingAddress, Keep } from "../contracts"
import { actionChannel, call, put, take, takeEvery } from "redux-saga/effects"
import { sendTransaction } from "./web3"
import {
  getContractsContext,
  getWeb3Context,
  identifyTaskByAddress,
  logErrorAndThrow,
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
import { messageType } from "../components/Message"
import { showMessage } from "../actions/messages"

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
          options: {
            onTransactionHashAction: (txHash) =>
              showModal({
                modalType: MODAL_TYPES.ThresholdAuthorizationLoadingModal,
                modalProps: {
                  txHash,
                },
              }),
          },
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
        options: {
          onTransactionHashAction: (txHash) =>
            showModal({
              modalType: MODAL_TYPES.ThresholdStakeConfirmationLoadingModal,
              modalProps: {
                txHash,
              },
            }),
        },
      },
    })
  } catch (err) {
    yield put(
      showModal({
        modalType: MODAL_TYPES.AuthorizedButNotStakedToTWarningModal,
      })
    )
    return
  }

  yield put(stakedToT(operator))
}

export function* watchAuthorizeAndStakeKeepToT() {
  yield takeEvery(STAKE_KEEP_TO_T, authorizeAndStakeKeepToT)
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

export function* watchFetchThresholdAuthDataSuccess() {
  yield takeEvery(
    FETCH_THRESHOLD_AUTH_DATA_SUCCESS,
    processThresholdAuthorizationNotification
  )
}

function* processThresholdAuthorizationNotification(action) {
  const stakesAvailableToStakeOnThreshold = action.payload.filter(
    (stake) => !stake.isStakedToT
  ).length

  const inflectionOfAStakeWord =
    stakesAvailableToStakeOnThreshold === 1 ? "stake" : "stakes"

  if (stakesAvailableToStakeOnThreshold > 0) {
    yield put(
      showMessage({
        messageType: messageType.STAKE_READY_TO_BE_STAKED_TO_T,
        messageProps: {
          sticky: true,
          title: `You have ${stakesAvailableToStakeOnThreshold} ${inflectionOfAStakeWord} to stake on Threshold and earn rewards`,
        },
      })
    )
  }
}

export function* fetchThresholdAuthDataRequest() {
  const {
    eth: { defaultAccount },
  } = yield getWeb3Context()

  yield put({
    type: FETCH_THRESHOLD_AUTH_DATA_REQUEST,
    payload: { address: defaultAccount },
  })
}
