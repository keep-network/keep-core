import web3Utils from "web3-utils"
import { createERC20Contract } from "../contracts"
import BigNumber from "bignumber.js"
import { toTokenUnit } from "../utils/token.utils"
import {
  getPairData,
  getKeepTokenPriceInUSD,
  getBTCPriceInUSD,
} from "./uniswap-api"
import moment from "moment"
/** @typedef {import("web3").default} Web3 */
/** @typedef {LiquidityRewards} LiquidityRewards */

// lp contract address -> wrapped ERC20 token as web3 contract instance
const LPRewardsToWrappedTokenCache = {}
const WEEKS_IN_YEAR = 52

class LiquidityRewards {
  constructor(_wrappedTokenContract, _LPRewardsContract) {
    this.wrappedToken = _wrappedTokenContract
    this.LPRewardsContract = _LPRewardsContract
  }

  get wrappedTokenAddress() {
    return this.wrappedToken.options.address
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
}

class SaddleLPRewards extends LiquidityRewards {
  calculateAPY = async (totalSupplyOfLPRewards) => {
    totalSupplyOfLPRewards = toTokenUnit(totalSupplyOfLPRewards)

    const wrappedTokenTotalSupply = await this.wrappedTokenTotalSupply()
    const BTCPriceInUSD = await getBTCPriceInUSD()

    // TODO fetch total Bitcoins deposited in the wrapped token pool
    const totalBitcoinDepositedInWrappedTokenPool = 0
    const wrappedTokenPoolInUSD = BTCPriceInUSD.multipliedBy(
      totalBitcoinDepositedInWrappedTokenPool
    )

    const rewardPoolPerWeek = await this.rewardPoolPerWeek()

    const lpRewardsPoolInUSD = totalSupplyOfLPRewards
      .multipliedBy(wrappedTokenPoolInUSD)
      .div(wrappedTokenTotalSupply)

    const keepTokenInUSD = await getKeepTokenPriceInUSD()

    const r = this._calculateR(
      keepTokenInUSD,
      rewardPoolPerWeek,
      lpRewardsPoolInUSD
    )

    return this._calculateAPY(r, WEEKS_IN_YEAR)
  }
}

const LiquidityRewardsPoolStrategy = {
  UNISWAP: UniswapLPRewards,
  SADDLE: SaddleLPRewards,
}

export class LiquidityRewardsFactory {
  /**
   *
   * @param {('UNISWAP' | 'SADDLE')} pool - The supported type of pools.
   * @param {Object} LPRewardsContract - The LPRewardsContract as web3 contract instance.
   * @param {Web3} web3 - web3
   * @return {LiquidityRewards} - The Liquidity Rewards Wrapper
   */
  static async initialize(pool, LPRewardsContract, web3) {
    const lpRewardsContractAddress = web3Utils.toChecksumAddress(
      LPRewardsContract.options.address
    )

    if (
      !LPRewardsToWrappedTokenCache.hasOwnProperty(lpRewardsContractAddress)
    ) {
      const wrappedTokenAddress = await LPRewardsContract.methods
        .wrappedToken()
        .call()
      LPRewardsToWrappedTokenCache[
        lpRewardsContractAddress
      ] = createERC20Contract(web3, wrappedTokenAddress)
    }

    const wrappedTokenContract =
      LPRewardsToWrappedTokenCache[lpRewardsContractAddress]
    const PoolStrategy = LiquidityRewardsPoolStrategy[pool]

    return new PoolStrategy(wrappedTokenContract, LPRewardsContract)
  }
}
