import React, {
  useContext,
  useReducer,
  useCallback,
  useMemo,
  useEffect,
} from "react"
import { Web3Context } from "../components/WithWeb3Context"
import { tokensPageService } from "../services/tokens-page.service"
import { tokenGrantsService } from "../services/token-grants.service"
import { useFetchData } from "../hooks/useFetchData"
import tokensPageReducer, {
  REFRESH_KEEP_TOKEN_BALANCE,
  REFRESH_GRANT_TOKEN_BALANCE,
  SET_STATE,
  GRANT_STAKED,
  GRANT_WITHDRAWN,
  SET_SELECTED_GRANT,
  GRANT_DEPOSITED,
} from "../reducers/tokens-page.reducer"
import { isEmptyObj } from "../utils/general.utils"
import { findIndexAndObject } from "../utils/array.utils"
import { add } from "../utils/arithmetics.utils"
import { usePrevious } from "../hooks/usePrevious"
import { fetchAvailableTopUps } from "../services/top-ups.service"
import { contracts } from "../contracts"

const tokensPageServiceInitialData = {
  delegations: [],
  undelegations: [],
  keepTokenBalance: "0",
  grantTokenBalance: "0",
  ownedTokensUndelegationsBalance: "0",
  ownedTokensDelegationsBalance: "0",
  initializationPeriod: "0",
  undelegationPeriod: "0",
  minimumStake: "0",
}

const TokensPageContext = React.createContext({
  refreshKeepTokenBalance: () => {},
  refreshGrantTokenBalance: () => {},
  dispatch: () => {},
  grants: [],
  ...tokensPageServiceInitialData,
})

const TokenPageContextProvider = (props) => {
  const web3Context = useContext(Web3Context)
  const [
    { data, isFetching: tokesnPageDataIsFetching },
    ,
    refreshData,
  ] = useFetchData(
    tokensPageService.fetchTokensPageData,
    tokensPageServiceInitialData
  )
  const [
    { data: grants, isFetching: grantsAreFetching },
    ,
    refreshGrants,
  ] = useFetchData(tokenGrantsService.fetchGrants, [])

  const [
    { data: availableTopUps, isFetching: topUpsAreFetching },
    ,
    ,
    setOperatorsForTopUp,
  ] = useFetchData(fetchAvailableTopUps, [])

  useEffect(() => {
    setOperatorsForTopUp([
      [...data.undelegations, ...data.delegations].map(
        ({ operatorAddress }) => operatorAddress
      ),
    ])
  }, [setOperatorsForTopUp, data.delegations, data.undelegations])

  const [state, dispatch] = useReducer(tokensPageReducer, {
    grants: [],
    delegations: [],
    undelegations: [],
    keepTokenBalance: "0",
    grantTokenBalance: "0",
    ownedTokensUndelegationsBalance: "0",
    ownedTokensDelegationsBalance: "0",
    initializationPeriod: "0",
    undelegationPeriod: "0",
    isFetching: true,
    grantsAreFetching: true,
    tokensContext: "granted",
    selectedGrant: {},
    availableTopUps: [],
    topUpsAreFetching: false,
    getGrantStakedAmount: () => {},
  })
  const previousSelectedGrant = usePrevious(state.selectedGrant)

  useEffect(() => {
    dispatch({
      type: SET_STATE,
      payload: { ...data, isFetching: tokesnPageDataIsFetching },
    })
  }, [data, tokesnPageDataIsFetching])

  useEffect(() => {
    dispatch({
      type: SET_STATE,
      payload: { grants, grantsAreFetching },
    })
    if (!isEmptyObj(state.selectedGrant)) {
      const { obj: updatedGrant } = findIndexAndObject(
        "id",
        state.selectedGrant.id,
        grants
      )
      dispatch({ type: SET_SELECTED_GRANT, payload: updatedGrant })
    }
  }, [grants, grantsAreFetching, state.selectedGrant])

  useEffect(() => {
    dispatch({
      type: SET_STATE,
      payload: { availableTopUps, topUpsAreFetching },
    })
  }, [availableTopUps, topUpsAreFetching])

  useEffect(() => {
    if (isEmptyObj(previousSelectedGrant) && state.grants.length > 0) {
      dispatch({ type: SET_SELECTED_GRANT, payload: state.grants[0] })
    }
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

    const grantTokenBalance = await grantContract.methods
      .balanceOf(yourAddress)
      .call()
    dispatch({ type: REFRESH_GRANT_TOKEN_BALANCE, payload: grantTokenBalance })
  }, [web3Context, dispatch])

  const grantStaked = useCallback(
    async (grantId, amount) => {
      const { grantContract } = web3Context

      const availableToStake = await grantContract.methods
        .availableToStake(grantId)
        .call()
      dispatch({
        type: GRANT_STAKED,
        payload: { grantId, amount, availableToStake },
      })
    },
    [web3Context, dispatch]
  )

  const grantWithdrawn = useCallback(
    async (grantId, amount, operator) => {
      const { grantContract } = web3Context

      const availableToStake = await grantContract.methods
        .availableToStake(grantId)
        .call()
      dispatch({
        type: GRANT_WITHDRAWN,
        payload: { grantId, amount, availableToStake, operator },
      })
    },
    [web3Context, dispatch]
  )

  const getGrantStakedAmount = useCallback(
    (grantId) => {
      if (!grantId) return 0

      return [...state.delegations, ...state.undelegations]
        .filter((delegation) => delegation.grantId === grantId)
        .map((grantDelegation) => grantDelegation.amount)
        .reduce(add, 0)
    },
    [state.delegations, state.undelegations]
  )

  const grantDeposited = useCallback(async (grantId, operator, amount) => {
    const availableToWithdrawEscrow = await contracts.tokenStakingEscrow.methods
      .withdrawable(operator)
      .call()
    const availableToWitdrawGrant = await contracts.grantContract.methods
      .withdrawable(grantId)
      .call()

    dispatch({
      type: GRANT_DEPOSITED,
      payload: {
        grantId,
        availableToWithdrawEscrow,
        availableToWitdrawGrant,
        amount,
        operator,
      },
    })
  }, [])

  return (
    <TokensPageContext.Provider
      value={{
        ...state,
        dispatch: contextValue.dispatch,
        refreshKeepTokenBalance,
        refreshGrantTokenBalance,
        refreshGrants,
        refreshData,
        grantWithdrawn,
        grantStaked,
        getGrantStakedAmount,
        grantDeposited,
      }}
    >
      {props.children}
    </TokensPageContext.Provider>
  )
}

export default TokenPageContextProvider

export const useTokensPageContext = () => {
  return useContext(TokensPageContext)
}
