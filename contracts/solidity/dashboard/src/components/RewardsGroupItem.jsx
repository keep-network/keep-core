import React from 'react'
import Button from './Button'
import AddressShortcut from './AddressShortcut'

export const RewardsGroupItem = (props) => {
  return (
    <div className='group-item'>
      <span>
        Group index: 1
      </span>
      <h4>
        Group public key:&nbsp;
        <AddressShortcut address={'0xcCFe2E36B3F10152D19dD7d14d651F213c9af4b0'} />
      </h4>
      <p>300 ETH</p>
      <Button className="btn btn-primary btn-large">
        Withdraw
      </Button>
    </div>
  )
}
