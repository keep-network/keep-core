import React from "react"
import { ModalHeader } from "../Modal"
import { SelectedWalletModal } from "./SelectedWalletModal"
import * as Icons from "../../Icons"
import { WALLETS } from "../../../constants/constants"
import { withBaseModal } from "../withBaseModal"

const WalletConnectModalBase = ({
  connector,
  connectAppWithWallet,
  onClose,
}) => {
  return (
    <>
      <ModalHeader>Connect wallet</ModalHeader>
      <SelectedWalletModal
        icon={
          <Icons.WalletConnect
            className="wallet-connect-logo wallet-connect-logo--black"
            width={30}
            height={28}
          />
        }
        walletName={WALLETS.WALLET_CONNECT.label}
        connector={connector}
        connectAppWithWallet={connectAppWithWallet}
        closeModal={onClose}
        connectWithWalletOnMount
      />
    </>
  )
}

export const WalletConnectModal = React.memo(
  withBaseModal(WalletConnectModalBase)
)
