import { ZERO_ADDRESS } from "../utils/ethereum.utils"
import {
  FETCH_OPERATOR_DELEGATIONS_START,
  FETCH_OPERATOR_DELEGATIONS_SUCCESS,
  FETCH_OPERATOR_DELEGATIONS_FAILURE,
  FETCH_OPERATOR_SLASHED_TOKENS_START,
  FETCH_OPERATOR_SLASHED_TOKENS_SUCCESS,
  FETCH_OPERATOR_SLASHED_TOKENS_FAILURE,
} from "../actions"

const initialState = {
  stakedBalance: "0",
  ownerAddress: ZERO_ADDRESS,
  beneficiaryAddress: ZERO_ADDRESS,
  authorizerAddress: ZERO_ADDRESS,
  isFetching: false,
  error: null,

  areSlashedTokensFetching: false,
  slashedTokens: [],
  slashedTokensError: null,
}

const operatorReducer = (state = initialState, action) => {
  switch (action.type) {
    case FETCH_OPERATOR_DELEGATIONS_START:
      return { ...state, isFetching: true }
    case FETCH_OPERATOR_DELEGATIONS_SUCCESS:
      return { ...state, isFetching: false, ...action.payload }
    case FETCH_OPERATOR_DELEGATIONS_FAILURE:
      return { ...state, error: action.payload.error, isFetching: false }
    case FETCH_OPERATOR_SLASHED_TOKENS_START:
      return { ...state, areSlashedTokensFetching: true }
    case FETCH_OPERATOR_SLASHED_TOKENS_SUCCESS:
      return {
        ...state,
        slashedTokens: action.payload,
        areSlashedTokensFetching: false,
      }
    case FETCH_OPERATOR_SLASHED_TOKENS_FAILURE:
      return {
        ...state,
        isFetching: false,
        slashedTokensError: action.payload.error,
      }
    default:
      return state
  }
}

export default operatorReducer
