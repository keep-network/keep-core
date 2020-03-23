import React, { useState } from 'react'
import { RewardsGroups } from '../components/RewardsGroups'
import { WithdrawalHistory } from '../components/WithdrawalHistory'

const RewardsPage = () => {
  const [totalRewardsBalance, setTotalRewardsBalance] = useState('0')

  return (
    <>
      <h2 className="mb-2">My Rewards</h2>
      <div className="rewards-wrapper flex row center">
        <section className="rewards-history tile flex column">
          <h3 className="text-grey-70 mb-1">Rewards</h3>
          <h2 className="balance">{totalRewardsBalance} ETH</h2>
          <hr/>
          <WithdrawalHistory />
        </section>
        <RewardsGroups setTotalRewardsBalance={setTotalRewardsBalance} />
      </div>
    </>
  )
}

export default React.memo(RewardsPage)
