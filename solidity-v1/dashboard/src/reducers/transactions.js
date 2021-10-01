const initialState = {
  transactionQueue: [],
}

const transactionsReducer = (state = initialState, action) => {
  switch (action.type) {
    case "transactions/transaction_added_to_queue":
      return {
        ...state,
        transactionQueue: [...state.transactionQueue, action.payload],
      }
    case "transactions/clear_queue":
      return initialState
    default:
      return state
  }
}

export default transactionsReducer
