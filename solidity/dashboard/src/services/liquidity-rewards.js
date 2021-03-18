import web3Utils from "web3-utils"
import {
  createERC20Contract,
  createSaddleSwapContract,
  CONTRACT_DEPLOY_BLOCK_NUMBER,
} from "../contracts"
import BigNumber from "bignumber.js"
import { toTokenUnit, fromTokenUnit } from "../utils/token.utils"
import {
  getPairData,
  getKeepTokenPriceInUSD,
  getBTCPriceInUSD,
} from "./uniswap-api"
import moment from "moment"
import { add } from "../utils/arithmetics.utils"
import { isEmptyArray } from "../utils/array.utils"
/** @typedef {import("web3").default} Web3 */
/** @typedef {LiquidityRewards} LiquidityRewards */

// lp contract address -> wrapped ERC20 token as address
const LPRewardsToWrappedTokenCache = {}
const WEEKS_IN_YEAR = 52

class LiquidityRewards {
  static async _getWrappedTokenAddress(LPRewardsContract) {
    return await LPRewardsContract.methods.wrappedToken().call()
  }

  constructor(_wrappedTokenContract, _LPRewardsContract, _web3) {
    this.wrappedToken = _wrappedTokenContract
    this.LPRewardsContract = _LPRewardsContract
    this.web3 = _web3
  }

  get wrappedTokenAddress() {
    return this.wrappedToken.options.address
  }

  get LPRewardsContractAddress() {
    return this.LPRewardsContract.options.address
  }

  get rewardClaimedEventName() {
    return "RewardPaid"
  }

  get depositWithdrawnEventName() {
    return "Withdrawn"
  }

  get withdrawTokensFnName() {
    return "exit"
  }

  withdrawTokensArgs() {
    return []
  }

  get stakedEventName() {
    return "Staked"
  }

  get stakeFnName() {
    return "stake"
  }

  stakeArgs(amount) {
    return [amount]
  }

  wrappedTokenBalance = async (address) => {
    return await this.wrappedToken.methods.balanceOf(address).call()
  }

  wrappedTokenTotalSupply = async () => {
    return await this.wrappedToken.methods.totalSupply().call()
  }

  wrappedTokenAllowance = async (owner, spender) => {
    return await this.wrappedToken.methods.allowance(owner, spender).call()
  }

  stakedBalance = async (address) => {
    return await this.LPRewardsContract.methods.balanceOf(address).call()
  }

  totalSupply = async () => {
    return await this.LPRewardsContract.methods.totalSupply().call()
  }

  rewardBalance = async (address) => {
    return await this.LPRewardsContract.methods.earned(address).call()
  }

  rewardRate = async () => {
    return await this.LPRewardsContract.methods.rewardRate().call()
  }

  rewardPoolPerWeek = async () => {
    const rewardRate = await this.rewardRate()
    return toTokenUnit(rewardRate).multipliedBy(
      moment.duration(7, "days").asSeconds()
    )
  }

  _calculateR = (
    keepTokenInUSD,
    rewardPoolPerInterval,
    totalLPTokensInLPRewardsInUSD
  ) => {
    return keepTokenInUSD
      .multipliedBy(rewardPoolPerInterval)
      .div(totalLPTokensInLPRewardsInUSD)
  }

  /**
   * Calculates the APY.
   *
   * @param {BigNumber} r Period rate.
   * @param {number | string | BigNumber} n Number of compounding periods.
   * @return {BigNumber} APY value.
   */
  _calculateAPY = (r, n = WEEKS_IN_YEAR) => {
    return r.plus(1).pow(n).minus(1)
  }

  calculateAPY = async (totalSupplyOfLPRewards) => {
    throw new Error("First, implement the `calculateAPY` function")
  }

  calculateLPTokenBalance = async (shareOfPoolInPercent) => {
    throw new Error("First, implement the `calculateLPTokenBalance` function")
  }

  calculateRewardMultiplier = async (address) => {
    throw new Error("First, implement the `calculateRewardMultiplier` function")
  }
}

class UniswapLPRewards extends LiquidityRewards {
  calculateAPY = async (totalSupplyOfLPRewards) => {
    totalSupplyOfLPRewards = toTokenUnit(totalSupplyOfLPRewards)

    const pairData = await getPairData(this.wrappedTokenAddress.toLowerCase())
    const rewardPoolPerWeek = await this.rewardPoolPerWeek()

    const lpRewardsPoolInUSD = totalSupplyOfLPRewards
      .multipliedBy(pairData.reserveUSD)
      .div(pairData.totalSupply)

    const ethPrice = new BigNumber(pairData.reserveUSD).div(pairData.reserveETH)

    let keepTokenInUSD = 0
    if (pairData.token0.symbol === "KEEP") {
      keepTokenInUSD = ethPrice.multipliedBy(pairData.token0.derivedETH)
    } else if (pairData.token1.symbol === "KEEP") {
      keepTokenInUSD = ethPrice.multipliedBy(pairData.token1.derivedETH)
    } else {
      keepTokenInUSD = await getKeepTokenPriceInUSD()
    }

    const r = this._calculateR(
      keepTokenInUSD,
      rewardPoolPerWeek,
      lpRewardsPoolInUSD
    )

    return this._calculateAPY(r, WEEKS_IN_YEAR)
  }

