import React from "react"
import { useEtherscanUrl } from "../hooks/useEtherscanUrl"

export const ViewInBlockExplorer = React.memo(
  ({ type, id, hashParam, text, ...restProps }) => {
    const etherscanDefaultUrl = useEtherscanUrl()

    return (
      <a
        href={`${etherscanDefaultUrl}/${type}/${id}${hashParam}`}
        className="arrow-link"
        {...restProps}
        rel="noopener noreferrer"
        target="_blank"
      >
        {text}
      </a>
    )
  }
)

ViewInBlockExplorer.defaultProps = {
  text: "View in Block Explorer",
  type: "address",
  hashParam: "",
}

export const ViewAddressInBlockExplorer = React.memo(
  ({ address, text, urlSuffix }) => {
    return (
      <ViewInBlockExplorer
        text={text}
        type="address"
        id={address}
        hashParam={urlSuffix}
      />
    )
  }
)

ViewAddressInBlockExplorer.defaultProps = {
  urlSuffix: "#code",
}
