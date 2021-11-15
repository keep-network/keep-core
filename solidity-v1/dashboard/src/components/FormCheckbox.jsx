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
  withoutLabel = false,
  className = "",
  ...inputProps
}) => {
  return (
    <>
      <div className={`form-checkbox ${className}`}>
        <input id={name} checked={checked} {...inputProps} />
        <OnlyIf condition={!withoutLabel}>
          <label htmlFor={name}>{children}</label>
        </OnlyIf>
      </div>
      <OnlyIf condition={withoutLabel}>
        <div className={"form-checkbox__text"}>{children}</div>
      </OnlyIf>
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
