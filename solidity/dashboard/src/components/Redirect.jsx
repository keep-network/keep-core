import React from "react"
import { Redirect as RedirectReadtRouterDOM } from "react-router-dom"
import useFinalPath from "../hooks/useFinalPath"

const Redirect = ({ to, ...props }) => {
  const finalPath = useFinalPath(to)
  return <RedirectReadtRouterDOM to={finalPath} />
}

export default Redirect
