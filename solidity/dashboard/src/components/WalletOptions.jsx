import React, { useCallback } from "react"
import * as Icons from "./Icons"
import { useWeb3Context } from "./WithWeb3Context"
import { useModal } from "../hooks/useModal"
import LedgerModal from "./LedgerModal"
import TrezorModal from "./TrezorModal"
import {
  TrezorConnector,
  LedgerConnector,
  LEDGER_DERIVATION_PATHS,
  InjectedConnector,
  WalletConnectConnector,
} from "../connectors"
import MetaMaskModal from "./MetaMaskModal"
import WallectConnectModal from "./WalletConnectModal"
import { WALLETS } from "../constants/constants"

const WALLETS_OPTIONS = [
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
    connector: new InjectedConnector(),
  },
  {
    label: "Ledger",
    icon: Icons.Ledger,
    providerName: "LEDGER",
    isHardwareWallet: true,
    connector: {
      LEDGER_LIVE: new LedgerConnector(LEDGER_DERIVATION_PATHS.LEDGER_LIVE),
      LEDGER_LEGACY: new LedgerConnector(LEDGER_DERIVATION_PATHS.LEDGER_LEGACY),
    },
  },
  {
    label: "Trezor",
    icon: Icons.Trezor,
    providerName: "TREZOR",
    isHardwareWallet: true,
    connector: new TrezorConnector(),
  },
  {
    label: "WalletConnect",
    // TODO icon
    icon: Icons.Trezor,
    providerName: "WALLET_CONNECT",
    isHardwareWallet: false,
    connector: new WalletConnectConnector(),
  },
]

const WalletOptions = () => {
  return (
    <ul className="wallet__options">{WALLETS_OPTIONS.map(renderWallet)}</ul>
  )
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
    const defaultProps = { connector, closeModal, connectAppWithWallet }
    switch (providerName) {
      case WALLETS.LEDGER.name:
        return <LedgerModal {...defaultProps} />
      case WALLETS.TREZOR.name:
        return <TrezorModal {...defaultProps} />
      case WALLETS.METAMASK.name:
        return <MetaMaskModal {...defaultProps} />
      case WALLETS.WALLET_CONNECT.name:
        return <WallectConnectModal {...defaultProps} />
      default:
        return null
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
      }}
    >
      <IconComponent className="wallet__item__icon" />
      {label}
    </li>
  )
}

export default WalletOptions
