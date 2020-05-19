import React, { useEffect, useContext } from "react"
import { useFetchData } from "../hooks/useFetchData"
import { operatorService } from "../services/token-staking.service"
import {
  isSameEthAddress,
  isEmptyObj,
  formatDate,
} from "../utils/general.utils"
import { displayAmount } from "../utils/token.utils"
import { LoadingOverlay } from "./Loadable"
import { Web3Context } from "./WithWeb3Context"
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"
import { DataTable, Column } from "./DataTable"
import { PENDING_STATUS } from "../constants/constants"
import moment from "moment"
import Tile from "./Tile"
import { usePrevious } from "../hooks/usePrevious"

const initialData = { undelegation: {} }

const PendingUndelegation = ({ latestUnstakeEvent }) => {
  const { stakingContract, yourAddress } = useContext(Web3Context)
  const [state, setData] = useFetchData(
    operatorService.fetchPendingUndelegation,
    initialData
  )
  const { isFetching, data } = state
  const { undelegationCompletedAt, undelegationPeriod } = data
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
    stakingContract.methods
      .getDelegationInfo(operator)
      .call()
      .then((data) => {
        const { amount } = data
        setData({
          ...state.data,
          undelegationCompletedAt,
          pendingUnstakeBalance: amount,
          undelegationStatus: "PENDING",
        })
      })
  })

  return (
    <LoadingOverlay isFetching={isFetching}>
      <Tile title="Token Undelegation" id="pending-undelegation">
        <DataTable data={[{ ...data }]} itemFieldId="pendingUnstakeBalance">
          <Column
            header="amount"
            field="pendingUnstakeBalance"
            renderContent={({ pendingUnstakeBalance }) =>
              pendingUnstakeBalance && `${displayAmount(pendingUnstakeBalance)}`
            }
          />
          <Column
            header="status"
            field="undelegationStatus"
            renderContent={({ undelegationStatus }) => {
              const statusText =
                undelegationStatus === PENDING_STATUS
                  ? `${undelegationStatus.toLowerCase()}, ${undelegationCompletedAt.fromNow(
                      true
                    )}`
                  : undelegationStatus

              return (
                undelegationStatus && (
                  <StatusBadge
                    className="self-start"
                    status={BADGE_STATUS[undelegationStatus]}
                    text={statusText.toLowerCase()}
                  />
                )
              )
            }}
          />
          <Column
            header="estimate"
            field="undelegationCompletedAt"
            renderContent={({ undelegationCompletedAt }) =>
              undelegationCompletedAt
                ? formatDate(undelegationCompletedAt)
                : "-"
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
    </LoadingOverlay>
  )
}

export default PendingUndelegation
