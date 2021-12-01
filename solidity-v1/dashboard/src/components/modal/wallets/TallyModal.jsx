import React from "react"
import { ModalHeader } from "../Modal"
import { SelectedWalletModal } from "./SelectedWalletModal"
import * as Icons from "../../Icons"
import { withBaseModal } from "../withBaseModal"

const rejectedConnectionErrorMsg =
  "You rejected the connection request in Tally. Please close and try again, and click confirm connection in Tally window."

const TallyModalBase = ({ connector, connectAppWithWallet, onClose }) => {
  return (
    <>
      <ModalHeader>Connect wallet</ModalHeader>

      <SelectedWalletModal
        icon={<Icons.Tally />}
        walletName="Tally"
        description={
          connector.getProvider()
            ? "The Tally login screen will open in an external window."
            : "Please install the Tally extension"
        }
        connector={connector}
        connectAppWithWallet={connectAppWithWallet}
        userRejectedConnectionRequestErrorMsg={rejectedConnectionErrorMsg}
        closeModal={onClose}
        connectWithWalletOnMount
      >
        {!connector.getProvider() && (
          <a
            href="https://tally.cash"
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

export const TallyModal = React.memo(withBaseModal(TallyModalBase))
