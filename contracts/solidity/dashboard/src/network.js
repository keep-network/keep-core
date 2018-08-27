import Web3 from 'web3'
import { sleep } from './utils'

const Network = {
  async web3() {
    const provider = await Network.provider()
    return new Web3(provider)
  },

  async eth() {
    const web3 = await Network.web3()
    return web3.eth
  },

  async getNetworkType() {
    const web3 = await Network.web3()
    return web3.eth.net.getNetworkType()
  },

  async provider() {
    let { web3 } = window

    while (web3 === undefined) {
      Network.log("Waiting for web3")
      await sleep(500)
      web3 = window.web3
    }

    return web3.currentProvider
  },

  getAccounts() {
    return new Promise((resolve, reject) => {
      Network.eth().then(eth => eth.getAccounts(Network._web3Callback(resolve, reject)))
    })
  },

  getCode(address) {
    return new Promise((resolve, reject) => {
      Network.eth().then(eth => eth.getCode(address, Network._web3Callback(resolve, reject)))
    })
  },

  _web3Callback(resolve, reject) {
    return (error, value) => {
      if (error) reject(error)
      else resolve(value)
    }
  },

  log(msg) {
    console.log(`[Network] ${msg}`)
  }
}

export default Network
