import web3Utils from "web3-utils"
import { KEEP } from "../utils/token.utils"

const REGEXP_NOT_BLANK_STRING = /^\s*$/
const REGEXP_ONLY_NUMBERS = /^\d+$/
const REGEX_DOT = /[.]/g

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

/**
 * Checks if the value is a valid decimal number.
 * To be a valid decimal number the value must fulfill three conditions:
 * - must have a dot (but not at the beginning)
 * - must have numbers only when we don't consider dots
 * - must have exactly one occurrence of a dot
 * @param {string} value given value
 * @return {boolean}
 */
const isValidDecimalNumber = (value) => {
  if (value.indexOf(".") <= 0) return false
  const valueWithoutDots = value.replace(".", "")
  return (
    REGEXP_ONLY_NUMBERS.test(valueWithoutDots) &&
    value.match(REGEX_DOT || []).length === 1
  )
}

export const validateAmountInRange = (
  value,
  maxValue,
  minValue = 0,
  token = KEEP
) => {
  let formattedValue = value
  try {
    formattedValue = value ? web3Utils.toWei(value) : 0
  } catch (err) {
    return "Invalid value"
  }
  const validateValueInBN = web3Utils.toBN(formattedValue)
  const maxValueInBN = web3Utils.toBN(maxValue || 0)
  const minValueInBN = web3Utils.toBN(minValue)

  if (isBlankString(value)) {
    return "Required"
  } else if (!REGEXP_ONLY_NUMBERS.test(value) && !isValidDecimalNumber(value)) {
    return "Invalid value"
  } else if (validateValueInBN.gt(maxValueInBN)) {
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