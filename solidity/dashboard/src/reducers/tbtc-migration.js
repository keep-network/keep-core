import { tbtcV2Migration } from "../actions"

const initialState = {
  tbtcV1Balance: "0",
  tbtcV2Balance: "0",
  unmintFee: "0",
  isFetching: false,
  error: null,
}

const tbtcV2MigrationReducer = (state = initialState, action) => {
  switch (action.type) {
    case tbtcV2Migration.TBTCV2_MIGRATION_FETCH_DATA_START:
      return {
        ...state,
        isFetching: true,
      }
    case tbtcV2Migration.TBTCV2_MIGRATION_FETCH_DATA_SUCCESS:
      return {
        ...state,
        tbtcV1Balance: action.payload.tbtcV1Balance,
        tbtcV2Balance: action.payload.tbtcV2Balance,
        unmintFee: action.payload.unmintFee,
        isFetching: false,
        error: null,
      }
    case tbtcV2Migration.TBTCV2_MIGRATION_FETCH_DATA_ERROR:
      return {
        ...state,
        isFetching: false,
        error: action.payload.error,
      }
    default:
      return state
  }
}

export default tbtcV2MigrationReducer
