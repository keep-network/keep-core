import React, { useCallback } from "react"
import * as Icons from "./Icons"
import Dropdown from "./Dropdown"
import SelectedGrantDropdown from "./SelectedGrantDropdown"
import { useTokensPageContext } from "../contexts/TokensPageContext"
import {
  SET_TOKENS_CONTEXT,
  SET_SELECTED_GRANT,
} from "../reducers/tokens-page.reducer"
import TokenAmount from "./TokenAmount"

const TokensContextSwitcher = (props) => {
  const {
    dispatch,
    selectedGrant,
    tokensContext,
    grants,
    keepTokenBalance,
  } = useTokensPageContext()

  const setTokensContext = useCallback(
    (contextName) => {
      dispatch({ type: SET_TOKENS_CONTEXT, payload: contextName })
    },
    [dispatch]
  )

  const onSelectGrant = useCallback(
    (grant) => {
      dispatch({ type: SET_SELECTED_GRANT, payload: grant })
    },
    [dispatch]
  )

  return (
    <div className="tokens-context-switcher-wrapper">
      <div
        className={`grants ${
          tokensContext === "granted" ? "active" : "inactive"
        }`}
        onClick={() => setTokensContext("granted")}
      >
        <div className="flex row">
          <Icons.GrantContextIcon />
          <div className="ml-1">
            <h2 className="text-grey-70">Grants</h2>
            <TokenAmount
              amount={selectedGrant.availableToStake}
              amountClassName="h4 text-primary"
              suffixClassName="h5"
            />
          </div>
        </div>
        <div className="grants-dropdown">
          <Dropdown
            onSelect={onSelectGrant}
            options={grants}
            valuePropertyName="id"
            labelPropertyName="id"
            selectedItem={selectedGrant}
            labelPrefix="Grant ID"
            noItemSelectedText="Select Grant"
            label="Choose Grant"
            selectedItemComponent={
              <SelectedGrantDropdown grant={selectedGrant} />
            }
          />
        </div>
      </div>
      <div
        className={`owned ${tokensContext === "owned" ? "active" : "inactive"}`}
        onClick={() => setTokensContext("owned")}
      >
        <Icons.MoneyWalletOpen />
        <h2 className="text-grey-70">Owned</h2>
        <TokenAmount
          amount={keepTokenBalance}
          amountClassName="h4 text-primary"
          suffixClassName="h5"
        />
      </div>
    </div>
  )
}

export default TokensContextSwitcher
