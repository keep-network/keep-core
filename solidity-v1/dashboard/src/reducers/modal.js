import { CLOSE_MODAL, OPEN_MODAL } from "../actions/modal"

const initialState = {
  isOpen: false,
  modalProps: {},
  modalType: null,
}

const modalReducer = (state = initialState, action) => {
  switch (action.type) {
    case OPEN_MODAL:
      return {
        isOpen: true,
        modalType: action.payload.modalType,
        modalProps: {
          ...action.payload.modalProps,
        },
      }
    case CLOSE_MODAL:
      return {
        ...state,
        isOpen: false,
        modalType: null,
        modalProps: {},
      }
    default:
      return state
  }
}
export default modalReducer
