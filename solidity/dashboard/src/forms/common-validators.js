import web3Utils from "web3-utils"
import { KEEP } from "../utils/token.utils"

const REGEXP_NOT_BLANK_STRING = /^\s*$/
const REGEXP_ONLY_NUMBERS = /^\d+$/

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

export const validateAmountInRange = (
  value,
  maxValue,
  minValue = 0,
  /** @type {import("../utils/token.utils").Token} */
  token = KEEP,
  isFloatingNumber = false
) => {
  /** @type {import("bignumber.js").BigNumber} */
  const formatedValue = token.fromTokenUnit(value)

  if (isBlankString(value)) {
    return "Required"
  } else if (
    (!isFloatingNumber && !REGEXP_ONLY_NUMBERS.test(value)) ||
    (isFloatingNumber && formatedValue.decimalPlaces() > 0)
  ) {
    return "Invalid value"
  }
  const validateValueInBN = web3Utils.toBN(formatedValue.toString())
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
