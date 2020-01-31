import React, { useContext } from 'react'
import { Web3Context } from './WithWeb3Context'
import { SubmitButton } from './Button'
import { useShowMessage, messageType } from './Message'

const Undelegation = ({ undelegation }) => {
  const { defaultContract, yourAddress } = useContext(Web3Context)
  const showMessage = useShowMessage()

  const submit = async (onTransactionHashCallback) => {
    try {
      await defaultContract.methods.recoverStake(undelegation.id)
        .send({ from: yourAddress })
        .on('transactionHash', onTransactionHashCallback)
      showMessage({ title: 'Success', content: 'Recover stake transaction has been successfully completed' })
    } catch (error) {
      showMessage({ type: messageType.ERROR, title: 'Error', content: 'Recover stake action has been failed' })
    }
  }

  return (
    <tr>
      <td>{undelegation.amount}</td>
      <td className="text-mute">{undelegation.availableAt}</td>
      <td>
        {undelegation.available ?
          <SubmitButton
            className="btn btn-priamry btn-sm"
            onSubmitAction={submit}
            pendingMessageTitle='Recover stake transaction is pending...'
          >
          Recover Stake
          </SubmitButton> : 'N/A'
        }
      </td>
    </tr>
  )
}

export default Undelegation
