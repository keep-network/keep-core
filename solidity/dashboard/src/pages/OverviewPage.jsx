import React, { useCallback, useMemo } from "react"
import DelegatedTokensTable from "../components/DelegatedTokensTable"
import Undelegations from "../components/Undelegations"
import TokenOverview from "../components/TokenOverview"
import { LoadingOverlay } from "../components/Loadable"
import { useTokensPageContext } from "../contexts/TokensPageContext"
import { add, sub } from "../utils/arithmetics.utils"
import { isEmptyArray } from "../utils/array.utils"
import DataTableSkeleton from "../components/skeletons/DataTableSkeleton"
import Tile from "../components/Tile"
import Button from "../components/Button"
import { Link } from "react-router-dom"
import PageWrapper from "../components/PageWrapper"
import TokenAmount from "../components/TokenAmount"
import Divider from "../components/Divider"
import ProgressBar from "../components/ProgressBar"
import Chip from "../components/Chip"
import * as Icons from "../components/Icons"
import { SpeechBubbleTooltip } from "../components/SpeechBubbleTooltip"
import TokenBalancesOverviewSkeleton from "../components/skeletons/TokenBalancesOverviewSkeleton"

const OverviewPage = (props) => {
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
    <PageWrapper {...props} headerClassName="header--overview">
      <OverviewFirstSection />
      <TokenOverview />
      {/* <LoadingOverlay
        isFetching={isFetching}
        skeletonComponent={<TokenBalancesOverviewSkeleton />}
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
      )} */}
    </PageWrapper>
  )
}

const OverviewFirstSection = () => {
  return (
    <Tile
      title="Make the most of your KEEP tokens by staking them and earning rewards with the token dashboard."
      titleClassName="h2 mb-2"
    >
      <div className="start-staking">
        <div className="start-staking__btn">
          <Link to="/delegate" className="btn btn-primary btn-lg">
            start staking
          </Link>
        </div>
        <a
          href="https://discordapp.com/invite/wYezN7v"
          className="arrow-link"
          rel="noopener noreferrer"
          target="_blank"
        >
          No KEEP tokens yet? Join our Discord
        </a>
      </div>
    </Tile>
  )
}

OverviewPage.route = {
  title: "Overview",
  path: "/overview",
  exact: true,
}

export default OverviewPage
