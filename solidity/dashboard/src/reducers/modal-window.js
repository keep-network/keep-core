const initialState = {
  displayModal: false,
  modalComponent: null,
}

const modalWindowReducer = (state = initialState, action) => {
  if (!action.payload) {
    return state
  }

  switch (action.type) {
    case "modal_window/display_modal":
      return {
        displayModal: true,
        modalComponent: action.payload,
      }
    case "modal_window/hide_modal":
      return initialState
    default:
      return state
  }
}

export default modalWindowReducer
