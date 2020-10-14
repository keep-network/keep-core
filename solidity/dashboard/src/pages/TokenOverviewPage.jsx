import React, { useCallback, useMemo } from "react"
import DelegatedTokensTable from "../components/DelegatedTokensTable"
import Undelegations from "../components/Undelegations"
import TokenOverview from "../components/TokenOverview"
import { LoadingOverlay } from "../components/Loadable"
import { useTokensPageContext } from "../contexts/TokensPageContext"
import { add, sub } from "../utils/arithmetics.utils"
import { isEmptyArray } from "../utils/array.utils"
import DataTableSkeleton from "../components/skeletons/DataTableSkeleton"
import TokenBalancesOverviewSkeleton from "../components/skeletons/TokenBalancesOverviewSkeleton"

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
    undelegationPeriod,
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
    const grantedBalance = grants
      .map(({ amount, released }) => sub(amount, released))
      .reduce(add, "0")
    return grantedBalance
  }, [grants])

  return (
    <>
      <LoadingOverlay
        isFetching={isFetching}
        skeletonComponent={<TokenBalancesOverviewSkeleton />}
      >
        <TokenOverview
          totalKeepTokenBalance={totalKeepTokenBalance}
          totalOwnedStakedBalance={totalOwnedStakedBalance}
          totalGrantedTokenBalance={totalGrantedTokenBalance}
          totalGrantedStakedBalance={totalGrantedStakedBalance}
        />
      </LoadingOverlay>
      <LoadingOverlay
        isFetching={isFetching}
        skeletonComponent={<DataTableSkeleton subtitleWidth={0} />}
      >
        <DelegatedTokensTable
          title="Delegation History"
          delegatedTokens={delegations}
          cancelStakeSuccessCallback={cancelStakeSuccessCallback}
          keepTokenBalance={keepTokenBalance}
          grants={grants}
          undelegationPeriod={undelegationPeriod}
        />
      </LoadingOverlay>
      {!isEmptyArray(undelegations) && (
        <Undelegations
          title="Undelegation History"
          undelegations={undelegations}
        />
      )}
    </>
  )
}

TokenOverviewPage.route = {
  title: "Overview",
  path: "/tokens/overview",
  exact: true,
}

export default TokenOverviewPage
