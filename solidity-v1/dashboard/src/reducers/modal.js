import {
  ADD_ADDITIONAL_DATA_TO_MODAL,
  CLEAR_ADDITIONAL_DATA_FROM_MODAL,
  CLOSE_MODAL,
  OPEN_MODAL,
} from "../actions/modal"

const initialState = {
  isOpen: false,
  modalComponentType: null,
  componentProps: null,
  modalProps: null,
  modalType: null,
}

// TODO: Clean up modal reducer.
const modalReducer = (state = initialState, action) => {
  switch (action.type) {
    case OPEN_MODAL:
      return {
        isOpen: true,
        modalComponentType: action.payload.modalComponentType,
        modalType: action.payload.modalType,
        componentProps: {
          ...state.componentProps,
          ...action.payload.componentProps,
        },
        modalProps: {
          ...state.modalProps,
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
    case ADD_ADDITIONAL_DATA_TO_MODAL:
      return {
        componentProps: {
          ...state.componentProps,
          ...action.payload.componentProps,
        },
        modalProps: {
          ...state.modalProps,
          ...action.payload.modalProps,
        },
      }
    case CLEAR_ADDITIONAL_DATA_FROM_MODAL:
      return {
        componentProps: initialState.componentProps,
        modalProps: initialState.modalProps,
      }
    default:
      return state
  }
}
export default modalReducer
