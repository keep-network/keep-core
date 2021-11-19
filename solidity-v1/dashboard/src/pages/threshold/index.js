import React from "react"
import HowItWorksPage from "./HowItWorksPage"
import ThresholdUpgradePage from "./ThresholdUpgradePage"
import PageWrapper from "../../components/PageWrapper"

const ThresholdPageContainer = ({ routes, withNewLabel }) => {
  return (
    <PageWrapper
      title={"Threshold Upgrade"}
      routes={routes}
      newPage={withNewLabel}
    />
  )
}

ThresholdPageContainer.route = {
  title: "Threshold",
  path: "/threshold",
  pages: [ThresholdUpgradePage, HowItWorksPage],
  withNewLabel: true,
}

export default ThresholdPageContainer
