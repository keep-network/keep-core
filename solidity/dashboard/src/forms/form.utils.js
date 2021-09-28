import BigNumber from "bignumber.js"

const NOT_NUMBERS = /[^0-9]+/g

export const normalizeAmount = (value) => {
  return value.replace(NOT_NUMBERS, "")
}

export const formatAmount = (value) => {
  if (value.includes(".") && value.indexOf(".") !== 0) {
    const valuesSplitByDot = value.split(".")
    value = valuesSplitByDot[0]
  }
  const newValue = value ? value.replace(NOT_NUMBERS, "") : 0
  return new BigNumber(newValue).toFormat(0)
}

export const normalizeFloatingAmount = (value) =>
  value
    .replace(/[^0-9.]/g, "") // remove chars except number, point.
    .replace(/(\..*)\./g, "$1") // remove multiple points.
    .replace(/(?!^)-/g, "") // remove middle hyphen.
    .replace(/^0+(\d)/gm, "$1") // remove multiple leading zeros.

export const formatFloatingAmount = (value) => {
  if (value.includes(".") && value.indexOf(".") !== 0) {
    const [firstPart, ...rest] = value.split(".")
    const restValue = rest.join(".").replace(NOT_NUMBERS, "")
    const newValue = firstPart ? firstPart.replace(NOT_NUMBERS, "") : 0
    const firstPartFormatted = new BigNumber(newValue).toFormat(0)
    return [firstPartFormatted, restValue].join(".")
  }
  const newValue = value ? value.replace(NOT_NUMBERS, "") : 0
  return new BigNumber(newValue).toFormat(0)
}
