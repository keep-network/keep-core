import React from "react"
import CoveragePoolPage from "./CoveragePoolPage"
import PageWrapper from "../../components/PageWrapper"
import HowItWorksPage from "./HowItWorksPage"

const CoveragePoolsPageContainer = ({ title, routes, withNewLabel }) => {
  return <PageWrapper title={title} routes={routes} newPage={withNewLabel} />
}

CoveragePoolsPageContainer.route = {
  title: "Coverage Pool",
  path: "/coverage-pools",
  pages: [HowItWorksPage, CoveragePoolPage],
  withNewLabel: true,
}

export default CoveragePoolsPageContainer
