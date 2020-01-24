import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Row, Form, FormGroup, FormControl, Col } from 'react-bootstrap'
import withWeb3Context from './WithWeb3Context'
import { SubmitButton } from './Button'

const ERRORS = {
  INVALID_AMOUNT: 'Invalid amount',
  SERVER: 'Sorry, your request cannot be completed at this time.',
}

const RESET_DELAY = 3000 // 3 seconds

class SigningForm extends Component {
  state = {
    messageToSign: this.props.defaultMessageToSign,
    signature: '',
    hasError: false,
    requestSent: false,
    requestSuccess: false,
    errorMsg: ERRORS.INVALID_AMOUNT,
  }

  onChange = (e) => {
    const name = e.target.name
    this.setState(
      { [name]: e.target.value },
    )
  }

  validateAddress() {
    const { web3 } = this.props
    if (web3.utils && web3.utils.isAddress(this.state.messageToSign)) return 'success'
    else return 'error'
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

  submit = async () => {
    const { messageToSign } = this.state
    const { web3 } = this.props

    const signature = await web3.eth.personal.sign(web3.utils.soliditySha3(messageToSign), web3.yourAddress, '')

    this.setState({
      signature: signature,
    })
  }

  render() {
    const { messageToSign, signature,
      hasError, errorMsg } = this.state
    const { description } = this.props

    const hidden = {
      display: signature ? 'block' : 'none',
    }

    return (
      <div>
        <p>{description}</p>
        <Form inline
          onSubmit={this.onSubmit}>
          <FormGroup>
            <FormControl
              type="text"
              name="messageToSign"
              value={messageToSign}
              onChange={this.onChange}
            />
          </FormGroup>
          <SubmitButton
            className="btn btn-primary btn-lg"
            onSubmitAction={this.submit}
            type="submit"
          >
            Sign
          </SubmitButton>
        </Form>
        { hasError &&
          <small className="error-message">{errorMsg}</small> }
        <Row style={ hidden }>
          <Col sm={12} >
            <p>Send the signature below to the stake owner to initiate stake delegation</p>
            <div className="well small">{signature}</div>
          </Col>
        </Row>
      </div>
    )
  }
}

SigningForm.propTypes = {
  btnText: PropTypes.string,
  action: PropTypes.string,
}

export default withWeb3Context(SigningForm)
