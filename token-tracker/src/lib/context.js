/** @typedef { import("./contract-helper.js").Contract } Contract */

import { initWeb3 } from "./ethereum-helper.js"

export default class Context {
  /** @property {Web3} web3 */

  /**
   * @typedef {Map<String,Contract>} Contracts
   * @type {Contracts}
   */
  #contracts

  /**
   * Deployment block is considered a starting
   * block number that can be used as 'fromBlock' for ethereum events lookup.
   * @type {Number}
   */
  #deploymentBlock

  /**
   * @param {Web3} web3
   */
  constructor(web3) {
    this.web3 = web3
    this.#contracts = new Map()
  }

  /**
   * Initializes context.
   *
   * @param {String} ethURL Ethereum API URL.
   * @param {String} ethPrivateKey Ethereum Account Private Key.
   * @return {Promise<Context>} Initialized Context instance.
   */
  static async initialize(ethURL, ethPrivateKey) {
    if (!ethURL) {
      throw new Error("ethereum rpc url not defined")
    }
    if (!ethPrivateKey) {
      throw new Error("ethereum account private key not defined")
    }

    const web3 = await initWeb3(ethURL, ethPrivateKey)

    return new Context(web3)
  }

  /**
   * @param {String} contractID Unique identifier of the contract.
   * @param {Contract} contract Contract.
   */
  addContract(contractID, contract) {
    if (this.#contracts.has(contractID)) {
      throw new Error(`contract with ID ${contractID} is already registered`)
    }

    this.#contracts.set(contractID, contract)
  }

  /**
   * @param {String} contractID
   * @return {Contract}
   */
  getContract(contractID) {
    if (!this.#contracts.has(contractID))
      throw new Error(`contract not registered: ${contractID}`)

    return this.#contracts.get(contractID)
  }

  /**
   * Sets deployment block to a number. Deployment block is considered a starting
   * block number that can be used as 'fromBlock' for ethereum events lookup.
   * @param {Number|String} value Block number.
   */
  set deploymentBlock(value) {
    this.#deploymentBlock = Number(value)
  }

  /**
   * Gets deployment block to a number. Deployment block is considered a starting
   * block number that can be used as 'fromBlock' for ethereum events lookup.
   * @return {Number} Block number.
   */
  get deploymentBlock() {
    return this.#deploymentBlock
  }
}
