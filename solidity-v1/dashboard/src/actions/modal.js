export const OPEN_MODAL = "modal/open"
export const CLOSE_MODAL = "modal/close"
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
