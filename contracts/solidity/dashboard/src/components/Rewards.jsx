import React from 'react'
import { RewardsGroups } from './RewardsGroups'
import rewardsService from '../services/rewards.service'
import { WithdrawalHistory } from './WithdrawalHistory'
import { useFetchData } from '../hooks/useFetchData'

// TODO implement update date hook when group reward are withdrawn
const useUpdateGroups = (data) => {
  const updateData = (groupPublicKey) => {

  }
}


export const Rewards = () => {
  const { isFetching, data: [groups, totalRewardsBalance] } = useFetchData(rewardsService.fetchAvailableRewards, [[], '0'])

  return (
    <div className="rewards-wrapper flex flex-row flex-row-center">
      <div className="rewards-history flex flex-column">
        <RewardsBalance balance={totalRewardsBalance} />
        <WithdrawalHistory />
      </div>
      <RewardsGroups isFetching={isFetching} groups={groups} />
    </div>
  )
}

const RewardsBalance = ({ balance }) => (
  <div className='rewards-balance tile'>
    <h2>{balance} ETH</h2>
    <h6>YOUR REWARDS BALANCE</h6>
  </div>
)
