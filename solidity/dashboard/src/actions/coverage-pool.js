export const COVERAGE_POOL_FETCH_TVL_REQUEST = "coverage_pool/fetch_tvl_request"
export const COVERAGE_POOL_FETCH_TVL_START = "coverage_pool/fetch_tvl_start"
export const COVERAGE_POOL_FETCH_TVL_SUCCESS = "coverage_pool/fetch_tvl_success"
export const COVERAGE_POOL_FETCH_TVL_ERROR = "coverage_pool/fetch_tvl_error"

export const COVERAGE_POOL_FETCH_SHARE_OF_POOL_REQUEST =
  "coverage_pool/fetch_share_of_pool_request"
export const COVERAGE_POOL_FETCH_SHARE_OF_POOL_START =
  "coverage_pool/fetch_share_of_pool_start"
export const COVERAGE_POOL_FETCH_SHARE_OF_POOL_SUCCESS =
  "coverage_pool/fetch_share_of_pool_success"
export const COVERAGE_POOL_FETCH_SHARE_OF_POOL_ERROR =
  "coverage_pool/fetch_share_of_pool_error"

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

export const fetchTvlRequest = () => {
  return {
    type: COVERAGE_POOL_FETCH_TVL_REQUEST,
  }
}

export const fetchShareOfPoolStart = () => {
  return {
    type: COVERAGE_POOL_FETCH_SHARE_OF_POOL_START,
  }
}

export const fetchShareOfPoolRequest = (address) => {
  return {
    type: COVERAGE_POOL_FETCH_SHARE_OF_POOL_REQUEST,
    payload: { address },
  }
}

/**
 * @param {Object} data Cov token info.
 * @param {string} data.shareOfPool The share of the pool.
 * @param {string} data.covBalance The user's token balance.
 * @param {string} data.covTotalSupply The total supply of the cov token.
 *
 * @return { { type: string, payload: object }}
 */
export const fetchShareOfPoolSuccess = (data) => {
  return {
    type: COVERAGE_POOL_FETCH_SHARE_OF_POOL_SUCCESS,
    payload: data,
  }
}
