import React from "react"
import { Rewards } from "../components/Rewards"
import PageWrapper from "../components/PageWrapper"

const RewardsPage = () => {
  return (
    <PageWrapper title="Random Beacon Rewards">
      <Rewards />
    </PageWrapper>
  )
}

export default React.memo(RewardsPage)
