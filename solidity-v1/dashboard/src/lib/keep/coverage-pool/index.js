import BigNumber from "bignumber.js"
import { KEEP } from "../../../utils/token.utils"
import { APYCalculator } from "../helper"
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
   * @param {BaseExchange} _rewardPoolContract
   * @param {BaseExchange} _riskManagerV1Contract
   * @param {BaseExchange} _exchangeService
   * @param {Web3LibWrapper} _web3
   */
  constructor(
    _assetPoolContract,
    _covTokenContract,
    _collateralToken,
    _rewardPoolContract,
    _riskManagerV1Contract,
    _exchangeService,
    _web3
  ) {
    this.assetPoolContract = _assetPoolContract
    this.covTokenContract = _covTokenContract
    this.collateralToken = _collateralToken
    this.rewardPoolContract = _rewardPoolContract
    this.riskManagerV1Contract = _riskManagerV1Contract
    this.exchangeService = _exchangeService
    this.web3 = _web3
  }

  /**
   * Calculates the share of the coverage pool.
   * @param {string} covTotalSupply Total supply of the cov token in the
   * smallest unit. It can be fetched via {@link CoveragePoolV1.covTotalSupply}.
   * @param {string} covBalanceOf Amount of tokens owned by user in the samllest
   * unit. It can be fetched via {@link CoveragePoolV1.covBalanceOf}.
   * @return {string} The share of pool. The value is between [0, 1].
   */
  shareOfPool = (covTotalSupply, covBalanceOf) => {
    if (new BigNumber(covTotalSupply).isZero()) {
      return "0"
    }
    return new BigNumber(covBalanceOf).div(covTotalSupply).toString()
  }

  /**
   * Returns the total supply of the cov token in the smallest unit- 18
   * ddecimals precision.
   * @return {Promise<string>} Total supply of the cov token.
   */
  covTotalSupply = async () => {
    return await this.covTokenContract.makeCall("totalSupply")
  }

  /**
   * Returns the amount of cov tokens owned by `address`.
   * @param {string} address User address.
   * @return {Promise<string>} Aamount of cov tokens owned by `address`.
   */
  covBalanceOf = async (address) => {
    return await this.covTokenContract.makeCall("balanceOf", address)
  }

  /**
   * Estimates the current reward balance earned for participating in the pool.
   * The collateral token is a reward token.
   * @param {string} address User address.
   * @param {string} shareOfPool The user's current share of the pool.
   * @return {Promise<string>} The current reward balance (in collateral token)
   * in the smallest unit- 18 decimals precision.
   */
  estimatedRewards = async (address, shareOfPool) => {
    if (shareOfPool <= 0) {
      return "0"
    }

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

    let deposited = "0"
    if (gt(curretlyDeposited, "0")) {
      deposited = curretlyDeposited
    }

    return new BigNumber(tvl)
      .multipliedBy(shareOfPool)
      .minus(deposited)
      .toFixed(0)
      .toString()
  }

  /**
   * Returns the current collateral token balance of the asset pool plus the
   * reward amount (in collateral token) earned by the asset pool and not yet
   * withdrawn to the asset pool.
   *
   * @return {Promise<string>}  The total value of asset pool in collateral
   * token in the smallest unit.
   */
  totalValueLocked = async () => {
    return await this.assetPoolContract.makeCall("totalValue")
  }

  /**
   * Estimates the collateral token balance based on the share of the pool.
   * @param {string | number} shareOfPool The user's current share of the
   * coverage pool. It can be calculated via {@link CoveragePoolV1#shareOfPool}
   * @return {Promise<string>} Estimated collateral token balance.
   */
  estimatedCollateralTokenBalance = async (shareOfPool) => {
    const balanceOfAssetPool = await this.totalValueLocked()

    return new BigNumber(balanceOfAssetPool)
      .multipliedBy(shareOfPool)
      .toFixed(0)
      .toString()
  }

  /**
   * Estimates collateral token balance from given cov token amount
   * @param {string | number} covToken cov token amount that we want to estimate
   * collateral token value from
   * @param {string | number} totalSupplyOfCovToken total supply of cov tokens
   * @param {string | number} tvl - total supply of collateral tokens places in
   * coverage pools
   * @return {string} - estimated collateral token amount of given cov tokens
   */
  estimatedBalanceFor = (covToken, totalSupplyOfCovToken, tvl) => {
    if (
      new BigNumber(covToken).isZero() ||
      new BigNumber(totalSupplyOfCovToken).isZero() ||
      new BigNumber(tvl).isZero()
    ) {
      return "0"
    }
    return (
      new BigNumber(covToken)
        .div(totalSupplyOfCovToken)
        .multipliedBy(tvl)
        // TODO: Make a coverage pool lib more abstract to be independent of the
        // KEEP token and pass collateral token and cov token object to
        // constructor.
        .decimalPlaces(KEEP.decimals)
        .toString()
    )
  }

  /**
   * Returns the `AssetPool` contract's balance of the collateral token.
   * @return {Promise<string>} The `AssetPool` contract's balance of the
   * collateral token.
   */
  assetPoolCollateralTokenBalance = async () => {
    return await this.collateralToken.makeCall(
      "balanceOf",
      this.assetPoolContract.address
    )
  }

  /**
   * Calculates the reward pool per week. The `RewardsPool` will earn a given
   * amount of collateral token per week. Reward tokens from the previous
   * interval that has not been yet unlocked, are added to the new interval
   * being created.
   * @return  {Promise<string>} The reward pool per week (in collateral token)
   * in the smallest unit.
   */
  rewardPoolPerWeek = async () => {
    const rewardRate = await this.rewardPoolRewardRate()

    return KEEP.toTokenUnit(rewardRate).multipliedBy(REWARD_DURATION)
  }

  /**
   * Returns the `RewardsPool` contract's rate per second with which reward
   * tokens are unlocked.
   * @return {Promise<string>} The rate per second with which reward tokens are
   * unlocked.
   */
  rewardPoolRewardRate = async () => {
    return await this.rewardPoolContract.makeCall("rewardRate")
  }

  /**
   * Calculates the APY of the coverage pool assuming that there are no calims.
   * @return {Promise<string>} APY.
   */
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

  /**
   * Calculates the total amount of the tokens that were allocated as a reward
   * to the `RewardsPool` contract.
   * @return {Promise<string>} The total amount of tokens allocated to
   * the `RewardsPool` contract.
   */
  totalAllocatedRewards = async () => {
    return (
      await this.rewardPoolContract.getPastEvents("RewardToppedUp")
    ).reduce((reducer, _) => add(reducer, _.returnValues.amount), "0")
  }

  /**
   * Calculates the amount of collateral token used from the coverage pool to
   * cover a loss since the start of the pool.
   * @return {Promise<string>} The total amount of collateral token used from the
   * coverage pool to cover a loss.
   */
  totalCoverageClaimed = async () => {
    return (
      await this.assetPoolContract.getPastEvents("CoverageClaimed")
    ).reduce((reducer, _) => add(reducer, _.returnValues.amount), "0")
  }

  covTokensAllowed = async (owner, spender) => {
    return await this.covTokenContract.makeCall("allowance", owner, spender)
  }

  withdrawalInitiatedTimestamp = async (address) => {
    return await this.assetPoolContract.makeCall(
      "withdrawalInitiatedTimestamp",
      address
    )
  }

  pendingWithdrawal = async (address) => {
    return await this.assetPoolContract.makeCall("pendingWithdrawal", address)
  }

  withdrawalDelays = async () => {
    const withdrawalDelay = await this.assetPoolContract.makeCall(
      "withdrawalDelay"
    )
    const withdrawalTimeout = await this.assetPoolContract.makeCall(
      "withdrawalTimeout"
    )
    return {
      withdrawalDelay,
      withdrawalTimeout,
    }
  }

  hasRiskManagerOpenAuctions = async () => {
    return await this.riskManagerV1Contract.makeCall("hasOpenAuctions")
  }
}

export default CoveragePoolV1
