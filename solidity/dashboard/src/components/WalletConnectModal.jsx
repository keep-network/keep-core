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
      icon={<Icons.MetaMask />}
      walletName={WALLETS.WALLET_CONNECT.label}
      iconDescription={null}
      // TODO
      description={"desc"}
      connector={connector}
      connectAppWithWallet={connectAppWithWallet}
      closeModal={closeModal}
      connectWithWalletOnMount
    />
  )
}

export default React.memo(WalletConnectModal)
