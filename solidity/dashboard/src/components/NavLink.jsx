import React from "react"
import { NavLink as NavLinkReactRouterDOM } from "react-router-dom"
import useFinalPath from "../hooks/useFinalPath"

const NavLink = ({ to, ...props }) => {
  const finalPath = useFinalPath(to)

  return <NavLinkReactRouterDOM to={finalPath} {...props} />
}

export default NavLink
