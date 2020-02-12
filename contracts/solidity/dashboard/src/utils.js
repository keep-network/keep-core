import Web3 from 'web3'
import BigNumber from 'bignumber.js'
import moment from 'moment'

moment.updateLocale('en', {
  relativeTime: {
    d: '1 day',
    dd: (number, withoutSuffix, key, isFuture) => {
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
      .toFormat(precision)
  }
}

export function formatAmount(amount, decimals = 18) {
  amount = new BigNumber(amount)
  return amount.times(new BigNumber(10).pow(new BigNumber(decimals)))
}

export const getWeb3 = () => {
  if (window.ethereum || window.web3) {
    return new Web3(window.ethereum || window.web3.currentProvider)
  }

  return null
}

export const getWeb3SocketProvider = () => {
  return new Web3(process.env.REACT_APP_ETH_NETWORK_WEB_SOCKET_ADDRESS)
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
  const now = moment()
  const date = moment(dateMillis)

  if (now.isSame(date, 'year')) {
    return date.format('MMM DD')
  }
  return date.format('MMM DD YYYY')
}
