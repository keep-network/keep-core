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
    default:
      return state
  }
}

export default stakingReducer
