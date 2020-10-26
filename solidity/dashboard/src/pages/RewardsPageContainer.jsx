import React from "react"
import TBTCRewardsPage from "./TBTCRewardsPage"
import RewardsPage from "./RewardsPage"
import PageWrapper from "../components/PageWrapper"

const RewardsPageContainer = ({ title, routes }) => {
  return <PageWrapper title={title} routes={routes} />
}

// TODO: Rename to `earings` in the followup PR, according to the Figma
// views.
RewardsPageContainer.route = {
  title: "Rewards",
  path: "/rewards",
  pages: [RewardsPage, TBTCRewardsPage],
}

export default RewardsPageContainer
