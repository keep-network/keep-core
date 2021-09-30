import { add, sub } from "../utils/arithmetics.utils"
import { keepBalanceActions } from "../actions"

export const keepBalanceInitialState = {
  value: "0",
  isFetching: false,
  error: "",
}

const keepBalance = (state = keepBalanceInitialState, action) => {
  switch (action.type) {
    case keepBalanceActions.KEEP_TOKEN_TRANSFERRED_FROM:
      return {
        ...state,
        value: sub(state.value, action.payload.value),
      }
    case keepBalanceActions.KEEP_TOKEN_TRANSFERRED_TO:
      return {
        ...state,
        value: add(state.value, action.payload.value),
      }
    case keepBalanceActions.KEEP_TOKEN_BALANCE_REQUEST:
      return { ...state, isFetching: true }
    case keepBalanceActions.KEEP_TOKEN_BALANCE_REQUEST_SUCCESS:
      return { ...state, isFetching: false, value: action.payload }
    case keepBalanceActions.KEEP_TOKEN_BALANCE_REQUEST_FAILURE:
      return {
        ...state,
        isFetching: false,
        error: action.payload.error,
      }
    default:
      return state
  }
}

export default keepBalance
