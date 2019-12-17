import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Form, FormGroup, FormControl } from 'react-bootstrap'
import WithWeb3Context from './WithWeb3Context'
import { formatAmount } from '../utils'
import { SubmitButton } from './Button'

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

  onSubmit = (e) => {
    e.preventDefault()
  }

  onKeyUp = (e) => {
    if (e.keyCode === 13) {
      this.submit()
    }
  }

  submit = async () => {
    const { amount } = this.state
    const { action, web3 } = this.props
    const stakingContractAddress = web3.stakingContract.options.address;
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
          <SubmitButton
            type="submit"
            className="btn btn-primary btn-lg"
            onSubmitAction={this.submit}
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
  action: PropTypes.string
}

export default WithWeb3Context(StakingForm);
