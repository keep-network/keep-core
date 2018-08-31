import BigNumber from "bignumber.js"

export function displayAmount(amount, decimals, precision) {
  amount = new BigNumber(amount)
  return amount.dividedBy(10**decimals).toFixed(precision)
}

export function formatAmount(amount, decimals) {
  return amount * (10 ** decimals)
}

export function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms))
}
