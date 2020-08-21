const FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_REQUEST =
  "FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_REQUEST"
const FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_SUCCESS =
  "FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_SUCCESS"
const FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_FAILURE =
  "FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_FAILURE"

export const INCREMENT_STEP = "INCREMENT_STEP"
export const DECREMENT_STEP = "DECREMENT_STEP"
export const SET_STRATEGY = "SET_STRATEGY"
export const SET_DELEGATION = "SET_DELEGATION"

export const copyStakeInitialData = {
  oldDelegations: [],
  oldDelegationsFetching: false,
  selectedOldDelegation: null,
  step: 0,
  // Copy stake to the new `TokenStaking` contract or only undelegate/recover delegation from an old `TokenStaking`
  selectedStrategy: null,
}

const copyStakeReducer = (state, action) => {
  switch (action.type) {
    case INCREMENT_STEP:
      return {
        ...state,
        step: state.step + 1,
      }
    case DECREMENT_STEP:
      return {
        ...state,
        step: state.step - 1,
      }
    case SET_STRATEGY:
      return { ...state, selectedStrategy: action.payload }
    case FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_SUCCESS:
      return {
        ...state,
        oldDelegations: action.payload,
        oldDelegationsFetching: false,
      }
    case FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_REQUEST:
      return {
        ...state,
        oldDelegationsFetching: true,
      }
    case FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_FAILURE:
      return {
        ...state,
        oldDelegationsFetching: false,
        error: action.payload,
      }
    case SET_DELEGATION: {
      return {
        ...state,
        selectedDelegations: action.payload,
      }
    }
    default:
      return state
  }
}

export default copyStakeReducer
