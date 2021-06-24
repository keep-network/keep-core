import React from "react"
import { withFormik } from "formik"
import { lte } from "../utils/arithmetics.utils"
import { getErrorsObj, validateAmountInRange } from "../forms/common-validators"
import {
  formatAmount as formatFormAmount,
  normalizeAmount,
} from "../forms/form.utils"
import FormInput from "./FormInput"
import MaxAmountAddon from "./MaxAmountAddon"
import { SubmitButton } from "./Button"
import { useCustomOnSubmitFormik } from "../hooks/useCustomOnSubmitFormik"

const WithdrawAmountForm = ({
  onCancel,
  onSubmit,
  submitBtnText,
  availableAmount,
  currentAmount,
  ...formikProps
}) => {
  const onSubmitBtn = useCustomOnSubmitFormik(onSubmit)

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
        normalize={normalizeAmount}
        format={formatFormAmount}
        inputAddon={
          <MaxAmountAddon
            onClick={() => console.log("on click addon")}
            text="Max Amount"
          />
        }
      />
      <SubmitButton
        className="btn btn-lg btn-primary w-100"
        onSubmitAction={onSubmitBtn}
      >
        withdraw
      </SubmitButton>
      <p
        className={
          "text-bold text-validation text-center withdraw-amount-form__button-text"
        }
      >
        <span className={"text-bold"}>14 days</span> cooldown period
      </p>
    </form>
  )
}

const WithdrawAmountFormWithFormik = withFormik({
  mapPropsToValues: () => ({
    amount: "",
  }),
  validate: ({ amount }, { availableAmount, minimumAmount }) => {
    const errors = {}

    if (lte(amount || 0, 0)) {
      errors.amount = "The value should be greater than zero"
    } else {
      errors.amount = validateAmountInRange(
        amount,
        availableAmount,
        minimumAmount
      )
    }

    return getErrorsObj(errors)
  },
  handleSubmit: (values, { props }) => props.onBtnClick(values),
  displayName: "CoveragePoolsWithdrawAmountForm",
})(WithdrawAmountForm)

export default WithdrawAmountFormWithFormik
