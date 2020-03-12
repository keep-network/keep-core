import React from 'react'
import AddressShortcut from './AddressShortcut'

export const WithdrawalHistoryItem = ({ date, groupPublicKey, amount }) => (
  <li className='withdrawal-history-item'>
    <div className="flex flex-1 text-smaller text-grey-70">
      {date}
    </div>
    <div className="flex flex-1">
      <AddressShortcut address={groupPublicKey} classNames="text-smaller" />
    </div>
    <div className="flex flex-2 text-smaller text-grey-70">
      {amount.toString()} ETH
    </div>
  </li>
)
