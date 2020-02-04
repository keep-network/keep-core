import React from 'react'
import { SubmitButton } from './Button'
import FormInput from './FormInput'
import { withFormik } from 'formik'
import { validateAmountInRange, validateEthAddress, getErrorsObj } from '../forms/common-validators'
import { useCustomOnSubmitFormik } from '../hooks/useCustomOnSubmitFormik'

const DelegateStakeForm = (props) => {
  const onSubmit = useCustomOnSubmitFormik(props.onSubmit)

  return (
    <form className="delegate-stake-form tile flex flex-column">
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
      <SubmitButton
        className="btn btn-primary btn-large"
        type="submit"
        onSubmitAction={onSubmit}
        withMessageActionIsPending={false}
        triggerManuallyFetch={true}
      >
        DELEGATE STAKE
      </SubmitButton>
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

    return getErrorsObj(errors)
  },
  displayName: 'DelegateStakeForm',
})(DelegateStakeForm)

export default connectedWithFormik
