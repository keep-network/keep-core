import React, { useEffect } from "react"
import { useUpdateHeaderData } from "./Header"

const PageWrapper = ({ title, subLinks, children }) => {
  const updateHeaderData = useUpdateHeaderData()

  useEffect(() => {
    updateHeaderData(title, subLinks)
  }, [title, subLinks, updateHeaderData])

  return <>{children}</>
}

export default PageWrapper
