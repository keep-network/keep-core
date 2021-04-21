import moment from "moment"
import web3Utils from "web3-utils"
import BigNumber from "bignumber.js"

moment.updateLocale("en", {
  relativeTime: {
    d: "1 day",
    dd: (number, withoutSuffix, key, isFuture) => {
      const weeks = Math.round(number / 7)
      if (number < 7) {
        return number + " days"
      } else {
        return weeks + " week" + (weeks === 1 ? "" : "s")
      }
    },
  },
})

export const shortenAddress = (address) => {
  if (!address) {
    return ""
  }
  const firstFourCharacters = address.substr(2, 4)
  const lastFourCharacters = address.substr(
    address.length - 4,
    address.length - 1
  )

  return "0x"
    .concat(firstFourCharacters)
    .concat("...")
    .concat(lastFourCharacters)
}

export const wait = (ms) => {
  return new Promise((resolve) => {
    return setTimeout(resolve, ms)
  })
}

export const formatDate = (dateMillis, format = "MM/DD/YYYY") => {
  const date = moment(dateMillis)

  return date.format(format)
}

export const isEmptyObj = (obj) =>
  Object.keys(obj).length === 0 && obj.constructor === Object

export const isSameEthAddress = (address1, address2) => {
  return (
    web3Utils.toChecksumAddress(address1) ===
    web3Utils.toChecksumAddress(address2)
  )
}

export const getBufferFromHex = (hex) => {
  const validHex = toValidHex(hex).toLowerCase()
  return new Buffer(validHex, "hex")
}

const toValidHex = (hex) => {
  hex = hex.substring(0, 2) === "0x" ? hex.substring(2) : hex
  if (hex === "") {
    return ""
  }
  return hex.length % 2 !== 0 ? `0${hex}` : hex
}

export const formatPercentage = (value, decimalPlaces = 2) => {
  if (!value) return 0

  value = BigNumber.isBigNumber(value) ? value : new BigNumber(value)

  return value.decimalPlaces(decimalPlaces, BigNumber.ROUND_DOWN).toNumber()
}

export const displayPercentageValue = (
  value,
  isFormattedValue = true,
  min = 0.01,
  max = 999
) => {
  if (!isFormattedValue) {
    value = formatPercentage(value)
  }

  let prefix = ""
  if (value > 0 && value <= min) {
    prefix = `<`
  } else if (value >= max) {
    prefix = `>`
  }
  return `${prefix}${value}%`
}

export const scaleInputForNumberRange = (inputX, xMin, xMax, yMin, yMax) => {
  inputX = BigNumber.isBigNumber(inputX) ? inputX : new BigNumber(inputX)
  xMin = BigNumber.isBigNumber(xMin) ? xMin : new BigNumber(xMin)
  xMax = BigNumber.isBigNumber(xMax) ? xMax : new BigNumber(xMax)
  yMin = BigNumber.isBigNumber(yMin) ? yMin : new BigNumber(yMin)
  yMax = BigNumber.isBigNumber(yMax) ? yMax : new BigNumber(yMax)

  const percent = inputX.minus(xMin).dividedBy(xMax.minus(xMin))
  return percent.multipliedBy(yMax.minus(yMin)).plus(yMin)
}
