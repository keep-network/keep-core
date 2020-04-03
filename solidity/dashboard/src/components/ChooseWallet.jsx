import React from 'react'
import * as Icons from './Icons'
import { useWeb3Context } from './WithWeb3Context'

// TODO change icons
const WALLETS = [
  {
    label: 'MetaMask',
    iconSrc: '/images/metamask.svg',
    providerName: 'METAMASK',
    type: 'web wallet',
    description: 'Use the MetaMask extension to access my wallet. Don’t have it? Download here.',
  },
  {
    label: 'Coinbase',
    iconSrc: '/images/trezor.svg',
    providerName: 'COINBASE',
    type: 'web wallet',
    description: 'Connect to your Coinbase wallet.',
  },
  {
    label: 'Ledger',
    iconSrc: '/images/trezor.svg',
    providerName: 'LEDGER',
    type: 'hardware wallet',
    description: 'You will need your Ledger device plugged into your USB drive.',
  },
  {
    label: 'Trezor',
    iconSrc: '/images/trezor.svg',
    providerName: 'TREZOR',
    type: 'hardware wallet',
    description: 'You will need your Trezor device plugged into your USB drive.',
  },
]

const ChooseWallet = () => {
  return (
    <>
      <h1 className="mb-1">Connect My Wallet</h1>
      <section className="tile">
        <h3>
          Choose a wallet below to get started.
        </h3>
        <div className="text-big text-grey-70">
          You’ll need a balance of KEEP tokens to stake with the token dashboard. Don’t have any? Request some.
        </div>
        <ul className="wallets-list" >
          {WALLETS.map(renderWallet)}
        </ul>
      </section>
    </>
  )
}

const renderWallet = (wallet) => <Wallet key={wallet.label} {...wallet} />

const Wallet = ({
  label,
  iconSrc,
  providerName,
  type,
  description,
}) => {
  const { connectAppWithWallet } = useWeb3Context()

  return (
    <li className="wallet" onClick={() => connectAppWithWallet(providerName)}>
      <img src={iconSrc} alt={label} />
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
  )
}

export default ChooseWallet
