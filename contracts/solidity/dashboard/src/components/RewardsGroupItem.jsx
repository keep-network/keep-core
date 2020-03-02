import React, { useContext, useState } from 'react'
import Button from './Button'
import AddressShortcut from './AddressShortcut'
import rewardsService from '../services/rewards.service'
import { Web3Context } from './WithWeb3Context'
import { useShowMessage, MessagesContext, messageType, useCloseMessage } from './Message'

export const RewardsGroupItem = ({ group, updateGroupsAfterWithdrawal }) => {
  const { groupPublicKey, reward } = group
  const [isFetching, withdrawAction] = useWithdrawAction(group)

  const withdraw = async () => {
    const groupToUpdate = await withdrawAction(group)
    updateGroupsAfterWithdrawal(groupToUpdate)
  }

  return (
    <li className='group-item center'>
      <div className='flex flex-1'>
        <AddressShortcut
          address={groupPublicKey}
          classNames="text-big"
        />
      </div>
      <div className='text-primary text-big flex flex-2'>
        {reward.toString()} ETH
      </div>
      <div className='flex flex-1'>
        <Button
          className='btn btn-lg btn-primary'
          onClick={withdraw}
          isFetching={isFetching}
        >
          withdraw
        </Button>
      </div>
    </li>
  )
}

const useWithdrawAction = (group) => {
  const { groupIndex, membersIndeces } = group
  const web3Context = useContext(Web3Context)
  const { utils } = web3Context
  const showMessage = useShowMessage(MessagesContext)
  const closeMessage = useCloseMessage(MessagesContext)
  const [isFetching, setIsFetching] = useState(false)

  const withdraw = async () => {
    try {
      setIsFetching(true)
      const message = showMessage({ type: messageType.PENDING_ACTION, sticky: true, title: 'Withdrawal in progress' })
      const result = await rewardsService.withdrawRewardFromGroup(groupIndex, membersIndeces, web3Context)
      setIsFetching(false)
      closeMessage(message)
      const unacceptedTransactions = result.filter((reward) => reward.isError)
      const errorTransactionCount = unacceptedTransactions.length

      if (errorTransactionCount === 0) {
        showMessage({ type: messageType.SUCCESS, title: 'Reward withdrawal completed' })
      } else if (errorTransactionCount === result.length) {
        showMessage({ type: messageType.ERROR, title: 'Reward withdrawal failed' })
      } else {
        showMessage({ type: messageType.INFO, title: `${errorTransactionCount} of ${result.length} transactions have been not approved` })
      }

      const updatedMemberIndices = {}
      let updatedRewardInGroup = utils.toBN('0')
      unacceptedTransactions.forEach((transactionDetails) => {
        updatedRewardInGroup = updatedRewardInGroup.add(utils.toBN(group.rewardPerMemberInWei).mul(utils.toBN(transactionDetails.memberIndices.length)))
        updatedMemberIndices[transactionDetails.memberAddress] = transactionDetails.memberIndices
      })

      const groupToUpdate = { ...group, membersIndeces: updatedMemberIndices, reward: utils.fromWei(updatedRewardInGroup, 'ether') }

      return groupToUpdate
    } catch (error) {
      showMessage({ type: messageType.ERROR, title: 'Reward withdrawal failed' })
      setIsFetching(false)
    }
  }

  return [isFetching, withdraw]
}
