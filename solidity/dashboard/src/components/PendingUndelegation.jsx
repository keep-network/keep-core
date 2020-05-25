import React, { useEffect, useContext } from "react"
import {
  isSameEthAddress,
  isEmptyObj,
  formatDate,
} from "../utils/general.utils"
import { displayAmount } from "../utils/token.utils"
import { Web3Context } from "./WithWeb3Context"
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"
import { DataTable, Column } from "./DataTable"
import moment from "moment"
import Tile from "./Tile"
import { usePrevious } from "../hooks/usePrevious"

const PendingUndelegation = ({ latestUnstakeEvent, data, setData }) => {
  const { yourAddress } = useContext(Web3Context)
  const { undelegationPeriod } = data
  const previousEvent = usePrevious(latestUnstakeEvent)

  useEffect(() => {
    if (isEmptyObj(latestUnstakeEvent)) {
      return
    } else if (
      previousEvent.transactionHash === latestUnstakeEvent.transactionHash
    ) {
      return
    }

    const {
      returnValues: { operator, undelegatedAt },
    } = latestUnstakeEvent
    if (!isSameEthAddress(yourAddress, operator)) {
      return
    }

    const undelegationCompletedAt = moment
      .unix(undelegatedAt)
      .add(undelegationPeriod, "seconds")
    setData({
      ...data,
      undelegationCompletedAt,
      delegatedStatus: "UNDELEGATED",
    })
  })

  return data.delegationStatus === "UNDELEGATED" ? (
    <Tile title="Pending Undelegations" id="pending-undelegation">
      <DataTable data={[{ ...data }]} itemFieldId="pendingUnstakeBalance">
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
