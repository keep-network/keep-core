import React, { useContext } from 'react'
import { Web3Context } from './WithWeb3Context'
import { SubmitButton } from './Button'
import AddressShortcut from './AddressShortcut'
import rewardsService from '../services/rewards.service'
import { useShowMessage, messageType, useCloseMessage } from './Message'

export const RewardsGroupItem = ({ group, updateGroupsAfterWithdrawal }) => {
  const { groupPublicKey, reward } = group
  const [withdrawAction] = useWithdrawAction(group)

  const withdraw = async () => {
    await withdrawAction(group)
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
        <SubmitButton
          className='btn btn-primary btn-lg flex-1'
          onSubmitAction={withdraw}
          successCallback={updateGroupsAfterWithdrawal}
        >
          withdraw
        </SubmitButton>
      </div>
    </li>
  )
}

const useWithdrawAction = (group) => {
  const { groupIndex, membersIndeces } = group
  const web3Context = useContext(Web3Context)
  const showMessage = useShowMessage()
  const closeMessage = useCloseMessage()

  const withdraw = async () => {
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
