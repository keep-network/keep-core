import React, { useEffect } from 'react'
import * as Icons from './Icons'
import { useWeb3Context } from './WithWeb3Context'
import { useModal } from '../hooks/useModal'
import SelectedWalletModal from './SelectedWalletModal'
import { LoadingOverlay } from './Loadable'

// TODO change icons
const WALLETS = [
  {
    label: 'MetaMask',
    icon: <Icons.MetaMask />,
    providerName: 'METAMASK',
    type: 'web wallet',
    description: 'Crypto wallet that’s a web browser plugin.',
    modalProps: {
      iconDescription: null,
      btnText: 'install extension',
      btnLink: 'https://metamask.io',
      description: 'The MetaMask login screen will open in an external window. If it doesn’t load right away, click below to install:',
    },
  },
  {
    label: 'Coinbase',
    icon: <Icons.Coinbase />,
    providerName: 'COINBASE',
    type: 'web wallet',
    description: 'Crypto wallet that’s a web browser plugin.',
    modalProps: {
      iconDescription: null,
      btnText: null,
      btnLink: null,
      description: 'Scan QR code to connect:',
    },
  },
  {
    label: 'Ledger',
    icon: <Icons.Ledger />,
    providerName: 'LEDGER',
    type: 'hardware wallet',
    description: 'Crypto wallet on a secure hardware device.',
    modalProps: {
      iconDescription: '/images/ledger-device.svg',
      btnText: 'install ledger live',
      btnLink: 'https://www.ledger.com/ledger-live/',
      description: 'Plug in Ledger device. Install Ledger Live below:',
    },
  },
  {
    label: 'Trezor',
    icon: <Icons.Trezor />,
    providerName: 'TREZOR',
    type: 'hardware wallet',
    description: 'Crypto wallet on a secure hardware device.',
    modalProps: {
      iconDescription: '/images/trezor-device.svg',
      btnText: 'go to trezor setup',
      btnLink: 'https://trezor.io/start/',
      description: 'Plug in your Trezor device. If the setup screen doesn’t load right away, go to Trezor setup:',
    },
  },
]

const ChooseWallet = () => {
  const { isFetching } = useWeb3Context()

  return (
    <LoadingOverlay isFetching={isFetching}>
      <h1 className="mb-1">Connect Wallet</h1>
      <section className="tile">
        <h3>
          Choose a wallet type.
        </h3>
        <ul className="wallets-list" >
          {WALLETS.map(renderWallet)}
        </ul>
      </section>
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
  const { ModalComponent, showModal, closeModal } = useModal()
  const { connectAppWithWallet, isFetching, web3, provider } = useWeb3Context()

  useEffect(() => {
    if (!isFetching && web3 && provider) {
      closeModal()
    }
  }, [isFetching, web3, provider])

  return (
    <>
      <ModalComponent title='Connect Wallet'>
        <SelectedWalletModal
          walletName={label}
          icon={icon}
          {...modalProps}
        />
      </ModalComponent>
      <li
        className="wallet"
        onClick={() => {
          showModal()
          connectAppWithWallet(providerName)
        }}
      >
        {icon}
        <div className="flex row center">
          <h4 className="mr-1">{label}</h4>
          <Icons.ArrowRight />
        </div>
        <h5 className="wallet-type">
          {type}
        </h5>
        <div className="wallet-description">
          {description}
        </div>
      </li>
    </>
  )
}

export default ChooseWallet
