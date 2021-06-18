export const COVERAGE_POOL_FETCH_TVL_REQUEST = "coverage_pool/fetch_tvl_request"
export const COVERAGE_POOL_FETCH_TVL_START = "coverage_pool/fetch_tvl_start"
export const COVERAGE_POOL_FETCH_TVL_SUCCESS = "coverage_pool/fetch_tvl_success"
export const COVERAGE_POOL_FETCH_TVL_ERROR = "coverage_pool/fetch_tvl_error"

export const fetchTvlStart = () => {
  return {
    type: COVERAGE_POOL_FETCH_TVL_START,
  }
}

export const fetchTvlSuccess = (data) => {
  return {
    type: COVERAGE_POOL_FETCH_TVL_SUCCESS,
    payload: data,
  }
}
