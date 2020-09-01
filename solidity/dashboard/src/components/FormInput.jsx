import React from "react"
import { useField } from "formik"
import SpeechBubbleTooltip from "./SpeechBubbleTooltip"
import { colors } from "../constants/colors"

const FormInput = ({ label, format, normalize, tooltipText, ...props }) => {
  const [field, meta, helpers] = useField(props.name, props.type)

  return (
    <div className="form-input flex flex-1 column">
      <label className="flex center">
        <span className="mr-1">{label}</span>
        {tooltipText && (
          <SpeechBubbleTooltip
            text={tooltipText}
            iconColor={colors.grey60}
            iconBackgroundColor="transparent"
          />
        )}
      </label>
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

FormInput.defaultProps = {
  tooltipText: null,
}

export default React.memo(FormInput)
