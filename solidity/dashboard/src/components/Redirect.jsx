import React from "react"
import { Redirect as RedirectReadtRouterDOM } from "react-router-dom"
import useWalletAddressFromUrl from "../hooks/useWalletAddressFromUrl"

const Redirect = ({ to, ...props }) => {
  const walletAddressFromUrl = useWalletAddressFromUrl()
  const finalPath = walletAddressFromUrl ? "/" + walletAddressFromUrl + to : to
  return <RedirectReadtRouterDOM to={finalPath} />
}

export default Redirect
