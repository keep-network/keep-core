import React, { useState } from "react"
import { shortenAddress } from "../utils/general.utils"
import copy from "copy-to-clipboard"

const AddressShortcut = ({ address, classNames }) => {
  const [copyStatus, setCopyStatus] = useState("Copy to clipboard")

  const copyToClipboard = () => {
    copy(address)
      ? setCopyStatus("Copied!")
      : setCopyStatus(`Cannot copy value: ${address}!`)
  }

  return (
    <span
      onClick={copyToClipboard}
      onMouseOut={() => setCopyStatus("Copy to clipboard")}
      className={`address-shortcut tooltip address ${classNames}`}
    >
      <span className="tooltip-text address bottom">{copyStatus}</span>
      {shortenAddress(address)}
    </span>
  )
}

export default React.memo(AddressShortcut)
