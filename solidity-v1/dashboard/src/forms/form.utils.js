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

export const normalizeFloatingAmount = (value, maxDecimals = null) => {
  let normalizedValue = value
    .replace(/[^0-9.]/g, "") // remove chars except number, point.
    .replace(/(\..*)\./g, "$1") // remove multiple points.
    .replace(/(?!^)-/g, "") // remove middle hyphen.
    .replace(/^0+(\d)/gm, "$1") // remove multiple leading zeros.
  // remove all decimals starting at {maxDecimals} index
  // eg if maxDecimals = 6 then we change 1.079321435 -> 1.079321
  if (
    normalizedValue.includes(".") &&
    normalizedValue.indexOf(".") !== 0 &&
    maxDecimals
  ) {
    const [firstPart, ...rest] = normalizedValue.split(".")
    let restValue = rest.join(".").replace(NOT_NUMBERS, "")
    if (maxDecimals) restValue = restValue.slice(0, maxDecimals)
    normalizedValue = [firstPart, restValue].join(".")
  }
  return normalizedValue
}

export const formatFloatingAmount = (value, maxDecimals = null) => {
  if (value.includes(".") && value.indexOf(".") !== 0) {
    const [firstPart, ...rest] = value.split(".")
    let restValue = rest.join(".").replace(NOT_NUMBERS, "")
    if (maxDecimals) restValue = restValue.slice(0, maxDecimals)
    const newValue = firstPart ? firstPart.replace(NOT_NUMBERS, "") : 0
    const firstPartFormatted = new BigNumber(newValue).toFormat(0)
    return [firstPartFormatted, restValue].join(".")
  }
  const newValue = value ? value.replace(NOT_NUMBERS, "") : 0
  return new BigNumber(newValue).toFormat(0)
}
