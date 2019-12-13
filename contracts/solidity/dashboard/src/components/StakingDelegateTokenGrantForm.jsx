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
    operatorAddress: "",
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

  validateAddress = (address) => {
    const { web3 } = this.props
    if (web3.utils && web3.utils.isAddress(address))
      return 'success'
    else
      return 'error'
  }

  validateAmount = () => {
    const { amount } = this.state
    const { web3, tokenBalance } = this.props
    if (web3.utils && tokenBalance && formatAmount(amount, 18).lte(tokenBalance))
      return 'success'
    else
      return 'error'
  }

  submit = async () => {
    const { grantId, amount, magpie, operatorAddress } = this.state
    const { web3 } = this.props
    const stakingContractAddress = web3.stakingContract.options.address;
    // Operator must sign grantee and token grant contract address since grant contract becomes the owner during grant staking.
    let delegation = Buffer.concat([
      Buffer.from(magpie.substr(2), 'hex'),
      Buffer.from(operatorAddress.substr(2), 'hex'),
    ]);

    await web3.grantContract.methods.stake(
      grantId,
      stakingContractAddress,
      web3.utils.toBN(formatAmount(amount, 18)).toString(),
      delegation
    ).send({from: web3.yourAddress})

  }

  render() {
    const { tokenBalance } = this.props
    const { grantId, amount, operatorAddress, magpie } = this.state

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

        <h4>Operator address</h4>
        <FormGroup validationState={this.validateAddress(operatorAddress)}>
          <FormControl
            type="textarea"
            name="operatorAddress"
            value={operatorAddress}
            onChange={this.onChange}
          />
          <FormControl.Feedback />
        </FormGroup>

        <h4>Magpie</h4>
        <p className="small">Address that receives earned rewards.</p>
        <FormGroup validationState={this.validateAddress(magpie)}>
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
