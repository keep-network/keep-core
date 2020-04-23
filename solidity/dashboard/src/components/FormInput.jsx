import React from "react"
import { useField } from "formik"

const FormInput = ({ label, format, normalize, ...props }) => {
  const [field, meta, helpers] = useField(props.name, props.type)

  return (
    <div className="form-input flex flex-1 column">
      <label>{label}</label>
      <input
        {...field}
        {...props}
        onChange={(event) => {
          const value = event && event.target ? event.target.value : event
          helpers.setValue(normalize ? normalize(value) : value)
        }}
        value={format ? format(field.value) : field.value}
      />
      {meta.touched && meta.error ? (
        <div className="form-error">{meta.error}</div>
      ) : null}
    </div>
  )
}

export default React.memo(FormInput)
