/** @typedef { import("../../web3").Web3LibWrapper} Web3LibWrapper */

class BaseContract {
  /**
   *
   * @param {any} _instance The Contract instance of a given web3 library.
   * @param {string} _deploymentTxnHash The Hash of the deployemnt transaction.
   * @param {Web3LibWrapper} _web3 The wrapper of the web3 library.
   * @param {string | null} _deployedAtBlock The block number, when the contract
   * was deployed. If null, the `deployedAtBlock` is set after the
   * `_getDeploymentBlock` is called.
   */
  constructor(_instance, _deploymentTxnHash, _web3, _deployedAtBlock = null) {
    this.instance = _instance
    this.deploymentTxnHash = _deploymentTxnHash
    this.web3Wrapper = _web3
    this.deployedAtBlock = _deployedAtBlock
  }

  /**
   * Calls the provided contract function with the provided arguments.
   *
   * @param {string} methodName Name of the contract function.
   * @param  {...any} args Arguments of function.
   *
   * @return {Promise<any>} response
   */
  makeCall = async (methodName, ...args) => {
    return await this._makeCall(methodName, ...args)
  }

  /**
   * Returns a promise combined event emitter. Resolves when the transaction
   * receipt is available. PromiEvents work like a normal promises with added
   * `on`, `once` and `off` functions. This way developers can watch for
   * additional events like on `receipt` or `transactionHash`.
   *
   * @param {string} methodName Name of the contract function.
   * @param  {...any} args Arguments of function.
   *
   * @return {any} PromiEvent object.
   */
  sendTransaction = async (methodName, ...args) => {
    return this._sendTransaction(methodName, ...args)
  }

  /**
   * Returns a past events for the provided event name.
   *
   * @param {string} eventName Name of the contract event.
   * @param {*} filter Let you filter events by indexed parameters.
   * @param {number} fromBlock The block number (greater than or equal to) from
   * which to get events on. By default, it refers to `this.deployedAtBlock`.
   *
   * @return {Promise<Array<any>>}
   */
  getPastEvents = async (
    eventName,
    filter,
    fromBlock = this.deployedAtBlock
  ) => {
    const _fromBlock = !fromBlock ? await this._getDeploymentBlock() : fromBlock

    return await this._getPastEvents(eventName, filter, _fromBlock)
  }

  /**
   * @return {string} The address used for this contract instance.
   */
  get address() {
    return this._address
  }

  /**
   * @return {Web3LibWrapper} The web3 lib wrapper.
   */
  get web3() {
    return this.web3Wrapper
  }

  set defaultAccount(defaultAccount) {
    this._defaultAccount = defaultAccount
  }

  get defaultAccount() {
    return this._defaultAccount
  }

  async _getDeploymentBlock() {
    if (!this.deployedAtBlock) {
      const transaction = await this.web3.getTransaction(this.deploymentTxnHash)
      this.deployedAtBlock = transaction.blockNumber.toString()
    }

    return this.deployedAtBlock
  }
}

/**
 * A wrapper for Web3.js Contract object {@link https://web3js.readthedocs.io/en/v1.3.4/web3-eth-contract.html}.
 */
class Web3jsContractWrapper extends BaseContract {
  constructor(_instance, _deploymentTxnHash, _web3, _deployedAtBlock = null) {
    super(_instance, _deploymentTxnHash, _web3, _deployedAtBlock)
  }

  async _makeCall(methodName, ...args) {
    return await this.instance.methods[methodName](...args).call()
  }

  _sendTransaction(methodName, ...args) {
    return this.instance.methods[methodName](...args).send()
  }

  async _getPastEvents(eventName, filter, fromBlock) {
    const searchFilter = { fromBlock, filter }
    return await this.instance.getPastEvents(eventName, searchFilter)
  }

  get _address() {
    return this.instance.options.address
  }

  /**
   * Contract's methods {@link https://web3js.readthedocs.io/en/v1.3.4/web3-eth-contract.html#id26}
   */
  get methods() {
    return this.instance.methods
  }

  /**
   * Contract's events {@link https://web3js.readthedocs.io/en/v1.3.4/web3-eth-contract.html#contract-events}
   */
  get events() {
    return this.instance.events
  }

  set _defaultAccount(defaultAccount) {
    this.instance.options.defaultAccount = defaultAccount
  }

  get _defaultAccount() {
    return this.instance.options.defaultAccount
  }
}

class ContractFactory {
  static createWeb3jsContract(instance, deploymentTxnHash, web3Wrapper) {
    return new Web3jsContractWrapper(instance, deploymentTxnHash, web3Wrapper)
  }
}

/**
 * @typedef {Object} EventData
 * @property {Object} returnValues
 * @property {Object} raw
 * @property {string} raw.data
 * @property {string[]} raw.topics
 * @property {string} event
 * @property {string} signature
 * @property {number} logIndex
 * @property {number} transactionIndex
 * @property {string} transactionHash
 * @property {string} blockHash
 * @property {number} blockNumber
 * @property {string} address
 */

export default ContractFactory

export { BaseContract, Web3jsContractWrapper }
