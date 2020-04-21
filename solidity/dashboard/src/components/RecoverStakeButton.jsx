import React, { useContext, useCallback } from 'react'
import { SubmitButton } from './Button'
import { Web3Context } from './WithWeb3Context'
import { useShowMessage, messageType } from './Message'
import { TOKEN_GRANT_CONTRACT_NAME, TOKEN_STAKING_CONTRACT_NAME } from '../constants/constants'

const RecoverStakeButton = ({ operatorAddress, ...props }) => {
  const web3Context = useContext(Web3Context)
  const { yourAddress } = web3Context
  const showMessage = useShowMessage()
  const { isFromGrant } = props
  const contract = web3Context[isFromGrant ? TOKEN_GRANT_CONTRACT_NAME : TOKEN_STAKING_CONTRACT_NAME]

  const recoverStake = useCallback(async (onTransactionHashCallback) => {

    try {
      await contract
        .methods
        .recoverStake(operatorAddress)
        .send({ from: yourAddress })
        .on('transactionHash', onTransactionHashCallback)
      showMessage({ type: messageType.SUCCESS, title: 'Success', content: 'Recover stake transaction successfully completed' })
    } catch (error) {
      showMessage({ type: messageType.ERROR, title: 'Recover stake action has been failed ', content: error.message })
      throw error
    }
  }, [operatorAddress, yourAddress, contract.methods, showMessage])

  return (
    <SubmitButton
      className={props.btnClassName}
      onSubmitAction={recoverStake}
      pendingMessageTitle='Recover stake transaction is pending...'
      successCallback={props.successCallback}
    >
      {props.btnText}
    </SubmitButton>
  )
}

RecoverStakeButton.defaultProps = {
  btnClassName: 'btn btn-sm btn-secondary',
  btnText: 'recover',
  successCallback: () => {},
  isFromGrant: false,
}

export default React.memo(RecoverStakeButton)
