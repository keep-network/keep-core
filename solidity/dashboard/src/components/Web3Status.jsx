import React from "react"
// import React, { useContext } from "react"
// import { Web3Context } from "./WithWeb3Context"

// Temporary cleanup. The PR that updates this compoent is in progress pls see
// https://github.com/keep-network/keep-core/pull/2107
export const Web3Status = () => {
  // const {} = useContext(Web3Context)

  const renderStatus = () => {
    return null
  }

  return <div className="web3">{renderStatus()}</div>
}
