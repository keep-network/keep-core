import { useFormik } from "formik"
import { getErrorsObj, validateAmountInRange } from "../forms/common-validators"
import { lte } from "../utils/arithmetics.utils"
import { KEEP } from "../utils/token.utils"

export const useAddAmountFormik = (
  availableAmount,
  onSubmit,
  token = KEEP,
  minAmount = token.fromTokenUnit(1)
) => {
  return useFormik({
    validateOnBlur: true,
    initialValues: {
      amount: "0",
    },
    validate: ({ amount }) => {
      const errors = {}
      if (lte(availableAmount || 0, 0)) {
        errors.amount = "Insufficient funds"
      } else {
        errors.amount = validateAmountInRange(
          amount,
          availableAmount,
          minAmount
        )
      }

      return getErrorsObj(errors)
    },
    onSubmit: (values) => {
      onSubmit(values)
    },
  })
}
