import { call, put } from "redux-saga/effects"
import { takeOnlyOnce } from "./effects"
import { operatorService } from "../services/token-staking.service"
import { logErrorAndThrow } from "./utils"
import {
  FETCH_OPERATOR_DELEGATIONS_RERQUEST,
  FETCH_OPERATOR_DELEGATIONS_START,
  FETCH_OPERATOR_DELEGATIONS_SUCCESS,
  FETCH_OPERATOR_DELEGATIONS_FAILURE,
} from "../actions"

function* fetchOperatorDelegations(action) {
  try {
    const {
      payload: { address },
    } = action
    yield put({ type: FETCH_OPERATOR_DELEGATIONS_START })
    const data = yield call(operatorService.fetchDelegatedTokensData, address)
    yield put({ type: FETCH_OPERATOR_DELEGATIONS_SUCCESS, payload: data })
  } catch (error) {
    yield* logErrorAndThrow(FETCH_OPERATOR_DELEGATIONS_FAILURE, error)
  }
}

export function* watchFetchOperatorDelegationRequest() {
  yield takeOnlyOnce(
    FETCH_OPERATOR_DELEGATIONS_RERQUEST,
    (action) => action.payload.address,
    fetchOperatorDelegations
  )
}
