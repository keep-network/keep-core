import React from "react"
import { useField } from "formik"
import OnlyIf from "./OnlyIf"

export const FormCheckboxBase = ({
  children,
  withError,
  hasError,
  errorMsg,
  name,
  checked,
  ...inputProps
}) => {
  return (
    <>
      <div className="form-checkbox">
        <input id={name} checked={checked} {...inputProps} />
        <label htmlFor={name}>{children}</label>
      </div>
      <OnlyIf condition={withError && hasError}>
        <div className="form-error">{errorMsg}</div>
      </OnlyIf>
    </>
  )
}

const FormCheckbox = ({ children, withError = false, ...props }) => {
  const [field, meta] = useField(props.name, props.type)

  return (
    <FormCheckboxBase
      withError={withError}
      hasError={meta.touched && meta.error}
      errorMsg={meta.error}
      checked={field.value}
      {...field}
      {...props}
    >
      {children}
    </FormCheckboxBase>
  )
}

export default React.memo(FormCheckbox)
