import React from "react"
import SelectedWalletModal from "./SelectedWalletModal"
import * as Icons from "./Icons"

const MetaMaskModal = ({ connector, connectAppWithWallet, closeModal }) => {
  return (
    <SelectedWalletModal
      icon={<Icons.MetaMask />}
      walletName="MetaMask"
      iconDescription={null}
      description={
        connector.getProvider()
          ? "The MetaMask login screen will open in an external window."
          : "Please install the MetaMask extension"
      }
      connector={connector}
      connectAppWithWallet={connectAppWithWallet}
      closeModal={closeModal}
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
  )
}

export default React.memo(MetaMaskModal)
