import React, {
  useContext,
  useReducer,
  useCallback,
  useMemo,
  useEffect,
} from "react"
import { Web3Context } from "../components/WithWeb3Context"
import tokensPageReducer, {
  REFRESH_KEEP_TOKEN_BALANCE,
  REFRESH_GRANT_TOKEN_BALANCE,
  GRANT_STAKED,
  GRANT_WITHDRAWN,
  SET_SELECTED_GRANT,
  GRANT_DEPOSITED,
} from "../reducers/tokens-page.reducer"
import { isEmptyObj } from "../utils/general.utils"
import { add } from "../utils/arithmetics.utils"
import { usePrevious } from "../hooks/usePrevious"
import { ContractsLoaded } from "../contracts"
import { useSelector } from "react-redux"

const TokensPageContext = React.createContext({
  refreshKeepTokenBalance: () => {},
  refreshGrantTokenBalance: () => {},
  dispatch: () => {},
  tokensContext: "granted",
  selectedGrant: {},
  getGrantStakedAmount: () => {},
})

const TokenPageContextProvider = (props) => {
  const web3Context = useContext(Web3Context)

  const [state, dispatch] = useReducer(tokensPageReducer, {
    tokensContext: "granted",
    selectedGrant: {},
  })

  const previousSelectedGrant = usePrevious(state.selectedGrant)
  const { grants } = useSelector((state) => state.tokenGrants)
  const { delegations, undelegations } = useSelector((state) => state.staking)

  useEffect(() => {
    if (isEmptyObj(previousSelectedGrant) && grants.length > 0) {
      dispatch({ type: SET_SELECTED_GRANT, payload: grants[0] })
    }
  })

  const refreshGrants = () => {}
  const refreshData = () => {}

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
      dispatch({
        type: GRANT_STAKED,
        payload: {
          grantId,
          amount,
        },
      })
    },
    [dispatch]
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

      return [...delegations, ...undelegations]
        .filter((delegation) => delegation.grantId === grantId)
        .map((grantDelegation) => grantDelegation.amount)
        .reduce(add, 0)
    },
    [delegations, undelegations]
  )

  const grantDeposited = useCallback(async (grantId, operator, amount) => {
    const { grantContract } = await ContractsLoaded
    const availableToWitdrawGrant = await grantContract.methods
      .withdrawable(grantId)
      .call()

    dispatch({
      type: GRANT_DEPOSITED,
      payload: {
        grantId,
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
