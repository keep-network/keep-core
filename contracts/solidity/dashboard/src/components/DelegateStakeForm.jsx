import React, { useCallback } from 'react'
import { SubmitButton } from './Button'
import FormInput from './FormInput'
import { withFormik, useFormikContext } from 'formik'
import { validateAmountInRange, validateEthAddress, getErrorsObj } from '../forms/common-validators'
import { useCustomOnSubmitFormik } from '../hooks/useCustomOnSubmitFormik'
import { displayAmount } from '../utils'

const DelegateStakeForm = ({ onSubmit, minStake, keepBalance, grantBalance, ...formikProps }) => {
  const onSubmitBtn = useCustomOnSubmitFormik(onSubmit)

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
        <div className="flex flex-column flex-column-center">
          <div className="text text-smaller" style={{ marginTop: '1.5rem' }}>
            Min Stake: {displayAmount(minStake)} KEEP
          </div>
          <div className="text text-smaller" style={{ marginTop: '1rem' }}>
            {formikProps.values.context === 'granted' ? displayAmount(grantBalance) : displayAmount(keepBalance)} KEEP available
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
        onSubmitAction={onSubmitBtn}
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
        id="granted"
        className={`tab ${getClassName('granted')}`}
        onClick={onClick}
      >
        GRANTED
      </div>
      <div
        id="owned"
        className={`tab ${getClassName('owned')}`}
        onClick={onClick}
      >
        OWNED
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
    const { keepBalance, grantBalance, minStake } = props
    const { beneficiaryAddress, operatorAddress, stakeTokens, authorizerAddress, context } = values
    const errors = {}

    errors.stakeTokens = validateAmountInRange(stakeTokens, context === 'granted' ? grantBalance : keepBalance, minStake)
    errors.beneficiaryAddress = validateEthAddress(beneficiaryAddress)
    errors.operatorAddress = validateEthAddress(operatorAddress)
    errors.authorizerAddress = validateEthAddress(authorizerAddress)

    return getErrorsObj(errors)
  },
  displayName: 'DelegateStakeForm',
})(DelegateStakeForm)

export default connectedWithFormik
