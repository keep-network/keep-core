import React, { useContext } from 'react'
import * as Icons from './Icons'
import { Web3Context } from './WithWeb3Context'
import { Web3Status } from './Web3Status'
import { MenuButton } from './MenuButton'
import AddressShortcut from './AddressShortcut'
import { NetworkStatus } from './NetworkStatus'

const Header = (props) => {
  const { yourAddress } = useContext(Web3Context)

  return (
    <header className='header'>
      <a href="/" className='logo'><Icons.Keep width='250px' height='80px'/></a>
      <Web3Status />
      <div className='account-address'>
        <h5 className="text-grey-50">
          <span>ADDRESS&nbsp;</span>
          <AddressShortcut classNames="text-small" address={yourAddress} />
        </h5>
        <NetworkStatus />
      </div>
      <MenuButton />
    </header>
  )
}

export default Header
