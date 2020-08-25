import { takeLatest, call, put } from "redux-saga/effects"
import {
  FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_REQUEST,
  FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_SUCCESS,
  FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_FAILURE,
} from "../actions"
import { fetchOldDelegations } from "../services/old-staking.service"
import { getContractsContext } from "./utils"
import { sendTransaction } from "./web3"

function* fetchOldStakingDelegations() {
  try {
    const { undelegationPeriod, intializationPeriod, delegations } = yield call(
      fetchOldDelegations
    )
    yield put({
      type: FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_SUCCESS,
      payload: { delegations, undelegationPeriod, intializationPeriod },
    })
  } catch (error) {
    yield put({
      type: FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_FAILURE,
      payload: error,
    })
  }
}

export function* watchFetchOldStakingContract() {
  yield takeLatest(
    FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_REQUEST,
    fetchOldStakingDelegations
  )
}

function* copyStake(action) {
  const operator = action.payload
  const { stakingPortBackerContract } = yield getContractsContext()

  try {
    yield call(sendTransaction, {
      payload: {
        contract: stakingPortBackerContract,
        methodName: "copyStake",
        args: [operator],
      },
    })
  } catch (error) {
    yield put({ type: "copy-stake/copy-stake_failure", payload: error })
  }
}

export function* watchCopyStakeRequest() {
  yield takeLatest("copy-stake/copy-stake-request", copyStake)
}

function* undelegateFromOldContract(action) {
  const delegation = action.payload
  const { operatorAddress } = delegation

  // TODO undelegate from grant
  const { oldTokenStakingContract } = yield getContractsContext()

  try {
    yield call(sendTransaction, {
      payload: {
        contract: oldTokenStakingContract,
        methodName: "undelegate",
        args: [operatorAddress],
      },
    })
  } catch (error) {
    yield put({ type: "copy-stake/undelegation_failure", payload: error })
  }
}
