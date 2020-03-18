import React, { useContext, useReducer, useCallback, useMemo } from 'react'
import { Web3Context } from '../components/WithWeb3Context'
import { findIndexAndObject, compareEthAddresses } from '../utils/array.utils'

const REFRESH_KEEP_TOKEN_BALANCE = 'REFRESH_KEEP_TOKEN_BALANCE'
const REFRESH_GRANT_TOKEN_BALANCE = 'REFRESH_GRANT_TOKEN_BALANCE'
const UPDATE_GRANT_BY_ID = 'UPDATE_GRANT_BY_ID'
const UPDATE_OWNED_UNDELEGATIONS_TOKEN_BALANCE = 'UPDATE_OWNED_UNDELEGATIONS_BALANCE'
const UPDATE_OWNED_DELEGATEN_TOKENS_BALANCE = 'UPDATE_DELEGATEN_TOKENS_BALANCE'
const ADD_DELEGATION = 'ADD_DELEGATION'
const REMOVE_DELEGATION = 'REMOVE_DELEGATION'
const ADD_UNDELEGATION = 'ADD_UNDELEGATION'
const REMOVE_UNDELEGATION = 'REMOVE_UNDELEGATION'

const TokensPageContext = React.createContext({
  refreshKeepTokenBalance: () => {},
  refreshGrantTokenBalance: () => {},
  dispatch: () => {},
  grants: [],
  delegations: [],
  undelegations: [],
  keepTokenBalance: '0',
  grantTokenBalance: '0',
  ownedTokensUndelegationsBalance: '0',
  ownedTokensDelegationsBalance: '0',
  initializationPeriod: '0',
  undelegationPeriod: '0',
})

const TokenPageContextProvider = () => {
  const web3Context = useContext(Web3Context)
  const [state, dispatch] = useReducer(tokensPageReducer, {
    grants: [],
    delegations: [],
    undelegations: [],
    keepTokenBalance: '0',
    grantTokenBalance: '0',
    ownedTokensUndelegationsBalance: '0',
    ownedTokensDelegationsBalance: '0',
    initializationPeriod: '0',
    undelegationPeriod: '0',
  })

  const contextValue = useMemo(() => {
    return { state, dispatch }
  }, [state, dispatch])

  const refreshKeepTokenBalance = useCallback(async () => {
    const { token, yourAddress } = web3Context

    const keepTokenBalance = await token.methods.balanceOf(yourAddress).call()
    dispatch({ type: REFRESH_KEEP_TOKEN_BALANCE, payload: keepTokenBalance })
  }, [web3Context, dispatch])

  const refreshGrantTokenBalance = useCallback(async () => {
    const { grantContract, yourAddress } = web3Context

    const grantTokenBalance = await grantContract.methods.balanceOf(yourAddress).call()
    dispatch({ type: REFRESH_GRANT_TOKEN_BALANCE, payload: grantTokenBalance })
  }, [web3Context, dispatch])

  return (
    <TokensPageContext.Provider value={{
      ...contextValue.state,
      dispatch: contextValue.dispatch,
      refreshKeepTokenBalance,
      refreshGrantTokenBalance,
    }}>
      {children}
    </TokensPageContext.Provider>
  )
}

const tokensPageReducer = (state, action) => {
  switch (action.type) {
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
  case UPDATE_GRANT_BY_ID:
    return {
      ...state,
      grants: updateGrants(state.grants, action.payload),
    }
  case UPDATE_OWNED_UNDELEGATIONS_TOKEN_BALANCE:
    return {
      ...state,
      ownedTokensUndelegationsBalance: action.payload.operation(state.ownedTokensUndelegationsBalance, action.payload.value),
    }
  case UPDATE_OWNED_DELEGATEN_TOKENS_BALANCE:
    return {
      ...state,
      ownedTokensDelegationsBalance: action.payload.operation(state.ownedTokensDelegationsBalance, action.payload.value),
    }
  case ADD_DELEGATION:
    return {
      ...state,
      delegations: [action.payload, ...state.delegations],
    }
  case REMOVE_DELEGATION:
    return {
      ...state,
      delegations: removeFromDelegationOrUndelegation(...state.delegations, action.payload),
    }
  case ADD_UNDELEGATION:
    return {
      ...state,
      undelegations: [action.payload, ...state.undelegations],
    }
  case REMOVE_UNDELEGATION:
    return {
      ...state,
      undelegations: removeFromDelegationOrUndelegation(...state.undelegations, action.payload),
    }
  default:
    return { ...state }
  }
}

const removeFromDelegationOrUndelegation = (array, id) => {
  const { indexInArray } = findIndexAndObject(id, 'operatorAddress', array, compareEthAddresses)
  if (indexInArray === null) {
    return array
  }
  array.splice(indexInArray, 1)

  return array
}

const updateGrants = (grants, { grantId, dataToUpdate }) => {
  const { indexInArray, obj } = findIndexAndObject(grantId, 'id', grants)
  if (indexInArray === null) {
    return grants
  }
  grants[indexInArray] = { ...obj, ...dataToUpdate }

  return grants
}
