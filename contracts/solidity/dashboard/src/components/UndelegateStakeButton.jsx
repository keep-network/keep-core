import React from 'react'
import { SubmitButton } from './Button'
import WithWeb3Context from './WithWeb3Context'

const UndelegateStakeButton = (props) => {

  const undelegate = async () => {
    const { web3, amount, operator} = props
    await web3.stakingContract.methods.initiateUnstake(amount, operator).send({from: web3.yourAddress})  
  }

  return (
    <SubmitButton className="btn btn-primary btn-sm" onSubmitAction={undelegate}>
      Undelegate
    </SubmitButton>
  )
}

export default WithWeb3Context(UndelegateStakeButton)
