import {
  DECREMENT_STEP,
  INCREMENT_STEP,
  SET_STRATEGY,
  SET_DELEGATION,
  FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_FAILURE,
  FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_SUCCESS,
  FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_REQUEST,
} from "../actions"

export const copyStakeInitialData = {
  oldDelegations: [],
  oldDelegationsFetching: false,
  selectedOldDelegation: null,
  step: 0,
  // Copy stake to the new `TokenStaking` contract or only undelegate/recover delegation from an old `TokenStaking`
  selectedStrategy: null,
  selectedDelegation: null,
}

const copyStakeReducer = (state = copyStakeInitialData, action) => {
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
        oldDelegations: action.payload.delegations,
        oldInitializationPeriod: action.payload.initializationPeriod,
        oldUndelegationPeriod: action.payload.undelegationPeriod,
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
        selectedDelegation: action.payload,
      }
    }
    default:
      return state
  }
}

export default copyStakeReducer
