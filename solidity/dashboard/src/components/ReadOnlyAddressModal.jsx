import React, { useEffect } from "react"
import SelectedWalletModal from "./SelectedWalletModal"
import * as Icons from "./Icons"
import { WALLETS } from "../constants/constants"
import { useState } from "react"
import ExplorerModeAddressForm from "./ExplorerModeAddressForm";

const ReadOnlyAddressModal = ({
  connector,
  connectAppWithWallet,
  closeModal,
}) => {
  const [selectedAddress, setSelectedAddress] = useState("")

  useEffect(() => {
    connector.setSelectedAccount(selectedAddress)
  }, [selectedAddress])

  const submitAction = (values) => {
    setSelectedAddress(values.address)
  }

  return selectedAddress ? (
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
  ) : (
    <div>
      <ExplorerModeAddressForm submitAction={submitAction} />
    </div>
  )
}

export default React.memo(ReadOnlyAddressModal)
