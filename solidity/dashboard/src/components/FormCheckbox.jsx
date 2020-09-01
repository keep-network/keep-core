import React from "react"
import { useField } from "formik"

const FormCheckbox = ({ label, ...props }) => {
  const [field, meta] = useField(props.name, props.type)

  return (
    <div className="form-input">
      <div className="flex row center">
        <label className="mr-1">{label}:</label>
        <input checked={field.value} {...field} {...props} />
      </div>
      {meta.touched && meta.error ? (
        <div className="form-error">{meta.error}</div>
      ) : null}
    </div>
  )
}

export default React.memo(FormCheckbox)
