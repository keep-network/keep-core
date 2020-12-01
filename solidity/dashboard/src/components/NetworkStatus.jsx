import React from "react"
import { useWeb3Context } from "./WithWeb3Context"
import * as Icons from "./Icons"

export const NetworkStatus = () => {
  const { networkType, error, isConnected } = useWeb3Context()

  let status = "disconnected"
  if (error) {
    status = "error"
  } else if (isConnected) {
    status = "connected"
  }

  return (
    <div className="network-status flex row center">
      <Icons.NetworkStatusIndicator
        className={`network-status__indicator--${status}`}
      />
      <span className={`network-status__text--${status}`}>
        {status === "disconnected" && "Network Disconnected"}
        {status === "error" && error}
        {status === "connected" && networkType}
      </span>
    </div>
  )
}
