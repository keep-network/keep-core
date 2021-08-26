import BigNumber from "bignumber.js"

const ONLY_NUMBERS = /[^0-9]+/g

export const normalizeAmount = (value) => value.replace(ONLY_NUMBERS, "")

export const formatAmount = (value) => {
  const newValue = value ? value.replace(ONLY_NUMBERS, "") : 0
  return new BigNumber(newValue).toFormat(0)
}

export const normalizeFloatingAmount = (value) =>
  value
    .replace(/[^0-9.]/g, "") // remove chars except number, point.
    .replace(/(\..*)\./g, "$1") // remove multiple points.
    .replace(/(?!^)-/g, "") // remove middle hyphen.
    .replace(/^0+(\d)/gm, "$1") // remove multiple leading zeros.
