import React, { useContext } from "react"
import { Web3Context } from "./WithWeb3Context"

export const NetworkStatus = () => {
  const { networkType, provider, error } = useContext(Web3Context)

  return (
    <div className="network-status flex row center">
      <div
        className={`network-indicator ${
          !error && provider !== null ? "connected" : "error"
        }`}
      />
      <h5 className="text-label">
        {!error && provider === null && "not connected"}
        {!error && provider !== null && `connected: ${networkType}`}
        {error && provider !== null && `wrong network`}
      </h5>
    </div>
  )
}
