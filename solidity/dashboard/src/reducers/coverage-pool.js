import {
  COVERAGE_POOL_FETCH_TVL_START,
  COVERAGE_POOL_FETCH_TVL_SUCCESS,
  COVERAGE_POOL_FETCH_TVL_ERROR,
  COVERAGE_POOL_FETCH_COV_POOL_DATA_SUCCESS,
  COVERAGE_POOL_FETCH_COV_POOL_DATA_ERROR,
  COVERAGE_POOL_FETCH_COV_POOL_DATA_START,
} from "../actions/coverage-pool"

export const coveragePoolInitialData = {
  // TVL data
  totalValueLocked: 0,
  isTotalValueLockedFetching: false,
  tvlError: null,

  isDataFetching: false,
  shareOfPool: 0,
  covBalance: 0,
  covTotalSupply: 0,
  error: null,
  estimatedRewards: 0,
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
        isTotalValueLockedFetching: false,
        tvlError: null,
      }
    case COVERAGE_POOL_FETCH_TVL_ERROR:
      return {
        ...state,
        isTotalValueLockedFetching: false,
        tvlError: action.payload.error,
      }

    case COVERAGE_POOL_FETCH_COV_POOL_DATA_START:
      return {
        ...state,
        isDataFetching: true,
      }
    case COVERAGE_POOL_FETCH_COV_POOL_DATA_SUCCESS:
      return {
        ...state,
        shareOfPool: action.payload.shareOfPool,
        covBalance: action.payload.covBalance,
        covTotalSupply: action.payload.covTotalSupply,
        estimatedRewards: action.payload.estimatedRewards,
        isDataFetching: false,
        error: null,
      }
    case COVERAGE_POOL_FETCH_COV_POOL_DATA_ERROR:
      return {
        ...state,
        isDataFetching: false,
        error: action.payload.error,
      }
    default:
      return state
  }
}

export default coveragePoolReducer
