import React, { useEffect, useCallback, useMemo } from "react"
import DelegatedTokensTable from "../components/DelegatedTokensTable"
import Undelegations from "../components/Undelegations"
import TokenOverview from "../components/TokenOverview"
import { LoadingOverlay } from "../components/Loadable"
import { add, sub } from "../utils/arithmetics.utils"
import { isEmptyArray } from "../utils/array.utils"
import DataTableSkeleton from "../components/skeletons/DataTableSkeleton"
import Tile from "../components/Tile"
import { Link } from "react-router-dom"
import PageWrapper from "../components/PageWrapper"
import { useSelector, useDispatch } from "react-redux"
import { useWeb3Context } from "../components/WithWeb3Context"
import DelegationPage from "./delegation"
import * as Icons from "../components/Icons"
import { useWeb3Address } from "../components/WithWeb3Context"
import * as CoveragePoolsComponents from "../components/coverage-pools"

const OverviewPage = (props) => {
  const { isConnected } = useWeb3Context()
  const address = useWeb3Address()
  const dispatch = useDispatch()

  useEffect(() => {
    if (isConnected) {
      dispatch({
        type: "staking/fetch_delegations_request",
        payload: { address },
      })
      dispatch({
        type: "token-grant/fetch_grants_request",
        payload: { address },
      })
    }
  }, [dispatch, isConnected, address])

  const keepToken = useSelector((state) => state.keepTokenBalance)
  const {
    delegations,
    undelegations,
    ownedTokensDelegationsBalance,
    ownedTokensUndelegationsBalance,
    isDelegationDataFetching,
    undelegationPeriod,
  } = useSelector((state) => state.staking)

  const { grants, isFetching: grantsAreFetching } = useSelector(
    (state) => state.tokenGrants
  )

  const cancelStakeSuccessCallback = useCallback(() => {
    // TODO
  }, [])

  const totalOwnedStakedBalance = useMemo(() => {
    return add(
      ownedTokensDelegationsBalance,
      ownedTokensUndelegationsBalance
    ).toString()
  }, [ownedTokensDelegationsBalance, ownedTokensUndelegationsBalance])

  const totalKeepTokenBalance = useMemo(() => {
    return add(totalOwnedStakedBalance, keepToken.value).toString()
  }, [keepToken.value, totalOwnedStakedBalance])

  const totalGrantedStakedBalance = useMemo(() => {
    return [...delegations, ...undelegations]
      .filter((delegation) => delegation.isFromGrant)
      .map(({ amount }) => amount)
      .reduce(add, "0")
      .toString()
  }, [delegations, undelegations])

  const totalGrantedTokenBalance = useMemo(() => {
    const grantedBalance = grants
      .map(({ amount, released }) => sub(amount, released))
      .reduce(add, "0")
      .toString()
    return grantedBalance
  }, [grants])

  return (
    <PageWrapper {...props} headerClassName="header--overview">
      <CoveragePoolsComponents.LearnMoreBanner />
      <OverviewFirstSection />
      <TokenOverview
        totalKeepTokenBalance={totalKeepTokenBalance}
        totalOwnedStakedBalance={totalOwnedStakedBalance}
        totalGrantedTokenBalance={totalGrantedTokenBalance}
        totalGrantedStakedBalance={totalGrantedStakedBalance}
        isFetching={
          keepToken.isFetching || grantsAreFetching || isDelegationDataFetching
        }
      />
      {isConnected && (
        <>
          <LoadingOverlay
            isFetching={isDelegationDataFetching}
            skeletonComponent={<DataTableSkeleton />}
          >
            <DelegatedTokensTable
              title="Delegation History"
              delegatedTokens={delegations}
              cancelStakeSuccessCallback={cancelStakeSuccessCallback}
              keepTokenBalance={keepToken.value}
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
      )}
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
          <Link
            to={DelegationPage.route.path}
            className="btn btn-primary btn-lg"
          >
            start staking
          </Link>
        </div>
        <div>
          <h4 className={"text-grey-40"}>
            Get KEEP tokens on the following exchanges:
          </h4>
          <Icons.BalancerLogo />
          &nbsp;
          <a
            target="_blank"
            rel="noopener noreferrer"
            href={
              "https://balancer.exchange/#/swap/ether/0x85eee30c52b0b379b046fb0f85f4f3dc3009afec"
            }
            className="text-black mr-2"
          >
            Balancer
          </a>
          &nbsp;
          <Icons.UniswapLogo style={{ verticalAlign: "text-top" }} />
          &nbsp;
          <a
            target="_blank"
            rel="noopener noreferrer"
            href={
              "https://app.uniswap.org/#/swap?inputCurrency=ETH&outputCurrency=0x85eee30c52b0b379b046fb0f85f4f3dc3009afec"
            }
            className="text-black"
          >
            Uniswap
          </a>
        </div>
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
