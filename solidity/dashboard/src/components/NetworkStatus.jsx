import React, { useContext } from "react"
import { Web3Context } from "./WithWeb3Context"

export const NetworkStatus = () => {
  const { networkType, provider, error, yourAddress } = useContext(Web3Context)

  return (
    <div className="network-status flex row center">
      <div
        className={`network-indicator ${
          yourAddress && !error && provider !== null ? "connected" : "error"
        }`}
      />
      <h5 className="text-label">
        {!yourAddress && "not connected"}
        {yourAddress &&
          !error &&
          provider !== null &&
          `connected: ${networkType}`}
        {yourAddress && error && provider !== null && `wrong network`}
      </h5>
    </div>
  )
}
