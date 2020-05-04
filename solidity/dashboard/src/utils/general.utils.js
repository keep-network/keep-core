import Web3 from "web3"
import moment from "moment"
import web3Utils from "web3-utils"

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

export const getWeb3 = () => {
  if (window.ethereum || window.web3) {
    return new Web3(window.ethereum || window.web3.currentProvider)
  }

  return null
}

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

export const formatDate = (dateMillis) => {
  const date = moment(dateMillis)

  return date.format("MM/DD/YYYY")
}

export const isEmptyObj = (obj) =>
  Object.keys(obj).length === 0 && obj.constructor === Object

export const isSameEthAddress = (address1, address2) => {
  return (
    web3Utils.toChecksumAddress(address1) ===
    web3Utils.toChecksumAddress(address2)
  )
}
