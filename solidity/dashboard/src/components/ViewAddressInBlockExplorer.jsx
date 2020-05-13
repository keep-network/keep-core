import React from "react"
import { useEtherscanUrl } from "../hooks/useEtherscanUrl"

const ViewAddressInBlockExplorer = ({ address, text, urlSuffix }) => {
  const etherscanDefaultUrl = useEtherscanUrl()

  return (
    <a
      href={`${etherscanDefaultUrl}/address/${address}${urlSuffix}`}
      rel="noopener noreferrer"
      target="_blank"
      className="arrow-link"
    >
      {text}
    </a>
  )
}

ViewAddressInBlockExplorer.defaultProps = {
  text: "View in Block Explorer",
  urlSuffix: "#code",
}

export default ViewAddressInBlockExplorer
