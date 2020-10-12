import React from "react"
import Header from "./Header"

const PageWrapper = ({ title, subLinks, children }) => {
  return (
    <>
      <Header title={title} subLinks={subLinks} />
      <main>{children}</main>
    </>
  )
}

export default PageWrapper
