import React from "react"
import { Rewards } from "../components/Rewards"
import PageWrapper from "../components/PageWrapper"

const RewardsPage = () => {
  return (
    <PageWrapper title="My Rewards">
      <Rewards />
    </PageWrapper>
  )
}

export default React.memo(RewardsPage)
