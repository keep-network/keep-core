import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Form, FormGroup, FormControl } from 'react-bootstrap'
import WithWeb3Context from './WithWeb3Context'
import { formatAmount, displayAmount } from '../utils'
import { SubmitButton } from './Button'


class StakingDelegateTokenGrantForm extends Component {

  state = {
    grantId: 0,
    amount: 0,
    operatorSignature1: "",
    addressFromSignature1: "",
    operatorSignature2: "",
    addressFromSignature2: "",
    magpie: "",
  }

  onChange = (e) => {
    const name = e.target.name
    this.setState(
      { [name]: e.target.value }
    )
  }

  onSubmit = (e) => {
    e.preventDefault()
  }

  onKeyUp = (e) => {
    if (e.keyCode === 13) {
      this.submit()
    }
  }

  async recoverAddress(signature, message, fieldName) {
    const { web3 } = this.props
    let recoveredAddress
    try {
      recoveredAddress = await web3.eth.personal.ecRecover(web3.utils.soliditySha3(message), signature)
    } catch {
      recoveredAddress = "0x0"
    }

    if (this.state[fieldName] !== recoveredAddress) {
      this.setState({
        [fieldName]: recoveredAddress
      });
    }
  }

  validateSignature(signature, message, fieldName) {
    if (signature.length === 132) {
      this.recoverAddress(signature, message, fieldName)
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

  submit = async() => {
    const { grantId, amount, magpie, operatorSignature1, operatorSignature2 } = this.state
    const { web3 } = this.props
    const stakingContractAddress = web3.stakingContract.options.address;
    // Operator must sign grantee and token grant contract address since grant contract becomes the owner during grant staking.
    let delegation = Buffer.concat([
      Buffer.from(magpie.substr(2), 'hex'),
      Buffer.from(operatorSignature2.substr(2), 'hex'),
      Buffer.from(operatorSignature1.substr(2), 'hex')
    ]);

    await web3.grantContract.methods.stake(
      grantId,
      stakingContractAddress,
      web3.utils.toBN(formatAmount(amount, 18)).toString(),
      delegation
    ).send({from: web3.yourAddress})

  }

  render() {
    const { web3, tokenBalance } = this.props
    const tokenGrantContractAddress =  web3.stakingContract.options.address 
    const { grantId, amount, operatorSignature1, operatorSignature2, addressFromSignature1, addressFromSignature2, magpie } = this.state

    let showRecoveredAddress1 = {
      display: addressFromSignature1 ? "block" : "none"
    }

    let showRecoveredAddress2 = {
      display: addressFromSignature2 ? "block" : "none"
    }

    return (
      <Form onSubmit={this.onSubmit}>

        <h4>Grant Id</h4>
        <FormGroup>
          <FormControl
            type="text"
            name="grantId"
            value={grantId}
            onChange={this.onChange}
            onKeyUp={this.onKeyUp}
          />
          <FormControl.Feedback />
        </FormGroup>

        <h4>Operator signature</h4>
        <p className="small">ECDSA signature of your address obtained from the operator. Please be aware that operator can
            unstake but can not transfer or withdraw your tokens, any misbehavior of the operator will
            result stake slashing of your token balance.</p>

        <FormGroup validationState={this.validateSignature(operatorSignature1, web3.yourAddress, 'addressFromSignature1')}>
          <FormControl
            type="textarea"
            name="operatorSignature1"
            value={operatorSignature1}
            onChange={this.onChange}
          />
          <FormControl.Feedback />
        </FormGroup>
        <div style={showRecoveredAddress1} className="alert alert-info small">
          Operator: <strong>{addressFromSignature1}</strong>. Please check that the address is correct.
        </div>

        <h4>Operator signature of Token Grant contract address</h4>
        <p className="small">ECDSA signature of the address of the Token Grant contract <strong>{tokenGrantContractAddress}</strong> obtained from the operator.</p>

        <FormGroup validationState={this.validateSignature(operatorSignature2, tokenGrantContractAddress, 'addressFromSignature2')}>
          <FormControl
            type="textarea"
            name="operatorSignature2"
            value={operatorSignature2}
            onChange={this.onChange}
          />
          <FormControl.Feedback />
        </FormGroup>
        <div style={showRecoveredAddress2} className="alert alert-info small">
          Operator: <strong>{addressFromSignature2}</strong>. Please check that the address is correct.
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
          type="submit"
          className="btn btn-primary btn-lg"
          onSubmitAction={this.submit}
        >
          Delegate stake
        </SubmitButton>
      </Form>
    )
  }
}

StakingDelegateTokenGrantForm.propTypes = {
  btnText: PropTypes.string,
  action: PropTypes.string
}

export default WithWeb3Context(StakingDelegateTokenGrantForm);