  /**
   * Calculates lp token balance for the given pair
   * The calculations were done based on
   * https://uniswap.org/docs/v2/advanced-topics/understanding-returns/#why-is-my-liquidity-worth-less-than-i-put-in
   *
   * @param {string} lpBalance Balance of liquidity token for a given uniswap pair deposited in
   * the LPRewards` contract.
   * @return {Promise<{token0: string, token1: string}>}
   */
  calculateLPTokenBalance = async (lpBalance) => {
    const pairData = await getPairData(this.wrappedTokenAddress.toLowerCase())

    return {
      token0: new BigNumber(lpBalance)
        .multipliedBy(pairData.reserve0)
        .dividedBy(pairData.totalSupply)
        .toString(),
      token1: new BigNumber(lpBalance)
        .multipliedBy(pairData.reserve1)
        .dividedBy(pairData.totalSupply)
        .toString(),
    }
  }

  calculateRewardMultiplier = async (address) => {
    return 0
  }
}

class SaddleLPRewards extends LiquidityRewards {
  BTC_POOL_TOKENS = [
    { name: "TBTC", decimals: 18 },
    { name: "WBTC", decimals: 8 },
    { name: "RENBTC", decimals: 8 },
    { name: "SBTC", decimals: 18 },
  ]

  constructor(_wrappedTokenContract, _LPRewardsContract, _web3) {
    super(_wrappedTokenContract, _LPRewardsContract, _web3)
    this.swapContract = createSaddleSwapContract(this.web3)
  }

  swapContract = null

  calculateAPY = async (totalSupplyOfLPRewards) => {
    totalSupplyOfLPRewards = toTokenUnit(totalSupplyOfLPRewards)

    const wrappedTokenTotalSupply = toTokenUnit(
      await this.wrappedTokenTotalSupply()
    )

    const BTCInPool = await this._getBTCInPool()
    const BTCPriceInUSD = await getBTCPriceInUSD()

    const wrappedTokenPoolInUSD = BTCPriceInUSD.multipliedBy(
      toTokenUnit(BTCInPool)
    )

    const keepTokenInUSD = await getKeepTokenPriceInUSD()

    const rewardPoolPerWeek = await this.rewardPoolPerWeek()

    const lpRewardsPoolInUSD = totalSupplyOfLPRewards
      .multipliedBy(wrappedTokenPoolInUSD)
      .div(wrappedTokenTotalSupply)

    const r = this._calculateR(
      keepTokenInUSD,
      rewardPoolPerWeek,
      lpRewardsPoolInUSD
    )

    return this._calculateAPY(r, WEEKS_IN_YEAR)
  }

  _getBTCInPool = async () => {
    return (
      await Promise.all(
        this.BTC_POOL_TOKENS.map(async (token, i) => {
          const balance = await this._getTokenBalance(i)
          return new BigNumber(10)
            .pow(18 - token.decimals) // cast all to 18 decimals
            .multipliedBy(balance)
        })
      )
    ).reduce(add, 0)
  }

  _getTokenBalance = async (index) => {
    return await this.swapContract.methods.getTokenBalance(index).call()
  }

  calculateLPTokenBalance = (shareOfTotalPoolInPercent) => {
    return {
      token0: "0",
      token1: "0",
    }
  }

  calculateRewardMultiplier = async (address) => {
    return 0
  }
}

class TokenGeyserLPRewards extends LiquidityRewards {
  static async _getWrappedTokenAddress(LPRewardsContract) {
    return await LPRewardsContract.methods.token().call()
  }

  get rewardClaimedEventName() {
    return "TokensClaimed"
  }

  get depositWithdrawnEventName() {
    return "Unstaked"
  }

  get withdrawTokensFnName() {
    return "unstake"
  }

  withdrawTokensArgs(amount) {
    return [amount, []]
  }

  stakeArgs(amount) {
    return [amount, []]
  }

  stakedBalance = async (address) => {
    return await this.LPRewardsContract.methods.totalStakedFor(address).call()
  }

  totalSupply = async () => {
    return await this.LPRewardsContract.methods.totalStaked().call()
  }

