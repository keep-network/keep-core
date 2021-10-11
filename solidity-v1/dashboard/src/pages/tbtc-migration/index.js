import React from "react"
import HowItWorksPage from "./HowItWorksPage"
import TokenUpgradePortalPage from "./TokenUpgradePortalPage"
import PageWrapper from "../../components/PageWrapper"

const TBTCMigrationsPageContainer = ({ title, routes, withNewLabel }) => {
  return <PageWrapper title={title} routes={routes} newPage={withNewLabel} />
}

TBTCMigrationsPageContainer.route = {
  title: "tBTC Token Upgrade Portal",
  path: "/tbtc-migration",
  pages: [HowItWorksPage, TokenUpgradePortalPage],
  withNewLabel: true,
}

export default TBTCMigrationsPageContainer
