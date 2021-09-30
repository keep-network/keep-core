import React from "react"
import { withFormik } from "formik"
import { lte } from "../utils/arithmetics.utils"
import { getErrorsObj, validateAmountInRange } from "../forms/common-validators"
import {
  formatFloatingAmount,
  normalizeFloatingAmount,
} from "../forms/form.utils"
import FormInput from "./FormInput"
import MaxAmountAddon from "./MaxAmountAddon"
import { SubmitButton } from "./Button"
import { useCustomOnSubmitFormik } from "../hooks/useCustomOnSubmitFormik"
import useSetMaxAmountToken from "../hooks/useSetMaxAmountToken"
import { covKEEP } from "../utils/token.utils"

const WithdrawAmountForm = ({
  onCancel,
  onSubmit,
  submitBtnText,
  withdrawAmount,
  withdrawalDelay, // <number> in seconds
  ...formikProps
}) => {
  const onSubmitBtn = useCustomOnSubmitFormik(onSubmit)
  const onAddonClick = useSetMaxAmountToken(
    "withdrawAmount",
    withdrawAmount,
    covKEEP,
    covKEEP.decimals
  )

  return (
    <form
      onSubmit={formikProps.handleSubmit}
      className={"withdraw-amount-form"}
    >
      <FormInput
        name="withdrawAmount"
        type="text"
        label="Withdraw Amount"
        placeholder="0"
        normalize={normalizeFloatingAmount}
        format={formatFloatingAmount}
        inputAddon={<MaxAmountAddon onClick={onAddonClick} text="Max Amount" />}
        leftIconComponent={
          <span className={"form-input__left-icon__cov-keep-amount"}>
            covKEEP
          </span>
        }
      />
      <SubmitButton
        className="btn btn-lg btn-primary w-100"
        onSubmitAction={onSubmitBtn}
        disabled={!(formikProps.isValid && formikProps.dirty)}
      >
        {submitBtnText}
      </SubmitButton>
    </form>
  )
}

const WithdrawAmountFormWithFormik = withFormik({
  validateOnChange: true,
  validateOnBlur: true,
  mapPropsToValues: () => ({
    withdrawAmount: "0",
  }),
  validate: (values, props) => {
    const { withdrawAmount } = values
    const errors = {}

    if (lte(props.withdrawAmount || 0, 0)) {
      errors.withdrawAmount = "The value should be greater than zero"
    } else {
      errors.withdrawAmount = validateAmountInRange(
        withdrawAmount,
        props.withdrawAmount,
        1,
        covKEEP
      )
    }

    return getErrorsObj(errors)
  },
  displayName: "CoveragePoolsWithdrawAmountForm",
})(WithdrawAmountForm)

export default WithdrawAmountFormWithFormik
