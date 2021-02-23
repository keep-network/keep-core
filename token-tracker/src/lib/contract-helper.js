/** @typedef { import("web3-eth-contract").Contract } Web3Contract */
/** @typedef { import("web3").default } Web3 */

import { EthereumHelpers } from "@keep-network/tbtc.js"

/**
 * @typedef {Object} ContractArtifact
 * @property {JSON} artifact
 * @property {Web3} web3
 */

export class Contract {
  /**
   * @param {JSON} artifact
   * @param {Web3} web3
   */
  constructor(artifact, web3) {
    this.artifact = artifact
    this.web3 = web3
  }

  /**
   * @return {Web3Contract}
   */
  async deployed() {
    const { getDeployedContract } = EthereumHelpers

    // const networkId = await web3.eth.net.getId()
    const networkId = 1 // FIXME: Workaround

    return getDeployedContract(this.artifact, this.web3, networkId)
  }

  /**
   * @param {string} address
   * @return {Web3Contract}
   */
  async at(address) {
    const { buildContract } = EthereumHelpers

    return buildContract(this.web3, this.artifact.abi, address)
  }
}

export async function getDeploymentBlockNumber(artifact, web3) {
  // const networkId = await web3.eth.net.getId()
  const networkId = 1 // FIXME: Workaround

  const transactionHash = artifact.networks[networkId].transactionHash

  console.log("transactionHash", transactionHash)

  console.log("web3", web3.version)

  const transaction = await web3.eth
    .getTransaction(transactionHash)
    .catch((err) => {
      console.error("DUPA", err)
    })

  return transaction.blockNumber
}
