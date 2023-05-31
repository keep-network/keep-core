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

export const COVERAGE_POOL_WITHDRAWAL_INITIATED_EVENT_EMITTED =
  "coverage_pool/withdrawal_initiated_event_emitted"

export const COVERAGE_POOL_WITHDRAWAL_COMPLETED_EVENT_EMITTED =
  "coverage_pool/withdrawal_completed_event_emitted"

export const RISK_MANAGER_AUCTION_CREATED_EVENT_EMITTED =
  "risk_manager/auction_created_event_emitted"

export const RISK_MANAGER_AUCTION_CLOSED_EVENT_EMITTED =
  "risk_manager/auction_closed_event_emitted"

export const COVERAGE_POOL_WITHDRAW_ASSET_POOL = "coverage_pool/withdraw"
export const COVERAGE_POOL_CLAIM_TOKENS_FROM_WITHDRAWAL =
  "coverage_pool/claim_tokens"

export const COVERAGE_POOL_COV_TOKEN_UPDATED = "coverage_pool/cov_token_updated"

export const COVERAGE_POOL_FETCH_APY_REQUEST = "coverage_pool/fetch_apy_request"
export const COVERAGE_POOL_FETCH_APY_START = "coverage_pool/fetch_apy_start"
export const COVERAGE_POOL_FETCH_APY_SUCCESS = "coverage_pool/fetch_apy_success"
export const COVERAGE_POOL_FETCH_APY_ERROR = "coverage_pool/fetch_apy_error"
export const COVERAGE_POOL_REINITAITE_WITHDRAW =
  "coverage_pool/reinitiate_withdraw"

export const COVERAGE_POOL_INCREASE_WITHDRAWAL =
  "coverage_pool/increase_withdrawal"

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

export const coveragePoolWithdrawalInitiatedEventEmitted = (event) => {
  return {
    type: COVERAGE_POOL_WITHDRAWAL_INITIATED_EVENT_EMITTED,
    payload: { event },
  }
}

export const coveragePoolWithdrawalCompletedEventEmitted = (event) => {
  return {
    type: COVERAGE_POOL_WITHDRAWAL_COMPLETED_EVENT_EMITTED,
    payload: { event },
  }
}

export const riskManagerAuctionCreatedEventEmitted = (event) => {
  return {
    type: RISK_MANAGER_AUCTION_CREATED_EVENT_EMITTED,
    payload: { event },
  }
}

export const riskManagerAuctionClosedEventEmitted = (event) => {
  return {
    type: RISK_MANAGER_AUCTION_CLOSED_EVENT_EMITTED,
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

export const withdrawAssetPool = (amount, meta) => {
  return {
    type: COVERAGE_POOL_WITHDRAW_ASSET_POOL,
    payload: {
      amount,
    },
    meta,
  }
}

export const claimTokensFromWithdrawal = (meta) => {
  return {
    type: COVERAGE_POOL_CLAIM_TOKENS_FROM_WITHDRAWAL,
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

export const reinitiateWithdraw = () => {
  return {
    type: COVERAGE_POOL_REINITAITE_WITHDRAW,
  }
}

export const increaseWithdrawal = (amount) => {
  return {
    type: COVERAGE_POOL_INCREASE_WITHDRAWAL,
    payload: { amount },
  }
}
