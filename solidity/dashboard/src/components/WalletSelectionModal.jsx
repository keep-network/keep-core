import React from "react"
import WalletOptions from "./WalletOptions"
import Tooltip from "./Tooltip"

export const WalletSelectionModal = () => {
  return (
    <div className="flex column center">
      <div className="flex full-center mb-3">
        <h3 className="ml-1">
          Select the wallet you want to connect with to proceed
        </h3>
      </div>
      <span className="text-center mt-1">
        {
          <Tooltip
            direction="top"
            simple
            className="empty-state__wallet-options-tooltip"
            triggerComponent={() => (
              <span
                className={`btn btn-primary btn-lg empty-state__connect-wallet-btn`}
              >
                Connect wallet
              </span>
            )}
          >
            <WalletOptions displayExplorerMode={false} />
          </Tooltip>
        }
      </span>
    </div>
  )
}
