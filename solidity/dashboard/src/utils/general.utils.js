import BigNumber from 'bignumber.js'
import moment from 'moment'
import { PENDING_STATUS, COMPLETE_STATUS } from '../constants/constants'
import web3Utils from 'web3-utils'

moment.updateLocale('en', {
  relativeTime: {
    d: '1 day',
    dd: (number) => {
      const weeks = Math.round(number / 7)
      if (number < 7) {
        return number + ' days'
      } else {
        return weeks + ' week' + (weeks === 1 ? '' : 's')
      }
    },
  },
})

export function displayAmount(amount, decimals = 18, precision = 0) {
  if (amount) {
    return new BigNumber(amount)
      .div(new BigNumber(10).pow(new BigNumber(decimals)))
      .toFormat(precision, BigNumber.ROUND_DOWN)
  }
}

export function formatAmount(amount, decimals = 18) {
  amount = new BigNumber(amount)
  return amount.times(new BigNumber(10).pow(new BigNumber(decimals)))
}

export const shortenAddress = (address) => {
  if (!address) {
    return ''
  }
  const firstFourCharacters = address.substr(2, 4)
  const lastFourCharacters = address.substr(address.length - 4, address.length - 1)

  return '0x'.concat(firstFourCharacters).concat('...').concat(lastFourCharacters)
}

export const wait = (ms) => {
  return new Promise((resolve) => {
    return setTimeout(resolve, ms)
  })
}

export const formatDate = (dateMillis) => {
  const date = moment(dateMillis)

  return date.format('MM/DD/YYYY')
}

export const isEmptyObj = (obj) => Object.keys(obj).length === 0 && obj.constructor === Object

export const getAvailableAtBlock = (blockNumber, status) => {
  if (status === PENDING_STATUS) {
    return `until ${blockNumber} block`
  } else if (status === COMPLETE_STATUS) {
    return `at ${blockNumber} block`
  }
}

export const isSameEthAddress = (address1, address2) => {
  return web3Utils.toChecksumAddress(address1) === web3Utils.toChecksumAddress(address2)
}
