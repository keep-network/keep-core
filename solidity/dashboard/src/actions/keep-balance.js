export const KEEP_TOKEN_TRANSFER_FROM_EVENT_EMITTED = "keep_token/transfer_from"
export const KEEP_TOKEN_TRANSFER_TO_EVENT_EMITTED = "keep_token/transfer_to"

export const keepTokenTransferFromEventEmitted = (event) => {
  return {
    type: KEEP_TOKEN_TRANSFER_FROM_EVENT_EMITTED,
    payload: { event },
  }
}

export const keepTokenTransferToEventEmitted = (event) => {
  return {
    type: KEEP_TOKEN_TRANSFER_TO_EVENT_EMITTED,
    payload: { event },
  }
}
