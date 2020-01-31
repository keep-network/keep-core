import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Form, FormGroup, FormControl } from 'react-bootstrap'
import withWeb3Context from './WithWeb3Context'
import { formatAmount } from '../utils'
import { SubmitButton } from './Button'
import { MessagesContext, messageType } from './Message'

const ERRORS = {
  INVALID_AMOUNT: 'Invalid amount',
  SERVER: 'Sorry, your request cannot be completed at this time.',
}

const RESET_DELAY = 3000 // 3 seconds

class StakingForm extends Component {
  static contextType = MessagesContext

  state = {
    amount: 0,
    hasError: false,
    requestSent: false,
    requestSuccess: false,
    errorMsg: ERRORS.INVALID_AMOUNT,
  }

  onChange = (e) => {
    this.setState(
      { amount: e.target.value },
    )
  }

  onRequestSuccess() {
    this.setState({
      hasError: false,
      requestSent: true,
      requestSuccess: true,
    })
    window.setTimeout(() => {
      this.setState(this.state)
    }, RESET_DELAY)
  }

  onSubmit = (e) => {
    e.preventDefault()
  }

  onKeyUp = (e) => {
    if (e.keyCode === 13) {
      this.submit()
    }
  }

  getCaptializedActionName = () => this.props.action.charAt(0).toUpperCase() + this.props.action.slice(1)

  submit = async (onTransactionHashCallback) => {
    const { amount } = this.state
    const { action, web3 } = this.props
    const stakingContractAddress = web3.stakingContract.options.address
    const actionName = this.getCaptializedActionName()

    try {
      if (action === 'stake') {
        const delegationData = '0x' + Buffer.concat([Buffer.from(web3.yourAddress.substr(2), 'hex'), Buffer.from(web3.yourAddress.substr(2), 'hex')]).toString('hex')
        await web3.token.methods
          .approveAndCall(stakingContractAddress, web3.utils.toBN(formatAmount(amount, 18)).toString(), delegationData)
          .send({ from: web3.yourAddress, gas: 150000 })
          .on('transactionHash', onTransactionHashCallback)
      } else if (action === 'undelegate') {
        await web3.stakingContract.methods.undelegate(web3.yourAddress)
          .send({ from: web3.yourAddress })
          .on('transactionHash', onTransactionHashCallback)
      }
      this.context.showMessage({ type: messageType.SUCCESS, title: 'Success', content: `${actionName} transaction has been successfully completed` })
    } catch (error) {
      this.context.showMessage({ type: messageType.ERROR, title: `${actionName} action has been failed`, content: error.message })
    }
  }

  render() {
    const { btnText } = this.props
    const { amount,
      hasError,
      errorMsg } = this.state

    return (
      <div className="staking-form">
        <Form inline
          onSubmit={this.onSubmit}>
          <FormGroup>
            <FormControl
              type="text"
              value={amount}
              onChange={this.onChange}
              onKeyUp={this.onKeyUp}
            />
          </FormGroup>
          <SubmitButton
            type="submit"
            className="btn btn-primary btn-lg"
            onSubmitAction={this.submit}
            pendingMessageTitle={`${this.getCaptializedActionName()} transaction is pending...`}
          >
            {btnText}
          </SubmitButton>
        </Form>
        { hasError &&
          <small className="error-message">{errorMsg}</small> }
      </div>
    )
  }
}

StakingForm.propTypes = {
  btnText: PropTypes.string,
  action: PropTypes.string,
}

export default withWeb3Context(StakingForm)
