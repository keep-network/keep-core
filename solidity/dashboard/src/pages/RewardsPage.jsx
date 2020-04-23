import React from "react"
import { RewardsGroups } from "../components/RewardsGroups"
import { WithdrawalHistory } from "../components/WithdrawalHistory"
import { useSubscribeToContractEvent } from "../hooks/useSubscribeToContractEvent"
import { OPERATOR_CONTRACT_NAME } from "../constants/constants"

const RewardsPage = () => {
  const { latestEvent } = useSubscribeToContractEvent(
    OPERATOR_CONTRACT_NAME,
    "GroupMemberRewardsWithdrawn"
  )

  return (
    <>
      <h2 className="mb-2">My Rewards</h2>
      <RewardsGroups latestWithdrawalEvent={latestEvent} />
      <WithdrawalHistory latestWithdrawalEvent={latestEvent} />
    </>
  )
}

export default React.memo(RewardsPage)
