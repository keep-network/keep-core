const initialState = {
  value: "0",
  isFetching: false,
  error: "",
}

const keepBalance = (state = initialState, action) => {
  switch (action.type) {
    case "keep-token/transfered":
      return {
        ...state,
        value: action.payload.arithmeticOpration(
          state.value,
          action.payload.value
        ),
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
