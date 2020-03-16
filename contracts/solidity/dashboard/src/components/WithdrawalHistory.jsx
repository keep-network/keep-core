import React, { useState, useContext } from 'react'
import { Web3Context } from './WithWeb3Context'
import { WithdrawalHistoryItem } from './WithdrawalHistoryItem'
import { SeeAllButton } from './SeeAllButton'
import { LoadingOverlay } from './Loadable'
import { useFetchData } from '../hooks/useFetchData'
import rewardsService from '../services/rewards.service'
import { useSubscribeToContractEvent } from '../hooks/useSubscribeToContractEvent'
import { OPERATOR_CONTRACT_NAME } from '../constants/constants'
import web3Utils from 'web3-utils'
import { formatDate, isSameEthAddress } from '../utils/general.utils'

const previewDataCount = 3
const initialData = []

export const WithdrawalHistory = (props) => {
  const [state, updateData] = useFetchData(rewardsService.fetchWithdrawalHistory, initialData)
  const { isFetching, data } = state
  const [showAll, setShowAll] = useState(false)
  const { yourAddress, eth, keepRandomBeaconOperatorContract } = useContext(Web3Context)

  const subscribeToEventCallback = async (event) => {
    const { blockNumber, returnValues: { groupIndex, amount, beneficiary } } = event
    if (!isSameEthAddress(yourAddress, beneficiary)) {
      return
    }
    const withdrawnAt = (await eth.getBlock(blockNumber)).timestamp
    const groupPublicKey = await keepRandomBeaconOperatorContract.methods.getGroupPublicKey(groupIndex).call()
    const withdrawal = {
      blockNumber,
      groupPublicKey,
      date: formatDate(withdrawnAt * 1000),
      amount: web3Utils.fromWei(amount, 'ether'),
    }
    updateData([withdrawal, ...data])
  }

  useSubscribeToContractEvent(
    OPERATOR_CONTRACT_NAME,
    'GroupMemberRewardsWithdrawn',
    subscribeToEventCallback
  )

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
