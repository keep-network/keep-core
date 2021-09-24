import web3Utils from "web3-utils"
import { KEEP } from "../utils/token.utils"

const REGEXP_NOT_BLANK_STRING = /^\s*$/

const isBlankString = (value) => {
  return !value || REGEXP_NOT_BLANK_STRING.test(value)
}

export const validateRequiredValue = (value) => {
  if (isBlankString(value)) {
    return "Required"
  }
}

export const validateEthAddress = (address, required = true) => {
  if (required && isBlankString(address)) {
    return "Required"
  } else if (!web3Utils.isAddress(address)) {
    return "Invalid eth address"
  }
}

export const isNumeric = (value) => {
  return !Number.isNaN(parseFloat(value)) && isFinite(value, "Is not finite")
}

export const validateAmountInRange = (
  value,
  maxValue,
  minValue = 0,
  token = KEEP
) => {
  const formattedValue = token.fromTokenUnit(value)

  if (isBlankString(value)) {
    return "Required"
  } else if (!isNumeric(value)) {
    return "Invalid value"
  }

  const validateValueInBN = web3Utils.toBN(formattedValue.toString())
  const maxValueInBN = web3Utils.toBN(maxValue.toString() || 0)
  const minValueInBN = web3Utils.toBN(minValue.toString())

  if (validateValueInBN.gt(maxValueInBN)) {
    return `The value should be less than or equal ${token.displayAmount(
      maxValueInBN.toString()
    )}`
  } else if (validateValueInBN.lt(minValueInBN)) {
    return `The value should be greater than or equal ${token.displayAmount(
      minValueInBN.toString()
    )}`
  }
}

export const getErrorsObj = (errors) => {
  return Object.keys(errors).every((name) => !errors[name]) ? {} : errors
}
