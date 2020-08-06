import { ADD_MESSAGE, REMOVE_MESSAGE } from "../actions/messages"

const initialState = []

const messages = (state = initialState, action) => {
  switch (action.type) {
    case ADD_MESSAGE:
      return [action.payload, ...state]
    case REMOVE_MESSAGE:
      return state.filter((message) => message.id !== action.payload)
    default:
      return state
  }
}

export default messages
