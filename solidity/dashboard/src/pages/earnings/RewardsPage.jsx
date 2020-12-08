import React from "react"
import { Rewards } from "../../components/Rewards"
import EmptyStatePage from "./EmptyStatePage"

const RewardsPage = () => {
  return <Rewards />
}

RewardsPage.route = {
  title: "Keep Random Beacon",
  path: "/earnings/random-beacon",
  exact: true,
  withConnectWalletGuard: true,
  emptyStateComponent: EmptyStatePage,
}

export default RewardsPage
