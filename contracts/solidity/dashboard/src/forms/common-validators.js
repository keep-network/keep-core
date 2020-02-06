import web3Utils from 'web3-utils'
import { displayAmount } from '../utils'

const REGEXP_NOT_BLANK_STRING = /^\s*$/
const REGEXP_ONLY_NUMBERS = /^\d+$/

const isBlankString = (value) => {
  return !value || REGEXP_NOT_BLANK_STRING.test(value)
}

export const validateEthAddress = (address, required = true) => {
  if (required && isBlankString(address)) {
    return 'Required'
  } else if (!web3Utils.isAddress(address)) {
    return 'Invalid eth address'
  }
}

export const validateAmountInRange = (value, maxValue, minValue = 0) => {
  const formatedValue = value ? web3Utils.toBN(value).mul(web3Utils.toBN(10).pow(web3Utils.toBN(18))).toString() : 0
  const validateValueInBN = web3Utils.toBN(formatedValue)
  const maxValueInBN = web3Utils.toBN(maxValue || 0)
  const minValueInBN = web3Utils.toBN(minValue)

  if (isBlankString(value)) {
    return 'Required'
  } else if (!REGEXP_ONLY_NUMBERS.test(value)) {
    return 'Invalid value'
  } else if (validateValueInBN.gt(maxValueInBN)) {
    return `The value should be less than or equals to ${displayAmount(maxValueInBN).toString()}`
  } else if (validateValueInBN.lt(minValueInBN)) {
    return `The value should be greater than or equlas to ${displayAmount(minValueInBN).toString()}`
  }
}

export const getErrorsObj = (errors) => {
  return Object.keys(errors).every((name) => !errors[name]) ? {} : errors
}
