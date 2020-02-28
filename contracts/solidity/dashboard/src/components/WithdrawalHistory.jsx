import React, { useState, useContext } from 'react'
import { Web3Context } from './WithWeb3Context'
import { WithdrawalHistoryItem } from './WithdrawalHistoryItem'
import { SeeAllButton } from './SeeAllButton'
import { LoadingOverlay } from './Loadable'
import { useFetchData } from '../hooks/useFetchData'
import rewardsService from '../services/rewards.service'
import { useSubscribeToContractEvent } from '../hooks/useSubscribeToContractEvent'
import { OPERATOR_CONTRACT_NAME_EVENTS } from '../constants/constants'
import web3Utils from 'web3-utils'
import { formatDate, isSameEthAddress } from '../utils'

const previewDataCount = 3

export const WithdrawalHistory = (props) => {
  const [state, updateData] = useFetchData(rewardsService.fetchWithdrawalHistory, [])
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
    OPERATOR_CONTRACT_NAME_EVENTS,
    'GroupMemberRewardsWithdrawn',
    subscribeToEventCallback
  )

  return (
    <LoadingOverlay isFetching={isFetching} >
      <ul className="withdrawal-history tile">
        <h6>Withdrawal History</h6>
        {showAll ? data.map(renderWithdrawalHistoryItem) : data.slice(0, previewDataCount).map(renderWithdrawalHistoryItem)}
        <SeeAllButton
          dataLength={data.length}
          previewDataCount={previewDataCount}
          onClickCallback={() => setShowAll(!showAll)}
          showAll={showAll}
        />
      </ul>
    </LoadingOverlay>
  )
}

const renderWithdrawalHistoryItem = (history, index) => <WithdrawalHistoryItem key={index} {...history} />
