import React from "react"
import PageWrapper from "../../components/PageWrapper"
import RandomBeaconRewardsPage from "./RandomBeaconRewardsPage"

const RewardsPage = (props) => {
  return <PageWrapper {...props} />
}

RewardsPage.route = {
  title: "Rewards",
  path: "/rewards",
  pages: [RandomBeaconRewardsPage],
}

export default RewardsPage
