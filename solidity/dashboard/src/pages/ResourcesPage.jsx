import React from "react"
import PageWrapper from "../components/PageWrapper"
import DocumentationSection from "../components/resources/DocumentationSection"
import TerminologyDataTable from "../components/resources/TerminologyDataTable"
import DelegationDiagram from "../components/resources/DelegationDiagram"

const ResourcesPage = () => {
  return (
    <PageWrapper title="Resources">
      <div className="resources-page-wrapper">
        <div>
          <DocumentationSection />
          <DelegationDiagram />
          <TerminologyDataTable />
        </div>
      </div>
    </PageWrapper>
  )
}

export default ResourcesPage
