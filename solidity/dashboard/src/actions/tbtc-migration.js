export const TBTCV2_MIGRATION_FETCH_DATA_REQUEST =
  "tbtcv2_migration/fetch_data_request"
export const TBTCV2_MIGRATION_FETCH_DATA_START =
  "tbtcv2_migration/fetch_data_start"
export const TBTCV2_MIGRATION_FETCH_DATA_SUCCESS =
  "tbtcv2_migration/fetch_data_success"
export const TBTCV2_MIGRATION_FETCH_DATA_ERROR =
  "tbtcv2_migration/fetch_data_error"

export const TBTCV2_TOKEN_MINTED_EVENT_EMITTED =
  "tbtcv2_migration/token_minted_event_emitted"
export const TBTCV2_TOKEN_UNMINTED_EVENT_EMITTED =
  "tbtcv2_migration/token_unminted_event_emitted"
export const TBTCV2_TOKEN_MINTED = "tbtcv2_migration/token_minted"
export const TBTCV2_TOKEN_UNMINTED = "tbtcv2_migration/token_unminted"

export const fetchDataStart = () => {
  return {
    type: TBTCV2_MIGRATION_FETCH_DATA_START,
  }
}

export const fetchDataSuccess = (data) => {
  return {
    type: TBTCV2_MIGRATION_FETCH_DATA_SUCCESS,
    payload: data,
  }
}

export const fetchDataRequest = (address) => {
  return {
    type: TBTCV2_MIGRATION_FETCH_DATA_REQUEST,
    payload: { address },
  }
}
