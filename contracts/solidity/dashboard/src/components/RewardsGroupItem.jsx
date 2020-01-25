import React from 'react'
import Button from './Button'
import AddressShortcut from './AddressShortcut'

export const RewardsGroupItem = ({ groupIndex, groupPublicKey, reward, isStale }) => {
  return (
    <li className='group-item'>
      <div className='group-key'>
        <AddressShortcut address={groupPublicKey} />
        <span>GROUP PUBLIC KEY</span>
      </div>
      <div className='group-reward'>
        <span className='reward-value'>
          {reward.toString()}
        </span>
        <span className='reward-currency'>ETH</span>
      </div>
      <Button
        className='btn btn-primary'
        disabled={!isStale}
      >
        WITHDRAW
      </Button>
    </li>
  )
}
