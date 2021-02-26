/** @typedef { import("web3").default } Web3 */
/** @typedef { import("web3-eth-contract").Contract } Web3Contract */
/** @typedef { import("@keep-network/tbtc.js/src/EthereumHelpers.js").TruffleArtifact} TruffleArtifact */
/** @typedef { import("@keep-network/tbtc.js/src/EthereumHelpers.js").Contract } ContractInstance */

import { EthereumHelpers } from "@keep-network/tbtc.js"

import { getChainID } from "./ethereum-helper.js"
import { logger } from "./winston.js"
export class Contract {
  /**
   * @param {TruffleArtifact} artifact The Truffle artifact for the deployed
   * contract.
   * @param {Web3} web3 The Web3 instance to instantiate the contract on.
   */
  constructor(artifact, web3) {
    this.artifact = artifact
    this.web3 = web3
  }

  /**
   * Gets the Web3 Contract instance deployed to the current network under
   * address provided in contract artifact. Throws if the artifact does not
   * contain deployment information for the specified network id.
   *
   * @return {Promise<ContractInstance>} A contract ready for usage with web3 for the
   * given network and artifact.
   */
  async deployed() {
    const { getDeployedContract } = EthereumHelpers

    return getDeployedContract(
      this.artifact,
      this.web3,
      await getChainID(this.web3)
    )
  }

  /**
   *
   * @param {string} address Address at which contract was deployed.
   * @return {Promise<ContractInstance>} A contract ready for usage.
   */
  async at(address) {
    const { buildContract } = EthereumHelpers

    return buildContract(this.web3, this.artifact.abi, address)
  }
}

/**
 * Fetches deployment block number from a contract artifact.
 *
 * @param {TruffleArtifact} artifact Truffle deployment contract artifact.
 * @param {Web3} web3 Web3 instance.
 * @return {Number} Block number at which contract was deployed.
 */
export async function getDeploymentBlockNumber(artifact, web3) {
  const chainID = await getChainID(web3)
  logger.debug(`chain id: ${chainID}`)

  if (!artifact.networks[chainID]) {
    throw new Error(`artifact does not define network ${chainID}`)
  }
  if (!artifact.networks[chainID].transactionHash) {
    throw new Error(`missing transaction hash for network ${chainID}`)
  }

  const transactionHash = artifact.networks[chainID].transactionHash

  const transaction = await web3.eth.getTransactionReceipt(transactionHash)
  if (!(transaction && transaction.blockNumber)) {
    throw new Error(
      `failed to fetch block number for transaction ${transactionHash}`
    )
  }

  return transaction.blockNumber
}
