import React from "react"
import TBTCRewardsPage from "./TBTCRewardsPage"
import RewardsPage from "./RewardsPage"
import PageWrapper from "../components/PageWrapper"

const RewardsPageContainer = ({ title, routes }) => {
  return <PageWrapper title={title} routes={routes} />
}

RewardsPageContainer.route = {
  title: "Earnings",
  path: "/earnings",
  pages: [RewardsPage, TBTCRewardsPage],
}

export default RewardsPageContainer
