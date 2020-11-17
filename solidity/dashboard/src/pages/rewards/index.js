import React from "react"
import PageWrapper from "../../components/PageWrapper"
import RandomBeaconRewardsPage from "./RandomBeaconRewardsPage"
import RewardsOverviewPage from "./RewardsOverviewPage"

const RewardsPage = (props) => {
  return <PageWrapper {...props} />
}

RewardsPage.route = {
  title: "Rewards",
  path: "/rewards",
  pages: [RewardsOverviewPage, RandomBeaconRewardsPage],
}

export default RewardsPage
