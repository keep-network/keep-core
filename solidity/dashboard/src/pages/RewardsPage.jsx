import React from 'react'
import { RewardsGroups } from '../components/RewardsGroups'
import { WithdrawalHistory } from '../components/WithdrawalHistory'

const RewardsPage = () => {
  return (
    <>
      <h2 className="mb-2">My Rewards</h2>
      <RewardsGroups />
      <WithdrawalHistory />
    </>
  )
}

export default React.memo(RewardsPage)
