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
import { displayAmount, fromTokenUnit } from "../utils/token.utils"
import { add } from "../utils/arithmetics.utils"

const AmountForm = ({
  onCancel,
  submitBtnText,
  availableAmount,
  currentAmount,
  ...formikProps
}) => {
  const { amount: formAmount } = formikProps.values || 0
  const newAmount =
    formAmount && displayAmount(add(fromTokenUnit(formAmount), currentAmount))

  return (
    <>
      <div className="flex row center mt-1">
        <div className="flex-1">
          <Tag text="Current" IconComponent={Icons.KeepToken} />
        </div>
        <h3 className="flex-2 text-primary">
          {displayAmount(currentAmount)} KEEP
        </h3>
      </div>
      <div className="flex row center mt-1">
        <div className="flex-1">
          <Tag text="New" IconComponent={Icons.KeepToken} />
        </div>
        {formAmount ? (
          <h3 className="flex-2 text-primary">
            {formAmount && `${newAmount} KEEP`}
          </h3>
        ) : (
          <span className="flex-2 text-big text-grey-40">
            Add an amount below
          </span>
        )}
      </div>
      <form onSubmit={formikProps.handleSubmit} className="mt-1">
        <FormInput
          name="amount"
          type="text"
          label="KEEP Amount"
          placeholder="0"
          normalize={normalizeAmount}
          format={formatFormAmount}
        />
        <ProgressBar
          styles={styles.progressBar}
          total={availableAmount}
          items={[
            {
              value: fromTokenUnit(formAmount),
              color: colors.primary,
            },
          ]}
        />
        <div className="text-caption text-grey-50">
          {displayAmount(availableAmount)} KEEP available
        </div>
        <div className="flex row center"></div>
        <Button
          className="btn btn-primary mt-1"
          type="submit"
          disabled={!(formikProps.isValid && formikProps.dirty)}
        >
          {submitBtnText}
        </Button>
        <span onClick={onCancel} className="mt-1 ml-1 text-link">
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
  validate: ({ amount }, { availableAmount, minimumAmount }) => {
    const errors = {}

    if (lte(amount || 0, 0)) {
      errors.amount = "Insufficient funds"
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

const styles = {
  progressBar: { margin: 0, marginTop: "-0.5rem" },
}

export default AmountFormWithFormik
