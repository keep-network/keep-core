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
import { showModal } from "../actions/modal"
import { MODAL_TYPES } from "../constants/constants"
import {
  FETCH_THRESHOLD_AUTH_DATA_FAILURE,
  FETCH_THRESHOLD_AUTH_DATA_REQUEST,
  FETCH_THRESHOLD_AUTH_DATA_START,
  FETCH_THRESHOLD_AUTH_DATA_SUCCESS,
} from "../actions"
import { thresholdAuthorizationService } from "../services/threshold-authorization.service"
import { takeOnlyOnce } from "./effects"

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
        },
      })
    )

    // const address = yield select(selectors.getUserAddress)
    // const {
    //   covBalance,
    //   totalValueLocked,
    //   covTotalSupply,
    //   withdrawalDelay,
    //   withdrawalTimeout,
    // } = yield select(selectors.getCoveragePool)
    //
    // if (!isSameEthAddress(address, underwriter)) {
    //   continue
    // }
    // // TODO: display modal with `WithdrawalOverview` component if a user
    // // increased existing withdrawal.
    // yield put(
    //   showModal({
    //     modalType: MODAL_TYPES.StakeOnThreshold,
    //     modalProps: {
    //       // amount: covAmount,
    //       transactionHash: event.transactionHash,
    //       authorizer,
    //       beneficiary,
    //       operator,
    //     },
    //   })
    // )
    //
    // yield put(
    //   covTokenUpdated({
    //     pendingWithdrawal: covAmount,
    //     withdrawalInitiatedTimestamp: timestamp,
    //     covTokensAvailableToWithdraw: sub(covBalance, covAmount).toString(),
    //   })
    // )
  }
}

function* authorizeAndStakeKeepToT(action) {
  const { payload } = action
  const { operator, isAuthorized } = payload
  const { stakingContract } = yield getContractsContext()
  const operatorContractAddress = getThresholdTokenStakingAddress()

  if (!isAuthorized) {
    yield call(sendTransaction, {
      payload: {
        contract: stakingContract,
        methodName: "authorizeOperatorContract",
        args: [operator, operatorContractAddress],
      },
    })

    yield put(thresholdContractAuthorized(operator))
  }

  yield call(sendTransaction, {
    payload: {
      contract: Keep.keepToTStaking.thresholdStakingContract.instance,
      methodName: "stakeKeep",
      args: [operator],
    },
  })

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
