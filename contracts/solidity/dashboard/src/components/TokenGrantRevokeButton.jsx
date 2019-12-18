import React, { useContext } from 'react'
import { Web3Context } from './WithWeb3Context'
import { SubmitButton } from './Button'
import { useShowMessage, messageType } from './Message'

const TokenGrantRevokeButton = ({ item }) => {
  const { grantContract, yourAddress } = useContext(Web3Context)
  const showMessage = useShowMessage()

  const submit = async (onTransactionHashCallback) => {
    try {
      await grantContract.methods.revoke(item.id).send({ from: yourAddress }).on('transactionHash', onTransactionHashCallback)
      showMessage({ title: 'Success', content: 'Revoke transaction has been successfully completed' })
    } catch(error) {
      showMessage({ type: messageType.ERROR, title: 'Revoke action has been failed', content: error.message })
    }
  }

  let button = 'Non revocable'

  if (item.revoked) {
    button = 'Revoked'
  }

  if (item.revocable && !item.revoked) {
    button = (
      <SubmitButton
        className="btn btn-primary btn-sm"
        onSubmitAction={submit}
        pendingMessageTitle="Revoke transaction is pending..."
      >
        Revoke
      </SubmitButton>
    )
  }

  return button
}

export default TokenGrantRevokeButton