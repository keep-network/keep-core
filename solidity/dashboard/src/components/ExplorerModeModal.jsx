import React, { useEffect } from "react"
import SelectedWalletModal from "./SelectedWalletModal"
import * as Icons from "./Icons"
import { WALLETS } from "../constants/constants"
import { useState } from "react"
import ExplorerModeAddressForm from "./ExplorerModeAddressForm"

const ExplorerModeModal = ({
  connector,
  connectAppWithWallet,
  closeModal,
  address = "",
}) => {
  const [selectedAddress, setSelectedAddress] = useState("")

  useEffect(() => {
    connector.setSelectedAccount(selectedAddress)
  }, [selectedAddress])

  useEffect(() => {
    if (address) setSelectedAddress(address)
  }, [address])

  const submitAction = (values) => {
    setSelectedAddress(values.address)
  }

  return selectedAddress ? (
    <SelectedWalletModal
      icon={
        <Icons.Explore
          className="wallet-connect-logo wallet-connect-logo--black"
          width={30}
          height={28}
        />
      }
      walletName={WALLETS.EXPLORER_MODE.label}
      connector={connector}
      connectAppWithWallet={connectAppWithWallet}
      closeModal={closeModal}
      connectWithWalletOnMount={true}
    />
  ) : (
    <ExplorerModeAddressForm
      submitAction={submitAction}
      onCancel={closeModal}
    />
  )
}

export default React.memo(ExplorerModeModal)
