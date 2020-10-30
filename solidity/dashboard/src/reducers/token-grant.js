const initialState = {
  isFetching: false,
  grants: [],
  error: null,
}

const tokenGrantsReducer = (state = initialState, action) => {
  switch (action.type) {
    case "token-grants/fetch_grants_start":
      return {
        ...state,
        isFetching: true,
        error: null,
      }
    case "token-grants/fetch_grants_success":
      return { ...state, isFetching: false, grants: action.payload }
    case "token-grants/fetch_grants_failure":
      return { ...state, isFetching: false, error: action.payload.error }
    default:
      return state
  }
}

export default tokenGrantsReducer
