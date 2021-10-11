import React from "react"
import WalletOptions from "./WalletOptions"
import Tooltip, { TOOLTIP_DIRECTION } from "./Tooltip"

const ConnectWalletBtn = ({
  text,
  btnClassName = "",
  displayExplorerMode = true,
}) => {
  return (
    <Tooltip
      direction={TOOLTIP_DIRECTION.TOP}
      simple
      className="empty-state__wallet-options-tooltip"
      triggerComponent={() => (
        <span
          className={`btn btn-primary btn-lg empty-state__connect-wallet-btn ${btnClassName}`}
        >
          {text}
        </span>
      )}
    >
      <WalletOptions displayExplorerMode={displayExplorerMode} />
    </Tooltip>
  )
}

export default ConnectWalletBtn
