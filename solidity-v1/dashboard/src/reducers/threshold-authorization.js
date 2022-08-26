import {
  FETCH_THRESHOLD_AUTH_DATA_START,
  FETCH_THRESHOLD_AUTH_DATA_SUCCESS,
  FETCH_THRESHOLD_AUTH_DATA_FAILURE,
  THRESHOLD_AUTHORIZED,
  THRESHOLD_STAKED_TO_T,
  REMOVE_STAKE_FROM_THRESHOLD_AUTH_DATA,
  ADD_STAKE_TO_THRESHOLD_AUTH_DATA,
} from "../actions"
import { findIndexAndObject, compareEthAddresses } from "../utils/array.utils"
import { isSameEthAddress } from "../utils/general.utils"
import { AUTH_CONTRACTS_LABEL } from "../constants/constants"

const initialState = {
  authData: [],
  isFetching: false,
  error: null,
}

const thresholdAuthorizationReducer = (state = initialState, action) => {
  switch (action.type) {
    case FETCH_THRESHOLD_AUTH_DATA_START:
      return {
        ...state,
        isFetching: true,
      }
    case FETCH_THRESHOLD_AUTH_DATA_SUCCESS:
      return {
        ...state,
        isFetching: false,
        authData: action.payload,
      }
    case FETCH_THRESHOLD_AUTH_DATA_FAILURE:
      return {
        ...state,
        isFetching: false,
        error: action.payload.error,
      }
    case THRESHOLD_AUTHORIZED:
      return {
        ...state,
        authData: authorizeThresholdContract([...state.authData], {
          ...action.payload,
        }),
      }
    case THRESHOLD_STAKED_TO_T:
      return {
        ...state,
        authData: updateThresholdAuthData([...state.authData], {
          ...action.payload,
        }),
      }
    case ADD_STAKE_TO_THRESHOLD_AUTH_DATA:
      return {
        ...state,
        authData: addStakeToAuthData([...state.authData], action.payload),
      }
    case REMOVE_STAKE_FROM_THRESHOLD_AUTH_DATA:
      return {
        ...state,
        authData: removeStakeFromAuthData([...state.authData], action.payload),
      }
    default:
      return state
  }
}

const authorizeThresholdContract = (authData, { operatorAddress }) => {
  const { indexInArray: operatorIndexInArray, obj: obsoleteOperator } =
    findIndexAndObject(
      "operatorAddress",
      operatorAddress,
      authData,
      compareEthAddresses
    )

  if (operatorIndexInArray === null) {
    return authData
  }

  const updatedContracts = {
    ...obsoleteOperator.contract,
    isAuthorized: true,
  }
  const updatedOperators = [...authData]
  updatedOperators[operatorIndexInArray] = {
    ...obsoleteOperator,
    contract: updatedContracts,
  }

  return updatedOperators
}

const updateThresholdAuthData = (authData, { operatorAddress }) => {
  const { indexInArray: operatorIndexInArray, obj: obsoleteOperator } =
    findIndexAndObject(
      "operatorAddress",
      operatorAddress,
      authData,
      compareEthAddresses
    )
  if (operatorIndexInArray === null) {
    return authData
  }

  const updatedOperators = [...authData]
  updatedOperators[operatorIndexInArray] = {
    ...obsoleteOperator,
    isStakedToT: true,
  }

  return updatedOperators
}

const addStakeToAuthData = (
  authData,
  {
    owner,
    authorizerAddress,
    beneficiary: beneficiaryAddress,
    operatorAddress,
    amount: stakeAmount,
    isFromGrant,
    operatorContractAddress,
  }
) => {
  const _isFromGrant = !!isFromGrant

  const newStake = {
    owner,
    authorizerAddress,
    beneficiaryAddress,
    operatorAddress,
    stakeAmount,
    contract: {
      contractName: AUTH_CONTRACTS_LABEL.THRESHOLD_TOKEN_STAKING,
      operatorContractAddress: operatorContractAddress,
      isAuthorized: false,
    },
    isStakedToT: false,
    isFromGrant: _isFromGrant,
    canBeMovedToT: !_isFromGrant,
    isInInitializationPeriod: true,
  }

  const updatedAuthData = [...authData, newStake]
  return updatedAuthData
}

const removeStakeFromAuthData = (authData, operatorAddress) => {
  return authData.filter((stake) => {
    return !isSameEthAddress(stake.operatorAddress, operatorAddress)
  })
}

export default thresholdAuthorizationReducer
