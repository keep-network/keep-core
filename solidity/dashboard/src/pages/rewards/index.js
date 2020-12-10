import React from "react"
import PageWrapper from "../../components/PageWrapper"
// import RandomBeaconRewardsPage from "./RandomBeaconRewardsPage"
import TBTCRewardsPage from "./TBTCRewardsPage"
import RewardsOverviewPage from "./RewardsOverviewPage"

const RewardsPage = (props) => {
  return <PageWrapper {...props} />
}

RewardsPage.route = {
  title: "Rewards",
  path: "/rewards",
  pages: [RewardsOverviewPage, TBTCRewardsPage],
}

export default RewardsPage
