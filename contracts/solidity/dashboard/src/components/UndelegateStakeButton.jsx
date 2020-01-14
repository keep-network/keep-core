import React, { useContext } from 'react'
import { Web3Context } from './WithWeb3Context'
import { useShowMessage } from './Message'
import { SubmitButton } from './Button'

const UndelegateStakeButton = (props) => {
  const web3 = useContext(Web3Context)
  const showMessage = useShowMessage()

  const undelegate = async (onTransactionHashCallback) => {
    const { amount, operator } = props

    try {
      await web3.stakingContract.methods.initiateUnstake(amount, operator).send({ from: web3.yourAddress }).on('transactionHash', onTransactionHashCallback)
      showMessage({ type: 'success', title: 'Success', content: 'Undelegate transaction successfully completed' })
    } catch (error) {
      showMessage({ type: 'error', title: 'Undelegate action has been failed ', content: error.message })
    }
  }

  return (
    <SubmitButton
      className="btn btn-primary btn-sm"
      onSubmitAction={undelegate}
      pendingMessageTitle='Undelegate transaction is pending...'
    >
      Undelegate
    </SubmitButton>
  )
}

export default UndelegateStakeButton
