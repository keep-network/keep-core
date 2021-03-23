import React from "react"
import SelectedWalletModal from "./SelectedWalletModal"
import * as Icons from "./Icons"
import { WALLETS } from "../constants/constants"

const WalletConnectModal = ({
  connector,
  connectAppWithWallet,
  closeModal,
}) => {
  return (
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
      closeModal={closeModal}
      connectWithWalletOnMount
    />
  )
}

export default React.memo(WalletConnectModal)
