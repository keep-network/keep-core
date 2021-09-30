import React, { useLayoutEffect, useRef, useState } from "react"
import { useField } from "formik"
import * as Icons from "./Icons"
import Tooltip from "./Tooltip"
import OnlyIf from "./OnlyIf"

const FormInput = ({
  label,
  format,
  normalize,
  tooltipText,
  additionalInfoText,
  leftIconComponent = null,
  inputAddon,
  tooltipProps = {},
  ...props
}) => {
  const [field, meta, helpers] = useField(props.name, props.type)
  const inputAddonRef = useRef(null)
  const [inputPaddingRight, setInputPaddingRight] = useState(0)

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
          {...field}
          {...props}
          onChange={(event) => {
            const value = event && event.target ? event.target.value : event
            helpers.setValue(normalize ? normalize(value) : value)
          }}
          value={format ? format(field.value) : field.value}
          style={{
            paddingRight: `${inputPaddingRight}px`,
            borderStyle: `${leftIconComponent ? "none" : "solid"}`,
          }}
        />
        <div ref={inputAddonRef} className="form-input__addon">
          {inputAddon}
        </div>
      </div>
      {meta.touched && meta.error ? (
        <div className="form-input__error">{meta.error}</div>
      ) : null}
    </div>
  )
}

FormInput.defaultProps = {
  tooltipText: null,
}

export default React.memo(FormInput)
