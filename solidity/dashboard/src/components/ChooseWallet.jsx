import React from 'react'
import * as Icons from './Icons'
import { useWeb3Context } from './WithWeb3Context'

// TODO change icons
const WALLETS = [
  {
    label: 'MetaMask',
    icon: <Icons.MetaMask />,
    providerName: 'METAMASK',
    type: 'web wallet',
    description: 'Crypto wallet that’s a web browser plugin.',
  },
  {
    label: 'Coinbase',
    icon: <Icons.Coinbase />,
    providerName: 'COINBASE',
    type: 'web wallet',
    description: 'Crypto wallet that’s a web browser plugin.',
  },
  {
    label: 'Ledger',
    icon: <Icons.Ledger />,
    providerName: 'LEDGER',
    type: 'hardware wallet',
    description: 'Crypto wallet on a secure hardware device.',
  },
  {
    label: 'Trezor',
    icon: <Icons.Trezor />,
    providerName: 'TREZOR',
    type: 'hardware wallet',
    description: 'Crypto wallet on a secure hardware device.',
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
  icon,
  providerName,
  type,
  description,
}) => {
  const { connectAppWithWallet } = useWeb3Context()

  return (
    <li className="wallet" onClick={() => connectAppWithWallet(providerName)}>
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
  )
}

export default ChooseWallet
