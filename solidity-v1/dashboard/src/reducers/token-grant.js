import { isSameEthAddress } from "../utils/general.utils"
import { add, sub, gt } from "../utils/arithmetics.utils"
import { findIndexAndObject } from "../utils/array.utils"

const initialState = {
  isFetching: false,
  grants: [],
  error: null,
}

const tokenGrantsReducer = (state = initialState, action) => {
  switch (action.type) {
    case "token-grant/fetch_grants_start":
      return {
        ...state,
        isFetching: true,
        error: null,
      }
    case "token-grant/fetch_grants_success":
      return { ...state, isFetching: false, grants: action.payload }
    case "token-grant/fetch_grants_failure":
      return { ...state, isFetching: false, error: action.payload.error }
    case "token-grant/grant_staked":
      return {
        ...state,
        ...findGrantAndUpdate(
          [...state.grants],
          { ...state.selectedGrant },
          action.payload,
          grantStaked
        ),
      }
    case "token-grant/grant_withdrawn":
      return {
        ...state,
        ...findGrantAndUpdate(
          [...state.grants],
          { ...state.selectedGrant },
          action.payload,
          grantWithdrawn
        ),
      }
    case "token-grant/grant_deposited":
      return {
        ...state,
        ...findGrantAndUpdate(
          [...state.grants],
          { ...state.selectedGrant },
          action.payload,
          grantDeposited
        ),
      }
    default:
      return state
  }
}

export default tokenGrantsReducer

const findGrantAndUpdate = (
  grants,
  selectedGrant,
  payload,
  updateGrantCallback
) => {
  const { grantId } = payload
  const { indexInArray, obj: grantToUpdate } = findIndexAndObject(
    "id",
    grantId,
    grants
  )
  if (indexInArray === null) {
    return { grants }
  }
  const updatedGrant = updateGrantCallback(grantToUpdate, payload)
  grants[indexInArray] = updatedGrant

  if (selectedGrant.id === grantId) {
    selectedGrant = updatedGrant
  }

  return { grants, selectedGrant }
}

const grantStaked = (grantToUpdate, { value }) => {
  grantToUpdate.staked = add(grantToUpdate.staked, value).toString()
  grantToUpdate.readyToRelease = sub(
    grantToUpdate.readyToRelease,
    value
  ).toString()

  grantToUpdate.readyToRelease = gt(grantToUpdate.readyToRelease, 0)
    ? grantToUpdate.readyToRelease
    : "0"
  grantToUpdate.availableToStake = sub(grantToUpdate.availableToStake, value)

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

const grantDeposited = (
  grantToUpdate,
  { amount, availableToWitdrawGrant, operator }
) => {
  grantToUpdate.staked = sub(grantToUpdate.staked, amount)
  grantToUpdate.withdrawableAmountGrantOnly = availableToWitdrawGrant
  grantToUpdate.readyToRelease = add(grantToUpdate.readyToRelease, amount)
  grantToUpdate.escrowOperatorsToWithdraw = [
    ...grantToUpdate.escrowOperatorsToWithdraw,
    operator,
  ]
  grantToUpdate.availableToStake = add(grantToUpdate.availableToStake, amount)

  return grantToUpdate
}
