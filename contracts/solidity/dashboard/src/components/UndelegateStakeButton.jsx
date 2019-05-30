import React, { Component } from 'react'
import { Button } from 'react-bootstrap'
import WithWeb3Context from './WithWeb3Context'

class UndelegateStakeButton extends Component {

  undelegate = async () => {
    const { web3, amount, operator} = this.props
    await web3.stakingContract.methods.initiateUnstake(amount, operator).send({from: web3.yourAddress})
  }

  render() {
    return (
      <Button bsSize="small" bsStyle="primary" onClick={this.undelegate}>Undelegate</Button>
    )
  }
}

export default WithWeb3Context(UndelegateStakeButton)
