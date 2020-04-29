import React, { useContext } from "react"
import { Web3Context } from "./WithWeb3Context"

export const NetworkStatus = (props) => {
  const { networkType, error } = useContext(Web3Context)

  return (
    <div className="network-status flex row center">
      <div className={`network-indicator ${!error ? "connected" : "error"}`} />
      <span className="text-label">
        {error ? "wrong network" : `connected: ${networkType}`}
      </span>
    </div>
  )
}
