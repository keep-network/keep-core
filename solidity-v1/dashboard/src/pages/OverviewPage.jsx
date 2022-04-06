import React, { useEffect, useCallback } from "react"
import DelegatedTokensTable from "../components/DelegatedTokensTable"
import Undelegations from "../components/Undelegations"
import TokenOverview from "../components/TokenOverview"
import { LoadingOverlay } from "../components/Loadable"
import { isEmptyArray } from "../utils/array.utils"
import DataTableSkeleton from "../components/skeletons/DataTableSkeleton"
import { Link } from "react-router-dom"
import PageWrapper from "../components/PageWrapper"
import { useSelector, useDispatch } from "react-redux"
import { useWeb3Context } from "../components/WithWeb3Context"
import * as Icons from "../components/Icons"
import { useWeb3Address } from "../components/WithWeb3Context"
import OnlyIf from "../components/OnlyIf"
import PendingWithdrawals from "../components/coverage-pools/PendingWithdrawals"
import {
  fetchAPYRequest,
  fetchCovPoolDataRequest,
  fetchTvlRequest,
} from "../actions/coverage-pool"
import useKeepBalanceInfo from "../hooks/useKeepBalanceInfo"
import useGrantedBalanceInfo from "../hooks/useGrantedBalanceInfo"
import ThresholdUpgradePage from "./threshold/ThresholdUpgradePage"
import NavLink from "../components/NavLink"

const OverviewPage = (props) => {
  const { isConnected } = useWeb3Context()
  const address = useWeb3Address()
  const dispatch = useDispatch()
  const { covTokensAvailableToWithdraw, withdrawalInitiatedTimestamp } =
    useSelector((state) => state.coveragePool)

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
      dispatch(fetchCovPoolDataRequest(address))
      dispatch(fetchTvlRequest())
      dispatch(fetchAPYRequest())
    }
  }, [dispatch, isConnected, address])

  const keepToken = useSelector((state) => state.keepTokenBalance)
  const {
    delegations,
    undelegations,
    isDelegationDataFetching,
    undelegationPeriod,
  } = useSelector((state) => state.staking)

  const { grants, isFetching: grantsAreFetching } = useSelector(
    (state) => state.tokenGrants
  )

  const cancelStakeSuccessCallback = useCallback(() => {
    // TODO
  }, [])

  const { totalOwnedStakedBalance, totalKeepTokenBalance } =
    useKeepBalanceInfo()

  const { totalGrantedStakedBalance, totalGrantedTokenBalance } =
    useGrantedBalanceInfo()

  return (
    <PageWrapper {...props} headerClassName="header--overview">
      <OverviewFirstSection />
      <OnlyIf condition={withdrawalInitiatedTimestamp > 0}>
        <PendingWithdrawals
          covTokensAvailableToWithdraw={covTokensAvailableToWithdraw}
        />
      </OnlyIf>
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
    <section className="upgrade-threshold-banner tile flex row center">
      <Icons.KeepTUpgrade className="threshold-upgrade-page__explanation__icon" />
      <div>
        <h4 className={"h2 mb-2 text-grey-70"}>
          Upgrade your KEEP to T and get started with the Threshold Network.
        </h4>
        <div className="upgrade-threshold-banner__content">
          <div className="upgrade-threshold-banner__upgrade-btn">
            <Link
              to={ThresholdUpgradePage.route.path}
              className="btn btn-primary btn-lg"
            >
              upgrade now
            </Link>
          </div>
          <div>
            <h4 className={"text-grey-50"}>
              Keep and NuCypher merged to form Threshold Network.
            </h4>
            <span className={"flex row center"}>
              <Icons.TTokenSymbol width={25} height={25} />
              &nbsp;
              <NavLink
                to="/threshold/how-it-works"
                className="upgrade-threshold-banner__learn-more text-black"
              >
                Learn more
              </NavLink>
              &nbsp;
            </span>
          </div>
        </div>
      </div>
    </section>
  )
}

OverviewPage.route = {
  title: "Overview",
  path: "/overview",
  exact: true,
}

export default OverviewPage
