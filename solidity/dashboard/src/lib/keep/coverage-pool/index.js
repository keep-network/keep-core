import BigNumber from "bignumber.js"
/** @typedef { import("../../web3").BaseContract} BaseContract */

class CoveragePoolV1 {
  /**
   *
   * @param {BaseContract} _assetPoolContract
   * @param {BaseContract} _rewardPoolContract
   * @param {BaseContract} _covTokenContract
   * @param {BaseContract} _corateralTokenContract
   */
  constructor(
    _assetPoolContract,
    _rewardPoolContract,
    _covTokenContract,
    _corateralTokenContract
  ) {
    this.assetPoolContract = _assetPoolContract
    this.rewardPoolContract = _rewardPoolContract
    this.covTokenContract = _covTokenContract
    this.corateralTokenContract = _corateralTokenContract
  }

  shareOfPool = (covTotalSupply, covBalanceOf) => {
    return new BigNumber(covBalanceOf).div(covTotalSupply).toString()
  }

  covTotalSupply = async () => {
    return await this.covTokenContract.makeCall("totalSupply")
  }

  covBalanceOf = async (address) => {
    return await this.covTokenContract.makeCall("balanceOf", address)
  }

  estimatedRewards = async (shareOfPool) => {
    const tokensInPool = await this.corateralTokenContract.makeCall(
      "balanceOf",
      this.assetPoolContract.address
    )

    const earned = await this.rewardPoolContract.makeCall("earned")

    return new BigNumber(tokensInPool)
      .plus(new BigNumber(earned))
      .multipliedBy(shareOfPool)
      .toString()
  }

  apy = async () => {}
}

export default CoveragePoolV1
