import React from "react"
import WalletOptions from "./WalletOptions"
import Tooltip from "./Tooltip"

const ConnectWalletBtn = ({
  text,
  btnClassName = "",
  displayExplorerMode = true,
}) => {
  return (
    <Tooltip
      direction="top"
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
