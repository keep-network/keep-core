import React from "react"
import * as Icons from "./Icons"
import { useWeb3Context } from "./WithWeb3Context"
import { useModal } from "../hooks/useModal"
import {
  TrezorConnector,
  LedgerConnector,
  LEDGER_DERIVATION_PATHS,
  injected,
  WalletConnectConnector,
} from "../connectors"
import { MODAL_TYPES, WALLETS } from "../constants/constants"
import { ExplorerModeConnector } from "../connectors/explorer-mode-connector"

const WALLETS_OPTIONS = [
  {
    label: "MetaMask",
    icon: Icons.MetaMask,
    isHardwareWallet: false,
    connector: injected,
    modalType: MODAL_TYPES.MetaMask,
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
    modalType: MODAL_TYPES.Ledger,
  },
  {
    label: "Trezor",
    icon: Icons.Trezor,
    isHardwareWallet: true,
    connector: new TrezorConnector(),
    modalType: MODAL_TYPES.Trezor,
  },
  {
    label: "WalletConnect",
    icon: Icons.WalletConnect,
    isHardwareWallet: false,
    connector: new WalletConnectConnector(),
    modalType: MODAL_TYPES.WalletConnect,
  },
  {
    label: "Explorer Mode",
    icon: Icons.Explore,
    isHardwareWallet: false,
    connector: new ExplorerModeConnector(),
    modalType: MODAL_TYPES.ExplorerMode,
  },
]

const WalletOptions = ({ displayExplorerMode = true }) => {
  return (
    <ul className="wallet__options">
      {WALLETS_OPTIONS.filter(
        (wallet) =>
          !(
            !displayExplorerMode && wallet.label === WALLETS.EXPLORER_MODE.label
          )
      ).map(renderWallet)}
    </ul>
  )
}

const renderWallet = (wallet) => <Wallet key={wallet.label} {...wallet} />

const Wallet = ({ label, icon: IconComponent, connector, modalType }) => {
  const { openModal } = useModal()
  const { connectAppWithWallet } = useWeb3Context()

  const openWalletModal = () => {
    openModal(modalType, {
      connector,
      connectAppWithWallet,
    })
  }

  return (
    <li className="wallet__item" onClick={openWalletModal}>
      <IconComponent className="wallet__item__icon" />
      {label}
    </li>
  )
}

export default WalletOptions
