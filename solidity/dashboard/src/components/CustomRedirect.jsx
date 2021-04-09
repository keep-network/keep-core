import React from "react"
import { Redirect } from "react-router-dom"
import useWalletAddressFromUrl from "../hooks/useWalletAddressFromUrl"

const CustomRedirect = ({ to, ...props }) => {
  const walletAddressFromUrl = useWalletAddressFromUrl()
  const finalPath = walletAddressFromUrl ? "/" + walletAddressFromUrl + to : to
  return <Redirect to={finalPath} />
}

export default CustomRedirect
