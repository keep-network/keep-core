import { CLOSE_MODAL, OPEN_MODAL } from "../actions/modal"

const initialState = {
  modalProps: {},
  modalType: null,
}

const modalReducer = (state = initialState, action) => {
  switch (action.type) {
    case OPEN_MODAL:
      return {
        modalType: action.payload.modalType,
        modalProps: {
          ...action.payload.modalProps,
        },
      }
    case CLOSE_MODAL:
      return {
        ...state,
        modalType: null,
        modalProps: {},
      }
    default:
      return state
  }
}
export default modalReducer
