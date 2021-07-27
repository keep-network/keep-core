import {
  FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_START,
  FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_SUCCESS,
  FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_FAILURE,
  KEEP_RANDOM_BEACON_AUTHORIZED,
} from "../actions"
import { findIndexAndObject, compareEthAddresses } from "../utils/array.utils"

const initialState = {
  authData: [],
  isFetching: false,
  error: null,
}

const authorizationReducer = (state = initialState, action) => {
  switch (action.type) {
    case FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_START:
      return {
        ...state,
        isFetching: true,
      }
    case FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_SUCCESS:
      return {
        ...state,
        isFetching: false,
        authData: action.payload,
      }
    case FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_FAILURE:
      return {
        ...state,
        isFetching: false,
        error: action.payload.error,
      }
    case KEEP_RANDOM_BEACON_AUTHORIZED:
      return {
        ...state,
        authData: updateAuthData([...state.authData], { ...action.payload }),
      }
    default:
      return state
  }
}

const updateAuthData = (authData, { contractName, operatorAddress }) => {
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

  const { indexInArray: contractIndexInArray, obj: obsoleteContract } =
    findIndexAndObject("contractName", contractName, obsoleteOperator.contracts)
  const updatedContracts = [...obsoleteOperator.contracts]
  updatedContracts[contractIndexInArray] = {
    ...obsoleteContract,
    isAuthorized: true,
  }
  const updatedOperators = [...authData]
  updatedOperators[operatorIndexInArray] = {
    ...obsoleteOperator,
    contracts: updatedContracts,
  }

  return updatedOperators
}

export default authorizationReducer
