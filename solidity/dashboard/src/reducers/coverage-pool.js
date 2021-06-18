import {
  COVERAGE_POOL_FETCH_TVL_START,
  COVERAGE_POOL_FETCH_TVL_SUCCESS,
  COVERAGE_POOL_FETCH_TVL_ERROR,
} from "../actions/coverage-pool"

export const coveragePoolInitialData = {
  // TVL data
  totalValueLocked: 0,
  isTotalValueLockedFetching: false,
  tvlError: null,

  // shareOfPool
  shareOfPool: 0,
  isShareOfPoolFetching: false,
  shareOfPoolError: null,

  weeklyRoi: 0,
  rewards: 0,
}

const coveragePoolReducer = (state = coveragePoolInitialData, action) => {
  switch (action.type) {
    case COVERAGE_POOL_FETCH_TVL_START:
      return {
        ...state,
        isTotalValueLockedFetching: true,
      }
    case COVERAGE_POOL_FETCH_TVL_SUCCESS:
      return {
        ...state,
        totalValueLocked: action.payload,
        tvlError: null,
      }
    case COVERAGE_POOL_FETCH_TVL_ERROR:
      return {
        ...state,
        isTotalValueLockedFetching: false,
        tvlError: action.payload.error,
      }
    default:
      return state
  }
}

export default coveragePoolReducer
