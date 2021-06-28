import BigNumber from "bignumber.js"
/** @typedef { import("../../web3").BaseContract} BaseContract */

class CoveragePoolV1 {
  /**
   *
   * @param {BaseContract} _assetPoolContract
   * @param {BaseContract} _covTokenContract
   * @param {BaseContract} _corateralTokenContract
   */
  constructor(_assetPoolContract, _covTokenContract, _corateralTokenContract) {
    this.assetPoolContract = _assetPoolContract
    this.covTokenContract = _covTokenContract
    this.corateralTokenContract = _corateralTokenContract
    this._rewardPoolContractAddress = null
  }

  shareOfPool = (covTotalSupply, covBalanceOf) => {
    if (new BigNumber(covTotalSupply).isZero()) {
      return 0
    }
    return new BigNumber(covBalanceOf).div(covTotalSupply).toString()
  }

  covTotalSupply = async () => {
    return await this.covTokenContract.makeCall("totalSupply")
  }

  covBalanceOf = async (address) => {
    return await this.covTokenContract.makeCall("balanceOf", address)
  }

  estimatedRewards = async (shareOfPool, estimatedKeepBalance) => {
    const tvl = await this.totalValueLocked()

    return new BigNumber(tvl)
      .multipliedBy(shareOfPool)
      .minus(estimatedKeepBalance)
      .toFixed(0)
      .toString()
  }

  totalValueLocked = async () => {
    return await this.assetPoolContract.makeCall("totalValue")
  }

  corateralTokenAllowance = async (owner, spender) => {
    return await this.corateralTokenContract.makeCall(
      "allowance",
      owner,
      spender
    )
  }

  estimatedCorateralTokenBalance = async (shareOfPool) => {
    const balanceOfAssetPool = await this.corateralTokenContract.makeCall(
      "balanceOf",
      this.assetPoolContract.address
    )

    return new BigNumber(balanceOfAssetPool)
      .multipliedBy(shareOfPool)
      .toFixed(0)
      .toString()
  }

  apy = async () => {}
}

export default CoveragePoolV1
