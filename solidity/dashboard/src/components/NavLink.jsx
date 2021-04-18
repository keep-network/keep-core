import React from "react"
import { NavLink as NavLinkReactRouterDOM } from "react-router-dom"
import useWalletAddressFromUrl from "../hooks/useWalletAddressFromUrl"

const NavLink = ({ to, ...props }) => {
  const walletAddressFromUrl = useWalletAddressFromUrl()
  const finalPath = walletAddressFromUrl ? "/" + walletAddressFromUrl + to : to

  return <NavLinkReactRouterDOM to={finalPath} {...props} />
}

export default NavLink
