import BigNumber from "bignumber.js"

const ONLY_NUMBERS = /[^0-9]+/g

export const normalizeAmount = (value) => value.replace(ONLY_NUMBERS, "")

export const formatAmount = (value) => {
  const newValue = value ? value.replace(ONLY_NUMBERS, "") : 0
  return new BigNumber(newValue).toFormat(0)
}
