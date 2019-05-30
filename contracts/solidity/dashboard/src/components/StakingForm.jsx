import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Button, Form, FormGroup,
  FormControl } from 'react-bootstrap'
import WithWeb3Context from './WithWeb3Context'
import { formatAmount } from '../utils'

const ERRORS = {
  INVALID_AMOUNT: 'Invalid amount',
  SERVER: 'Sorry, your request cannot be completed at this time.'
}

const RESET_DELAY = 3000 // 3 seconds

class StakingForm extends Component {

  state = {
    amount: 0,
    hasError: false,
    requestSent: false,
    requestSuccess: false,
    errorMsg: ERRORS.INVALID_AMOUNT
  }

  onChange = (e) => {
    this.setState(
      { amount: e.target.value }
    )
  }

  onRequestSuccess() {
    this.setState({
      hasError: false,
      requestSent: true,
      requestSuccess: true
    })
    window.setTimeout(() => {
      this.setState(this.state)
    }, RESET_DELAY)
  }

  onClick = (e) => {
    this.submit()
  }

  onSubmit = (e) => {
    e.preventDefault()
  }

  onKeyUp = (e) => {
    if (e.keyCode === 13) {
      this.submit()
    }
  }

  async submit() {
    const { amount } = this.state
    const { action, web3, stakingContractAddress } = this.props

    if (action === 'stake') {
      await web3.token.methods.approveAndCall(stakingContractAddress, formatAmount(amount, 18), "", {from: web3.yourAddress, gas: 150000})
    } else if (action === 'unstake') {
      await web3.stakingContract.methods.initiateUnstake(web3.utils.toBN(formatAmount(amount, 18)).toString(), web3.yourAddress).send({from: web3.yourAddress})
    }
  }

  render() {
    const { btnText } = this.props
    const { amount,
        hasError,
        errorMsg} = this.state

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
          <Button
            bsStyle="primary"
            bsSize="large"
            onClick={this.onClick}>
            {btnText}
          </Button>
        </Form>
        { hasError &&
          <small className="error-message">{errorMsg}</small> }
      </div>
    )
  }
}

StakingForm.propTypes = {
  btnText: PropTypes.string,
  action: PropTypes.string
}

export default WithWeb3Context(StakingForm);
