import Web3 from 'web3'
import BigNumber from "bignumber.js"

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

export function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms))
}

let web3;

export const getWeb3 = () => {

  if (web3) return web3

  if (window.ethereum || window.web3)
    return new Web3(window.ethereum || window.web3.currentProvider)
  
  return null;
}

export const getWeb3SocketProvider = () => {
  return new Web3(process.env.REACT_APP_ETH_NETWORK_WEB_SOCKET_ADDRESS)
}