import React from 'react'
import * as Icons from './Icons'
import { useWeb3Context } from './WithWeb3Context'

const WALLETS = [
  { label: 'MetaMask', iconSrc: '/images/metamask.svg', providerName: 'METAMASK' },
  { label: 'Trezor', iconSrc: '/images/trezor.svg', providerName: 'TREZOR' },
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

const Wallet = ({
  label,
  iconSrc,
  providerName,
}) => {
  const { connectAppWithWallet } = useWeb3Context()

  return (
    <li key={label} className="wallet" onClick={() => connectAppWithWallet(providerName)}>
      <img src={iconSrc} alt={label} />
      <h4>{label}</h4>
      <Icons.ArrowRight />
    </li>
  )
}

export default ChooseWallet
