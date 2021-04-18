import React from "react"
import { NavLink as NavLinkReactRouterDOM } from "react-router-dom"
import useWalletAddressFromUrl from "../hooks/useWalletAddressFromUrl"
import useFinalPath from "../hooks/useFinalPath";

const NavLink = ({ to, ...props }) => {
  const walletAddressFromUrl = useWalletAddressFromUrl()
  const finalPath = useFinalPath(to, walletAddressFromUrl)

  return <NavLinkReactRouterDOM to={finalPath} {...props} />
}

export default NavLink
