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

export const getWeb3 = async () => {

  if (web3) return web3

  // Modern dapp browsers...
  if (window.ethereum) {
    web3 = new Web3(window.ethereum)
    try {
      // Request account access if needed
      await window.ethereum.enable()
      // Acccounts now exposed
      return web3;
    } catch (error) {
      return error;
    }
  }
  // Legacy dapp browsers...
  else if (window.web3) {
    // Use Mist/MetaMask's provider.
    web3 = new Web3(window.web3.currentProvider)
    console.log("Injected web3 detected.")
    return web3;
  }
  
  return null;
}

export const getWeb3SocketProvider = () => {
  return new Web3(process.env.REACT_APP_ETH_NETWORK_WEB_SOCKET_ADDRESS)
}