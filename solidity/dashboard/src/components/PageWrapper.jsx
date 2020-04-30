import React from "react"

const PageWrapper = ({ title, children }) => (
  <>
    <h1 className="mb-2">{title}</h1>
    {children}
  </>
)

export default PageWrapper
