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
      const message = showMessage({ type: messageType.PENDING_ACTION, sticky: true, title: 'Withdraw acion is pending' })
      const result = await rewardsService.withdrawRewardFromGroup(groupIndex, membersIndeces, web3Context)
      setIsFetching(false)
      closeMessage(message)
      const unacceptedTransactions = result.filter((reward) => reward.isError)
      const errorTransactionCount = unacceptedTransactions.length

      if (errorTransactionCount === 0) {
        showMessage({ type: messageType.SUCCESS, title: 'Withdraw action has been successfully completed' })
      } else if (errorTransactionCount === result.length) {
        showMessage({ type: messageType.ERROR, title: 'Withdraw action has been failed' })
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
      showMessage({ type: messageType.ERROR, title: 'Something goes wrong...' })
      setIsFetching(false)
    }
  }

  return [isFetching, withdraw]
}
