import React from "react"
import { formatDate } from "../utils/general.utils"
import { displayAmount } from "../utils/token.utils"
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"
import { DataTable, Column } from "./DataTable"
import moment from "moment"
import Tile from "./Tile"

const PendingUndelegation = ({ data }) => {
  return data.delegationStatus === "UNDELEGATED" ? (
    <Tile id="pending-undelegation">
      <DataTable
        data={[{ ...data }]}
        itemFieldId="pendingUnstakeBalance"
        title="Pending Undelegations"
        noDataMessage="No undelegated tokens."
      >
        <Column
          header="amount"
          field="stakedBalance"
          renderContent={({ stakedBalance }) =>
            stakedBalance && `${displayAmount(stakedBalance)}`
          }
        />
        <Column
          header="status"
          field="delegationStatus"
          renderContent={() => {
            return (
              <StatusBadge status={BADGE_STATUS.PENDING} text="processing" />
            )
          }}
        />
        <Column
          header="estimate"
          field="undelegationCompletedAt"
          renderContent={({ undelegationCompletedAt }) =>
            undelegationCompletedAt ? formatDate(undelegationCompletedAt) : "-"
          }
        />
        <Column
          header="undelegation period"
          field="undelegationPeriod"
          renderContent={({ undelegationPeriod }) => {
            const undelegationPeriodRelativeTime = moment()
              .add(undelegationPeriod, "seconds")
              .fromNow(true)
            return undelegationPeriodRelativeTime
          }}
        />
      </DataTable>
    </Tile>
  ) : null
}

export default PendingUndelegation
