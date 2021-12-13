import React from "react"
import { withFormik } from "formik"
import { lte } from "../../utils/arithmetics.utils"
import {
  getErrorsObj,
  validateAmountInRange,
} from "../../forms/common-validators"
import {
  formatFloatingAmount,
  normalizeFloatingAmount,
} from "../../forms/form.utils"
import FormInput from "../FormInput"
import MaxAmountAddon from "../MaxAmountAddon"
import Button from "../Button"
import useSetMaxAmountToken from "../../hooks/useSetMaxAmountToken"
import { covKEEP } from "../../utils/token.utils"
import { COV_POOLS_FORMS_MAX_DECIMAL_PLACES } from "../../pages/coverage-pools/CoveragePoolPage"

const WithdrawAmountForm = ({
  onCancel,
  submitBtnText,
  withdrawAmount,
  withdrawalDelay, // <number> in seconds
  ...formikProps
}) => {
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
        normalize={(value) => {
          return normalizeFloatingAmount(
            value,
            COV_POOLS_FORMS_MAX_DECIMAL_PLACES
          )
        }}
        format={(value) => {
          return formatFloatingAmount(value, COV_POOLS_FORMS_MAX_DECIMAL_PLACES)
        }}
        inputAddon={<MaxAmountAddon onClick={onAddonClick} text="Max Amount" />}
        leftIconComponent={
          <span className={"form-input__left-icon__cov-keep-amount"}>
            covKEEP
          </span>
        }
      />
      <Button
        className="btn btn-lg btn-primary w-100"
        onClick={formikProps.handleSubmit}
      >
        {submitBtnText}
      </Button>
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
      errors.withdrawAmount = "Insufficient funds"
    } else {
      errors.withdrawAmount = validateAmountInRange(
        withdrawAmount,
        props.withdrawAmount,
        1000000000000,
        covKEEP,
        COV_POOLS_FORMS_MAX_DECIMAL_PLACES
      )
    }

    return getErrorsObj(errors)
  },
  handleSubmit: (values, { props, resetForm }) => {
    props.onSubmit(values)
    resetForm({ withdrawAmount: "0" })
  },
  displayName: "CoveragePoolsWithdrawAmountForm",
})(WithdrawAmountForm)

export default WithdrawAmountFormWithFormik
