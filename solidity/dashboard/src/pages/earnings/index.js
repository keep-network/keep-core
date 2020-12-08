import React from "react"
import TBTCRewardsPage from "./TBTCRewardsPage"
import KeepRandomBeaconRewardsPage from "./RewardsPage"
import PageWrapper from "../../components/PageWrapper"

const RewardsPage = ({ title, routes }) => {
  return <PageWrapper title={title} routes={routes} />
}

RewardsPage.route = {
  title: "Earnings",
  path: "/earnings",
  pages: [KeepRandomBeaconRewardsPage, TBTCRewardsPage],
}

export default RewardsPage
