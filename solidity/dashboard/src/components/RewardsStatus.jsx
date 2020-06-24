import React from "react"
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"
import {
  COMPLETE_STATUS,
  PENDING_STATUS,
  REWARD_STATUS,
} from "../constants/constants"
import { ViewInBlockExplorer } from "./ViewInBlockExplorer"

const RewardsStatus = ({ status, transactionHash }) => {
  switch (status) {
    case PENDING_STATUS:
      return <StatusBadge text={status} status={BADGE_STATUS[PENDING_STATUS]} />
    case REWARD_STATUS.AVAILABLE:
      return (
        <StatusBadge text={status} status={BADGE_STATUS[COMPLETE_STATUS]} />
      )
    case REWARD_STATUS.ACCUMULATING:
      return <StatusBadge text={status} status={BADGE_STATUS.ACTIVE} />

    case REWARD_STATUS.WITHDRAWN:
      return (
        <>
          <StatusBadge text={status} status={BADGE_STATUS.DISABLED} />
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
    default:
      return null
  }
}

export default React.memo(RewardsStatus)
