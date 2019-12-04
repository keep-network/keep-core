import React, { useContext } from 'react'
import { Button } from 'react-bootstrap'
import { Web3Context } from './WithWeb3Context'
import { useShowMessage } from './Message'

const UndelegateStakeButton = (props) => {
  const web3 = useContext(Web3Context)
  const showMessage = useShowMessage()

  const undelegate = async () => {
    const { amount, operator} = props
    try {
      await web3.stakingContract.methods.initiateUnstake(amount, operator).send({from: web3.yourAddress})
      showMessage({ type: 'success', title: 'Success', content: 'Undelegate transaction successfully completed' })
    } catch (error) {
      showMessage({ type: 'error', title: 'Undelegate action has been failed ', content: error.message })
    }
  }

  return (
    <Button bsSize="small" bsStyle="primary" onClick={undelegate}>Undelegate</Button>
  )
}

export default UndelegateStakeButton
