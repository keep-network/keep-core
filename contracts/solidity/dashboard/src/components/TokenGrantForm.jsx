import moment from 'moment'
import React, { Component } from 'react'
import { Form, FormGroup, FormControl, ControlLabel, Col, HelpBlock, Checkbox } from 'react-bootstrap'
import WithWeb3Context from './WithWeb3Context'
import { formatAmount } from '../utils'
import { SubmitButton } from './Button'
import { MessagesContext, messageType } from './Message'

const ERRORS = {
  INVALID_AMOUNT: 'Invalid amount.',
  SERVER: 'Sorry, your request cannot be completed at this time.'
}

class TokenGrantForm extends Component {
  static contextType = MessagesContext

  state = {
    amount: 0,
    grantee: "0x0",
    duration: 1,
    start: moment().unix(),
    cliff: 1,
    revocable: false,
    formErrors: {
      grantee: '',
      amount: ''
    },
    hasError: false,
    requestSent: false,
    requestSuccess: false,
    errorMsg: ERRORS.INVALID_AMOUNT
  }

  onChange = (e) => {
    const name = e.target.name
    const value = e.target.type === 'checkbox' ? e.target.checked : e.target.value
    this.setState({[name]: value})
  }

  onSubmit = (e) => {
    e.preventDefault()
  }

  validateGrantee() {
    const { web3 } = this.props
    if (web3.utils && web3.utils.isAddress(this.state.grantee)) return 'success'
    else return 'error'
  }

  submit = async (onTransationHashCallback) => {
    const { amount, grantee, duration, start, cliff, revocable} = this.state
    const { web3 } = this.props

    const tokenGrantContractAddress = web3.grantContract.options.address;
    const extraData = Buffer.concat([
      Buffer.from(grantee.substr(2), 'hex'),
      web3.utils.toBN(duration).toBuffer('be', 32),
      web3.utils.toBN(start).toBuffer('be', 32),
      web3.utils.toBN(cliff).toBuffer('be', 32),
      Buffer.from(revocable ? "01" : "00", 'hex'),
    ])
    
    try {
      await web3.token.methods.approveAndCall(
        tokenGrantContractAddress,
        web3.utils.toBN(formatAmount(amount, 18)).toString(),
        extraData)
        .send({from: web3.yourAddress})
        .on('transactionHash', onTransationHashCallback)
      this.context.showMessage({ title: 'Success', content: 'Grant tokens transaction successfully completed' })
    } catch(error) {
      this.context.showMessage({ type: messageType.ERROR, title: 'Grant tokens action has been failed ', content: error.message })
    }
  }

  render() {
    const { amount, grantee, duration, start, cliff, revocable,
        hasError,
        errorMsg} = this.state

    return (
      <div className="token-grant-form">
        <Form horizontal onSubmit={this.onSubmit}>
          <FormGroup validationState={this.validateGrantee()}>
            <Col componentClass={ControlLabel} sm={2}>
              Grantee:
            </Col>
            <Col sm={8}>
              <FormControl
                type="text"
                name="grantee"
                value={grantee}
                onChange={this.onChange}
              />
              <FormControl.Feedback />
              <HelpBlock className="small">Address to which granted tokens are going to be withdrawn.</HelpBlock>
            </Col>
          </FormGroup>

          <FormGroup>
            <Col componentClass={ControlLabel} sm={2}>
              Amount:
            </Col>
            <Col sm={8}>
              <FormControl
                type="text"
                name="amount"
                value={amount}
                onChange={this.onChange}
              />
              <FormControl.Feedback />
              <HelpBlock className="small">Amount to be granted.</HelpBlock>
            </Col>
          </FormGroup>

          <FormGroup>
            <Col componentClass={ControlLabel} sm={2}>
             Duration:
            </Col>
            <Col sm={8}>
              <FormControl
                type="text"
                name="duration"
                value={duration}
                onChange={this.onChange}
              />
              <FormControl.Feedback />
              <HelpBlock className="small">Duration in seconds of the period in which the tokens will vest.</HelpBlock>
            </Col>
          </FormGroup>

          <FormGroup>
            <Col componentClass={ControlLabel} sm={2}>
              Start:
            </Col>
            <Col sm={8}>
              <FormControl
                type="text"
                name="start"
                value={start}
                onChange={this.onChange}
              />
              <FormControl.Feedback />
              <HelpBlock className="small">Timestamp at which vesting will start.</HelpBlock>
            </Col>
          </FormGroup>

          <FormGroup>
            <Col componentClass={ControlLabel} sm={2}>
              Cliff:
            </Col>
            <Col sm={8}>
              <FormControl
                type="text"
                name="cliff"
                value={cliff}
                onChange={this.onChange}
              />
              <FormControl.Feedback />
              <HelpBlock className="small">Duration in seconds of the cliff after which tokens will begin to vest.</HelpBlock>
            </Col>
          </FormGroup>

          <FormGroup>
            <Col componentClass={ControlLabel} sm={2}>
              Revocable:
            </Col>
            <Col sm={8}>
              <Checkbox
                name="revocable"
                checked={revocable}
                onChange={this.onChange}></Checkbox>
              <HelpBlock className="small">Whether the token grant is revocable or not.</HelpBlock>
            </Col>
          </FormGroup>

          <SubmitButton
            className="btn btn-primary btn-lg"
            type="submit"
            onSubmitAction={this.submit}
            pendingMessageTitle="Grant tokens transaction is pending..."
          >
            Grant tokens
          </SubmitButton>
        </Form>
        { hasError &&
          <small className="error-message">{errorMsg}</small> }
      </div>
    )
  }
}

export default WithWeb3Context(TokenGrantForm)
