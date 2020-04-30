import React, { useContext } from "react"
import { Web3Context } from "./WithWeb3Context"

export const NetworkStatus = (props) => {
  const { networkType, error } = useContext(Web3Context)

  return (
    <div className="network-status flex row center">
      <div className={`network-indicator ${!error ? "connected" : "error"}`} />
      <h5 className="text-grey-50">
        {error ? "wrong network" : `connected: ${networkType}`}
      </h5>
    </div>
  )
}
