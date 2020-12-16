import React from "react"
import PageWrapper from "../../components/PageWrapper"
import RandomBeaconRewardsPage from "./RandomBeaconRewardsPage"
import TBTCRewardsPage from "./TBTCRewardsPage"

const RewardsPage = (props) => {
  return <PageWrapper {...props} />
}

RewardsPage.route = {
  title: "Rewards",
  path: "/rewards",
  pages: [RandomBeaconRewardsPage, TBTCRewardsPage],
}

export default RewardsPage
