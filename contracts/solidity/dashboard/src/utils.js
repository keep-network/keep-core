import Web3 from 'web3'
import BigNumber from "bignumber.js"

export function displayAmount(amount, decimals, precision) {
  amount = new BigNumber(amount)
  return amount.dividedBy(10**decimals).toFixed(precision)
}

export function formatAmount(amount, decimals) {
  return amount * (10 ** decimals)
}

export function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms))
}

export function getWeb3() {
  return new Promise(resolve => {
    window.addEventListener('load', function() {
      let web3 = window.web3
      if (typeof web3 !== 'undefined') {
        resolve(new Web3(web3.currentProvider));
      }
      resolve(null)
    })
  })
}
