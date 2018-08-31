import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Button, Form, FormGroup,
  FormControl } from 'react-bootstrap'
import { getKeepToken, getTokenStaking } from '../contracts'
import Network from '../network'
import { formatAmount } from '../utils'

const ERRORS = {
  INVALID_AMOUNT: 'Invalid amount',
  SERVER: 'Sorry, your request cannot be completed at this time.'
}

const RESET_DELAY = 3000 // 3 seconds

class StakingForm extends Component {
  constructor(props) {
    super(props)
    this.state = this.getInitialState()
  }

  getInitialState() {
    return {
      amount: 0,
      hasError: false,
      requestSent: false,
      requestSuccess: false,
      errorMsg: ERRORS.INVALID_AMOUNT,
    }
  }

  onChange(e) {
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
      this.setState(this.getInitialState())
    }, RESET_DELAY)
  }

  onClick(e) {
    this.submit()
  }

  onKeyUp(e) {
    if (e.keyCode === 13) {
      this.submit()
    }
  }

  async submit() {
    const { amount } = this.state
    const { action } = this.props
    const { stakingContractAddress } = this.props

    const accounts = await Network.getAccounts()
    const token = await getKeepToken(process.env.REACT_APP_TOKEN_ADDRESS)
    const stakingContract = await getTokenStaking(stakingContractAddress)

    if (action === 'stake') {
      token.approveAndCall(stakingContractAddress, formatAmount(amount, 18), "", {from: accounts[0], gas: 150000})
    } else if (action === 'unstake') {
      stakingContract.initiateUnstake(formatAmount(amount, 18), {from: accounts[0], gas: 150000})
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
          onSubmit={(e) => { e.preventDefault() }}>
          <FormGroup>
            <FormControl
              type="text"
              value={amount}
              onChange={this.onChange.bind(this)}
              onKeyUp={this.onKeyUp.bind(this)}
              />
          </FormGroup>
          <Button
            bsStyle="primary"
            bsSize="large"
            onClick={this.onClick.bind(this)}>
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

export default StakingForm
