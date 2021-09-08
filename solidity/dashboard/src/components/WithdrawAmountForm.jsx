import React from "react"
import { withFormik } from "formik"
import { lte } from "../utils/arithmetics.utils"
import { getErrorsObj, validateAmountInRange } from "../forms/common-validators"
import { normalizeFloatingAmount } from "../forms/form.utils"
import FormInput from "./FormInput"
import MaxAmountAddon from "./MaxAmountAddon"
import { SubmitButton } from "./Button"
import { useCustomOnSubmitFormik } from "../hooks/useCustomOnSubmitFormik"
import { KEEP } from "../utils/token.utils"
import useSetMaxAmountToken from "../hooks/useSetMaxAmountToken"

const WithdrawAmountForm = ({
  onCancel,
  onSubmit,
  submitBtnText,
  withdrawAmount,
  withdrawalDelay, // <number> in seconds
  allowDecimals = false,
  ...formikProps
}) => {
  const onSubmitBtn = useCustomOnSubmitFormik(onSubmit)
  const onAddonClick = useSetMaxAmountToken("withdrawAmount", withdrawAmount)

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
        inputAddon={<MaxAmountAddon onClick={onAddonClick} text="Max Amount" />}
      />
      <SubmitButton
        className="btn btn-lg btn-primary w-100"
        onSubmitAction={onSubmitBtn}
      >
        {submitBtnText}
      </SubmitButton>
    </form>
  )
}

const WithdrawAmountFormWithFormik = withFormik({
  mapPropsToValues: () => ({
    withdrawAmount: "",
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
        KEEP.fromTokenUnit(1)
      )
    }

    return getErrorsObj(errors)
  },
  displayName: "CoveragePoolsWithdrawAmountForm",
})(WithdrawAmountForm)

export default WithdrawAmountFormWithFormik
