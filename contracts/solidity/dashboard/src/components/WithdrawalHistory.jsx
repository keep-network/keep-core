import React, { useState }from 'react'
import rewardsService from '../services/rewards.service'
import { useFetchData } from '../hooks/useFetchData'
import { WithdrawalHistoryItem } from './WithdrawalHistoryItem'
import { SeeAllButton } from './SeeAllButton'
import { LoadingOverlay } from './Loadable'

export const WithdrawalHistory = (props) => {
  const { isFetching, data } = useFetchData(rewardsService.fetchWithdrawalHistory, [])
  const [showAll, setShowAll] = useState(false)

  return (
    <LoadingOverlay isFetching={isFetching} >
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
    </LoadingOverlay>
  )
}
