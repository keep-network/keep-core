import React from 'react'
import Button from './Button'
import AddressShortcut from './AddressShortcut'

export const RewardsGroupItem = ({ groupIndex, groupPubKey, reward, isStale }) => {
  return (
    <div className='group-item'>
      <div className='group-key'>
        <AddressShortcut address={groupPubKey} />
        <span>GROUP PUBLIC KEY</span>
      </div>
      <div className='group-reward'>
        <span className='reward-value'>
          {reward.toString()}
        </span>
        <span className='reward-currency'>ETH</span>
      </div>
      <Button
        className='btn btn-primary btn-sm'
        disabled={!isStale}
      >
        Withdraw
      </Button>
    </div>
  )
}
