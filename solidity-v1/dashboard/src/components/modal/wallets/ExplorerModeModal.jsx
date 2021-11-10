import React, { useEffect, useState } from "react"
import { SelectedWalletModal } from "./SelectedWalletModal"
import { ModalBody, ModalHeader } from "../Modal"
import { withBaseModal } from "../withBaseModal"
import ExplorerModeAddressForm from "../../ExplorerModeAddressForm"
import * as Icons from "../../Icons"
import { WALLETS } from "../../../constants/constants"

const ExplorerModeModalBase = ({
  connector,
  connectAppWithWallet,
  onClose,
  address = "",
}) => {
  const [selectedAddress, setSelectedAddress] = useState("")

  useEffect(() => {
    connector.setSelectedAccount(selectedAddress)
  }, [selectedAddress, connector])

  useEffect(() => {
    if (address) setSelectedAddress(address)
  }, [address])

  const submitAction = (values) => {
    setSelectedAddress(values.address)
  }

  return (
    <>
      <ModalHeader>Connect Ethereum Address</ModalHeader>
      {selectedAddress ? (
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
          closeModal={onClose}
          connectWithWalletOnMount={true}
        />
      ) : (
        <ModalBody>
          <ExplorerModeAddressForm
            submitAction={submitAction}
            onCancel={onClose}
          />
        </ModalBody>
      )}
    </>
  )
}

export const ExplorerModeModal = React.memo(
  withBaseModal(ExplorerModeModalBase)
)
