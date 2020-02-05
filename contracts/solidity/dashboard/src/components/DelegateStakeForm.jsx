import React, { useCallback } from 'react'
import { SubmitButton } from './Button'
import FormInput from './FormInput'
import { withFormik, useFormikContext } from 'formik'
import { validateAmountInRange, validateEthAddress, getErrorsObj } from '../forms/common-validators'
import { useCustomOnSubmitFormik } from '../hooks/useCustomOnSubmitFormik'

const DelegateStakeForm = (props) => {
  const onSubmit = useCustomOnSubmitFormik(props.onSubmit)

  return (
    <form className="delegate-stake-form tile flex flex-column">
      <h5>Delegate Stake</h5>
      <ContextSwitch />
      <div className="input-wrapper flex flex-row">
        <FormInput
          name="stakeTokens"
          type="text"
          label="Token Amount"
        />
        <div className="flex flex-1 flex-column flex-column-center">
          <div className="text text-smaller" style={{ marginTop: '1.5rem' }}>
            Min Stake: {props.minStake} KEEP
          </div>
          <div className="text text-smaller" style={{ marginTop: '1rem' }}>
            {props.availableTokens} KEEP available
          </div>
        </div>
      </div>
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
      <FormInput
        name="authorizerAddress"
        type="text"
        label="Authorizer Address"
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

const ContextSwitch = (props) => {
  const { setFieldValue, values } = useFormikContext()

  const getClassName = useCallback((contextName) => {
    return values.context === contextName ? 'active' : 'inactive'
  }, [values.context])

  const onClick = useCallback((event) => {
    setFieldValue('context', event.target.id, false)
  }, [])

  return (
    <div className="tabs flex">
      <div
        id="owned"
        className={`tab ${getClassName('owned')}`}
        onClick={onClick}
      >
        OWNED
      </div>
      <div
        id="granted"
        className={`tab ${getClassName('granted')}`}
        onClick={onClick}
      >
        GRANTED
      </div>
    </div>
  )
}

const connectedWithFormik = withFormik({
  mapPropsToValues: () => ({
    beneficiaryAddress: '',
    stakeTokens: '',
    operatorAddress: '',
    authorizerAddress: '',
    context: 'granted',
  }),
  validate: (values, props) => {
    const { availableTokens, minStake } = props
    const { beneficiaryAddress, operatorAddress, stakeTokens, authorizerAddress } = values
    const errors = {}

    errors.stakeTokens = validateAmountInRange(stakeTokens, availableTokens, minStake)
    errors.beneficiaryAddress = validateEthAddress(beneficiaryAddress)
    errors.operatorAddress = validateEthAddress(operatorAddress)
    errors.authorizerAddress = validateEthAddress(authorizerAddress)

    return getErrorsObj(errors)
  },
  displayName: 'DelegateStakeForm',
})(DelegateStakeForm)

export default connectedWithFormik
