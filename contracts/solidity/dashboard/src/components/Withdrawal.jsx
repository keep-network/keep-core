import React, { Component } from 'react'
import { Button } from 'react-bootstrap'
import WithWeb3Context from './WithWeb3Context'

class Withdrawal extends Component {

  async finishUnstake(withdrawalId) {
    const { web3 } = this.props
    web3.stakingContract.finishUnstake(withdrawalId, {from: web3.yourAddress, gas: 110000})
  }

  render() {
    const { withdrawal } = this.props
    let action = 'N/A'
    if (withdrawal.available) {
      action = <Button bsSize="small" bsStyle="primary" onClick={()=>this.finishUnstake(withdrawal.id)}>Finish Unstake</Button>
    }

    return (
      <tr>
        <td>{withdrawal.amount}</td>
        <td className="text-mute">{withdrawal.availableAt}</td>
        <td>{action}</td>
      </tr>
    )
  }
}

export default WithWeb3Context(Withdrawal)
