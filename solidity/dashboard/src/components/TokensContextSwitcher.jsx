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
import Skeleton from "./skeletons/Skeleton"
import TokenAmountSkeleton from "./skeletons/TokenAmountSkeleton"

const TokensContextSwitcher = (props) => {
  const {
    dispatch,
    selectedGrant,
    tokensContext,
    grants,
    keepTokenBalance,
    isFetching,
    grantsAreFetching,
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
        {grantsAreFetching ? (
          <GrantedTokensLoadingComponent />
        ) : (
          <>
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
          </>
        )}
      </div>
      <div
        className={`owned ${tokensContext === "owned" ? "active" : "inactive"}`}
        onClick={() => setTokensContext("owned")}
      >
        {isFetching ? (
          <OwnedTokensLoadingComponent />
        ) : (
          <>
            <Icons.MoneyWalletOpen />
            <h2 className="text-grey-70">Owned</h2>

            <TokenAmount
              amount={keepTokenBalance}
              amountClassName="h4 text-primary"
              suffixClassName="h5"
            />
          </>
        )}
      </div>
    </div>
  )
}

const OwnedTokensLoadingComponent = () => {
  return (
    <>
      <Skeleton className="h2 ml-1" styles={{ width: "35%" }} />
      <TokenAmountSkeleton
        wrapperStyles={{ width: "35%", marginLeft: "auto" }}
        textStyles={{ width: "35%" }}
      />
    </>
  )
}

const GrantedTokensLoadingComponent = () => {
  return (
    <>
      <div className="flex column" style={{ width: "40%" }}>
        <Skeleton className="h2" />
        <TokenAmountSkeleton
          wrapperStyles={{ width: "100%", marginTop: "1rem" }}
          textStyles={{ width: "35%" }}
        />
      </div>
      <div className="grants-dropdown">
        <Skeleton styles={{ padding: "1.5rem 4rem" }} />
      </div>
    </>
  )
}

export default TokensContextSwitcher
