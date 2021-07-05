import BigNumber from "bignumber.js"
import { KEEP } from "../../../utils/token.utils"
import { APYCalculator } from "../helper"
import { RewardsPoolArtifact } from "../contracts"
import { add, sub, gt } from "../../../utils/arithmetics.utils"

/** @typedef { import("../../web3").BaseContract} BaseContract */
/** @typedef { import("../../web3").Web3LibWrapper} Web3LibWrapper */
/** @typedef { import("../exchange-api").BaseExchange} BaseExchange */

const REWARD_DURATION = 604800 // 7 days in seconds
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
    this._rewardPoolContract = undefined
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

  estimatedRewards = async (address, shareOfPool) => {
    const tvl = await this.totalValueLocked()
    const toAssetPool = (
      await this.collateralToken.getPastEvents("Transfer", {
        from: address,
        to: this.assetPoolContract.address,
      })
    ).reduce((reducer, _) => add(reducer, _.returnValues.value), "0")
    const fromAssetPool = (
      await this.collateralToken.getPastEvents("Transfer", {
        from: this.assetPoolContract.address,
        to: address,
      })
    ).reduce((reducer, _) => add(reducer, _.returnValues.value), "0")

    const curretlyDeposited = sub(toAssetPool, fromAssetPool)

    let deposited = 0
    if (gt(curretlyDeposited, "0")) {
      deposited = curretlyDeposited
    }

    return new BigNumber(tvl)
      .multipliedBy(shareOfPool)
      .minus(deposited)
      .toFixed(0)
      .toString()
  }

  totalValueLocked = async () => {
    return await this.assetPoolContract.makeCall("totalValue")
  }

  estimatedCollateralTokenBalance = async (shareOfPool) => {
    const balanceOfAssetPool = await this.assetPoolCollateralTokenBalance()

    return new BigNumber(balanceOfAssetPool)
      .multipliedBy(shareOfPool)
      .toFixed(0)
      .toString()
  }

  assetPoolCollateralTokenBalance = async () => {
    return await this.collateralToken.makeCall(
      "balanceOf",
      this.assetPoolContract.address
    )
  }

  rewardPoolPerWeek = async () => {
    const rewardRate = await this.rewardPoolRewardRate()

    return KEEP.toTokenUnit(rewardRate).multipliedBy(REWARD_DURATION)
  }

  rewardPoolRewardRate = async () => {
    const rewardPoolContract = await this.getRewardPoolContract()
    return await rewardPoolContract.makeCall("rewardRate")
  }

  /**
   * @return {Promise<BaseContract>} The reward pool contract.
   */
  getRewardPoolContract = async () => {
    if (!this._rewardPoolContract) {
      const rewardPoolAddress = await this.assetPoolContract.makeCall(
        "rewardsPool"
      )
      this._rewardPoolContract = this.web3.createContractInstance(
        RewardsPoolArtifact.abi,
        rewardPoolAddress,
        // The `RewardsPool` contract is created in the same transaction as the
        // `AssetPool` contract (in the `AssetPool` constructor). In thah case
        // we can pass `deploymentTxnHash` and `deployedAtBlock` from the
        // `AssetPool` contract.
        this.assetPoolContract.deploymentTxnHash,
        this.assetPoolContract.deployedAtBlock
      )
    }
    return this._rewardPoolContract
  }

  apy = async () => {
    const totalSupply = await this.assetPoolCollateralTokenBalance()
    const rewardPoolPerWeek = await this.rewardPoolPerWeek()

    // We know that the collateral token is KEEP. TODO: consider a more abstract
    // solution to fetch the collateral token price in USD.
    const collateralTokenPriceInUSD =
      await this.exchangeService.getKeepTokenPriceInUSD()

    const totalSupplyInUSD = KEEP.toTokenUnit(totalSupply).multipliedBy(
      collateralTokenPriceInUSD
    )

    const rewardRate = APYCalculator.calculatePoolRewardRate(
      collateralTokenPriceInUSD,
      rewardPoolPerWeek,
      totalSupplyInUSD
    )

    return APYCalculator.calculateAPY(rewardRate).toString()
  }

  totalAllocatedRewards = async () => {
    const rewardPoolContract = await this.getRewardPoolContract()

    return (await rewardPoolContract.getPastEvents("RewardToppedUp")).reduce(
      (reducer, _) => add(reducer, _.returnValues.amount),
      "0"
    )
  }
}

export default CoveragePoolV1
