import React from 'react'
import AddressShortcut from './AddressShortcut'

export const WithdrawalHistoryItem = ({ date, groupPublicKey, amount }) => (
  <li className='group-item'>
    <div>
      {date}
    </div>
    <div className='group-key'>
      <AddressShortcut address={groupPublicKey} />
    </div>
    <div className>
      {amount.toString()} ETH
    </div>
  </li>
)
