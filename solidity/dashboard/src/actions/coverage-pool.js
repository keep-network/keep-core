export const COVERAGE_POOL_FETCH_TVL_REQUEST = "coverage_pool/fetch_tvl_request"
export const COVERAGE_POOL_FETCH_TVL_START = "coverage_pool/fetch_tvl_start"
export const COVERAGE_POOL_FETCH_TVL_SUCCESS = "coverage_pool/fetch_tvl_success"
export const COVERAGE_POOL_FETCH_TVL_ERROR = "coverage_pool/fetch_tvl_error"

export const COVERAGE_POOL_FETCH_COV_POOL_DATA_REQUEST =
  "coverage_pool/fetch_cov_pool_data_request"
export const COVERAGE_POOL_FETCH_COV_POOL_DATA_START =
  "coverage_pool/fetch_cov_pool_data_start"
export const COVERAGE_POOL_FETCH_COV_POOL_DATA_SUCCESS =
  "coverage_pool/fetch_cov_pool_data_success"
export const COVERAGE_POOL_FETCH_COV_POOL_DATA_ERROR =
  "coverage_pool/fetch_cov_pool_data_error"

export const COVERAGE_POOL_COV_TOKEN_TRANSFER_EVENT_EMITTED =
  "coverage_pool/cov_token_transfer_event_emitted"

export const COVERAGE_POOL_DEPOSIT_ASSET_POOL = "coverage_pool/deposit"

export const COVERAGE_POOL_COV_TOKEN_UPDATED = "coverage_pool/cov_token_updated"

export const COVERAGE_POOL_FETCH_APY_REQUEST = "coverage_pool/fetch_apy_request"
export const COVERAGE_POOL_FETCH_APY_START = "coverage_pool/fetch_apy_start"
export const COVERAGE_POOL_FETCH_APY_SUCCESS = "coverage_pool/fetch_apy_success"
export const COVERAGE_POOL_FETCH_APY_ERROR = "coverage_pool/fetch_apy_error"

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

export const fetchCovPoolDataStart = () => {
  return {
    type: COVERAGE_POOL_FETCH_COV_POOL_DATA_START,
  }
}

export const fetchCovPoolDataRequest = (address) => {
  return {
    type: COVERAGE_POOL_FETCH_COV_POOL_DATA_REQUEST,
    payload: { address },
  }
}

/**
 * @param {Object} data Cov token info.
 * @param {string} data.shareOfPool The share of the pool.
 * @param {string} data.covBalance The user's token balance.
 * @param {string} data.covTotalSupply The total supply of the cov token.
 * @param {string} data.estimatedRewards The estimated rewards.
 *
 * @return { { type: string, payload: object }}
 */
export const fetchCovPoolDataSuccess = (data) => {
  return {
    type: COVERAGE_POOL_FETCH_COV_POOL_DATA_SUCCESS,
    payload: data,
  }
}

export const covTokenTransferEventEmitted = (event) => {
  return {
    type: COVERAGE_POOL_COV_TOKEN_TRANSFER_EVENT_EMITTED,
    payload: { event },
  }
}

/**
 * @param {Object} data Cov token info.
 * @param {string} data.shareOfPool The share of the pool. Tha value should be
 * between [0, 1].
 * @param {string} data.covBalance The amount of tokens owned by user in the
 * smallest unit (18 decimals precision).
 * @param {string} data.covTotalSupply The total supply of the cov token in the
 * samallest unit (18 decimals precision).
 *
 * @return { { type: string, payload: object }}
 */
export const covTokenUpdated = (data) => {
  return {
    type: COVERAGE_POOL_COV_TOKEN_UPDATED,
    payload: data,
  }
}

export const depositAssetPool = (amount, meta) => {
  return {
    type: COVERAGE_POOL_DEPOSIT_ASSET_POOL,
    payload: {
      amount,
    },
    meta,
  }
}

export const fetchAPYRequest = () => {
  return {
    type: COVERAGE_POOL_FETCH_APY_REQUEST,
  }
}

export const fetchAPYStart = () => {
  return {
    type: COVERAGE_POOL_FETCH_APY_START,
  }
}

export const fetchAPYSuccess = (apy) => {
  return {
    type: COVERAGE_POOL_FETCH_APY_SUCCESS,
    payload: apy,
  }
}
