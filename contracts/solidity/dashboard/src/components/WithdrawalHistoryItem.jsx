import React from 'react'
import AddressShortcut from './AddressShortcut'

export const WithdrawalHistoryItem = ({ date, groupPublicKey, amount }) => (
  <li className='withdrawal-history-item'>
    <span className="small-title">
      {date}
    </span>
    <AddressShortcut address={groupPublicKey} />
    <span className="small-title">{amount.toString()} ETH</span>
  </li>
)
