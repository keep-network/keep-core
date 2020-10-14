import React, { useCallback } from "react"
import * as Icons from "./Icons"
import { useWeb3Context } from "./WithWeb3Context"
import { useModal } from "../hooks/useModal"
import SelectedWalletModal from "./SelectedWalletModal"
import LedgerModal from "./LedgerModal"
import TrezorModal from "./TrezorModal"
import { TrezorProvider } from "../connectors/trezor"
import { LedgerProvider, LEDGER_DERIVATION_PATHS } from "../connectors/ledger"

const WALLETS = [
  {
    label: "MetaMask",
    icon: Icons.MetaMask,
    providerName: "METAMASK",
    modalProps: {
      iconDescription: null,
      btnText: "install extension",
      btnLink: "https://metamask.io",
      description:
        "The MetaMask login screen will open in an external window. If it doesn’t load right away, click below to install:",
    },
    isHardwareWallet: false,
  },
  {
    label: "Ledger",
    icon: Icons.Ledger,
    providerName: "LEDGER",
    isHardwareWallet: true,
    connector: {
      LEDGER_LIVE: new LedgerProvider(LEDGER_DERIVATION_PATHS.LEDGER_LIVE),
      LEDGER_LEGACY: new LedgerProvider(LEDGER_DERIVATION_PATHS.LEDGER_LEGACY),
    },
  },
  {
    label: "Trezor",
    icon: Icons.Trezor,
    providerName: "TREZOR",
    isHardwareWallet: true,
    connector: new TrezorProvider(),
  },
]

const WalletOptions = () => {
  return <ul className="wallet__options">{WALLETS.map(renderWallet)}</ul>
}

const renderWallet = (wallet) => <Wallet key={wallet.label} {...wallet} />

const Wallet = ({
  label,
  icon: IconComponent,
  providerName,
  modalProps,
  connector,
}) => {
  const { openModal, closeModal } = useModal()
  const { connectAppWithWallet, abortWalletConnection } = useWeb3Context()

  const customCloseModal = useCallback(() => {
    abortWalletConnection()
    closeModal()
  }, [abortWalletConnection, closeModal])

  const renderModalContent = () => {
    switch (providerName) {
      case "LEDGER":
        return (
          <LedgerModal
            connector={connector}
            closeModal={closeModal}
            connectAppWithWallet={connectAppWithWallet}
          />
        )
      case "TREZOR":
        return (
          <TrezorModal
            connector={connector}
            closeModal={closeModal}
            connectAppWithWallet={connectAppWithWallet}
          />
        )
      case "METAMASK":
      case "COINBASE":
      default:
        return (
          <SelectedWalletModal
            walletName={label}
            icon={<IconComponent />}
            {...modalProps}
          />
        )
    }
  }

  return (
    <li
      className="wallet__item"
      onClick={async () => {
        openModal(renderModalContent(), {
          title: "Connect Wallet",
          closeModal: customCloseModal,
        })
        if (providerName === "METAMASK") {
          await connectAppWithWallet(window.ethereum, providerName)
          closeModal()
        }
      }}
    >
      <IconComponent className="wallet__item__icon" />
      {label}
    </li>
  )
}

export default WalletOptions
