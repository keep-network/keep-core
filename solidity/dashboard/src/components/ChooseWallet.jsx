import React, { useState, useCallback } from "react"
import * as Icons from "./Icons"
import { useWeb3Context } from "./WithWeb3Context"
import { useModal } from "../hooks/useModal"
import SelectedWalletModal from "./SelectedWalletModal"
import { LoadingOverlay } from "./Loadable"
import Tile from "./Tile"
import PageWrapper from "./PageWrapper"
import LedgerModal from "./LedgerModal"
import TrezorModal from "./TrezorModal"

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
  },
  {
    label: "Ledger",
    icon: <Icons.Ledger />,
    providerName: "LEDGER",
    type: "hardware wallet",
    description: "Crypto wallet on a secure hardware device.",
  },
  {
    label: "Trezor",
    icon: <Icons.Trezor />,
    providerName: "TREZOR",
    type: "hardware wallet",
    description: "Crypto wallet on a secure hardware device.",
  },
]

const ChooseWallet = () => {
  const { isFetching } = useWeb3Context()

  return (
    <LoadingOverlay isFetching={isFetching}>
      <PageWrapper title="Connect Wallet">
        <Tile title="Choose a wallet type.">
          <ul className="wallets-list">{WALLETS.map(renderWallet)}</ul>
        </Tile>
      </PageWrapper>
    </LoadingOverlay>
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
}) => {
  const { ModalComponent, openModal, closeModal } = useModal()
  const {
    connectAppWithWallet,
    setAccount,
    abortWalletConnection,
  } = useWeb3Context()
  const [accounts, setAccounts] = useState(null)

  const onSelectProvider = async (providerName) => {
    const firstAccountAsSelected = providerName === "METAMASK"
    const availableAccounts = await connectAppWithWallet(
      providerName,
      firstAccountAsSelected
    )
    setAccounts(availableAccounts)
  }

  const onSelectAccount = (account) => {
    setAccount([account])
    closeModal()
  }

  const customCloseModal = useCallback(() => {
    abortWalletConnection()
    closeModal()
  }, [abortWalletConnection, closeModal])

  const renderModalContent = () => {
    switch (providerName) {
      case "LEDGER":
        return (
          <LedgerModal
            onSelectProvider={onSelectProvider}
            onSelectAccount={onSelectAccount}
            accounts={accounts}
          />
        )
      case "TREZOR":
        return (
          <TrezorModal accounts={accounts} onSelectAccount={onSelectAccount} />
        )
      case "METAMASK":
      case "COINBASE":
      default:
        return (
          <SelectedWalletModal walletName={label} icon={icon} {...modalProps} />
        )
    }
  }

  return (
    <>
      <ModalComponent title="Connect Wallet" closeModal={customCloseModal}>
        {renderModalContent()}
      </ModalComponent>
      <li
        title={providerName === "COINBASE" && "Coinbase not yet supported"}
        className={`wallet${providerName === "COINBASE" ? " disabled" : ""}`}
        onClick={async () => {
          if (providerName === "COINBASE") {
            return
          }
          openModal()
          if (providerName === "LEDGER") {
            return
          }
          await onSelectProvider(providerName)
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
    </>
  )
}

export default ChooseWallet
