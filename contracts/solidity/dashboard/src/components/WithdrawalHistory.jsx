import React, { useState } from 'react'
import { WithdrawalHistoryItem } from './WithdrawalHistoryItem'
import { SeeAllButton } from './SeeAllButton'
import { LoadingOverlay } from './Loadable'
import { useFetchData } from '../hooks/useFetchData'
import rewardsService from '../services/rewards.service'

const previewDataCount = 3

export const WithdrawalHistory = (props) => {
  const [state] = useFetchData(rewardsService.fetchWithdrawalHistory, [])
  const { isFetching, data } = state
  const [showAll, setShowAll] = useState(false)

  return (
    <LoadingOverlay isFetching={isFetching} >
      <section className="mt-1">
        <h5 className="mb-1 text-grey-50">Rewards History</h5>
        <div className="flex row center">
          <div className="flex-1 text-label">
            date
          </div>
          <div className="flex-1 text-label">
            group key
          </div>
          <div className="flex-2 text-label">
            amount
          </div>
        </div>
        <ul className="withdrawal-history">
          {showAll ? data.map(renderWithdrawalHistoryItem) : data.slice(0, previewDataCount).map(renderWithdrawalHistoryItem)}
          <SeeAllButton
            dataLength={data.length}
            previewDataCount={previewDataCount}
            onClickCallback={() => setShowAll(!showAll)}
            showAll={showAll}
          />
        </ul>
      </section>
    </LoadingOverlay>
  )
}

const renderWithdrawalHistoryItem = (history, index) => <WithdrawalHistoryItem key={index} {...history} />
