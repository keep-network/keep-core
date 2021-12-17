import React from "react"
import { ModalHeader } from "../Modal"
import { SelectedWalletModal } from "./SelectedWalletModal"
import * as Icons from "../../Icons"
import { withBaseModal } from "../withBaseModal"

const rejectedConnectionErrorMsg =
  "You rejected the connection request in MetaMask. Please close and try again, then click confirm connection in MetaMask window."

const MetaMaskModalBase = ({ connector, connectAppWithWallet, onClose }) => {
  return (
    <>
      <ModalHeader>Connect wallet</ModalHeader>

      <SelectedWalletModal
        icon={<Icons.MetaMask />}
        walletName="MetaMask"
        description={
          connector.getProvider()
            ? "The MetaMask login screen will open in an external window."
            : "Please install the MetaMask extension"
        }
        connector={connector}
        connectAppWithWallet={connectAppWithWallet}
        userRejectedConnectionRequestErrorMsg={rejectedConnectionErrorMsg}
        closeModal={onClose}
        connectWithWalletOnMount
      >
        {!connector.getProvider() && (
          <a
            href="https://metamask.io"
            className="btn btn-lg btn-primary mt-1 mb-1"
            target="_blank"
            rel="noopener noreferrer"
          >
            install extension
          </a>
        )}
      </SelectedWalletModal>
    </>
  )
}

export const MetaMaskModal = React.memo(withBaseModal(MetaMaskModalBase))
