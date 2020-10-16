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
        connector
          ? "The MetaMask login screen will open in an external window."
          : "Please install the MetaMask extension"
      }
      providerName="METAMASK"
      connector={connector}
      connectAppWithWallet={connectAppWithWallet}
      closeModal={closeModal}
      connectWithWalletOnMount
    >
      {!connector && (
        <a
          href="https://metamask.io"
          className="btn bt-lg btn-primary mt-3"
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
