import * as KeepBalanceActions from "./keep-balance"
import * as modal from "./modal"

// COPY STAKE
export const FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_REQUEST =
  "copy-stake/fetch_old_delegations_request"
export const FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_SUCCESS =
  "copy-stake/fetch_old_delegations_success"
export const FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_FAILURE =
  "copy-stake/fetch_old_delegations_failure"
export const INCREMENT_STEP = "copy-stake/increment_step"
export const DECREMENT_STEP = "copy-stake/decrement_step"
export const RESET_COPY_STAKE_FLOW = "copy-stake/reset_flow"
export const SET_STRATEGY = "copy-stake/set_strategy"
export const SET_DELEGATION = "copy-stake/set_delegation"
export const COPY_STAKE_REQUEST = "copy-stake/copy_stake_request"

// OPERATOR
export const FETCH_OPERATOR_DELEGATIONS_RERQUEST =
  "operator/fetch_delegations_request"
export const FETCH_OPERATOR_DELEGATIONS_START =
  "operator/fetch_delegations_start"
export const FETCH_OPERATOR_DELEGATIONS_SUCCESS =
  "operator/fetch_delegations_success"
export const FETCH_OPERATOR_DELEGATIONS_FAILURE =
  "operator/fetch_delegations_failure"
export const FETCH_OPERATOR_SLASHED_TOKENS_RERQUEST =
  "operator/fetch_slashed_tokens_request"
export const FETCH_OPERATOR_SLASHED_TOKENS_START =
  "operator/fetch_slashed_tokens_start"
export const FETCH_OPERATOR_SLASHED_TOKENS_SUCCESS =
  "operator/fetch_slashed_tokens_success"
export const FETCH_OPERATOR_SLASHED_TOKENS_FAILURE =
  "operator/fetch_slashed_tokens_failure"
export const OPERATOR_DELEGATION_UNDELEGATED = "operator/delegation_undelegated"
export const OPERATR_DELEGATION_CANCELED = "operator/delegation_canceled"

// AUTHORIZATION RANDOM BEACON
export const FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_REQUEST =
  "authorization_beacon/fetch_auth_data_request"
export const FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_START =
  "authorization_beacon/fetch_auth_data_start"
export const FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_SUCCESS =
  "authorization_beacon/fetch_auth_data_success"
export const FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_FAILURE =
  "authorization_beacon/fetch_auth_data_failure"
export const KEEP_RANDOM_BEACON_AUTHORIZED =
  "authorization_beacon/beacon_authorized"

// AUTHORIZATION THRESHOLD
export const FETCH_THRESHOLD_AUTH_DATA_REQUEST =
  "threshold/fetch_auth_data_request"
export const FETCH_THRESHOLD_AUTH_DATA_START = "threshold/fetch_auth_data_start"
export const FETCH_THRESHOLD_AUTH_DATA_SUCCESS =
  "threshold/fetch_auth_data_success"
export const FETCH_THRESHOLD_AUTH_DATA_FAILURE =
  "threshold/fetch_auth_data_failure"
export const THRESHOLD_AUTHORIZED = "threshold/contract_authorized"
export const THRESHOLD_STAKED_TO_T = "threshold/staked_to_t"
export const ADD_STAKE_TO_THRESHOLD_AUTH_DATA =
  "threshold/add_stake_to_threshold_auth_data"
export const REMOVE_STAKE_FROM_THRESHOLD_AUTH_DATA =
  "threshold/remove_stake_from_threshold_auth_data"

export const keepBalanceActions = {
  ...KeepBalanceActions,
}

export const modalActions = {
  ...modal,
}
