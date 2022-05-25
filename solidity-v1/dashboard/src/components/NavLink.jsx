import React from "react"
import { NavLink as NavLinkReactRouterDOM } from "react-router-dom"
import useFinalPath from "../hooks/useFinalPath"

const NavLink = ({ to, ...props }) => {
  // NavLink from react-router-dom can also accept `to` prop as an object so we
  // have to extract the pathname from it if that will be the case here
  const path = typeof to === "object" ? to.pathname : to
  const finalPath = useFinalPath(path)
  if (typeof to === "object") {
    to.pathname = finalPath
  }

  return (
    <NavLinkReactRouterDOM
      to={typeof to === "object" ? to : finalPath}
      {...props}
    />
  )
}

export default NavLink
