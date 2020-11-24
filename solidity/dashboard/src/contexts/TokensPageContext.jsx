import React, {
  useContext,
  useReducer,
  useCallback,
  useMemo,
  useEffect,
} from "react"
import tokensPageReducer, {
  SET_SELECTED_GRANT,
} from "../reducers/tokens-page.reducer"
import { isEmptyObj } from "../utils/general.utils"
import { add } from "../utils/arithmetics.utils"
import { usePrevious } from "../hooks/usePrevious"
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

  return (
    <TokensPageContext.Provider
      value={{
        ...state,
        dispatch: contextValue.dispatch,
        refreshGrants,
        refreshData,
        getGrantStakedAmount,
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
