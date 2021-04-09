import React from "react"
import { NavLink } from "react-router-dom"
import useWalletAddressFromUrl from "../hooks/useWalletAddressFromUrl"

const CustomNavLink = ({ to, ...props }) => {
  const walletAddressFromUrl = useWalletAddressFromUrl()
  const finalPath = walletAddressFromUrl ? "/" + walletAddressFromUrl + to : to

  return <NavLink to={finalPath} {...props} />
}

export default CustomNavLink
