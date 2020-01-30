import React, { useCallback, useContext, useState } from 'react'
import InlineForm from './InlineForm'
import { useShowMessage, messageType } from './Message'
import { Web3Context } from './WithWeb3Context'
import { formatAmount } from '../utils'

const UndelegateForm = (props) => {
  const showMessage = useShowMessage()
  const { stakingContract, yourAddress, utils } = useContext(Web3Context)
  const [amount, setAmount] = useState('')

  const onSubmitAction = useCallback(async (onTransactionHashCallback) => {
    try {
      await stakingContract.methods.initiateUnstake(utils.toBN(formatAmount(amount)).toString(), yourAddress)
        .send({ from: yourAddress })
        .on('transactionHash', onTransactionHashCallback)

      setAmount('')
    } catch (error) {
      showMessage({ type: messageType.ERROR, title: `Undelegate action has been failed`, content: error.message })
    }
  }, [amount, yourAddress])

  return (
    <div>
      <InlineForm
        inputProps={{ type: 'text', value: amount, onChange: (event) => setAmount(event.target.value), placeholder: 'Amount' }}
        classNames="undelegation-form"
        onSubmit={onSubmitAction}
      />
      <p className="text-warning">Initiating an undelegation resets undelegation period.</p>
    </div>
  )
}

export default UndelegateForm
