import React, { useContext, useReducer, useCallback, useMemo, useEffect } from 'react'
import { Web3Context } from '../components/WithWeb3Context'
import { findIndexAndObject, compareEthAddresses } from '../utils/array.utils'
import { tokensPageService } from '../services/tokens-page.service'
import { tokenGrantsService } from '../services/token-grants.service'
import { useFetchData } from '../hooks/useFetchData'
import { add, sub, gte } from '../utils/arithmetics.utils'

export const REFRESH_KEEP_TOKEN_BALANCE = 'REFRESH_KEEP_TOKEN_BALANCE'
export const REFRESH_GRANT_TOKEN_BALANCE = 'REFRESH_GRANT_TOKEN_BALANCE'
export const UPDATE_OWNED_UNDELEGATIONS_TOKEN_BALANCE = 'UPDATE_OWNED_UNDELEGATIONS_BALANCE'
export const UPDATE_OWNED_DELEGATED_TOKENS_BALANCE = 'UPDATE_OWNED_DELEGATED_TOKENS_BALANCE'
export const ADD_DELEGATION = 'ADD_DELEGATION'
export const REMOVE_DELEGATION = 'REMOVE_DELEGATION'
export const ADD_UNDELEGATION = 'ADD_UNDELEGATION'
export const REMOVE_UNDELEGATION = 'REMOVE_UNDELEGATION'
export const GRANT_STAKED = 'GRANT_STAKED'
export const GRANT_WITHDRAWN = 'GRANT_WITHDRAWN'
const SET_STATE = 'SET_STATE'

const tokesnPageServiceInitialData = {
  delegations: [],
  undelegations: [],
  keepTokenBalance: '0',
  grantTokenBalance: '0',
  ownedTokensUndelegationsBalance: '0',
  ownedTokensDelegationsBalance: '0',
  initializationPeriod: '0',
  undelegationPeriod: '0',
  minimumStake: '0',
}

const TokensPageContext = React.createContext({
  refreshKeepTokenBalance: () => {},
  refreshGrantTokenBalance: () => {},
  dispatch: () => {},
  grants: [],
  ...tokesnPageServiceInitialData,
})

const TokenPageContextProvider = (props) => {
  const web3Context = useContext(Web3Context)
  const [{ data, isFetching: tokesnPageDataIsFetching }] = useFetchData(tokensPageService.fetchTokensPageData, tokesnPageServiceInitialData)
  const [{ data: grants, isFetching: grantsAreFetching }, , refreshGrants] = useFetchData(tokenGrantsService.fetchGrants, [])

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
    isFetching: true,
  })

  useEffect(() => {
    dispatch({ type: SET_STATE, payload: { ...data, isFetching: tokesnPageDataIsFetching } })
  }, [data, tokesnPageDataIsFetching])

  useEffect(() => {
    dispatch({ type: SET_STATE, payload: { grants, isFetching: grantsAreFetching } })
  }, [grants, grantsAreFetching])

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
      ...state,
      dispatch: contextValue.dispatch,
      refreshKeepTokenBalance,
      refreshGrantTokenBalance,
      refreshGrants,
    }}>
      {props.children}
    </TokensPageContext.Provider>
  )
}

export default TokenPageContextProvider

export const useTokensPageContext = () => {
  return useContext(TokensPageContext)
}

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
      ownedTokensUndelegationsBalance: action.payload.operation(state.ownedTokensUndelegationsBalance, action.payload.value),
    }
  case UPDATE_OWNED_DELEGATED_TOKENS_BALANCE:
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
      delegations: removeFromDelegationOrUndelegation([...state.delegations], action.payload),
    }
  case ADD_UNDELEGATION:
    return {
      ...state,
      undelegations: [action.payload, ...state.undelegations],
    }
  case REMOVE_UNDELEGATION:
    return {
      ...state,
      undelegations: removeFromDelegationOrUndelegation([...state.undelegations], action.payload),
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
  default:
    return { ...state }
  }
}

const removeFromDelegationOrUndelegation = (array, id) => {
  const { indexInArray } = findIndexAndObject('operatorAddress', id, array, compareEthAddresses)
  if (indexInArray === null) {
    return array
  }
  array.splice(indexInArray, 1)

  return array
}

const grantStaked = (grants, { grantId, amount }) => {
  const { indexInArray, obj: grantToUpdate } = findIndexAndObject('id', grantId, grants)
  if (indexInArray === null) {
    return grants
  }
  grantToUpdate.staked = add(grantToUpdate.staked, amount)
  grantToUpdate.readyToRelease = sub(grantToUpdate.readyToRelease, amount)
  grantToUpdate.readyToRelease = gte(grantToUpdate.readyToRelease, 0) ? grantToUpdate.readyToRelease : '0'
  grantToUpdate.availableToStake = sub(grantToUpdate.availableToStake, amount)
  grantToUpdate.availableToStake = gte(grantToUpdate.availableToStake, 0) ? grantToUpdate.availableToStake : '0'
  grants[indexInArray] = grantToUpdate

  return grants
}

const grantWithdrawn = (grants, { grantId, amount }) => {
  const { indexInArray, obj: grantToUpdate } = findIndexAndObject('id', grantId, grants)
  if (indexInArray === null) {
    return grants
  }
  grantToUpdate.readyToRelease = '0'
  grantToUpdate.released = add(grantToUpdate.released, amount)
  grantToUpdate.vested = add(grantToUpdate.released, grantToUpdate.staked)
  grantToUpdate.availableToStake = sub(grantToUpdate.availableToStake, amount)
  grantToUpdate.availableToStake = gte(grantToUpdate.availableToStake, 0) ? grantToUpdate.availableToStake : '0'
  grants[indexInArray] = grantToUpdate

  return grants
}
