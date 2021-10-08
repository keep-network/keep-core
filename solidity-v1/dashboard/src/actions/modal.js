export const OPEN_MODAL = "OPEN_MODAL"
export const CLOSE_MODAL = "CLOSE_MODAL"
export const ADD_ADDITIONAL_DATA_TO_MODAL = "ADD_ADDITIONAL_DATA_TO_MODAL"
export const CLEAR_ADDITIONAL_DATA_FROM_MODAL =
  "CLEAR_ADDITIONAL_DATA_FROM_MODAL"

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
