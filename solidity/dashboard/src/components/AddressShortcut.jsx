import React, { useState } from "react"
import { shortenAddress } from "../utils/general.utils"
import copy from "copy-to-clipboard"
import Tooltip from "./Tooltip"

const AddressShortcut = ({ address, classNames }) => {
  const [copyStatus, setCopyStatus] = useState("Copy to clipboard")

  const copyToClipboard = () => {
    copy(address)
      ? setCopyStatus("Copied!")
      : setCopyStatus(`Cannot copy value: ${address}!`)
  }

  return (
    <Tooltip
      simple
      triggerComponent={() => {
        return (
          <span
            onClick={copyToClipboard}
            onMouseOut={() => setCopyStatus("Copy to clipboard")}
            className={`address-shortcut ${classNames}`}
          >
            {shortenAddress(address)}
          </span>
        )
      }}
    >
      {copyStatus}
    </Tooltip>
  )
}

export default React.memo(AddressShortcut)
