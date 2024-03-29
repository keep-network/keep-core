import React, { useCallback } from "react"
import { ModalHeader } from "../Modal"
import { SelectedWalletModal } from "./SelectedWalletModal"
import * as Icons from "../../Icons"
import { withBaseModal } from "../withBaseModal"

const TrezorModalBase = ({ connector, connectAppWithWallet, onClose }) => {
  const fetchAccounts = useCallback(
    async (numberOfAccounts, accountsOffset) => {
      try {
        const accounts = await connector.getAccounts(
          numberOfAccounts,
          accountsOffset
        )
        return accounts
      } catch (error) {
        throw error
      }
    },
    [connector]
  )
  return (
    <>
      <ModalHeader>Connect wallet</ModalHeader>

      <SelectedWalletModal
        icon={<Icons.Trezor className="trezor-logo trezor-logo--black" />}
        walletName="Trezor"
        descriptionIcon={<Icons.TrezorDevice />}
        description="Plug in your Trezor device and unlock. If the setup screen doesn’t
        load right away, go to Trezor setup:"
        connector={connector}
        connectAppWithWallet={connectAppWithWallet}
        closeModal={onClose}
        fetchAvailableAccounts={fetchAccounts}
      >
        <a
          href="https://trezor.io/start/"
          className="btn btn-lg btn-primary mt-1 mb-1"
          target="_blank"
          rel="noopener noreferrer"
        >
          go to trezor setup
        </a>
      </SelectedWalletModal>
    </>
  )
}

export const TrezorModal = withBaseModal(TrezorModalBase)
