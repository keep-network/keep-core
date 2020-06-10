import React, { useCallback, useMemo } from "react"
import PageWrapper from "../components/PageWrapper"
import DelegatedTokensTable from "../components/DelegatedTokensTable"
import Undelegations from "../components/Undelegations"
import TokenOverview from "../components/TokenOverview"
import { LoadingOverlay } from "../components/Loadable"
import { useTokensPageContext } from "../contexts/TokensPageContext"
import { add } from "../utils/arithmetics.utils"
import { isEmptyArray } from "../utils/array.utils"

const TokenOverviewPage = () => {
  const {
    delegations,
    undelegations,
    keepTokenBalance,
    ownedTokensDelegationsBalance,
    ownedTokensUndelegationsBalance,
    grants,
    refreshGrants,
    refreshData,
    isFetching,
  } = useTokensPageContext()
  const cancelStakeSuccessCallback = useCallback(() => {
    refreshGrants()
    refreshData()
  }, [refreshGrants, refreshData])

  const totalOwnedStakedBalance = useMemo(() => {
    return add(ownedTokensDelegationsBalance, ownedTokensUndelegationsBalance)
  }, [ownedTokensDelegationsBalance, ownedTokensUndelegationsBalance])

  const totalKeepTokenBalance = useMemo(() => {
    return add(totalOwnedStakedBalance, keepTokenBalance)
  }, [keepTokenBalance, totalOwnedStakedBalance])

  const totalGrantedStakedBalance = useMemo(() => {
    return [...delegations, ...undelegations]
      .filter((delegation) => delegation.isFromGrant)
      .map(({ amount }) => amount)
      .reduce(add, "0")
  }, [delegations, undelegations])

  const totalGrantedTokenBalance = useMemo(() => {
    const grantedBalance = grants.map(({ amount }) => amount).reduce(add, "0")
    return add(grantedBalance, totalGrantedStakedBalance)
  }, [grants, totalGrantedStakedBalance])

  return (
    <LoadingOverlay isFetching={isFetching}>
      <PageWrapper title="Token Overview">
        <TokenOverview
          totalKeepTokenBalance={totalKeepTokenBalance}
          totalOwnedStakedBalance={totalOwnedStakedBalance}
          totalGrantedTokenBalance={totalGrantedTokenBalance}
          totalGrantedStakedBalance={totalGrantedStakedBalance}
        />
        <DelegatedTokensTable
          title="Delegation History"
          delegatedTokens={delegations}
          cancelStakeSuccessCallback={cancelStakeSuccessCallback}
        />
        {!isEmptyArray(undelegations) && (
          <Undelegations
            title="Undelegation History"
            undelegations={undelegations}
          />
        )}
      </PageWrapper>
    </LoadingOverlay>
  )
}

export default TokenOverviewPage
