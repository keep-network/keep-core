import React, { useMemo } from "react"
import { shortenAddress } from "../utils/general.utils"
import Tooltip from "./Tooltip"
import CopyToClipboard from "./CopyToClipboard"

const AddressShortcut = ({ address, classNames }) => {
  const addr = useMemo(() => shortenAddress(address), [address])

  return (
    <CopyToClipboard
      toCopy={address}
      render={({ reset, copyStatus, copyToClipboard }) => {
        return (
          <Tooltip
            simple
            triggerComponent={() => {
              return (
                <button
                  onClick={copyToClipboard}
                  onMouseOut={reset}
                  className={`address-shortcut ${classNames}`}
                >
                  {addr}
                </button>
              )
            }}
          >
            {copyStatus}
          </Tooltip>
        )
      }}
    />
  )
}

export default React.memo(AddressShortcut)
