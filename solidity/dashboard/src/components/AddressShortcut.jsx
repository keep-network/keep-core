import React, { useMemo } from "react"
import { shortenAddress } from "../utils/general.utils"
import CopyToClipboard from "./CopyToClipboard"

const AddressShortcut = ({ address, classNames }) => {
  const addr = useMemo(() => shortenAddress(address), [address])

  return (
    <CopyToClipboard
      toCopy={address}
      render={({ reset, copyStatus, copyToClipboard }) => {
        return (
          <span
            onClick={copyToClipboard}
            onMouseOut={reset}
            className={`address-shortcut tooltip address ${classNames}`}
          >
            <span className="tooltip-text address bottom">{copyStatus}</span>
            {addr}
          </span>
        )
      }}
    />
  )
}

export default React.memo(AddressShortcut)
