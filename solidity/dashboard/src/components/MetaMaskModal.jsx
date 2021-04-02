import React from "react"
import SelectedWalletModal from "./SelectedWalletModal"
import * as Icons from "./Icons"

const rejectedConnectionErrorMsg =
  "You rejected the connection request in MetaMask. Please close and try again, and click confirm connection in MetaMask window."

const MetaMaskModal = ({
  connector,
  connectAppWithWallet,
  closeModal,
  payload = null,
}) => {
  return (
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
      closeModal={closeModal}
      connectWithWalletOnMount
      payload={payload}
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
  )
}

export default React.memo(MetaMaskModal)
