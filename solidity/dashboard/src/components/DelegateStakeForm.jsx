import React, { useState } from 'react'
import { SubmitButton } from './Button'
import FormInput from './FormInput'
import { withFormik, useFormikContext } from 'formik'
import {
  validateAmountInRange,
  validateEthAddress,
  getErrorsObj,
} from '../forms/common-validators'
import { useCustomOnSubmitFormik } from '../hooks/useCustomOnSubmitFormik'
import { displayAmount, formatAmount } from '../utils/general.utils'
import ProgressBar from './ProgressBar'
import { colors } from '../constants/colors'
import {
  normalizeAmount,
  formatAmount as formatFormAmount,
} from '../forms/form.utils.js'
import { lte } from '../utils/arithmetics.utils'
import * as Icons from './Icons'

const DelegateStakeForm = ({ onSubmit, minStake, availableToStake, ...formikProps }) => {
  const onSubmitBtn = useCustomOnSubmitFormik(onSubmit)
  const stakeTokensValue = formatAmount(formikProps.values.stakeTokens)

  return (
    <form className="delegate-stake-form flex column">
     <TokensAmountField
      availableToStake={availableToStake}
      minStake={minStake}
      stakeTokensValue={stakeTokensValue}
     />
      <AddressField
        name="authorizerAddress"
        type="text"
        label="Authorizer Address"
        placeholder="0x0"
        icon={<Icons.AuthorizerFormIcon />}
      />
      <AddressField
        name="operatorAddress"
        type="text"
        label="Operator Address"
        placeholder="0x0"
        icon={<Icons.OperatorFormIcon />}
      />      
      <AddressField
        name="beneficiaryAddress"
        type="text"
        label="Beneficiary Address"
        placeholder="0x0"
        icon={<Icons.BeneficiaryFormIcon />}
      />
      <SubmitButton
        className="btn btn-primary"
        type="submit"
        onSubmitAction={onSubmitBtn}
        withMessageActionIsPending={false}
        triggerManuallyFetch={true}
        disabled={true}
      > 
        delegate stake
      </SubmitButton>
    </form>
  )
}

const AddressField = ({ icon, ...formInputProps}) => {
  const [focused, setFocused] = useState(false)
  const { setFieldTouched, touched } = useFormikContext()
  const isTouched = focused || touched[formInputProps.name]
  console.log('isTouched', focused, formInputProps.name, touched, touched[formInputProps.name])

  const onFocus = () => {
    setFocused(true)
    if (formInputProps.name === 'operatorAddress' && !touched.authorizerAddress) {
      setFieldTouched('authorizerAddress', true, false)
    } else if (formInputProps.name === 'beneficiaryAddress' && (!touched.authorizerAddress || !touched.operatorAddress)) {
      setFieldTouched('authorizerAddress', true, false)
      setFieldTouched('operatorAddress', true, false)
    }
  }

  return (
    <div className={`address-field-wrapper${isTouched ? ' touched' : ''}`}>
      <Icons.DashedLine />
      {icon}
      <FormInput
        {...formInputProps}
        onFocus={onFocus}
      />
    </div>
  )
}

const TokensAmountField = ({ availableToStake, minStake, stakeTokensValue }) => {
  return (
    <div className="token-amount-wrapper">
      <Icons.KeepGreenOutline />
      <div className="token-amount-field">
        <FormInput
            name="stakeTokens"
            type="text"
            label="Token Amount"
            normalize={normalizeAmount}
            format={formatFormAmount}
            placeholder="0"
        />
        <ProgressBar
          total={availableToStake}
          items={[{ value: stakeTokensValue, color: colors.primary }]}
        />
        <div className="text-small text-grey-50">
          {availableToStake} KEEP available
        </div>
        <div className="text-smaller text-grey-30">
          Min Stake: {displayAmount(minStake)} KEEP
        </div>
      </div>
  </div>
  )
}

const connectedWithFormik = withFormik({
  mapPropsToValues: () => ({
    beneficiaryAddress: null,
    stakeTokens: null,
    operatorAddress: null,
    authorizerAddress: null,
  }),
  validate: (values, props) => {
    const { beneficiaryAddress, operatorAddress, authorizerAddress } = values
    const errors = {}

    errors.stakeTokens = getStakeTokensError(props, values)
    errors.beneficiaryAddress = validateEthAddress(beneficiaryAddress)
    errors.operatorAddress = validateEthAddress(operatorAddress)
    errors.authorizerAddress = validateEthAddress(authorizerAddress)

    return getErrorsObj(errors)
  },
  displayName: 'DelegateStakeForm',
})(DelegateStakeForm)

const getStakeTokensError = (props, { stakeTokens }) => {
  const { availableToStake, minStake } = props

  if (lte(availableToStake, 0)) {
    return 'Insufficient funds'
  } else {
    return validateAmountInRange(stakeTokens, availableToStake, minStake)
  }
}

export default connectedWithFormik
