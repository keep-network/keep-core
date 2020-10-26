import React, { useMemo } from "react"
import CopyToClipboard from "./CopyToClipboard"
import { shortenAddress } from "../utils/general.utils"

const AddressShortcut = ({ address, classNames }) => {
  const addr = useMemo(() => shortenAddress(address), [address])

  return (
    <CopyToClipboard
      toCopy={address}
      render={({ reset, copyStatus, copyToClipboard }) => {
        return (
          <button
            onClick={copyToClipboard}
            onMouseOut={reset}
            className={`address-shortcut tooltip address ${classNames}`}
          >
            <span className="tooltip-text address bottom">{copyStatus}</span>
            {addr}
          </button>
        )
      }}
    />
  )
}

export default React.memo(AddressShortcut)
