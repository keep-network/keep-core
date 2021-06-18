import contracts from "./contracts"
import CoveragePoolV1 from "./coverage-pool"

/** @typedef { import("../web3").Web3LibWrapper} Web3LibWrapper */

class Keep {
  static initialize(web3) {
    const keep = new Keep(web3)
    keep.initializeContracts()
    keep.initializeServices()

    return keep
  }

  /**
   * @param {Web3LibWrapper} _web3 The web3 lib wrapper.
   */
  constructor(_web3) {
    this.web3 = _web3
  }

  initializeContracts = () => {
    for (const [contractName, { artifact }] of Object.entries(contracts)) {
      this[contractName] = this.web3.createContractInstance(artifact)
    }
  }

  initializeServices = () => {
    this.coveragePoolV1 = new CoveragePoolV1(
      this.assetPoolContract,
      this.rewardPoolContract,
      this.covTokenContract,
      this.keepTokenContract
    )
  }
}

export default Keep
