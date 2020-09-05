import React, { useMemo, useCallback } from "react"
import Undelegations from "../components/Undelegations"
import DelegatedTokensTable from "../components/DelegatedTokensTable"
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"
import { useTokensPageContext } from "../contexts/TokensPageContext"
import { formatDate, isSameEthAddress } from "../utils/general.utils"
import moment from "moment"
import { LoadingOverlay } from "./Loadable"
import DataTableSkeleton from "./skeletons/DataTableSkeleton"
import TopUpsDataTable from "./TopUpsDataTable"
import Tile from "./Tile"

const filterByOwned = (delegation) => !delegation.grantId
const filterBySelectedGrant = (selectedGrant) => (delegation) =>
  selectedGrant.id && delegation.grantId === selectedGrant.id

const DelegationOverview = () => {
  const {
    undelegations,
    delegations,
    refreshData,
    tokensContext,
    selectedGrant,
    isFetching,
    grantsAreFetching,
    keepTokenBalance,
    availableTopUps,
    topUpsAreFetching,
    grants,
    initializationPeriod,
    undelegationPeriod,
  } = useTokensPageContext()

  const ownedDelegations = useMemo(() => {
    return delegations.filter(filterByOwned)
  }, [delegations])

  const ownedUndelegations = useMemo(() => {
    return undelegations.filter(filterByOwned)
  }, [undelegations])

  const grantDelegations = useMemo(() => {
    return delegations.filter(filterBySelectedGrant(selectedGrant))
  }, [delegations, selectedGrant])

  const grantUndelegations = useMemo(() => {
    return undelegations.filter(filterBySelectedGrant(selectedGrant))
  }, [undelegations, selectedGrant])

  const getDelegations = useCallback(() => {
    if (tokensContext === "granted") {
      return grantDelegations
    }
    return ownedDelegations
  }, [tokensContext, grantDelegations, ownedDelegations])

  const getUndelegations = useCallback(() => {
    if (tokensContext === "granted") {
      return grantUndelegations
    }

    return ownedUndelegations
  }, [tokensContext, grantUndelegations, ownedUndelegations])

  const cancelStakeSuccessCallback = useCallback(() => {
    refreshData()
  }, [refreshData])

  const filteredTopUps = useMemo(() => {
    const topUps = []
    for (const topUp of availableTopUps) {
      const { operatorAddress: lookupOperator } = topUp
      const isUndelegation = getUndelegations().some(({ operatorAddress }) =>
        isSameEthAddress(lookupOperator, operatorAddress)
      )

      const isDelegation = getDelegations().some(({ operatorAddress }) =>
        isSameEthAddress(lookupOperator, operatorAddress)
      )

      if (isDelegation || isUndelegation) {
        topUp.isInUndelegation = isUndelegation
        topUps.push(topUp)
      }
    }
    return topUps
  }, [availableTopUps, getDelegations, getUndelegations])

  return (
    <section>
      <div className="flex wrap self-center mt-3 mb-2">
        <h2 className="text-grey-60">
          {`${tokensContext === "granted" ? "Grant " : ""}Delegation Overview`}
        </h2>
        {tokensContext === "granted" && (
          <>
            <span className="flex self-center ml-2">
              <StatusBadge
                className="self-center"
                status={BADGE_STATUS.DISABLED}
                text="grant id"
              />
              <span className="self-center h4 text-grey-50 ml-1">
                {selectedGrant.id}
              </span>
            </span>
            <span className="flex self-center ml-2">
              <StatusBadge
                className="self-center"
                status={BADGE_STATUS.DISABLED}
                text="issued"
              />
              <span className="h4 text-grey-50 ml-1">
                {selectedGrant.start &&
                  formatDate(moment.unix(selectedGrant.start))}
              </span>
            </span>
          </>
        )}
      </div>
      <LoadingOverlay
        isFetching={
          tokensContext === "granted" ? grantsAreFetching : isFetching
        }
        skeletonComponent={<DataTableSkeleton />}
      >
        <DelegatedTokensTable
          delegatedTokens={getDelegations()}
          cancelStakeSuccessCallback={cancelStakeSuccessCallback}
          keepTokenBalance={keepTokenBalance}
          grants={grants}
          undelegationPeriod={undelegationPeriod}
        />
      </LoadingOverlay>
      <LoadingOverlay
        isFetching={
          tokensContext === "granted" ? grantsAreFetching : isFetching
        }
        skeletonComponent={<DataTableSkeleton />}
      >
        <Undelegations undelegations={getUndelegations()} />
      </LoadingOverlay>
      <LoadingOverlay
        isFetching={topUpsAreFetching}
        skeletonComponent={<DataTableSkeleton columns={3} />}
      >
        <Tile>
          <TopUpsDataTable
            topUps={filteredTopUps}
            initializationPeriod={initializationPeriod}
          />
        </Tile>
      </LoadingOverlay>
    </section>
  )
}

export default DelegationOverview
