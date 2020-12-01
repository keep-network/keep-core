import React from "react"
import PageWrapper from "../components/PageWrapper"
import DocumentationSection from "../components/resources/DocumentationSection"
import TerminologyDataTable from "../components/resources/TerminologyDataTable"
import DelegationDiagram from "../components/resources/DelegationDiagram"

const ResourcesPage = ({ title, routes }) => {
  return <PageWrapper title={title} routes={routes} />
}

DocumentationSection.route = {
  title: "Documentation",
  path: "/resources/documentation",
  exact: true,
}

DelegationDiagram.route = {
  title: "Diagrams",
  path: "/resources/diagram",
  exact: true,
}

TerminologyDataTable.route = {
  title: "Quick Terminology",
  path: "/resources/quick-terminology",
  exact: true,
}

ResourcesPage.route = {
  title: "Resources",
  path: "/resources",
  pages: [DocumentationSection, TerminologyDataTable, DelegationDiagram],
}

export default ResourcesPage
