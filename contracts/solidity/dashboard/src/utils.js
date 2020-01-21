import Web3 from 'web3'
import BigNumber from 'bignumber.js'

export function displayAmount(amount, decimals, precision) {
  if (amount) {
    amount = new BigNumber(amount)
    return amount.div(new BigNumber(10).pow(new BigNumber(decimals))).toFixed(precision)
  }
}

export function formatAmount(amount, decimals) {
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
  const firstFourCharacters = address.substr(2, 4)
  const lastFourCharacters = address.substr(address.length - 4, address.length - 1)

  return '0x'.concat(firstFourCharacters).concat('...').concat(lastFourCharacters)
}

export const wait = (ms) => {
  return new Promise((resolve) => {
    return setTimeout(resolve, ms)
  })
}
