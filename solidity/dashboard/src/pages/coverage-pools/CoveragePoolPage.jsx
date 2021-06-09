import React from "react"
import PageWrapper from "../../components/PageWrapper"
import { CheckListBanner } from "../../components/coverage-pools"

const CoveragePoolPage = ({ title, withNewLabel }) => {
  return (
    <PageWrapper title={title} newPage={withNewLabel}>
      <CheckListBanner />
    </PageWrapper>
  )
}

export default CoveragePoolPage
