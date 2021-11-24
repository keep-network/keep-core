import React, { useLayoutEffect, useRef, useState } from "react"
import { useField } from "formik"
import * as Icons from "./Icons"
import Tooltip from "./Tooltip"
import OnlyIf from "./OnlyIf"

const FormInputBase = ({
  label,
  format,
  normalize,
  tooltipText,
  additionalInfoText,
  leftIconComponent = null,
  inputAddon,
  tooltipProps = {},
  hasError,
  errorMsg,
  ...props
}) => {
  const [inputPaddingRight, setInputPaddingRight] = useState(0)
  const inputAddonRef = useRef(null)

  useLayoutEffect(() => {
    const inputAddonStyles = window.getComputedStyle(inputAddonRef.current)
    const finalWidth =
      parseInt(inputAddonStyles.right, 10) +
      parseInt(inputAddonStyles.width, 10) +
      10
    setInputPaddingRight(finalWidth)
  }, [])

  return (
    <div className="form-input flex flex-1 column">
      <label className="input__label">
        <span className="input__label__text">{label}</span>
        <div
          className={`input__label__info-container ${
            additionalInfoText ? "align-right" : ""
          }`}
        >
          {tooltipText && (
            <Tooltip
              simple
              delay={0}
              triggerComponent={Icons.MoreInfo}
              className="input__label__info-container__tooltip"
              {...tooltipProps}
            >
              {tooltipText}
            </Tooltip>
          )}
          {additionalInfoText && (
            <span className="input__label__info-container__additional-info-text">
              {additionalInfoText}
            </span>
          )}
        </div>
      </label>
      <div
        className="form-input__wrapper"
        style={{
          border: `${leftIconComponent ? "1px solid #0A0806" : "none"}`,
        }}
      >
        <OnlyIf condition={leftIconComponent}>
          <div className={"form-input__left-icon-container"}>
            {leftIconComponent}
          </div>
        </OnlyIf>
        <input
          {...props}
          onChange={(event) => {
            const value = event && event.target ? event.target.value : event
            props.onChange(event, normalize ? normalize(value) : value)
          }}
          value={format ? format(props.value) : props.value}
          style={{
            paddingRight: `${inputPaddingRight}px`,
            borderStyle: `${leftIconComponent ? "none" : "solid"}`,
          }}
        />
        <div ref={inputAddonRef} className="form-input__addon">
          {inputAddon}
        </div>
      </div>
      {hasError ? <div className="form-input__error">{errorMsg}</div> : null}
    </div>
  )
}

FormInputBase.defaultProps = {
  tooltipText: null,
}

const FormInput = (props) => {
  const [field, meta, helpers] = useField(props.name, props.type)

  return (
    <FormInputBase
      {...props}
      {...field}
      onChange={(_, value) => helpers.setValue(value)}
      hasError={meta.touched && meta.error}
      errorMsg={meta.error}
    />
  )
}

export default React.memo(FormInput)

export { FormInputBase }
