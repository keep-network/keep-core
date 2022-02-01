import React from "react"
import TBTCApplicationPage from "./TBTCApplicationPage"
import KeepRandomBeaconApplicationPage from "./KeepRandomBeaconApplicationPage"
import PageWrapper from "../../components/PageWrapper"
import ThresholdApplicationPage from "./ThresholdApplicationPage"

const ApplicationsPageContainer = ({ title, routes, withNewLabel }) => {
  return <PageWrapper title={title} routes={routes} newPage={withNewLabel} />
}

ApplicationsPageContainer.route = {
  title: "Applications",
  path: "/applications",
  pages: [
    ThresholdApplicationPage,
    TBTCApplicationPage,
    KeepRandomBeaconApplicationPage,
  ],
  withNewLabel: true,
}

export default ApplicationsPageContainer
