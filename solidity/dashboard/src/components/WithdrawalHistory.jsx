import React, { useState, useContext } from 'react'
import { Web3Context } from './WithWeb3Context'
import { SeeAllButton } from './SeeAllButton'
import { LoadingOverlay } from './Loadable'
import { useFetchData } from '../hooks/useFetchData'
import rewardsService from '../services/rewards.service'
import { useSubscribeToContractEvent } from '../hooks/useSubscribeToContractEvent'
import { OPERATOR_CONTRACT_NAME } from '../constants/constants'
import web3Utils from 'web3-utils'
import { formatDate, isSameEthAddress } from '../utils/general.utils'
import { DataTable, Column } from './DataTable'
import AddressShortcut from './AddressShortcut'

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

  console.log('data history',data)

  return (
    <LoadingOverlay isFetching={isFetching} >
      <section className="tile">
        <h5 className="mb-1 text-grey-50">Rewards History</h5>
        <DataTable data={showAll ? data : data.slice(0, previewDataCount)}>
          <Column
            header="amount"
            field="amount"
            renderContent={({ amount }) => `${amount.toString()} ETH`}
          />
          <Column
            header="date"
            field="date"
          />
          <Column
            header="group key"
            field="groupPublicKey"
            renderContent={({ groupPublicKey }) => <AddressShortcut address={groupPublicKey} classNames="text-smaller" />}
          />
        </DataTable>
        <SeeAllButton
          dataLength={data.length}
          previewDataCount={previewDataCount}
          onClickCallback={() => setShowAll(!showAll)}
          showAll={showAll}
        />
      </section>
    </LoadingOverlay>
  )
}
