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
import ExplorerModeModal from "./ExplorerModeModal"
import { ExplorerModeConnector } from "../connectors/explorer-mode-connector"

const WALLETS_OPTIONS = [
  {
    label: "MetaMask",
    icon: Icons.MetaMask,
    isHardwareWallet: false,
    connector: new InjectedConnector(),
  },
  {
    label: "Ledger",
    icon: Icons.Ledger,
    isHardwareWallet: true,
    connector: {
      // Add additional filed `name` to handle `switch` statement in
      // `renderModalContent` function in `Wallet` component. The `LedgerModal`
      // component is able to handle this pseudo connector object correctly and
      // choose the right connector from the object depending on user
      // preferences.
      name: WALLETS.LEDGER.name,
      LEDGER_LIVE: new LedgerConnector(LEDGER_DERIVATION_PATHS.LEDGER_LIVE),
      LEDGER_LEGACY: new LedgerConnector(LEDGER_DERIVATION_PATHS.LEDGER_LEGACY),
    },
  },
  {
    label: "Trezor",
    icon: Icons.Trezor,
    isHardwareWallet: true,
    connector: new TrezorConnector(),
  },
  {
    label: "WalletConnect",
    icon: Icons.WalletConnect,
    isHardwareWallet: false,
    connector: new WalletConnectConnector(),
  },
  {
    label: "Explorer Mode",
    icon: Icons.Wallet,
    isHardwareWallet: false,
    connector: new ExplorerModeConnector(),
  },
]

const WalletOptions = ({ displayExplorerMode = true }) => {
  return (
    <ul className="wallet__options">
      {WALLETS_OPTIONS.map((wallet) => {
        // do not display explorer mode option inside WalletSelectionModal
        if (
          !displayExplorerMode &&
          wallet.label === WALLETS.EXPLORER_MODE.label
        )
          return null
        return renderWallet(wallet)
      })}
    </ul>
  )
}

const renderWallet = (wallet) => <Wallet key={wallet.label} {...wallet} />

const Wallet = ({ label, icon: IconComponent, connector }) => {
  const { openModal, closeModal } = useModal()
  const {
    connectAppWithWallet,
    abortWalletConnection,
    connector: currentConnector,
  } = useWeb3Context()

  const customCloseModal = useCallback(() => {
    if (currentConnector?.name !== WALLETS.EXPLORER_MODE.name) {
      abortWalletConnection()
    }
    closeModal()
  }, [abortWalletConnection, closeModal])

  const renderModalContent = () => {
    const defaultProps = { connector, closeModal, connectAppWithWallet }
    switch (connector.name) {
      case WALLETS.LEDGER.name:
        return <LedgerModal {...defaultProps} />
      case WALLETS.TREZOR.name:
        return <TrezorModal {...defaultProps} />
      case WALLETS.METAMASK.name:
        return <MetaMaskModal {...defaultProps} />
      case WALLETS.WALLET_CONNECT.name:
        return <WallectConnectModal {...defaultProps} />
      case WALLETS.EXPLORER_MODE.name:
        return <ExplorerModeModal {...defaultProps} />
      default:
        return null
    }
  }

  const modalTitle =
    label === WALLETS.EXPLORER_MODE.label
      ? "Connect Ethereum Address"
      : "Connect Wallet"

  console.log('LABEL', label)

  return (
    <li
      className="wallet__item"
      onClick={async () => {
        openModal(renderModalContent(), {
          title: modalTitle,
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
