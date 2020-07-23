import React from "react"
import { withFormik } from "formik"
import FormInput from "./FormInput"
import Button from "./Button"
import { validateAmountInRange, getErrorsObj } from "../forms/common-validators"
import { lte } from "../utils/arithmetics.utils"
import {
  normalizeAmount,
  formatAmount as formatFormAmount,
} from "../forms/form.utils.js"

const AmountForm = ({ onCancel, submitBtnText, ...formikProps }) => {
  return (
    <form onSubmit={formikProps.handleSubmit}>
      <FormInput
        name="amount"
        type="text"
        label="Amount"
        placeholder="0"
        normalize={normalizeAmount}
        format={formatFormAmount}
      />
      <Button
        className="btn btn-primary"
        type="submit"
        disabled={!(formikProps.isValid && formikProps.dirty)}
      >
        {submitBtnText}
      </Button>
      <span onClick={onCancel} className="ml-1 text-link">
        Cancel
      </span>
    </form>
  )
}

const AmountFormWithFormik = withFormik({
  mapPropsToValues: () => ({
    amount: "",
  }),
  validate: ({ amount }, { availableAmount }) => {
    const errors = {}

    if (lte(amount || 0, 0)) {
      errors.amount = "Insufficient funds"
    } else {
      errors.amount = validateAmountInRange(amount, availableAmount)
    }

    return getErrorsObj(errors)
  },
  handleSubmit: (values, { props }) => props.onBtnClick(values),
  displayName: "KEEPTokenAmountForm",
})(AmountForm)

export default AmountFormWithFormik
