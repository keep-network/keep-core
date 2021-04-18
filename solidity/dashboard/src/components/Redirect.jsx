import React from "react"
import { Redirect as RedirectReadtRouterDOM } from "react-router-dom"
import useWalletAddressFromUrl from "../hooks/useWalletAddressFromUrl"
import useFinalPath from "../hooks/useFinalPath"

const Redirect = ({ to, ...props }) => {
  const walletAddressFromUrl = useWalletAddressFromUrl()
  const finalPath = useFinalPath(to, walletAddressFromUrl)
  return <RedirectReadtRouterDOM to={finalPath} />
}

export default Redirect
