import React, { useContext, useState } from 'react'
import Button from './Button'
import AddressShortcut from './AddressShortcut'
import rewardsService from '../services/rewards.service'
import { Web3Context } from './WithWeb3Context'
import { useShowMessage, MessagesContext, messageType, useCloseMessage } from './Message'

const useWithdrawAction = (groupIndex, membersIndeces) => {
  const web3Context = useContext(Web3Context)
  const showMessage = useShowMessage(MessagesContext)
  const closeMessage = useCloseMessage(MessagesContext)
  const [isFetching, setIsFetching] = useState(false)

  const withdraw = async () => {
    try {
      setIsFetching(true)
      const message = showMessage({ type: messageType.PENDING_ACTION, sticky: true, title: 'Withdraw acion is pending' })
      const result = await rewardsService.withdrawRewardFromGroup(groupIndex, membersIndeces, web3Context)
      setIsFetching(false)
      closeMessage(message)
      const errorTransactionCount = result.filter((transaction) => transaction.isError).length

      if (errorTransactionCount === 0) {
        showMessage({ type: messageType.SUCCESS, title: 'Withdraw action has been successfully completed' })
      } else if (errorTransactionCount === result.length) {
        showMessage({ type: messageType.ERROR, title: 'Withdraw action has been failed' })
      } else {
        showMessage({ type: messageType.ERROR, title: `${errorTransactionCount} of ${result.length} transactionshave been failed` })
      }
    } catch (error) {
      showMessage({ type: messageType.ERROR, title: 'Something goes wrong...' })
      setIsFetching(false)
    }
  }

  return [isFetching, withdraw]
}


export const RewardsGroupItem = ({ groupIndex, groupPublicKey, membersIndeces, reward }) => {
  const [isFetching, withdraw] = useWithdrawAction(groupIndex, membersIndeces)

  return (
    <li className='group-item'>
      <div className='group-key'>
        <AddressShortcut address={groupPublicKey} />
        <span>GROUP PUBLIC KEY</span>
      </div>
      <div className='group-reward'>
        <span className='reward-value'>
          {reward.toString()}
        </span>
        <span className='reward-currency'>ETH</span>
      </div>
      <Button
        className='btn btn-primary'
        onClick={withdraw}
        isFetching={isFetching}
      >
        WITHDRAW
      </Button>
    </li>
  )
}
