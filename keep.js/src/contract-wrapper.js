import { lookupArtifactAddress } from "./utils.js"

/**
 * A wrapper for Web3.Contract {@link https://web3js.readthedocs.io/en/v1.2.9/web3-eth-contract.html#}
 */
export class ContractWrapper {
  constructor(instance, deployedAtBlock) {
    this.deployedAtBlock = deployedAtBlock
    this.instance = instance
  }

  /**
   * Calls the provided contract function with the provided arguments.
   *
   * @param {string} contractMethodName Name of the contract function.
   * @param  {...any} args Arguments of function.
   *
   * @return {Promise<any>} response
   */
  async makeCall(contractMethodName, ...args) {
    return await this.instance.methods[contractMethodName](...args).call()
  }

  /**
   * Returns send object of the provided contract function with the provided arguments.
   *
   * @param {string} contractMethodName Name of the contract function.
   * @param  {...any} args Arguments of function.
   *
   * @return {any} Send object from Web3
   */
  sendTransaction(contractMethodName, ...args) {
    return this.instance.methods[contractMethodName](...args).send
  }

  /**
   * Returns a past events for the provided event.
   *
   * @param {string} eventName Name of the contract event.
   * @param {*} filter Let you filter events by indexed parameters.
   * @param {number} fromBlock The block number (greater than or equal to) from which to get events on.
   *
   * @return {Promise<Arra<EventData>>}
   */
  async getPastEvents(eventName, filter, fromBlock = this.deployedAtBlock) {
    const searchFilter = { fromBlock, filter }
    return await this.instance.getPastEvents(eventName, searchFilter)
  }

  /**
   * @return {string} The address used for this contract instance.
   */
  get address() {
    return this.instance.options.address
  }

  /**
   * {@link https://web3js.readthedocs.io/en/v1.2.9/web3-eth-contract.html#id26}
   * @retun {any} Contract methods.
   */
  get methods() {
    return this.instance.methods
  }
}

class ContractFactory {
  /**
   * @typedef {Object} Config
   * @property {Object} web3 Web3 instance {@link https://web3js.readthedocs.io/en/v1.2.9/index.html}
   * @property {number} networkId Network id.
   */

  /**
   * Creates a contract wrapper instance.
   *
   * @param {*} artifact Artifacts of the provided contract.
   * @param {Config} config
   *
   * @return {ContractWrapper} A contract wrapper for the provided contract.
   */
  static async createContractInstance(artifact, config) {
    const { web3, networkId } = config

    const contractDeployedAtBlock = async () => {
      const deployTransactionHash = artifact.networks[networkId].transactionHash
      const transaction = await web3.eth.getTransaction(deployTransactionHash)

      return transaction.blockNumber.toString()
    }

    const address = lookupArtifactAddress(artifact)
    const deployedAtBlock = await contractDeployedAtBlock()
    const instance = new web3.eth.Contract(artifact.abi, address)
    instance.options.from = web3.eth.defaultAccount

    return new ContractWrapper(instance, deployedAtBlock)
  }

  static new(instance) {
    return new ContractWrapper(instance)
  }
}

/**
 * @typedef {Object} EventData
 * @property {Object} returnValues
 * @property {Object} raw
 * @property {string} raw.data
 * @property {string[]} raw.topics
 * @property {string} event
 * @property {string}signature
 * @property {number} logIndex
 * @property {number} transactionIndex
 * @property {string} transactionHash
 * @property {string} blockHash
 * @property {number} blockNumber
 * @property {string} address
 */

export default ContractFactory
