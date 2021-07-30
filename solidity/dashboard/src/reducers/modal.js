const initialState = {
  isOpen: false,
  emittedEvent: null,
  transactionHash: null,
}

const modalReducer = (state = initialState, action) => {
  switch (action.type) {
    case "modal/is_opened":
      return {
        isOpen: true,
        emittedEvent: action.payload.emittedEvent,
        transactionHash: action.payload.transactionHash,
      }
    case "modal/is_closed":
      return {
        isOpen: false,
        emittedEvent: null,
        transactionHash: null,
      }
    default:
      return state
  }
}

export default modalReducer
