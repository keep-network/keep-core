class KeepToTStaking {
  /**
   * @param {BaseContract} _thresholdStakingContract
   * @param {Web3LibWrapper} _web3
   */
  constructor(_thresholdStakingContract, _web3) {
    this.thresholdStakingContract = _thresholdStakingContract
    this.web3 = _web3
  }
}

export default KeepToTStaking
