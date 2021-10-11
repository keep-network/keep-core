export const KEEP_TOKEN_TRANSFER_FROM_EVENT_EMITTED =
  "keep_token/transfer_from_event_emitted"
export const KEEP_TOKEN_TRANSFER_TO_EVENT_EMITTED =
  "keep_token/transfer_to_event_emitted"

export const KEEP_TOKEN_TRANSFERRED_FROM = "keep-token/transferred_from"
export const KEEP_TOKEN_TRANSFERRED_TO = "keep-token/transferred_to"

export const KEEP_TOKEN_BALANCE_REQUEST = "keep-token/balance_request"
export const KEEP_TOKEN_BALANCE_REQUEST_SUCCESS =
  "keep-token/balance_request_success"
export const KEEP_TOKEN_BALANCE_REQUEST_FAILURE =
  "keep-token/balance_request_failure"

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
