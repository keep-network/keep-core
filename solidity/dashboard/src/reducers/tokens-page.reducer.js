import { add, sub, gte } from "../utils/arithmetics.utils"
import { findIndexAndObject, compareEthAddresses } from "../utils/array.utils"
import { isSameEthAddress } from "../utils/general.utils"

export const REFRESH_KEEP_TOKEN_BALANCE = "REFRESH_KEEP_TOKEN_BALANCE"
export const REFRESH_GRANT_TOKEN_BALANCE = "REFRESH_GRANT_TOKEN_BALANCE"
export const UPDATE_OWNED_UNDELEGATIONS_TOKEN_BALANCE =
  "UPDATE_OWNED_UNDELEGATIONS_BALANCE"
export const UPDATE_OWNED_DELEGATED_TOKENS_BALANCE =
  "UPDATE_OWNED_DELEGATED_TOKENS_BALANCE"
export const ADD_DELEGATION = "ADD_DELEGATION"
export const REMOVE_DELEGATION = "REMOVE_DELEGATION"
export const ADD_UNDELEGATION = "ADD_UNDELEGATION"
export const REMOVE_UNDELEGATION = "REMOVE_UNDELEGATION"
export const GRANT_STAKED = "GRANT_STAKED"
export const GRANT_WITHDRAWN = "GRANT_WITHDRAWN"
export const SET_STATE = "SET_STATE"
export const SET_SELECTED_GRANT = "SET_SELECTED_GRANT"
export const SET_TOKENS_CONTEXT = "SET_TOKENS_CONTEXT"
export const TOP_UP_INITIATED = "TOP_UP_INITIATED"
export const TOP_UP_COMPLETED = "TOP_UP_COMPLETED"

const tokensPageReducer = (state, action) => {
  switch (action.type) {
    case SET_STATE:
      return {
        ...state,
        ...action.payload,
      }
    case REFRESH_KEEP_TOKEN_BALANCE:
      return {
        ...state,
        keepTokenBalance: action.payload,
      }
    case REFRESH_GRANT_TOKEN_BALANCE:
      return {
        ...state,
        grantTokenBalance: action.payload,
      }
    case UPDATE_OWNED_UNDELEGATIONS_TOKEN_BALANCE:
      return {
        ...state,
        ownedTokensUndelegationsBalance: action.payload.operation(
          state.ownedTokensUndelegationsBalance,
          action.payload.value
        ),
      }
    case UPDATE_OWNED_DELEGATED_TOKENS_BALANCE:
      return {
        ...state,
        ownedTokensDelegationsBalance: action.payload.operation(
          state.ownedTokensDelegationsBalance,
          action.payload.value
        ),
      }
    case ADD_DELEGATION:
      return {
        ...state,
        delegations: [action.payload, ...state.delegations],
      }
    case REMOVE_DELEGATION:
      return {
        ...state,
        delegations: removeFromDelegationOrUndelegation(
          [...state.delegations],
          action.payload
        ),
      }
    case ADD_UNDELEGATION:
      return {
        ...state,
        undelegations: [action.payload, ...state.undelegations],
      }
    case REMOVE_UNDELEGATION:
      return {
        ...state,
        undelegations: removeFromDelegationOrUndelegation(
          [...state.undelegations],
          action.payload
        ),
      }
    case GRANT_STAKED:
      return {
        ...state,
        grants: grantStaked([...state.grants], action.payload),
      }
    case GRANT_WITHDRAWN:
      return {
        ...state,
        grants: grantWithdrawn([...state.grants], action.payload),
      }
    case SET_SELECTED_GRANT:
      return {
        ...state,
        selectedGrant: action.payload,
      }
    case SET_TOKENS_CONTEXT:
      return {
        ...state,
        tokensContext: action.payload,
      }
    case TOP_UP_INITIATED:
      return {
        ...state,
        availableTopUps: topUpInitiated(
          [...state.availableTopUps],
          action.payload
        ),
      }
    case TOP_UP_COMPLETED:
      return {
        ...state,
        availableTopUps: state.availableTopUps.filter(
          ({ operatorAddress }) =>
            !isSameEthAddress(operatorAddress, action.payload.operator)
        ),
        delegations: updateDelegationAmount(
          [...state.delegations],
          action.payload
        ),
      }
    default:
      return state
  }
}

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

const grantStaked = (grants, { grantId, amount, availableToStake }) => {
  const { indexInArray, obj: grantToUpdate } = findIndexAndObject(
    "id",
    grantId,
    grants
  )
  if (indexInArray === null) {
    return grants
  }
  grantToUpdate.staked = add(grantToUpdate.staked, amount)
  grantToUpdate.readyToRelease = sub(grantToUpdate.readyToRelease, amount)
  grantToUpdate.readyToRelease = gte(grantToUpdate.readyToRelease, 0)
    ? grantToUpdate.readyToRelease
    : "0"
  grantToUpdate.availableToStake = availableToStake
  grants[indexInArray] = grantToUpdate

  return grants
}

const grantWithdrawn = (grants, { grantId, amount, availableToStake }) => {
  const { indexInArray, obj: grantToUpdate } = findIndexAndObject(
    "id",
    grantId,
    grants
  )
  if (indexInArray === null) {
    return grants
  }
  grantToUpdate.readyToRelease = "0"
  grantToUpdate.released = add(grantToUpdate.released, amount)
  grantToUpdate.unlocked = add(grantToUpdate.released, grantToUpdate.staked)
  grantToUpdate.availableToStake = availableToStake
  grantToUpdate.withdrawn = add(grantToUpdate.withdrawn, amount)
  grants[indexInArray] = grantToUpdate

  return grants
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
      { operatorAddress: operator, availableTopUpAmount: topUp },
      ...topUps,
    ]
  }

  topUpToUpdate.availableTopUpAmount = add(
    topUpToUpdate.availableTopUpAmount,
    topUp
  )
  topUps[indexInArray] = topUpToUpdate

  return topUps
}

const updateDelegationAmount = (delegations, { operator, newAmount }) => {
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

export default tokensPageReducer
