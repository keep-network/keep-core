import { add, sub } from "../utils/arithmetics.utils"

export const keepBalanceInitialState = {
  value: "0",
  isFetching: false,
  error: "",
}

const keepBalance = (state = keepBalanceInitialState, action) => {
  switch (action.type) {
    case "keep-token/transferred_from":
      return {
        ...state,
        value: sub(state.value, action.payload.value),
      }
    case "keep-token/transferred_to":
      return {
        ...state,
        value: add(state.value, action.payload.value),
      }
    case "keep-token/balance_request":
      return { ...state, isFetching: true }
    case "keep-token/balance_request_success":
      return { ...state, isFetching: false, value: action.payload }
    case "keep-token/balance_request_failure":
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
