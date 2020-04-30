import React from "react"
import { RewardsGroups } from "../components/RewardsGroups"
import { WithdrawalHistory } from "../components/WithdrawalHistory"
import { useSubscribeToContractEvent } from "../hooks/useSubscribeToContractEvent"
import { OPERATOR_CONTRACT_NAME } from "../constants/constants"
import PageWrapper from "../components/PageWrapper"

const RewardsPage = () => {
  const { latestEvent } = useSubscribeToContractEvent(
    OPERATOR_CONTRACT_NAME,
    "GroupMemberRewardsWithdrawn"
  )

  return (
    <PageWrapper title="My Rewards">
      <RewardsGroups latestWithdrawalEvent={latestEvent} />
      <WithdrawalHistory latestWithdrawalEvent={latestEvent} />
    </PageWrapper>
  )
}

export default React.memo(RewardsPage)
