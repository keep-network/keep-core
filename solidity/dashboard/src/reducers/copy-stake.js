import {
  DECREMENT_STEP,
  INCREMENT_STEP,
  SET_STRATEGY,
  SET_DELEGATION,
  FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_FAILURE,
  FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_SUCCESS,
  FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_REQUEST,
  RESET_COPY_STAKE_FLOW,
} from "../actions"
import { isSameEthAddress } from "../utils/general.utils"

export const copyStakeInitialData = {
  oldDelegations: [],
  oldDelegationsFetching: false,
  selectedOldDelegation: null,
  step: 0,
  // Copy stake to the new `TokenStaking` contract or only undelegate/recover delegation from an old `TokenStaking`
  selectedStrategy: null,
  selectedDelegation: null,
  oldUndelegationPeriod: 0,
  isProcessing: false,
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
    case RESET_COPY_STAKE_FLOW: {
      return {
        ...state,
        step: 0,
        selectedDelegation: null,
        selectedStrategy: null,
      }
    }
    case "copy-stake/copy-stake_request":
    case "copy-stake/undelegate_request":
    case "copy-stake/recover_request": {
      return {
        ...state,
        isProcessing: true,
      }
    }
    case "copy-stake/copy-stake_success":
    case "copy-stake/recover_success":
    case "copy-stake/undelegation_success":
      return {
        ...state,
        isProcessing: false,
      }
    case "copy-stake/copy-stake_failure":
    case "copy-stake/recover_failure":
    case "copy-stake/undelegation_failure":
      return {
        ...state,
        isProcessing: false,
        error: action.payload,
      }
    case "copy-stake/remove_old_delegation": {
      return {
        ...state,
        oldDelegations: state.oldDelegations.filter(
          (delegation) =>
            !isSameEthAddress(delegation.operatorAddress, action.payload)
        ),
      }
    }
    default:
      return state
  }
}

export default copyStakeReducer
