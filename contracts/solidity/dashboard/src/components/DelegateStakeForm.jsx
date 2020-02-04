import React from 'react'
import Button from './Button'
import FormInput from './FormInput'
import { withFormik } from 'formik'
import { validateAmountInRange, validateEthAddress } from '../forms/common-validators'

const DelegateStakeForm = ({ handleSubmit }) => {
  return (
    <form className="delegate-stake-form tile flex flex-column" onSubmit={handleSubmit}>
      <FormInput
        name="stakeTokens"
        type="text"
        label="Stake Tokens"
      />
      <FormInput
        name="beneficiaryAddress"
        type="text"
        label="Beneficiary Address"
      />
      <FormInput
        name="operatorAddress"
        type="text"
        label="Operator Address"
      />
      <Button
        className="btn btn-primary btn-large"
        type="submit"
      >
        DELEGATE STAKE
      </Button>
    </form>
  )
}

const connectedWithFormik = withFormik({
  mapPropsToValues: () => ({ beneficiaryAddress: '', stakeTokens: '', operatorAddress: '' }),
  validate: (values, props) => {
    const { availableTokens, minStake } = props
    const { beneficiaryAddress, operatorAddress, stakeTokens } = values
    const errors = {}

    errors.stakeTokens = validateAmountInRange(stakeTokens, availableTokens, minStake)
    errors.beneficiaryAddress = validateEthAddress(beneficiaryAddress)
    errors.operatorAddress = validateEthAddress(operatorAddress)

    return errors
  },
  handleSubmit: (values, { props }) => {
    console.log('submit', values, props)
  },

  displayName: 'DelegateStakeForm',
})(DelegateStakeForm)

export default connectedWithFormik
