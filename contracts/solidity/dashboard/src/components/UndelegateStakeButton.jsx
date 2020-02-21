import React, { useContext } from 'react'
import { Web3Context } from './WithWeb3Context'
import { useShowMessage, messageType } from './Message'
import { SubmitButton } from './Button'

const UndelegateStakeButton = (props) => {
  const { yourAddress, stakingContract } = useContext(Web3Context)
  const showMessage = useShowMessage()

  const undelegate = async (onTransactionHashCallback) => {
    const { operator, isInInitializationPeriod } = props

    try {
      await stakingContract
        .methods[isInInitializationPeriod ? 'cancelStake' : 'undelegate'](operator)
        .send({ from: yourAddress })
        .on('transactionHash', onTransactionHashCallback)
      showMessage({ type: messageType.SUCCESS, title: 'Success', content: 'Undelegate transaction successfully completed' })
    } catch (error) {
      showMessage({ type: messageType.ERROR, title: 'Undelegate action has been failed ', content: error.message })
    }
  }

  return (
    <SubmitButton
      className={props.btnClassName}
      onSubmitAction={undelegate}
      pendingMessageTitle='Undelegate transaction is pending...'
    >
      {props.isInInitializationPeriod ? 'cancel' :props.btnText }
    </SubmitButton>
  )
}

UndelegateStakeButton.defaultProps = {
  btnClassName: 'btn btn-primary btn-sm',
  btnText: 'undelegate',
  isInInitializationPeriod: false,
}

export default UndelegateStakeButton
