import React from "react"
import TBTCApplicationPage from "./TBTCApplicationPage"
import KeepRandomBeaconApplicationPage from "./KeepRandomBeaconApplicationPage"
import PageWrapper from "../components/PageWrapper"

const ApplicationsPageContainer = ({ title, routes }) => {
  return <PageWrapper title={title} routes={routes} />
}

ApplicationsPageContainer.route = {
  title: "Applications",
  path: "/applications",
  pages: [KeepRandomBeaconApplicationPage, TBTCApplicationPage],
}

export default ApplicationsPageContainer
