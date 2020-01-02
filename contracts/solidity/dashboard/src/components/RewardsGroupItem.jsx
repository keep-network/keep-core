import React from 'react'
import Button from './Button'
import AddressShortcut from './AddressShortcut'

export const RewardsGroupItem = ({ groupIndex, groupPubKey, reward, isStale }) => {
  return (
    <div className='group-item'>
      <span>
        Group index: {groupIndex}
      </span>
      <h4>
        Group public key:&nbsp;
        <AddressShortcut address={groupPubKey} />
      </h4>
      <p>{reward.toString()} ETH</p>
      <Button
        className="btn btn-primary btn-large"
        disabled={!isStale}
      >
        Withdraw
      </Button>
    </div>
  )
}
