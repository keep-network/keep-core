import React, { useMemo, useCallback } from "react"
import Undelegations from "../components/Undelegations"
import DelegatedTokensTable from "../components/DelegatedTokensTable"
import { formatDate, isSameEthAddress } from "../utils/general.utils"
import moment from "moment"
import { LoadingOverlay } from "./Loadable"
import DataTableSkeleton from "./skeletons/DataTableSkeleton"
import TopUpsDataTable from "./TopUpsDataTable"
import Tile from "./Tile"
import Tag from "./Tag"
import * as Icons from "./Icons"

const DelegationOverview = ({
  delegationsWithTAuthData,
  undelegations,
  isFetching,
  topUps: availableTopUps,
  areTopUpsFetching,
  undelegationPeriod,
  initializationPeriod,
  keepTokenBalance,
  grants = [],
  selectedGrant = null,
  context = "wallet",
}) => {
  const cancelStakeSuccessCallback = useCallback(() => {
    // TODO
  }, [])

  const filteredTopUps = useMemo(() => {
    const topUps = []
    for (const topUp of availableTopUps) {
      const { operatorAddress: lookupOperator } = topUp
      const isUndelegation = undelegations.some(({ operatorAddress }) =>
        isSameEthAddress(lookupOperator, operatorAddress)
      )

      const isDelegation = delegationsWithTAuthData.some(
        ({ operatorAddress }) =>
          isSameEthAddress(lookupOperator, operatorAddress)
      )

      if (isDelegation || isUndelegation) {
        topUp.isInUndelegation = isUndelegation
        topUps.push(topUp)
      }
    }
    return topUps
  }, [availableTopUps, delegationsWithTAuthData, undelegations])

  return (
    <section>
      <DelegationHeader type={context} selectedGrant={selectedGrant} />
      <LoadingOverlay
        isFetching={areTopUpsFetching}
        skeletonComponent={<DataTableSkeleton columns={3} />}
      >
        <Tile>
          <TopUpsDataTable
            topUps={filteredTopUps}
            initializationPeriod={initializationPeriod}
          />
        </Tile>
      </LoadingOverlay>
      <LoadingOverlay
        isFetching={isFetching}
        skeletonComponent={<DataTableSkeleton />}
      >
        <DelegatedTokensTable
          delegationsWithTAuthData={delegationsWithTAuthData}
          cancelStakeSuccessCallback={cancelStakeSuccessCallback}
          keepTokenBalance={keepTokenBalance}
          grants={grants}
          undelegationPeriod={undelegationPeriod}
        />
      </LoadingOverlay>
      <LoadingOverlay
        isFetching={isFetching}
        skeletonComponent={<DataTableSkeleton />}
      >
        <Undelegations undelegations={undelegations} />
      </LoadingOverlay>
    </section>
  )
}

const DelegationHeader = ({ type, selectedGrant }) => {
  switch (type) {
    case "wallet":
    default:
      return <h2 className="h2--alt text-grey-60 mb-2">Activity</h2>
    case "granted":
      return (
        <header className="flex row center mb-2">
          <h2 className="h2--alt text-grey-60">Grant Activity</h2>
          <div className="flex row center ml-a">
            <Tag
              IconComponent={Icons.Grant}
              text="Grant ID"
              className="grant-id-tag"
            />
            <span className="ml-1 mr-2">
              {selectedGrant && selectedGrant.id}
            </span>
            <Tag IconComponent={Icons.Time} text="Issued" />
            <span className="ml-1">
              {selectedGrant && selectedGrant.start
                ? formatDate(moment.unix(selectedGrant.start))
                : null}
            </span>
          </div>
        </header>
      )
  }
}

export default DelegationOverview
