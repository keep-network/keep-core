import React, { useState, useContext } from 'react'
import { SeeAllButton } from './SeeAllButton'
import { LoadingOverlay } from './Loadable'
import { useFetchData } from '../hooks/useFetchData'
import rewardsService from '../services/rewards.service'
import Dropdown from './Dropdown'
import { DataTable, Column } from './DataTable'
import AddressShortcut from './AddressShortcut'
import SubmitButton from './Button'
import { useShowMessage, messageType, useCloseMessage } from './Message'
import { Web3Context } from './WithWeb3Context'
import { findIndexAndObject } from '../utils/array.utils'

const previewDataCount = 3
const initialData = [[], '0']

export const RewardsGroups = React.memo(({ setTotalRewardsBalance }) => {
  const [state, updateData] = useFetchData(rewardsService.fetchAvailableRewards, initialData)
  const { isFetching, data: [groups, totalRewardsBalance] } = state
  const [showAll, setShowAll] = useState(false)
  const [selectedReward, setSelectedReward] = useState({})
  const [withdrawAction] = useWithdrawAction()

  const setWithdrawalStatus = (group) => {
    const { groupIndex } = group
    const { indexInArray, obj } = findIndexAndObject('groupIndex', groupIndex, groups)
    const updateGroups = [...groups]
    updateGroups[indexInArray] = { ...obj, status: 'PENDING' }

    updateData([updateGroups, totalRewardsBalance])
  }

  const updateGroupsAfterWithdrawal = () => {
    
  }

  return (
    <>
    <LoadingOverlay isFetching={isFetching} >
      <section className="tile total-rewards-section">
        <div className="total-rewards-balance">
          <h3 className='text-grey-70 pb-2'>Total Balance</h3>
          <h2 className="balance">{`${totalRewardsBalance} ETH`}</h2>
        </div>
        <div className="withdraw-dropdown">
          <h4 className="text-grey-70 text-normal">Withdraw</h4>
          <Dropdown
            options={[]}
            onSelect={(reward) => setSelectedReward(reward)}
            valuePropertyName='address'
            labelPropertyName='address'
            selectedItem={selectedReward}
            labelPrefix='Operator:'
            noItemSelectedText='Select Operator'
            label=''
          />
          <SubmitButton
            className='btn btn-primary btn-lg flex-1'
            onSubmitAction={() => withdrawAction(selectedReward, setWithdrawalStatus)}
            successCallback={updateGroupsAfterWithdrawal}
          >
            withdraw
          </SubmitButton>
        </div>
      </section>
    </LoadingOverlay>
      <LoadingOverlay isFetching={isFetching} classNames='group-items self-start'>
        <section className="group-items tile">
          <h3 className='text-grey-70 mb-2'>Totals</h3>
          <DataTable data={showAll ? groups : groups.slice(0, previewDataCount)} >
            <Column
              header="amount"
              field="reward"
              renderContent={({ reward }) => `${reward.toString()} ETH`}
            />
            <Column
              header="status"
              field="status"
              renderContent={({ status }) => `status here`}
            />
            <Column
              header="group key"
              field="groupPubKey"
              renderContent={({ groupPubKey }) => <AddressShortcut address={groupPubKey} /> }
            />
          </DataTable>
          <SeeAllButton
            dataLength={groups.length}
            previewDataCount={previewDataCount}
            onClickCallback={() => setShowAll(!showAll)}
            showAll={showAll}
          />
        </section>
      </LoadingOverlay>
    </>
  )
})

const useWithdrawAction = () => {
  const web3Context = useContext(Web3Context)
  const showMessage = useShowMessage()
  const closeMessage = useCloseMessage()

  const withdraw = async (group) => {
    const { groupIndex, membersIndeces } = group
    try {
      const message = showMessage({ type: messageType.PENDING_ACTION, sticky: true, title: 'Withdrawal in progress' })
      const result = await rewardsService.withdrawRewardFromGroup(groupIndex, membersIndeces, web3Context)
      closeMessage(message)
      const unacceptedTransactions = result.filter((reward) => reward.isError)
      const errorTransactionCount = unacceptedTransactions.length

      if (errorTransactionCount === 0) {
        showMessage({ type: messageType.SUCCESS, title: 'Reward withdrawal completed' })
      } else if (errorTransactionCount === result.length) {
        throw new Error('Reward withdrawal failed')
      } else {
        showMessage({ type: messageType.INFO, title: `${errorTransactionCount} of ${result.length} transactions have been not approved` })
      }
    } catch (error) {
      showMessage({ type: messageType.ERROR, title: error.message })
      throw error
    }
  }

  return [withdraw]
}
