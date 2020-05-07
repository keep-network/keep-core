import React from "react"

const PageWrapper = ({ title, children, ...titleProps }) => (
  <>
    <h1 className="mb-2" {...titleProps}>
      {title}
    </h1>
    {children}
  </>
)

export default PageWrapper
