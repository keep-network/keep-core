import React from "react"
import ConnectWalletBtn from "../../ConnectWalletBtn"
import { ModalBody, ModalHeader } from "../Modal"
import { withBaseModal } from "../withBaseModal"

const WalletSelectionModalBase = () => {
  return (
    <>
      <ModalHeader>Select Wallet</ModalHeader>
      <ModalBody>
        <div className="flex column center">
          <div className="flex full-center mb-3">
            <h3 className="ml-1">
              {
                "You're viewing the Dashboard in Explorer Mode. Connect a wallet to proceed."
              }
            </h3>
          </div>
          <span className="text-center mt-1">
            <ConnectWalletBtn
              text={"Connect wallet"}
              displayExplorerMode={false}
            />
          </span>
        </div>
      </ModalBody>
    </>
  )
}

export const WalletSelectionModal = withBaseModal(WalletSelectionModalBase)
