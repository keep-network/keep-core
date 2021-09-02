import BigNumber from "bignumber.js"

const NOT_NUMBERS = /[^0-9]+/g

export const normalizeAmount = (value, allowDecimals = false) => {
  if (allowDecimals && value.includes(".") && value.indexOf(".") !== 0) {
    const [firstPart, ...rest] = value.split(".")
    const restValue = rest.join(".").replace(NOT_NUMBERS, "")
    return [firstPart, restValue].join(".")
  }
  return value.replace(NOT_NUMBERS, "")
}

export const formatAmount = (value, allowDecimals = false) => {
  if (allowDecimals && value.includes(".") && value.indexOf(".") !== 0) {
    const [firstPart, ...rest] = value.split(".")
    const restValue = rest.join(".").replace(NOT_NUMBERS, "")
    return [firstPart, restValue].join(".")
  }
  const newValue = value ? value.replace(NOT_NUMBERS, "") : 0
  return new BigNumber(newValue).toFormat(0)
}
