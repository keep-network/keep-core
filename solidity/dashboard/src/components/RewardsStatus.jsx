import React from "react"
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"
import { COMPLETE_STATUS, PENDING_STATUS } from "../constants/constants"
import { ViewInBlockExplorer } from "./ViewInBlockExplorer"

const RewardsStatus = ({ isStale, status, transactionHash }) => {
  if (status && status === PENDING_STATUS) {
    return <StatusBadge text="pending" status={BADGE_STATUS[PENDING_STATUS]} />
  } else if (status === "WITHDRAWN") {
    return (
      <>
        <StatusBadge text="WITHDRAWN" status={BADGE_STATUS.DISABLED} />
        <div>
          <ViewInBlockExplorer
            className="text-smaller text-grey-50 arrow-link grey"
            text="View transaction"
            type="tx"
            id={transactionHash}
          />
        </div>
      </>
    )
  } else if (isStale) {
    return (
      <StatusBadge text="available" status={BADGE_STATUS[COMPLETE_STATUS]} />
    )
  } else {
    return (
      <>
        <StatusBadge
          text="active"
          status={BADGE_STATUS[PENDING_STATUS]}
          bgClassName="bg-success-light"
        />
        <div className="text-smaller">Signing group still working.</div>
      </>
    )
  }
}

export default React.memo(RewardsStatus)
