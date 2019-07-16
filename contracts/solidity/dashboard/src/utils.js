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
  // Fallback to localhost; use dev console port by default...
  else {
    const provider = new Web3.providers.HttpProvider(
      "http://127.0.0.1:9545"
    );
    web3 = new Web3(provider);
    console.log("No web3 instance injected, using Local web3.")
    return web3;
  }
}
