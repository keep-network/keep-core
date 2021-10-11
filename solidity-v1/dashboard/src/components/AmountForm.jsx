import React from "react"
import { withFormik } from "formik"
import FormInput from "./FormInput"
import Button from "./Button"
import { validateAmountInRange, getErrorsObj } from "../forms/common-validators"
import {
  normalizeAmount,
  formatAmount as formatFormAmount,
} from "../forms/form.utils.js"
import Divider from "./Divider"
import { lte } from "../utils/arithmetics.utils"
import TokenAmount from "./TokenAmount"

const styles = {
  divider: { margin: "2rem -2rem 0", padding: "2rem 2rem 0" },
  availableAmountWrapper: { marginTop: "-1rem", alignItems: "baseline" },
}

const AmountForm = ({
  onCancel,
  submitBtnText,
  availableAmount,
  currentAmount,
  ...formikProps
}) => {
  return (
    <>
      <form onSubmit={formikProps.handleSubmit} className="mt-1">
        <FormInput
          name="amount"
          type="text"
          label="KEEP Amount"
          placeholder="0"
          normalize={normalizeAmount}
          format={formatFormAmount}
        />
        <div className="flex row" style={styles.availableAmountWrapper}>
          <TokenAmount
            wrapperClassName="ml-a"
            amountClassName="text-caption--green-theme"
            symbolClassName="text-caption--green-theme"
            amount={availableAmount}
            withMetricSuffix
          />
          <span className="text-caption--green-theme">&nbsp;available.</span>
        </div>
        <Divider style={styles.divider} />
        <div className="flex row center">
          <Button
            className="btn btn-lg btn-primary"
            type="submit"
            disabled={!(formikProps.isValid && formikProps.dirty)}
          >
            {submitBtnText}
          </Button>
          <span onClick={onCancel} className="ml-2 text-link text-grey-70">
            Cancel
          </span>
        </div>
      </form>
    </>
  )
}

const AmountFormWithFormik = withFormik({
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
  displayName: "KEEPTokenAmountForm",
})(AmountForm)

export default AmountFormWithFormik
