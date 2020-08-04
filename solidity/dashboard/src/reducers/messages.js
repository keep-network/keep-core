import { ADD_MESSAGE, REMOVE_MESSAGE } from "../actions/messages"

const initialState = {
  messages: [],
}

const messages = (state = initialState, action) => {
  switch (action.type) {
    case ADD_MESSAGE:
      return {
        ...state,
        messages: [action.payload, ...state.messages],
      }
    case REMOVE_MESSAGE:
      return {
        ...state,
        messages: state.messages.filter(
          (message) => message.id !== action.payload
        ),
      }
    default:
      return state
  }
}

export default messages
