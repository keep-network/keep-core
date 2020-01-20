import React from 'react'
import rewardsService from '../services/rewards.service'
import { useFetchData } from '../hooks/useFetchData'
import { WithdrawalHistoryItem } from './WithdrawalHistoryItem'

export const WithdrawalHistory = (props) => {
  const { isFetching, data } = useFetchData(rewardsService.fetchWithdrawalHistory, [])

  return (
    <ul className="withdrawals-history">
      Withdrawal History
      {data.map((history) => <WithdrawalHistoryItem key={history.groupPubKey} {...history} /> )}
    </ul>
  )
}
