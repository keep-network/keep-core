import BigNumber from "bignumber.js"
import { AMOUNT_UNIT } from "../constants/constants"
import * as Icons from "../components/Icons"
import { isZero, lt } from "./arithmetics.utils"

const metrics = [
  { divisor: 1, symbol: "" },
  { divisor: 1e3, symbol: "K" },
  { divisor: 1e6, symbol: "M" },
]

export function displayAmount(
  amount,
  withCommaSeparator = true,
  unit = AMOUNT_UNIT.WEI
) {
  if (amount) {
    const readableFormat =
      unit === AMOUNT_UNIT.WEI ? toTokenUnit(amount) : new BigNumber(amount)
    return withCommaSeparator
      ? readableFormat.toFormat(0, BigNumber.ROUND_DOWN)
      : readableFormat.toString()
  }
  return 0
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

export const displayNumberWithMetricSuffix = (number) => {
  return getNumberWithMetricSuffix(number).formattedValue
}

export class Token {
  static toTokenUnit = (amount, decimals) => {
    if (!amount) {
      return new BigNumber(0)
    }
    return new BigNumber(amount).div(new BigNumber(10).pow(decimals))
  }

  static fromTokenUnit(amount, decimals = 18) {
    amount = new BigNumber(amount)
    return amount.times(new BigNumber(10).pow(decimals))
  }

  constructor(
    _name,
    _decimals,
    _symbol,
    _smallestPrecisionUnit,
    _smallestPrecisionDecimals,
    _icon,
    _decimalsToDisplay = 4
  ) {
    this.name = _name
    this.decimals = _decimals
    this.symbol = _symbol
    this.smallestPrecisionUnit = _smallestPrecisionUnit
    this.smallestPrecisionDecimals = _smallestPrecisionDecimals
    this.icon = _icon
    this.decimalsToDisplay = _decimalsToDisplay
  }

  toTokenUnit = (amount, decimals = this.decimals) => {
    return this.constructor.toTokenUnit(amount, decimals)
  }

  fromTokenUnit(amount, decimals = this.decimals) {
    return this.constructor.fromTokenUnit(amount, decimals)
  }

  displayAmount = (amount) => {
    if (!amount || isZero(amount)) {
      return "0"
    }

    const MIN_AMOUNT = new BigNumber(10)
      .pow(this.smallestPrecisionDecimals)
      .toString()

    if (lt(amount, MIN_AMOUNT)) {
      return `<${this.toTokenUnit(MIN_AMOUNT).toString()}`
    }

    return this.toFormat(this.toTokenUnit(amount))
  }

  toFormat = (amountInBn, decimalPlaces = this.decimalsToDisplay) => {
    return amountInBn.decimalPlaces() < decimalPlaces
      ? amountInBn.toFormat(undefined, BigNumber.ROUND_DOWN)
      : amountInBn.toFormat(decimalPlaces, BigNumber.ROUND_DOWN)
  }

  displayAmountWithSymbol = (amount) => {
    return `${this.displayAmount(amount)} ${this.symbol}`
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
  getNumberWithMetricSuffix = (number) => {
    const bigNumber = new BigNumber(number)
    let metric
    for (let i = metrics.length - 1; i >= 0; i--) {
      metric = metrics[i]
      if (bigNumber.gte(metric.divisor)) {
        break
      }
    }

    const value = bigNumber.div(metric.divisor)

    return {
      value: value.toString(),
      suffix: metric.symbol,
      formattedValue: `${this.toFormat(
        value,
        metric.divisor === 1 ? this.decimalsToDisplay : 2
      )}${metric.symbol}`,
    }
  }

  displayAmountWithMetricSuffix = (amount) => {
    if (!amount || isZero(amount)) {
      return "0"
    }

    const MIN_AMOUNT = new BigNumber(10)
      .pow(this.smallestPrecisionDecimals)
      .toString()

    if (lt(amount, MIN_AMOUNT)) {
      return `<${this.toTokenUnit(MIN_AMOUNT).toString()}`
    }

    return this.getNumberWithMetricSuffix(this.toTokenUnit(amount))
      .formattedValue
  }
}

export const KEEP = new Token(
  "Keep Token",
  18,
  "KEEP",
  "KEEP",
  18,
  Icons.KeepOutline,
  0
)

export const ETH = new Token("Ether", 18, "ETH", "gwei", 14, Icons.ETH)
export const TBTC = new Token("tBTC", 18, "TBTC", "tSats", 8, Icons.TBTC)
