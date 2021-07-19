import React from "react"
import { useField } from "formik"
import OnlyIf from "./OnlyIf"

const FormCheckbox = ({ label, withError = false, ...props }) => {
  const [field, meta] = useField(props.name, props.type)

  return (
    <>
      <label>
        <input checked={field.value} {...field} {...props} />
        &nbsp;{label}
      </label>
      <OnlyIf condition={withError && meta.touched && meta.error}>
        <div className="form-error">{meta.error}</div>
      </OnlyIf>
    </>
  )
}

export default React.memo(FormCheckbox)
