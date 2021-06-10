import React from "react"
import PageWrapper from "../../components/PageWrapper"
import {
  CheckListBanner,
  HowDoesItWorkBanner,
} from "../../components/coverage-pools"

const CoveragePoolPage = ({ title, withNewLabel }) => {
  return (
    <PageWrapper title={title} newPage={withNewLabel}>
      <CheckListBanner />
      <HowDoesItWorkBanner />
    </PageWrapper>
  )
}

export default CoveragePoolPage
