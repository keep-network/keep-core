const initialState = {
  displayModal: false,
  payload: null,
}

const modalWindowReducer = (state = initialState, action) => {
  // if (!action.payload) {
  //   return state
  // }

  switch (action.type) {
    case "modal_window/display_modal":
      return {
        displayModal: true,
        payload: action.payload,
      }
    case "modal_window/hide_modal":
      return initialState
    default:
      return state
  }
}

export default modalWindowReducer
