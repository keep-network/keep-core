import BigNumber from "bignumber.js"

const metrics = [
  { divisor: 1, symbol: "" },
  { divisor: 1e3, symbol: "K" },
  { divisor: 1e6, symbol: "M" },
]

export function displayAmount(amount, withCommaSeparator = true) {
  if (amount) {
    const readableFormat = toTokenUnit(amount)
    return withCommaSeparator
      ? readableFormat.toFormat(0, BigNumber.ROUND_DOWN)
      : readableFormat.toString()
  }
  return 0
}

export function displayAmountHigherOrderFn(withCommaSeparator = true, formatDecimalPlaces = 0) {
  return (amount) => {
    if (amount) {
      const readableFormat = toTokenUnit(amount)
      return withCommaSeparator
        ? readableFormat.toFormat(formatDecimalPlaces, BigNumber.ROUND_DOWN)
        : readableFormat.toString()
    }
    return 0
  }
}

/**
 * Convert wei amount to token units
 * @param {*} amount amount in wei
 *
 * @return {BigNumber} amount in token units
 */
export const toTokenUnit = (amount) => {
  if (!amount) {
    return new BigNumber(0)
  }
  return new BigNumber(amount).div(new BigNumber(10).pow(new BigNumber(18)))
}

/**
 * Convert sats amount to weitoshi, then to tBTC
 * @param {*} amount _utxoValue: "The size of the utxo in sat."
 *
 * @return {BigNumber} amount in weitoshi, the 18th decimal of BTC (https://docs.keep.network/tbtc/index.html)
 */
export const satsToTBtcViaWeitoshi = (amount) => {
  return toTokenUnit(fromTokenUnit(amount, 10))
}

/**
 * Convert token unit amount to wei.
 * @param {*} amount amount in token units
 * @param {*} decimals decimals
 *
 * @return {BigNumber} amount in wei
 */
export function fromTokenUnit(amount, decimals = 18) {
  amount = new BigNumber(amount)
  return amount.times(new BigNumber(10).pow(new BigNumber(decimals)))
}

/**
 * Returns a number with a metric suffix eg.:
 * * 10000000 => { value: 10, suffix: 'M', formattedValue: '10M' }
 * * 1000.4 => { value: 1.4, suffix: 'K', formattedValue: '1.4K'}
 * * 10000.4 => { value: 10, suffix: 'K', formattedValue: '10K' }
 * * 1 => { value: 1, suffix: '', formattedValue: '1' }
 *
 * @param {*} number
 *
 * @typedef {Object} NumebrWithSuffix
 * @property {string} value - number in string
 * @property {string} suffix - metric suffix
 * @property {string} formattedValue - number with suffix
 *
 * @return {NumebrWithSuffix}
 */
export const getNumberWithMetricSuffix = (number) => {
  const bigNumber = new BigNumber(number)
  let metric
  for (let i = metrics.length - 1; i >= 0; i--) {
    metric = metrics[i]
    if (bigNumber.gte(metric.divisor)) {
      break
    }
  }

  let value = bigNumber.div(metric.divisor)
  const beforeDecimal = value.toFraction(1)[0]
  const precision = beforeDecimal.toString().length === 1 ? 1 : 0
  value = bigNumber
    .div(metric.divisor)
    .toFormat(precision, BigNumber.ROUND_DOWN)

  return {
    value,
    suffix: metric.symbol,
    formattedValue: `${value}${metric.symbol}`,
  }
}

export const displayAmountWithMetricSuffix = (amount) => {
  return getNumberWithMetricSuffix(toTokenUnit(amount)).formattedValue
}
