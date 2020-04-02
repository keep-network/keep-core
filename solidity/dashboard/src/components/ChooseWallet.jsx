import React from 'react'
import * as Icons from './Icons'

const WALLETS = [
  { label: 'MetaMask', iconSrc: '/images/metamask.svg' },
  { label: 'Trezor', iconSrc: '/images/trezor.svg' },
]

const ChooseWallet = () => {
  return (
    <React.Fragment>
      <h3>
        Connect a wallet to get started:
      </h3>
      <ul className="wallets-list" >
        {WALLETS.map(renderWallet)}
      </ul>
    </React.Fragment>
  )
}

const renderWallet = (wallet) => <Wallet {...wallet} />

const Wallet = ({ label, iconSrc }) => {
  return (
    <li key={label} className="wallet">
      <img src={iconSrc} alt={label} />
      <h4>{label}</h4>
      <Icons.ArrowRight />
    </li>
  )
}

export default ChooseWallet
