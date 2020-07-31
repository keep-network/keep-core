import { add, sub, gt } from "../utils/arithmetics.utils"
import { findIndexAndObject, compareEthAddresses } from "../utils/array.utils"
import { isSameEthAddress } from "../utils/general.utils"
import moment from "moment"

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
export const GRANT_DEPOSITED = "UPDATE_GRANT_DATA"

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
        grants: findGrantAndUpdate(
          [...state.grants],
          action.payload,
          grantStaked
        ),
      }
    case GRANT_WITHDRAWN:
      return {
        ...state,
        grants: findGrantAndUpdate(
          [...state.grants],
          action.payload,
          grantWithdrawn
        ),
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
        delegations: topUpCompleted([...state.delegations], action.payload),
      }
    case GRANT_DEPOSITED:
      return {
        ...state,
        grants: findGrantAndUpdate(
          [...state.grants],
          action.payload,
          grantDeposited
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

const grantDeposited = (
  grantToUpdate,
  { amount, availableToWithdrawEscrow, availableToWitdrawGrant, operator }
) => {
  grantToUpdate.staked = sub(grantToUpdate.staked, amount)
  grantToUpdate.withdrawableAmountGrantOnly = availableToWitdrawGrant
  grantToUpdate.readyToRelease = add(
    grantToUpdate.readyToRelease,
    availableToWithdrawEscrow
  )
  grantToUpdate.escrowOperatorsToWithdraw = [
    ...grantToUpdate.escrowOperatorsToWithdraw,
    operator,
  ]

  return grantToUpdate
}

const grantStaked = (grantToUpdate, { amount, availableToStake }) => {
  grantToUpdate.staked = add(grantToUpdate.staked, amount)
  grantToUpdate.readyToRelease = sub(grantToUpdate.readyToRelease, amount)
  grantToUpdate.readyToRelease = gt(grantToUpdate.readyToRelease, 0)
    ? grantToUpdate.readyToRelease
    : "0"
  grantToUpdate.availableToStake = availableToStake

  return grantToUpdate
}

const grantWithdrawn = (
  grantToUpdate,
  { amount, availableToStake, operator }
) => {
  grantToUpdate.readyToRelease = sub(grantToUpdate.readyToRelease, amount)
  grantToUpdate.released = add(grantToUpdate.released, amount)
  const unlocked = add(grantToUpdate.released, grantToUpdate.staked)
  if (!gt(unlocked, grantToUpdate.amount)) {
    grantToUpdate.unlocked = unlocked
  }
  grantToUpdate.availableToStake = availableToStake
  if (operator) {
    grantToUpdate.escrowOperatorsToWithdraw = [
      ...grantToUpdate.escrowOperatorsToWithdraw,
    ].filter((escrowOperator) => !isSameEthAddress(operator, escrowOperator))
  }

  return grantToUpdate
}

const findGrantAndUpdate = (grants, payload, updateGrantCallback) => {
  const { grantId } = payload
  const { indexInArray, obj: grantToUpdate } = findIndexAndObject(
    "id",
    grantId,
    grants
  )
  if (indexInArray === null) {
    return grants
  }

  grants[indexInArray] = updateGrantCallback(grantToUpdate, payload)

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

export default tokensPageReducer
