import { ZERO_ADDRESS } from "../utils/ethereum.utils"
import {
  FETCH_OPERATOR_DELEGATIONS_START,
  FETCH_OPERATOR_DELEGATIONS_SUCCESS,
  FETCH_OPERATOR_DELEGATIONS_FAILURE,
} from "../actions"

const initialState = {
  stakedBalance: "0",
  ownerAddress: ZERO_ADDRESS,
  beneficiaryAddress: ZERO_ADDRESS,
  authorizerAddress: ZERO_ADDRESS,
  isFetching: false,
  error: null,
}

const operatorReducer = (state = initialState, action) => {
  switch (action.type) {
    case FETCH_OPERATOR_DELEGATIONS_START:
      return { ...state, isFetching: true }
    case FETCH_OPERATOR_DELEGATIONS_SUCCESS:
      return { ...state, isFetching: false, ...action.payload }
    case FETCH_OPERATOR_DELEGATIONS_FAILURE:
      return { ...state, error: action.payload.error, isFetching: false }
    default:
      return state
  }
}

export default operatorReducer
