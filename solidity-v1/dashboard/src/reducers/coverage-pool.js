import {
  COVERAGE_POOL_FETCH_TVL_START,
  COVERAGE_POOL_FETCH_TVL_SUCCESS,
  COVERAGE_POOL_FETCH_TVL_ERROR,
  COVERAGE_POOL_FETCH_COV_POOL_DATA_SUCCESS,
  COVERAGE_POOL_FETCH_COV_POOL_DATA_ERROR,
  COVERAGE_POOL_FETCH_COV_POOL_DATA_START,
  COVERAGE_POOL_COV_TOKEN_UPDATED,
  COVERAGE_POOL_FETCH_APY_START,
  COVERAGE_POOL_FETCH_APY_SUCCESS,
  COVERAGE_POOL_FETCH_APY_ERROR,
} from "../actions/coverage-pool"

export const coveragePoolInitialData = {
  // TVL data
  totalValueLocked: 0,
  totalValueLockedInUSD: 0,
  isTotalValueLockedFetching: false,
  tvlError: null,
  totalAllocatedRewards: 0,
  totalCoverageClaimed: 0,

  // APY
  apy: 0,
  isApyFetching: false,
  apyError: null,

  isDataFetching: false,
  shareOfPool: 0,
  covBalance: 0,
  covTokensAvailableToWithdraw: 0,
  covTotalSupply: 0,
  error: null,
  estimatedRewards: 0,
  estimatedKeepBalance: 0,
  withdrawalDelay: 0,
  withdrawalTimeout: 0,
  pendingWithdrawal: 0,
  withdrawalInitiatedTimestamp: 0,

  // riskManager
  hasRiskManagerOpenAuctions: false,
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
        totalValueLocked: action.payload.tvl,
        totalValueLockedInUSD: action.payload.tvlInUSD,
        totalAllocatedRewards: action.payload.totalAllocatedRewards,
        totalCoverageClaimed: action.payload.totalCoverageClaimed,
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
    case COVERAGE_POOL_COV_TOKEN_UPDATED:
      return {
        ...state,
        ...action.payload,
        isDataFetching: false,
        error: null,
      }
    case COVERAGE_POOL_FETCH_COV_POOL_DATA_ERROR:
      return {
        ...state,
        isDataFetching: false,
        error: action.payload.error,
      }
    case COVERAGE_POOL_FETCH_APY_START:
      return {
        ...state,
        isApyFetching: true,
      }
    case COVERAGE_POOL_FETCH_APY_SUCCESS:
      return {
        ...state,
        isApyFetching: false,
        apy: action.payload,
        apyError: null,
      }
    case COVERAGE_POOL_FETCH_APY_ERROR:
      return {
        ...state,
        isApyFetching: false,
        apyError: action.payload.error,
      }
    default:
      return state
  }
}

export default coveragePoolReducer
