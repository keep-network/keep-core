import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Form, FormGroup, FormControl } from 'react-bootstrap'
import WithWeb3Context from './WithWeb3Context'
import { formatAmount, displayAmount } from '../utils'
import { SubmitButton } from './Button'

const ERRORS = {
  INVALID_AMOUNT: 'Invalid amount',
  SERVER: 'Sorry, your request cannot be completed at this time.'
}

const RESET_DELAY = 3000 // 3 seconds

class StakingDelegateForm extends Component {

  state = {
    amount: 100,
    operatorSignature: "",
    operatorAddress: "",
    magpie: "",
    hasError: false,
    requestSent: false,
    requestSuccess: false,
    errorMsg: ERRORS.INVALID_AMOUNT
  }

  onChange = (e) => {
    const name = e.target.name
    this.setState(
      { [name]: e.target.value }
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

  async recoverOperatorAddress(signature) {
    const { web3 } = this.props
    let recoveredOperator
    try {
      recoveredOperator = await web3.eth.personal.ecRecover(web3.utils.soliditySha3(web3.yourAddress), signature)
    } catch {
      recoveredOperator = "0x0"
    }

    if (this.state.operatorAddress !== recoveredOperator) {
      this.setState({
        operatorAddress: recoveredOperator
      });
    }
  }

  validateOperatorSignature() {
    const { web3 } = this.props
    if (web3.utils && this.state.operatorSignature.length === 132) {
      this.recoverOperatorAddress(this.state.operatorSignature)
      return 'success'
    } else return 'error'
  }

  validateMagpie() {
    const { web3 } = this.props
    if (web3.utils && web3.utils.isAddress(this.state.magpie)) return 'success'
    else return 'error'
  }

  validateAmount() {
    const { amount } = this.state
    const { web3, tokenBalance } = this.props
    if (web3.utils && tokenBalance && formatAmount(amount, 18).lte(tokenBalance)) return 'success'
    else return 'error'
  }

  submit = async () => {
    const { amount, magpie, operatorSignature } = this.state
    const { web3 } = this.props;
    const stakingContractAddress = web3.stakingContract.options.address;
    let delegationData = '0x' + Buffer.concat([Buffer.from(magpie.substr(2), 'hex'), Buffer.from(operatorSignature.substr(2), 'hex')]).toString('hex');
    await web3.token.methods.approveAndCall(stakingContractAddress, web3.utils.toBN(formatAmount(amount, 18)).toString(), delegationData).send({from: web3.yourAddress})
  }

  render() {
    const { tokenBalance } = this.props
    const { amount, operatorSignature, operatorAddress, magpie,
        hasError,
        errorMsg} = this.state
    
    let operatorRecovered = {
      display: operatorAddress ? "block" : "none"
    }
    
    return (
      <Form onSubmit={this.onSubmit}>
        <h4>Operator signature</h4>
        <p className="small">ECDSA signature of your address obtained from the operator. Please be aware that operator can
            unstake but can not transfer or withdraw your tokens, any misbehavior of the operator will
            result stake slashing of your token balance.</p>

        <FormGroup validationState={this.validateOperatorSignature()}>
          <FormControl
            type="textarea"
            name="operatorSignature"
            value={operatorSignature}
            onChange={this.onChange}
          />
          <FormControl.Feedback />
        </FormGroup>
        <div style={operatorRecovered} className="alert alert-info small">
          Operator: <strong>{operatorAddress}</strong>. Please check that the address is correct.
        </div>
        <h4>Magpie</h4>
        <p className="small">Address that receives earned rewards.</p>
        <FormGroup validationState={this.validateMagpie()}>
          <FormControl
            type="text"
            name="magpie"
            value={magpie}
            onChange={this.onChange}
          />
          <FormControl.Feedback />
        </FormGroup>

        <p className="small"> You can stake up to { displayAmount(tokenBalance, 18, 3) } KEEP</p>

        <FormGroup validationState={this.validateAmount()}>
          <FormControl
            type="text"
            name="amount"
            value={amount}
            onChange={this.onChange}
            onKeyUp={this.onKeyUp}
          />
          <FormControl.Feedback />
        </FormGroup>

        <SubmitButton
          type='submit'
          className="btn btn-primary btn-lg"
          onSubmitAction={this.submit}>
          Delegate stake
        </SubmitButton>
        { hasError &&
          <small className="error-message">{errorMsg}</small> }
      </Form>
    )
  }
}

StakingDelegateForm.propTypes = {
  btnText: PropTypes.string,
  action: PropTypes.string
}

export default WithWeb3Context(StakingDelegateForm);
