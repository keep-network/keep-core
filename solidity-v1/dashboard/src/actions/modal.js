export const OPEN_MODAL = "modal/open"
export const CLOSE_MODAL = "modal/close"
export const ADD_ADDITIONAL_DATA_TO_MODAL = "modal/add_additional_data"
export const CLEAR_ADDITIONAL_DATA_FROM_MODAL = "modal/clear_additional_data"
export const CONFIRM = "modal/confirm"
export const CANCEL = "modal/cancel"

export const showModal = (options) => {
  return {
    type: OPEN_MODAL,
    payload: options,
  }
}

export const hideModal = (options) => {
  return {
    type: CLOSE_MODAL,
    payload: options,
  }
}

export const addAdditionalDataToModal = (additionalData) => {
  return {
    type: ADD_ADDITIONAL_DATA_TO_MODAL,
    payload: additionalData,
  }
}

export const clearAdditionalDataFromModal = () => {
  return {
    type: CLEAR_ADDITIONAL_DATA_FROM_MODAL,
  }
}

export const openConfirmationModal = (modalType, modalProps = {}) => {
  return {
    type: OPEN_MODAL,
    payload: {
      modalProps: {
        ...modalProps,
        isConfirmationModal: true,
      },
      modalType,
    },
  }
}
