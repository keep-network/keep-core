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
import keepEthTokenJson from "../../artifacts/KEEP-ETH-UNI-V2-Token.json"
// https://etherscan.io/address/0x38c8ffee49f286f25d25bad919ff7552e5daf081#code
import keepTbtcTokenJson from "../../artifacts/KEEP-TBTC-UNI-V2-Token.json"

import LPRewardsKEEPETHJson from "@keep-network/keep-ecdsa/artifacts/LPRewardsKEEPETH.json"
import LPRewardsKEEPTBTCJson from "@keep-network/keep-ecdsa/artifacts/LPRewardsKEEPTBTC.json"

// https://info.uniswap.org/pair/0xe6f19dab7d43317344282f803f8e8d240708174a
const KEEPETH_PAIR = "0xe6f19dab7d43317344282f803f8e8d240708174a"
// https://info.uniswap.org/pair/0x38c8ffee49f286f25d25bad919ff7552e5daf081
const KEEPTBTC_PAIR = "0x38c8ffee49f286f25d25bad919ff7552e5daf081"

// https://etherscan.io/tx/0xc64ac175846e719bb4f7f9b17a0b04bc365db3dda9d97ef70d7ede8f9c1a265b
const KEEPETH_CREATION_BLOCK = "10100034"
// https://etherscan.io/tx/0x1592f9b235c602c87a5b8cc5f896164dc43d16b92664cb9c8b420d28b64ca4a0
const KEEPTBTC_CREATION_BLOCK = "11452642"

const KEEP_IN_LP_KEEPETH_BALANCES_PATH =
  "./tmp/keep-in-lp-keepeth-token-balances.json"

const KEEP_IN_LP_KEEPTBTC_BALANCES_PATH =
  "./tmp/keep-in-lp-keeptbtc-token-balances.json"

export class LPTokenTruthSource extends ITruthSource {
  /**
   * @param {Context} context
   * @param {Number} finalBlock
   */
  constructor(context, finalBlock) {
    super(context, finalBlock)
  }

  async initialize() {
    const keepEthPairTokenAbi = JSON.parse(keepEthTokenJson.result)
    const keepEthTokenContract = new this.context.web3.eth.Contract(
      keepEthPairTokenAbi,
      KEEPETH_PAIR
    )
    const lpRewardKeepEth = new Contract(
      LPRewardsKEEPETHJson,
      this.context.web3
    )

    const KEEPTBTCPairTokenAbi = JSON.parse(keepTbtcTokenJson.result)
    const keepTbtcTokenContract = new this.context.web3.eth.Contract(
      KEEPTBTCPairTokenAbi,
      KEEPTBTC_PAIR
    )
    const lpRewardKeepTbtc = new Contract(
      LPRewardsKEEPTBTCJson,
      this.context.web3
    )

    this.liquidityStakingObjects = {
      KEEPETH: {
        lpTokenContract: keepEthTokenContract,
        lpRewardsContract: await lpRewardKeepEth.deployed(),
        lpCreationBlock: KEEPETH_CREATION_BLOCK,
        keepInLpTokenFilePath: KEEP_IN_LP_KEEPETH_BALANCES_PATH,
        lpPairAddress: KEEPETH_PAIR,
      },
      KEEPTBTC: {
        lpTokenContract: keepTbtcTokenContract,
        lpRewardsContract: await lpRewardKeepTbtc.deployed(),
        lpCreationBlock: KEEPTBTC_CREATION_BLOCK,
        keepInLpTokenFilePath: KEEP_IN_LP_KEEPTBTC_BALANCES_PATH,
        lpPairAddress: KEEPTBTC_PAIR,
      },
    }
  }

