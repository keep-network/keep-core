import { call, put } from "redux-saga/effects"
import { takeOnlyOnce } from "./effects"
import { operatorService } from "../services/token-staking.service"
import { logErrorAndThrow, identifyTaskByAddress } from "./utils"
import {
  FETCH_OPERATOR_DELEGATIONS_RERQUEST,
  FETCH_OPERATOR_DELEGATIONS_START,
  FETCH_OPERATOR_DELEGATIONS_SUCCESS,
  FETCH_OPERATOR_DELEGATIONS_FAILURE,
  FETCH_OPERATOR_SLASHED_TOKENS_START,
  FETCH_OPERATOR_SLASHED_TOKENS_SUCCESS,
  FETCH_OPERATOR_SLASHED_TOKENS_FAILURE,
  FETCH_OPERATOR_SLASHED_TOKENS_RERQUEST,
} from "../actions"
import { slashedTokensService } from "../services/slashed-tokens.service"

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
    identifyTaskByAddress,
    fetchOperatorDelegations
  )
}

function* fetchOperatorSlashedTokens(action) {
  try {
    const {
      payload: { address },
    } = action
    yield put({ type: FETCH_OPERATOR_SLASHED_TOKENS_START })
    const data = yield call(slashedTokensService.fetchSlashedTokens, address)
    yield put({ type: FETCH_OPERATOR_SLASHED_TOKENS_SUCCESS, payload: data })
  } catch (error) {
    yield* logErrorAndThrow(FETCH_OPERATOR_SLASHED_TOKENS_FAILURE, error)
  }
}

export function* watchFetchOperatorSlashedTokensRequest() {
  yield takeOnlyOnce(
    FETCH_OPERATOR_SLASHED_TOKENS_RERQUEST,
    identifyTaskByAddress,
    fetchOperatorSlashedTokens
  )
}
