import { useFormik } from "formik"
import { getErrorsObj } from "../forms/common-validators"

export const useAcceptTermToConfirmFormik = (onConfirm) => {
  return useFormik({
    initialValues: {
      checked: false,
    },
    validate: (values) => {
      const errors = {}

      if (!values.checked) {
        errors.checked = "Required"
      }

      return getErrorsObj(errors)
    },
    onSubmit: (values) => {
      onConfirm(values)
    },
  })
}
