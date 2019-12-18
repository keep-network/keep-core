import React, { useContext } from 'react'
import { Web3Context } from './WithWeb3Context'
import { SubmitButton } from './Button'
import { useShowMessage, messageType } from './Message'

const Withdrawal = ({ withdrawal }) => {
  const { defaultContract, yourAddress } = useContext(Web3Context)
  const showMessage = useShowMessage()

  const submit = async (onTransactionHashCallback) => {
    try {
      await defaultContract.methods.finishUnstake(withdrawal.id)
        .send({ from: yourAddress })
        .on('transactionHash', onTransactionHashCallback)
      showMessage({ title: 'Success', content: 'Finish unstake transaction successfully completed' })
    } catch(error) {
      showMessage({ type: messageType.ERROR, title: 'Error', content: 'Finish unstake action has been failed' })
    }
   
  }

  return (
    <tr>
      <td>{withdrawal.amount}</td>
      <td className="text-mute">{withdrawal.availableAt}</td>
      <td>
      {withdrawal.available ? 
        <SubmitButton
          className="btn btn-priamry btn-sm"
          onSubmitAction={submit}
          pendingMessageTitle='Finish unstake transaction is pending...'
        >
          Finish Unstake
        </SubmitButton> : 'N/A'
      }
      </td>
    </tr>
  )
}

export default Withdrawal