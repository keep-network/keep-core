import React from 'react'
import AddressShortcut from './AddressShortcut'

export const WithdrawalHistoryItem = ({ date, groupPublicKey, amount }) => (
  <li className='withdrawal-history-item'>
    <span className="text-small">
      {date}
    </span>
    <AddressShortcut address={groupPublicKey} />
    <span className="text-small">{amount.toString()} ETH</span>
  </li>
)
