/** @typedef { import("../lib/context").Context } Context */
/** @typedef { import("../lib/ethereum-helper").Address } Address */

import BN from "bn.js"

import { ITruthSource } from "./truth-source.js"
import { Contract } from "../lib/contract-helper.js"
import { getPastEvents } from "../lib/ethereum-helper.js"
import { dumpDataToFile } from "../lib/file-helper.js"
import { logger } from "../lib/winston.js"
import { getPairData } from "../lib/uniswap.js"

// https://etherscan.io/address/0xe6f19dab7d43317344282f803f8e8d240708174a#code
import KEEPETHTokenJson from "../../artifacts/KEEP-ETH-UNI-V2-Token.json"
// https://etherscan.io/address/0x47a5f2ffdf66d13ed7e317581f458d09b49d6f44
import LPRewardsKEEPETHJson from "@keep-network/keep-ecdsa/artifacts/LPRewardsKEEPETH.json"

const KEEPETH_PAIR = "0xe6f19dab7d43317344282f803f8e8d240708174a" // https://info.uniswap.org/pair/0xe6f19dab7d43317344282f803f8e8d240708174a
const KEEPETH_CREATION_BLOCK = "10100034" // https://etherscan.io/tx/0xc64ac175846e719bb4f7f9b17a0b04bc365db3dda9d97ef70d7ede8f9c1a265b

const LP_KEEPETH_TOKEN_HISTORIC_STAKERS_DUMP_PATH =
  "./tmp/lp-keepeth-token-stakers.json"
const KEEP_IN_LP_KEEPETH_BALANCES_DUMP_PATH =
  "./tmp/keep-in-lp-keepeth-token-balances.json"

export class LPKeepEthTokenTruthSource extends ITruthSource {
  /**
   * @param {Context} context
   * @param {Number} finalBlock
   */
  constructor(context, finalBlock) {
    super(context, finalBlock)
  }

  async initialize() {
    const KEEPETHPairTokenAbi = JSON.parse(KEEPETHTokenJson.result)
    this.keepEthTokenContract = new this.context.web3.eth.Contract(
      KEEPETHPairTokenAbi,
      KEEPETH_PAIR
    )

    const lpRewardKeepEth = new Contract(
      LPRewardsKEEPETHJson,
      this.context.web3
    )
    this.lpRewardKeepEthContract = await lpRewardKeepEth.deployed()
  }

  /**
   * Finds all historic stakers of LP KEEP-ETH pair token based on Transfer events
   * 
   * @return {Set<Address>} All historic LP KEEP-ETH token stakers.
   * */
  async findStakers() {
    const lpRewardKeepEthContractAddress = this.lpRewardKeepEthContract.options
      .address

    logger.info(
      `looking for Transfer events emitted from ${this.keepEthTokenContract.options.address} ` +
        `to LP KEEP-ETH pair ${lpRewardKeepEthContractAddress} ` +
        `between blocks ${KEEPETH_CREATION_BLOCK} and ${this.finalBlock}`
    )

    const events = await getPastEvents(
      this.context.web3,
      this.keepEthTokenContract,
      "Transfer",
      KEEPETH_CREATION_BLOCK,
      this.finalBlock
    )
    logger.info(`found ${events.length} lp keep-eth token transfer events`)

    const lpTokenStakersSet = new Set()
    events.forEach(function (event) {
      // include accounts that staked in LPRewards contracts only
      if (event.returnValues.to == lpRewardKeepEthContractAddress) {
        lpTokenStakersSet.add(event.returnValues.from)
      }
    })

    logger.info(`found ${lpTokenStakersSet.size} unique historic stakers`)

    dumpDataToFile(
      lpTokenStakersSet,
      LP_KEEPETH_TOKEN_HISTORIC_STAKERS_DUMP_PATH
    )

    return Array.from(lpTokenStakersSet)
  }

  /**
   * Retrieves balances of LP KEEP-ETH pair for stakers in LPRewardsKEEPETH contract
   *
   * @param {Array<Address>} lpStakers LP KEEP-ETH stakers
   *
   * @return {Map<Address,BN>} LP Balances by lp stakers
   */
  async lpKeepEthStakersBalances(lpStakers) {
    let expectedTotalSupply = new BN(0)
    const lpBalanceByStaker = new Map()

    for (let i = 0; i < lpStakers.length; i++) {
      const lpBalance = new BN(
        await this.lpRewardKeepEthContract.methods
          .balanceOf(lpStakers[i])
          .call({}, this.finalBlock)
      )
      if (!lpBalance.isZero()) {
        lpBalanceByStaker.set(lpStakers[i], lpBalance)
        expectedTotalSupply = expectedTotalSupply.add(lpBalance)
      }
    }
    const actualTotalSupply = new BN(await this.lpRewardKeepEthContract.methods
      .totalSupply()
      .call({}, this.finalBlock))

    if (!expectedTotalSupply.eq(actualTotalSupply)) {
      logger.error(
        `Sum of LP staker balances ${expectedTotalSupply} does not match the total supply ${actualTotalSupply}`
      )
    }

    logger.info(`Total supply of LP Token: ${expectedTotalSupply.toString()}`)

    return lpBalanceByStaker
  }

  /**
   * Calculates KEEP for all LP KEEP-ETH stakers.
   *
   * @param {Map<Address, BN>} stakersBalances LP KEEP-ETH Token amounts by stakers
   *
   * @return {Map<Address,BN>} KEEP Tokens in LP KEEP-ETH at the final block
   */
  async calcKeepInStakersBalances(stakersBalances) {
    logger.info(`check token stakers at block ${this.finalBlock}`)

    const keepInLpByStakers = new Map()

    // Retrieve current pair data
    const pairData = await getPairData(KEEPETH_PAIR)
    for (const [stakerAddress, lpBalance] of stakersBalances.entries()) {
      const keepInLPToken = await this.calcKeepTokenfromLPToken(
        lpBalance,
        pairData
      )
      keepInLpByStakers.set(stakerAddress, keepInLPToken)

      logger.info(
        `Staker: ${stakerAddress} - LP Balance: ${lpBalance} - KEEP in LP: ${keepInLPToken}`
      )
    }

    logger.info(
      `found ${keepInLpByStakers.size} stakers at block ${this.finalBlock}`
    )

    dumpDataToFile(keepInLpByStakers, KEEP_IN_LP_KEEPETH_BALANCES_DUMP_PATH)

    return keepInLpByStakers
  }

  /**
   * Calculating amount of KEEP token which makes a KEEP-ETH Uniswap pair.
   * Math is based on https://uniswap.org/docs/v2/advanced-topics/understanding-returns/
   *
   * @param {BN} lpBalance LP amount staked by a staker
   * @param {PairData} pairData KEEP-ETH pair data fetched from Uniswap
   *
   * @return {BN} KEEP token amounts in a LP token balance
   */
  async calcKeepTokenfromLPToken(lpBalance, pairData) {
    const uniswapTotalSupply = new BN(
      this.context.web3.utils.toWei(pairData.totalSupply.toString())
    )
    const keepLiquidityPool = new BN(
      this.context.web3.utils.toWei(pairData.reserve0.toString())
    )

    return lpBalance.mul(keepLiquidityPool).div(uniswapTotalSupply)
  }

  /**
   * @return {Map<Address,BN>} KEEP token amounts staked by stakers at the final block
   */
  async getTokenHoldingsAtFinalBlock() {
    await this.initialize()

    const lpStakers = await this.findStakers()
    const stakersBalances = await this.lpKeepEthStakersBalances(lpStakers)

    return await this.calcKeepInStakersBalances(stakersBalances)
  }
}
