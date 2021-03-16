import React, { useLayoutEffect, useRef, useState } from "react"
import { useField } from "formik"
import * as Icons from "./Icons"
import Tooltip from "./Tooltip"

const iconDefaultValues = {
  width: 32,
  heigh: 32,
  marginRight: 1.5,
}

const FormInput = ({
  label,
  format,
  normalize,
  tooltipText,
  additionalInfoText,
  leftIcon,
  inputAddon,
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

  const leftIconComponent =
    leftIcon && React.isValidElement(leftIcon)
      ? React.cloneElement(leftIcon, {
          width: iconDefaultValues.width,
          height: iconDefaultValues.height,
          style: { marginRight: `${iconDefaultValues.marginRight}rem` },
        })
      : null

  const alignToInput = leftIconComponent
    ? {
        marginLeft: `calc(${iconDefaultValues.width}px + ${iconDefaultValues.marginRight}rem)`,
      }
    : {}
  return (
    <div className="form-input flex flex-1 column">
      <label className="input__label" style={alignToInput}>
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
      <div className="form-input__wrapper">
        {leftIconComponent}
        <input
          {...field}
          {...props}
          onChange={(event) => {
            const value = event && event.target ? event.target.value : event
            helpers.setValue(normalize ? normalize(value) : value)
          }}
          value={format ? format(field.value) : field.value}
          style={{ paddingRight: `${inputPaddingRight}px` }}
        />
        <div ref={inputAddonRef} className="form-input__addon">
          {inputAddon}
        </div>
      </div>
      {meta.touched && meta.error ? (
        <div className="form-input__error" style={alignToInput}>
          {meta.error}
        </div>
      ) : null}
    </div>
  )
}

FormInput.defaultProps = {
  tooltipText: null,
}

export default React.memo(FormInput)
