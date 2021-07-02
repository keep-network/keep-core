import BigNumber from "bignumber.js"
/** @typedef { import("../../web3").BaseContract} BaseContract */
/** @typedef { import("../../web3").Web3LibWrapper} Web3LibWrapper */
/** @typedef { import("../exchange-api").BaseExchange} BaseExchange */

class CoveragePoolV1 {
  /**
   * @param {BaseContract} _assetPoolContract
   * @param {BaseContract} _covTokenContract
   * @param {BaseContract} _collateralToken
   * @param {BaseExchange} _exchangeService
   * @param {Web3LibWrapper} _web3
   */
  constructor(
    _assetPoolContract,
    _covTokenContract,
    _collateralToken,
    _exchangeService,
    _web3
  ) {
    this.assetPoolContract = _assetPoolContract
    this.covTokenContract = _covTokenContract
    this.collateralToken = _collateralToken
    this.exchangeService = _exchangeService
    this.web3 = _web3
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

  estimatedRewards = async (shareOfPool, estimatedCollateralTokenBalance) => {
    const tvl = await this.totalValueLocked()

    return new BigNumber(tvl)
      .multipliedBy(shareOfPool)
      .minus(estimatedCollateralTokenBalance)
      .toFixed(0)
      .toString()
  }

  totalValueLocked = async () => {
    return await this.assetPoolContract.makeCall("totalValue")
  }

  estimatedCollateralTokenBalance = async (shareOfPool) => {
    const balanceOfAssetPool = await this.collateralToken.makeCall(
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
