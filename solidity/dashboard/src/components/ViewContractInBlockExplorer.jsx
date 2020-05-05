import React from "react"
import { useEtherscanUrl } from "../hooks/useEtherscanUrl"

const ViewContractInBlockExplorer = ({ contractAddress, text }) => {
  const etherscanDefaultUrl = useEtherscanUrl()

  return (
    <a
      href={`${etherscanDefaultUrl}/address/${contractAddress}#code`}
      rel="noopener noreferrer"
      target="_blank"
    >
      {text}
    </a>
  )
}

ViewContractInBlockExplorer.defaultProps = {
  text: "View in Block Explorer",
}

export default ViewContractInBlockExplorer
