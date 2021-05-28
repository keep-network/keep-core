import React from "react"
import { withFormik } from "formik"
import { getErrorsObj } from "../../forms/common-validators"
import { SubmitButton } from "../../components/Button"
import FormInput from "../../components/FormInput"
import { useCustomOnSubmitFormik } from "../../hooks/useCustomOnSubmitFormik"
import { normalizeAmount, formatAmount as formatFormAmount, } from "../../forms/form.utils.js"
import { colors } from "../../constants/colors"

const TransferKeepTokensForm = ({
  onSubmit,
  ...formikProps
}) => {
  const onSubmitBtn = useCustomOnSubmitFormik(onSubmit)

  return (
    <form>
      <FormInput
        name="address"
        type="text"
        label="Address"
        placeholder="0x0"
      />
      <FormInput
        name="amount"
        type="text"
        label="Amount"
        placeholder="0"
        normalize={normalizeAmount}
        format={formatFormAmount}
        placeholder="0"
      />
      <div
        className="flex row center mt-2"
        style={{
          borderTop: `1px solid ${colors.grey20}`,
          margin: "0 -2rem",
          padding: "2rem 2rem 0",
        }}
      >
        <SubmitButton
          className="btn btn-primary"
          type="submit"
          onSubmitAction={onSubmitBtn}
          withMessageActionIsPending={false}
          triggerManuallyFetch={true}
          disabled={!formikProps.dirty}
        >
          Transfer
      </SubmitButton>
      </div>
    </form>
  )
}

export const TransferKeepTokensFormik = withFormik({
  validateOnChange: false,
  validateOnBlur: false,
  mapPropsToValues: () => ({
    amount: "0",
    address: ""
  }),
  validate: (values) => {
    const { amount, address } = values
    const errors = {}

    if (!amount) {
      errors.amount = "Required"
    }

    if (!address) {
      errors.address = "Required"
    }

    return getErrorsObj(errors)
  },
  displayName: "MintBondTokensForm",
})(TransferKeepTokensForm)
