import React from "react"
import PageWrapper from "../components/PageWrapper"
import DocumentationSection from "../components/glossary/DocumentationSection"
import TerminologyDataTable from "../components/glossary/TerminologyDataTable"
import DelegationDiagram from "../components/glossary/DelegationDiagram"
import Navigation from "../components/glossary/Navigation"

const GlossaryPage = () => {
  return (
    <PageWrapper title="Glossary">
      <div className="glossary-page-wrapper">
        <div>
          <DocumentationSection />
          <TerminologyDataTable />
          <DelegationDiagram />
        </div>
        <Navigation />
      </div>
    </PageWrapper>
  )
}

export default GlossaryPage
