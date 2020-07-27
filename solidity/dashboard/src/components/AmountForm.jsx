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
import * as Icons from "./Icons"
import Tag from "./Tag"
import ProgressBar from "./ProgressBar"
import { colors } from "../constants/colors"
import { displayAmount } from "../utils/token.utils"

const AmountForm = ({
  onCancel,
  submitBtnText,
  availableAmount,
  ...formikProps
}) => {
  return (
    <>
      <div className="flex row center">
        <Tag text="Current" IconComponent={Icons.KeepToken} />
        <h3 className="balance">10000 KEEP</h3>
      </div>
      <div className="flex row center mt-1">
        <Tag text="New" IconComponent={Icons.KeepToken} />
        <h3 className="balance">15000 KEEP</h3>
      </div>

      <form onSubmit={formikProps.handleSubmit}>
        <FormInput
          name="amount"
          type="text"
          label="Amount"
          placeholder="0"
          normalize={normalizeAmount}
          format={formatFormAmount}
        />
        <ProgressBar
          total={availableAmount}
          items={[{ value: formikProps.values.amount, color: colors.primary }]}
        />
        <div className="text-caption text-grey-50">
          {displayAmount(availableAmount)} KEEP available
        </div>
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
    </>
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
