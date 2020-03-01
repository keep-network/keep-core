import React, { useState } from 'react'
import { RewardsGroups } from './RewardsGroups'
import { WithdrawalHistory } from './WithdrawalHistory'

export const Rewards = () => {
  const [totalRewardsBalance, setTotalRewardsBalance] = useState('0')

  return (
    <div className="rewards-wrapper flex row center">
      <div className="rewards-history flex column">
        <RewardsBalance balance={totalRewardsBalance} />
        <WithdrawalHistory />
      </div>
      <RewardsGroups setTotalRewardsBalance={setTotalRewardsBalance} />
    </div>
  )
}

const RewardsBalance = ({ balance }) => (
  <div className='rewards-balance tile'>
    <h2>{balance} ETH</h2>
    <h6>YOUR REWARDS BALANCE</h6>
  </div>
)
