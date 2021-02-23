/** @typedef { import("./contract-helper").ContractArtifact } ContractArtifact */

import ProviderEngine from "web3-provider-engine"
import CacheSubprovider from "web3-provider-engine/subproviders/cache.js"
import FilterSubprovider from "web3-provider-engine/subproviders/filters.js"
import NonceSubprovider from "web3-provider-engine/subproviders/nonce-tracker.js"

import WebsocketSubprovider from "web3-provider-engine/subproviders/websocket.js"

/** @typedef { import("web3").default } Web3 */
import Web3 from "web3"

import { getDeploymentBlockNumber, Contract } from "./contract-helper.js"

import KeepTokenJson from "@keep-network/keep-core/artifacts/KeepToken.json"

/**
 * @typedef {Object} ContractsArtifacts
 * @property {Number} deploymentBlock
 * @property {Contract} KeepToken
 * @property {Contract} TokenStaking
 */

/**
 * @typedef {Object} Context
 * @property {Web3} web3 Description.
 * @property {ContractsArtifacts} contracts Description.
 */

export default class Context {
  /**
   * @param {Web3} web3 Description.
   * @param {ContractsArtifacts} contracts Description.
   */
  constructor(web3, contracts) {
    this.web3 = web3
    this.contracts = contracts
  }

  static async initialize() {
    if (!process.env.ETH_HOSTNAME) {
      throw new Error("ETH_HOSTNAME not defined")
    }
    if (!process.env.ETH_ACCOUNT_PRIVATE_KEY) {
      throw new Error("ETH_ACCOUNT_PRIVATE_KEY not defined")
    }

    const ethUrl = process.env.ETH_HOSTNAME
    const ethPrivateKey = process.env.ETH_ACCOUNT_PRIVATE_KEY

    const web3 = await initWeb3(ethUrl, ethPrivateKey)

    const KeepToken = new Contract(KeepTokenJson, web3)

    const deploymentBlock = await getDeploymentBlockNumber(KeepTokenJson, web3)

    console.debug("deploymentBlock", deploymentBlock)

    const contracts = {
      KeepToken: KeepToken,
      deploymentBlock: deploymentBlock,
    }

    return new Context(web3, contracts)
  }

  /**
   * @param {String} contractID Unique identifier of the contract instance.
   * @param {Contract} contract Contract.
   */
  addContract(contractID, contract) {
    if (this.contracts[contractID]) {
      throw new Error(`contract with ID ${contractID} is already registered`)
    }

    this.contracts[contractID] = contract
  }
}

/**
 * @param {string} url
 * @param {string} ethPrivateKey
 */
async function initWeb3(url, ethPrivateKey) {
  const engine = new ProviderEngine()
  const web3 = new Web3(engine)

  // cache layer
  engine.addProvider(new CacheSubprovider())

  // filters
  engine.addProvider(new FilterSubprovider())

  // pending nonce
  engine.addProvider(new NonceSubprovider())

  // // vm
  // engine.addProvider(new VmSubprovider())

  engine.addProvider(
    new WebsocketSubprovider({
      rpcUrl: url,
      // origin: undefined,
    })
  )

  // network connectivity error
  engine.on("error", function (err) {
    // report connectivity errors
    console.error(err.stack)
  })

  engine.start()

  // const web3Provider = new Web3.providers.WebsocketProvider(url)
  // const web3 = new Web3(web3Provider)

  const account = web3.eth.accounts.privateKeyToAccount(ethPrivateKey)
  web3.eth.accounts.wallet.add(account)

  web3.eth.defaultAccount = account.address

  return web3
}
