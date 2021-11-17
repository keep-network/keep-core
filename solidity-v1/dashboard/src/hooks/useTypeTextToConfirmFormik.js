import { useFormik } from "formik"
import { getErrorsObj } from "../forms/common-validators"

export const useTypeTextToConfirmFormik = (confirmationText, onConfirm) => {
  return useFormik({
    validateOnBlur: true,
    initialValues: {
      confirmationText: "",
    },
    validate: (values) => {
      const errors = {}

      if (values.confirmationText !== confirmationText) {
        errors.confirmationText = "Does not match"
      }

      return getErrorsObj(errors)
    },
    onSubmit: (values) => {
      onConfirm(values)
    },
  })
}
