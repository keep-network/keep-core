import web3Utils from 'web3-utils'

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
  const validateValueInBN = web3Utils.toBN(value || 0)
  const maxValueInBN = web3Utils.toBN(maxValue || 0)
  const minValueInBN = web3Utils.toBN(minValue)

  if (isBlankString(value)) {
    return 'Required'
  } else if (!REGEXP_ONLY_NUMBERS.test(value)) {
    return 'Invalid value'  
  } else if (validateValueInBN.gte(maxValueInBN)) {
    return 'You do not have enough KEEP tokens'
  } else if (minValueInBN.gte(validateValueInBN)) {
    return `You have to stake more than ${minValueInBN.toString()} KEEP tokens`
  }
}

export const getErrorsObj = (errors) => {
  return Object.keys(errors).every((name) => !errors[name]) ? {} : errors
}
