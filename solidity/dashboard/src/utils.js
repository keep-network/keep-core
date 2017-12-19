export function displayAmount(amount, decimals) {
  amount = amount / (10 ** decimals)
  return Math.round(amount * 10000) / 10000
}

export function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms))
}
