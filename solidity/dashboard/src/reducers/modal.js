const initialState = {
  isOpen: false,
  emittedEvent: null,
  transactionHash: null,
  additionalData: null,
}

const modalReducer = (state = initialState, action) => {
  switch (action.type) {
    case "modal/is_opened":
      return {
        ...state,
        isOpen: true,
        emittedEvent: action.payload.emittedEvent || state.emittedEvent,
        transactionHash: action.payload.transactionHash,
        additionalData: action.payload.additionalData,
      }
    case "modal/is_closed":
      return {
        ...state,
        isOpen: false,
        emittedEvent: null,
        transactionHash: null,
      }
    case "modal/set_emitted_event":
      return {
        ...state,
        emittedEvent: action.payload,
      }
    default:
      return state
  }
}

export default modalReducer
