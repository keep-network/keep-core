import React, { useState }from 'react'
import rewardsService from '../services/rewards.service'
import { useFetchData } from '../hooks/useFetchData'
import { WithdrawalHistoryItem } from './WithdrawalHistoryItem'
import { SeeAllButton } from './SeeAllButton'

export const WithdrawalHistory = (props) => {
  const { isFetching, data } = useFetchData(rewardsService.fetchWithdrawalHistory, [])
  const [showAll, setShowAll] = useState(false)

  return (
    <ul className="withdrawal-history tile">
      <h6>Withdrawal History</h6>
      {data.map((history) => <WithdrawalHistoryItem key={history.groupPubKey} {...history} /> )}
      <SeeAllButton
        dataLength={data.length}
        previewDataCount={3}
        onClickCallback={() => setShowAll(!showAll)}
        showAll={showAll}
      />
    </ul>
  )
}
