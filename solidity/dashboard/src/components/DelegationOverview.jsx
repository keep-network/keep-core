import React, { useMemo, useCallback } from "react"
import Undelegations from "../components/Undelegations"
import DelegatedTokensTable from "../components/DelegatedTokensTable"
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"
import { useTokensPageContext } from "../contexts/TokensPageContext"
import { formatDate } from "../utils/general.utils"
import moment from "moment"

const filterByOwned = (delegation) => !delegation.grantId
const filterBySelectedGrant = (delegation, selectedGrant) =>
  selectedGrant.id && delegation.grantId === selectedGrant.id

const DelegationOverview = () => {
  const {
    undelegations,
    delegations,
    refreshData,
    refreshGrants,
    tokensContext,
    selectedGrant,
  } = useTokensPageContext()

  const ownedDelegations = useMemo(() => {
    return delegations.filter(filterByOwned)
  }, [delegations])

  const ownedUndelegations = useMemo(() => {
    return undelegations.filter(filterByOwned)
  }, [undelegations])

  const grantDelegations = useMemo(() => {
    return delegations.filter((delegation) =>
      filterBySelectedGrant(delegation, selectedGrant)
    )
  }, [delegations, selectedGrant])

  const grantUndelegations = useMemo(() => {
    return undelegations.filter((undelegation) =>
      filterBySelectedGrant(undelegation, selectedGrant)
    )
  }, [undelegations, selectedGrant])

  const getDelegations = () => {
    if (tokensContext === "granted") {
      return grantDelegations
    }
    return ownedDelegations
  }

  const getUndelegations = () => {
    if (tokensContext === "granted") {
      return grantUndelegations
    }

    return ownedUndelegations
  }

  const cancelStakeSuccessCallback = useCallback(() => {
    refreshGrants()
    refreshData()
  }, [refreshGrants, refreshData])

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
      <DelegatedTokensTable
        delegatedTokens={getDelegations()}
        cancelStakeSuccessCallback={cancelStakeSuccessCallback}
      />
      <Undelegations undelegations={getUndelegations()} />
    </section>
  )
}

export default DelegationOverview
