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

export const formatValue = (
  value,
  decimalPlaces = 2,
  roundingType = BigNumber.ROUND_DOWN
) => {
  if (!value) return 0

  value = BigNumber.isBigNumber(value) ? value : new BigNumber(value)

  return value.decimalPlaces(decimalPlaces, roundingType).toNumber()
}

export const displayPercentageValue = (
  value,
  isFormattedValue = true,
  min = 0.01,
  max = 999
) => {
  if (!isFormattedValue) {
    value = formatValue(value)
  }

  if (value > 0 && value <= min) {
    return `<${min}%`
  } else if (value >= max) {
    return `>${max}%`
  }
  return `${value}%`
}

/**
 * Returns a number that is in the same spot on <yMin, yMax> number line that
 * inputX is on <xMin, xMax> number line
 *
 * e.g 2 on the number line of <1, 3> is exactly in the middle of the number
 * line of <1, 10> so it will return 5.5 (because 5.5 is exactly in the middle
 * of the number line of <1, 10>.
 *
 * @param {BigNumber|String|Number} inputX Number on number line
 * @param {BigNumber|String|Number} xMin Minimal value of the current number
 * line
 * @param {BigNumber|String|Number} xMax Maximal value of the current number
 * line
 * @param {BigNumber|String|Number} yMin Minimal value of the new number line
 * @param {BigNumber|String|Number} yMax Maximal value of the new number line
 * @return {BigNumber} Number from the new number line
 */
export const scaleInputForNumberRange = (inputX, xMin, xMax, yMin, yMax) => {
  inputX = BigNumber.isBigNumber(inputX) ? inputX : new BigNumber(inputX)
  xMin = BigNumber.isBigNumber(xMin) ? xMin : new BigNumber(xMin)
  xMax = BigNumber.isBigNumber(xMax) ? xMax : new BigNumber(xMax)
  yMin = BigNumber.isBigNumber(yMin) ? yMin : new BigNumber(yMin)
  yMax = BigNumber.isBigNumber(yMax) ? yMax : new BigNumber(yMax)

  const percent = inputX.minus(xMin).dividedBy(xMax.minus(xMin))
  return percent.multipliedBy(yMax.minus(yMin)).plus(yMin)
}

/**
 * Removes a substring from the input that occurs between the two same
 * characters on given index
 *
 * e.g to remove test3 from the input "test1/test2/test3/test4" we should call
 * this function like this:
 * removeSubstringBetweenCharacter("test1/test2/test3/test4", "/", 2)
 *
 * NOTE: if we keep the occurrenceIndex as 0 it will remove the word between
 * start of the input and first character occurrence. By analogy when we add max
 * possible index it will remove the word between last character occurrence and
 * the end of the input
 *
 * @param {string} input Input that we want to operate on
 * @param {string} character Character between which we want to cut out text
 * @param {Number} occurrenceIndex If there are more than two same character
 * in the input then we should specify the index to cut it out from (starts
 * from 0)
 * @return {String} Input with the word cut out
 */
export const removeSubstringBetweenCharacter = (
  input,
  character,
  occurrenceIndex = 0
) => {
  if (occurrenceIndex < 0) return input
  occurrenceIndex =
    input.charAt(0) === character ? occurrenceIndex + 1 : occurrenceIndex
  const splittedInput = input.split(character)
  splittedInput.splice(occurrenceIndex, 1)
  const final = splittedInput.join(character)
  return final
}

export const isString = (value) =>
  typeof value === "string" || value instanceof String

export const getSamePercentageValue = (givenValue, givenValueTotal, searchedValueTotal) => {
  if (
    new BigNumber(givenValue).isZero() ||
    new BigNumber(givenValueTotal).isZero() ||
    new BigNumber(searchedValueTotal).isZero()) {
    return "0"
  }
  return new BigNumber(givenValue).div(givenValueTotal)
    .multipliedBy(searchedValueTotal).toString()
}