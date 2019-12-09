import React, { Component } from 'react'
import WithWeb3Context from './WithWeb3Context'
import { SubmitButton } from './Button'

class Withdrawal extends Component {

  finishUnstake = async () => {
    const { web3, withdrawal } = this.props
    await web3.defaultContract.methods.finishUnstake(withdrawal.id).send({from: web3.yourAddress})
  }

  render() {
    const { withdrawal } = this.props
    let action = 'N/A'
    if (withdrawal.available) {
      action = <SubmitButton className="btn btn-priamry btn-sm" onSubmitAction={this.finishUnstake}>Finish Unstake</SubmitButton>
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
