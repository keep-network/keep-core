import React from "react"
import PageWrapper from "../components/PageWrapper"
import DocumentationSection from "../components/glossary/DocumentationSection"
import TerminologyDataTable from "../components/glossary/TerminologyDataTable"
import DelegationDiagram from "../components/glossary/DelegationDiagram"

const GlossaryPage = () => {
  return (
    <PageWrapper title="Glossary">
      <DocumentationSection />
      <TerminologyDataTable />
      <DelegationDiagram />
    </PageWrapper>
  )
}

export default GlossaryPage