  rewardBalance = async (address, amount) => {
    try {
      // The `TokenGeyser.unstakeQuery` throws an error in case when eg. the
      // amount param is greater than the real user's stake or when
      // the user stakes KEEP in block `X` and call unstakeQuery in block `X`
      // (`SafeMath: division by zero` error is thrown.). The web3 parses the
      // error message in the wrong way when the `hanleRevert` option is enabled
      // [1]. So here we clone the rewards contract instance and disable the
      // `hanldeRevert` option.
      // References: [1]:
      // https://github.com/ChainSafe/web3.js/issues/3742
      const clonedLPRewardsContract = this.LPRewardsContract.clone()
      clonedLPRewardsContract.handleRevert = false
      return await clonedLPRewardsContract.methods.unstakeQuery(amount).call()
    } catch (error) {
      return 0
    }
  }

  calculateAPY = async (totalSupplyOfLPRewards) => {
    totalSupplyOfLPRewards = toTokenUnit(totalSupplyOfLPRewards)

    const rewardPoolPerWeek = await this.rewardPoolPerWeek()
    const keepTokenInUSD = await getKeepTokenPriceInUSD()

    const lpRewardsPoolInUSD = totalSupplyOfLPRewards.multipliedBy(
      keepTokenInUSD
    )

    const r = this._calculateR(
      keepTokenInUSD,
      rewardPoolPerWeek,
      lpRewardsPoolInUSD
    )

    return this._calculateAPY(r, WEEKS_IN_YEAR)
  }

  rewardPoolPerWeek = async () => {
    const tokensLockedEvents = await this.LPRewardsContract.getPastEvents(
      "TokensLocked",
      {
        fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.keepTokenGeyserContract,
      }
    )

    // The KEEP-only pool will earn 100k KEEP per month.
    let rewardPoolPerMonth = fromTokenUnit(10e4)
    const weeksInMonth = new BigNumber(
      moment.duration(1, "months").asSeconds()
    ).div(moment.duration(7, "days").asSeconds())

    if (!isEmptyArray(tokensLockedEvents)) {
      rewardPoolPerMonth = new BigNumber(
        tokensLockedEvents.reverse()[0].returnValues.amount
      )
    }

    return toTokenUnit(rewardPoolPerMonth.div(weeksInMonth))
  }

  calculateLPTokenBalance = (shareOfTotalPoolInPercent) => {
    return {
      token0: "0",
      token1: "0",
    }
  }

  /**
   * Calculates reward multiplier for KEEP-ONLY pool for a given user
   *
   * @param {string} address - address of the user's wallet
   * @return {Promise<string>}
   */
  calculateRewardMultiplier = async (address) => {
    const stakedBalanceOfUser = await this.stakedBalance(address)
    const rewardBalance = await this.rewardBalance(stakedBalanceOfUser)

    const stakedBalanceOfUserBN = new BigNumber(stakedBalanceOfUser)
    const rewardBalanceBN = new BigNumber(rewardBalance)

    const rewardMultiplier = stakedBalanceOfUserBN
      .plus(rewardBalanceBN)
      .dividedBy(stakedBalanceOfUserBN)
      .toString()

    return rewardMultiplier
  }
}

const LiquidityRewardsPoolStrategy = {
  UNISWAP: UniswapLPRewards,
  SADDLE: SaddleLPRewards,
  TOKEN_GEYSER: TokenGeyserLPRewards,
}

export class LiquidityRewardsFactory {
  /**
   *
   * @param {('UNISWAP' | 'SADDLE' | 'TOKEN_GEYSER')} pool - The supported type of pools.
   * @param {Object} LPRewardsContract - The LPRewardsContract as web3 contract instance.
   * @param {Web3} web3 - web3
   * @return {LiquidityRewards} - The Liquidity Rewards Wrapper
   */
  static async initialize(pool, LPRewardsContract, web3) {
    const PoolStrategy = LiquidityRewardsPoolStrategy[pool]

    const lpRewardsContractAddress = web3Utils.toChecksumAddress(
      LPRewardsContract.options.address
    )

    if (
      !LPRewardsToWrappedTokenCache.hasOwnProperty(lpRewardsContractAddress)
    ) {
      const wrappedTokenAddress = await PoolStrategy._getWrappedTokenAddress(
        LPRewardsContract
      )
      LPRewardsToWrappedTokenCache[
        lpRewardsContractAddress
      ] = wrappedTokenAddress
    }

    const wrappedTokenContract = createERC20Contract(
      web3,
      LPRewardsToWrappedTokenCache[lpRewardsContractAddress]
    )

    return new PoolStrategy(wrappedTokenContract, LPRewardsContract, web3)
  }
}
