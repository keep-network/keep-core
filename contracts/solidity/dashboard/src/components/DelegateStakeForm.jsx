import React from 'react'
import Button from './Button'
import FormInput from './FormInput'
import { useFormik } from 'formik'

const DelegateStakeForm = () => {
  const { handleSubmit, getFieldMeta, getFieldProps } = useFormik({
    initialValues: {
      beneficiaryAddress: '',
      stakeTokens: '',
      operatorAddress: '',
    }, onSubmit: (values) => {
      console.log('values', values)
    },
    validate: (values) => {
      console.log('validation')
    },
  })

  return (
    <form className="delegate-stake-form tile flex flex-column" onSubmit={handleSubmit}>
      <FormInput
        getFieldProps={getFieldProps}
        getFieldMeta={getFieldMeta}
        name="stakeTokens"
        type="text"
        label="Stake Tokens"
      />
      <FormInput
        getFieldProps={getFieldProps}
        getFieldMeta={getFieldMeta}
        name="beneficiaryAddress"
        type="text"
        label="Beneficiary Address"
      />
      <FormInput
        getFieldProps={getFieldProps}
        getFieldMeta={getFieldMeta}
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

export default DelegateStakeForm