  /**
   * Finds all historic stakers of LP KEEP-ETH / KEEP-TBTC pair token based on
   * "Transfer" events
   *
   * @return {Set<Address>} All historic LP KEEP-ETH / KEEP-TBTC token stakers
   * */
  async findStakers(pairName) {
    const lpRewardsContractAddress = this.liquidityStaking.lpRewardsContract
      .options.address

    logger.info(
      `looking for Transfer events emitted from ${lpRewardsContractAddress} ` +
        `to ${pairName} pair ${this.liquidityStaking.lpTokenContract.options.address} ` +
        `between blocks ${this.liquidityStaking.lpCreationBlock} and ${this.finalBlock}`
    )

    const events = await getPastEvents(
      this.context.web3,
      this.liquidityStaking.lpTokenContract,
      "Transfer",
      this.liquidityStaking.lpCreationBlock,
      this.finalBlock
    )
    logger.info(`found ${events.length} lp ${pairName} token transfer events`)

    const lpTokenStakersSet = new Set()
    events.forEach((event) => {
      // include accounts that staked in LPReward contract only
      if (event.returnValues.to == lpRewardsContractAddress) {
        lpTokenStakersSet.add(event.returnValues.from)
      }
    })

    logger.info(`found ${lpTokenStakersSet.size} unique historic stakers`)

    return Array.from(lpTokenStakersSet)
  }

  /**
   * Retrieves balances of LP KEEP-ETH / KEEP-TBTC pair for stakers in LPRewards* contract
   *
   * @param {Array<Address>} lpStakers LP KEEP-ETH / KEEP-TBTC stakers
   *
   * @return {Map<Address,BN>} LP Balances by lp stakers
   */
  async getLpTokenStakersBalances(lpStakers) {
    const lpBalanceByStaker = new Map()
    let expectedTotalSupply = new BN(0)

    for (let i = 0; i < lpStakers.length; i++) {
      const lpBalance = new BN(
        await this.liquidityStaking.lpRewardsContract.methods
          .balanceOf(lpStakers[i])
          .call({}, this.finalBlock)
      )
      if (!lpBalance.isZero()) {
        lpBalanceByStaker.set(lpStakers[i], lpBalance)
        expectedTotalSupply = expectedTotalSupply.add(lpBalance)
      }
    }
    const actualTotalSupply = new BN(
      await this.liquidityStaking.lpRewardsContract.methods
        .totalSupply()
        .call({}, this.finalBlock)
    )

    if (!expectedTotalSupply.eq(actualTotalSupply)) {
      logger.error(
        `Sum of LP staker balances ${expectedTotalSupply} does not match the total supply ${actualTotalSupply}`
      )
    }

    logger.info(`Total supply of LP Token: ${expectedTotalSupply.toString()}`)

    return lpBalanceByStaker
  }

  /**
   * Calculates KEEP for all LP KEEP-ETH / KEEP-TBTC stakers.
   *
   * @param {Map<Address, BN>} stakersBalances LP KEEP-ETH / KEEP-TBTC Token amounts by stakers
   *
   * @return {Map<Address,BN>} KEEP Tokens in LP KEEP-ETH / KEEP-TBTC at the final block
   */
  async calcKeepInStakersBalances(stakersBalances) {
    logger.info(`check token stakers at block ${this.finalBlock}`)

    const keepInLpByStakers = new Map()

    // Retrieve current pair data
    const pairData = await getPairData(this.liquidityStaking.lpPairAddress)
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

    dumpDataToFile(
      keepInLpByStakers,
      this.liquidityStaking.keepInLpTokenFilePath
    )

    return keepInLpByStakers
  }

  /**
   * Calculating amount of KEEP token which makes a KEEP-ETH / KEEP-TBTC Uniswap pair.
   * Math is based on https://uniswap.org/docs/v2/advanced-topics/understanding-returns/
   *
   * @param {BN} lpBalance LP amount staked by a staker
   * @param {PairData} pairData KEEP-ETH / KEEP-TBTC pair data fetched from Uniswap
   *
   * @return {BN} KEEP token amounts in LP token balance
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

    for (const [pairName, pairObj] of Object.entries(
      this.liquidityStakingObjects
    )) {
      this.liquidityStaking = pairObj

      const lpStakers = await this.findStakers(pairName)
      const stakersBalances = await this.getLpTokenStakersBalances(lpStakers)
      await this.calcKeepInStakersBalances(stakersBalances)
    }

    return {}
  }
}
