import React, { useCallback } from "react"
import * as Icons from "./Icons"
import { useWeb3Context } from "./WithWeb3Context"
import { useModal } from "../hooks/useModal"
import Tile from "./Tile"
import PageWrapper from "./PageWrapper"
import LedgerModal from "./LedgerModal"
import TrezorModal from "./TrezorModal"
import {
  TrezorProvider,
  LedgerProvider,
  LEDGER_DERIVATION_PATHS,
  InjectedProvider,
} from "../connectors"
import MetaMaskModal from "./MetaMaskModal"

const WALLETS = [
  {
    label: "MetaMask",
    icon: <Icons.MetaMask />,
    providerName: "METAMASK",
    type: "web wallet",
    description: "Crypto wallet that’s a web browser plugin.",
    modalProps: {
      iconDescription: null,
      btnText: "install extension",
      btnLink: "https://metamask.io",
      description:
        "The MetaMask login screen will open in an external window. If it doesn’t load right away, click below to install:",
    },
    isHardwareWallet: false,
    connector: InjectedProvider,
  },
  {
    label: "Coinbase",
    icon: <Icons.Coinbase />,
    providerName: "COINBASE",
    type: "web wallet",
    description: "Crypto wallet that’s a web browser plugin.",
    modalProps: {
      iconDescription: null,
      btnText: null,
      btnLink: null,
      description: "Scan QR code to connect:",
    },
    isHardwareWallet: false,
  },
  {
    label: "Ledger",
    icon: <Icons.Ledger />,
    providerName: "LEDGER",
    type: "hardware wallet",
    description: "Crypto wallet on a secure hardware device.",
    isHardwareWallet: true,
    connector: {
      LEDGER_LIVE: new LedgerProvider(LEDGER_DERIVATION_PATHS.LEDGER_LIVE),
      LEDGER_LEGACY: new LedgerProvider(LEDGER_DERIVATION_PATHS.LEDGER_LEGACY),
    },
  },
  {
    label: "Trezor",
    icon: <Icons.Trezor />,
    providerName: "TREZOR",
    type: "hardware wallet",
    description: "Crypto wallet on a secure hardware device.",
    isHardwareWallet: true,
    connector: new TrezorProvider(),
  },
]

const ChooseWallet = () => {
  return (
    <PageWrapper title="Connect Wallet">
      <Tile title="Choose a wallet type.">
        <ul className="wallets-list">{WALLETS.map(renderWallet)}</ul>
      </Tile>
    </PageWrapper>
  )
}

const renderWallet = (wallet) => <Wallet key={wallet.label} {...wallet} />

const Wallet = ({
  label,
  icon,
  providerName,
  type,
  description,
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
      case "LEDGER":
        return <LedgerModal {...defaultProps} />
      case "TREZOR":
        return <TrezorModal {...defaultProps} />
      case "METAMASK":
        return (
          <MetaMaskModal
            {...defaultProps}
            walletName={label}
            providerName={providerName}
            icon={icon}
            {...modalProps}
          />
        )
      default:
        return null
    }
  }

  return (
    <li
      title={providerName === "COINBASE" ? "Coinbase not yet supported" : ""}
      className={`wallet${providerName === "COINBASE" ? " disabled" : ""}`}
      onClick={async () => {
        if (providerName === "COINBASE") {
          return
        }
        openModal(renderModalContent(), {
          title: "Connect Wallet",
          closeModal: customCloseModal,
        })
      }}
    >
      {icon}
      <div className="flex row center">
        <h4 className="mr-1">{label}</h4>
        <Icons.ArrowRight />
      </div>
      <h5 className="wallet-type">{type}</h5>
      <div className="wallet-description">{description}</div>
    </li>
  )
}

export default ChooseWallet
