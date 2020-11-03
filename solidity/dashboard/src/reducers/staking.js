import moment from "moment"
import { add } from "../utils/arithmetics.utils"
import { findIndexAndObject, compareEthAddresses } from "../utils/array.utils"
import { isSameEthAddress } from "../utils/general.utils"

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
  delegationsFetchingStatus: null,
}

const stakingReducer = (state = initialState, action) => {
  switch (action.type) {
    case "staking/fetch_delegations_start":
      return {
        ...state,
        isDelegationDataFetching: true,
        error: null,
        delegationsFetchingStatus: null,
      }
    case "staking/fetch_delegations_success":
      return {
        ...state,
        isDelegationDataFetching: false,
        ...action.payload,
        delegationsFetchingStatus: "completed",
      }
    case "staking/fetch_delegations_failure":
      return {
        ...state,
        isDelegationDataFetching: false,
        error: action.payload.error,
        delegationsFetchingStatus: "failure",
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
    case "staking/add_undelegation":
      return {
        ...state,
        undelegations: [action.payload, ...state.undelegations],
      }
    case "staking/remove_undelegation":
      return {
        ...state,
        undelegations: removeFromDelegationOrUndelegation(
          [...state.undelegations],
          action.payload
        ),
      }
    case "staking/top_up_initiated":
      return {
        ...state,
        topUps: topUpInitiated([...state.topUps], action.payload),
      }
    case "staking/top_up_completed":
      return {
        ...state,
        topUps: state.topUps.filter(
          ({ operatorAddress }) =>
            !isSameEthAddress(operatorAddress, action.payload.operator)
        ),
        delegations: topUpCompleted([...state.delegations], action.payload),
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

const topUpInitiated = (topUps, { operator, topUp }) => {
  const { indexInArray, obj: topUpToUpdate } = findIndexAndObject(
    "operatorAddress",
    operator,
    topUps,
    compareEthAddresses
  )
  if (indexInArray === null) {
    return [
      {
        operatorAddress: operator,
        availableTopUpAmount: topUp,
        createdAt: moment.unix(),
      },
      ...topUps,
    ]
  }

  topUpToUpdate.availableTopUpAmount = add(
    topUpToUpdate.availableTopUpAmount,
    topUp
  )
  topUpToUpdate.createdAt = moment.unix()
  topUps[indexInArray] = topUpToUpdate

  return topUps
}

const topUpCompleted = (delegations, { operator, newAmount }) => {
  const { indexInArray, obj: delegationsToUpdate } = findIndexAndObject(
    "operatorAddress",
    operator,
    delegations,
    compareEthAddresses
  )
  if (indexInArray === null) {
    return delegations
  }

  delegations[indexInArray] = { ...delegationsToUpdate, amount: newAmount }

  return delegations
}
