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
    default:
      return state
  }
}

export default stakingReducer
