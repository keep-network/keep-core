import React, { useEffect } from "react"
import SelectedWalletModal from "./SelectedWalletModal"
import * as Icons from "./Icons"
import { WALLETS } from "../constants/constants"

const ReadOnlyAddressModal = ({
  connector,
  connectAppWithWallet,
  closeModal,
}) => {
  useEffect(() => {
    connector.setSelectedAccount("0x857173e7c7d76e051e80d30FCc3EA6A9C2b53756")
  }, [])

  return (
    <SelectedWalletModal
      icon={
        <Icons.Wallet
          className="wallet-connect-logo wallet-connect-logo--black"
          width={30}
          height={28}
        />
      }
      walletName={WALLETS.READ_ONLY_ADDRESS.label}
      connector={connector}
      connectAppWithWallet={connectAppWithWallet}
      closeModal={closeModal}
      connectWithWalletOnMount
    />
  )
}

export default React.memo(ReadOnlyAddressModal)
