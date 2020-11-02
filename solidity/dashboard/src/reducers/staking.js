import { findIndexAndObject, compareEthAddresses } from "../utils/array.utils"

const initialState = {
  isDelegationDataFetching: false,
  delegations: [],
  undelegations: [],
  ownedTokensDelegationsBalance: "0",
  ownedTokensUndelegationsBalance: "0",
  minimumStake: "0",
  initializationPeriod: "0",
  undelegationPeriod: "0",
  error: null,
  topUps: [],
  areTopUpsFetching: false,
  topUpsFetchingError: null,
}

const stakingReducer = (state = initialState, action) => {
  switch (action.type) {
    case "staking/fetch_delegations_start":
      return {
        ...state,
        isDelegationDataFetching: true,
        error: null,
      }
    case "staking/fetch_delegations_success":
      return { ...state, isDelegationDataFetching: false, ...action.payload }
    case "staking/fetch_delegations_failure":
      return {
        ...state,
        isDelegationDataFetching: false,
        error: action.payload.error,
      }
    case "staking/fetch_top_ups_start":
      return {
        ...state,
        areTopUpsFetching: true,
        topUpsFetchingError: null,
      }
    case "staking/fetch_top_ups_success":
      return {
        ...state,
        areTopUpsFetching: false,
        topUps: action.payload,
      }
    case "staking/fetch_top_ups_failure":
      return {
        ...state,
        areTopUpsFetching: false,
        topUpsFetchingError: action.payload.error,
      }
    case "staking/add_delegation":
      return {
        ...state,
        delegations: [action.payload, ...state.delegations],
      }
    case "staking/remove_delegation":
      return {
        ...state,
        undelegations: removeFromDelegationOrUndelegation(
          [...state.undelegations],
          action.payload
        ),
      }
    case "staking/update_owned_undelegations_tokens_balance":
      return {
        ...state,
        ownedTokensUndelegationsBalance: action.payload.operation(
          state.ownedTokensUndelegationsBalance,
          action.payload.value
        ),
      }
    case "staking/update_owned_delegated_tokens_balance":
      return {
        ...state,
        ownedTokensDelegationsBalance: action.payload.operation(
          state.ownedTokensDelegationsBalance,
          action.payload.value
        ),
      }
    default:
      return state
  }
}

export default stakingReducer

const removeFromDelegationOrUndelegation = (array, id) => {
  const { indexInArray } = findIndexAndObject(
    "operatorAddress",
    id,
    array,
    compareEthAddresses
  )
  if (indexInArray === null) {
    return array
  }
  array.splice(indexInArray, 1)

  return array
}
