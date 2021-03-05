/** @typedef { import("../lib/context").Context } Context */
/** @typedef { import("../lib/ethereum-helper").Address } Address */

import BN from "bn.js"

import { ITruthSource } from "./truth-source.js"
import { getPastEvents, getChainID } from "../lib/ethereum-helper.js"
import { dumpDataToFile } from "../lib/file-helper.js"
import { logger } from "../lib/winston.js"
import { EthereumHelpers } from "@keep-network/tbtc.js"
const { callWithRetry } = EthereumHelpers

// https://etherscan.io/address/0xe6f19dab7d43317344282f803f8e8d240708174a#code
import keepEthTokenJson from "../../artifacts/KEEP-ETH-UNI-V2-Token.json"
// https://etherscan.io/address/0x38c8ffee49f286f25d25bad919ff7552e5daf081#code
import keepTbtcTokenJson from "../../artifacts/KEEP-TBTC-UNI-V2-Token.json"

import LPRewardsKEEPETHJson from "@keep-network/keep-ecdsa/artifacts/LPRewardsKEEPETH.json"
import LPRewardsKEEPTBTCJson from "@keep-network/keep-ecdsa/artifacts/LPRewardsKEEPTBTC.json"

// https://info.uniswap.org/pair/0xe6f19dab7d43317344282f803f8e8d240708174a
const KEEPETH_PAIR_ADDRESS = "0xe6f19dab7d43317344282f803f8e8d240708174a"
// https://info.uniswap.org/pair/0x38c8ffee49f286f25d25bad919ff7552e5daf081
const KEEPTBTC_PAIR_ADDRESS = "0x38c8ffee49f286f25d25bad919ff7552e5daf081"

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
   * @param {Number} targetBlock
   */
  constructor(context, targetBlock) {
    super(context, targetBlock)
  }

  async initialize() {
    const chainID = await getChainID(this.context.web3)

    const keepEthTokenContract = EthereumHelpers.buildContract(
      this.context.web3,
      JSON.parse(keepEthTokenJson.result),
      KEEPETH_PAIR_ADDRESS
    )

    const lpRewardKeepEthContract = EthereumHelpers.buildContract(
      this.context.web3,
      LPRewardsKEEPETHJson.abi,
      LPRewardsKEEPETHJson.networks[chainID].address
    )

    const keepTbtcTokenContract = EthereumHelpers.buildContract(
      this.context.web3,
      JSON.parse(keepTbtcTokenJson.result),
      KEEPTBTC_PAIR_ADDRESS
    )

    const lpRewardKeepTbtcContract = EthereumHelpers.buildContract(
      this.context.web3,
      LPRewardsKEEPTBTCJson.abi,
      LPRewardsKEEPTBTCJson.networks[chainID].address
    )

    this.liquidityStakingObjects = {
      KEEPETH: {
        lpTokenContract: keepEthTokenContract,
        lpRewardsContract: lpRewardKeepEthContract,
        lpCreationBlock: KEEPETH_CREATION_BLOCK,
        keepInLpTokenFilePath: KEEP_IN_LP_KEEPETH_BALANCES_PATH,
        lpPairAddress: KEEPETH_PAIR_ADDRESS,
      },
      KEEPTBTC: {
        lpTokenContract: keepTbtcTokenContract,
        lpRewardsContract: lpRewardKeepTbtcContract,
        lpCreationBlock: KEEPTBTC_CREATION_BLOCK,
        keepInLpTokenFilePath: KEEP_IN_LP_KEEPTBTC_BALANCES_PATH,
        lpPairAddress: KEEPTBTC_PAIR_ADDRESS,
      },
    }
  }

  /**
   * Finds all historic stakers of LP KEEP-ETH / KEEP-TBTC pair token based on
   * "Transfer" events
   *
   * @param {String} pairName LP pair name
   * @param {Object} pairObj LP pair object
   *
   * @return {Set<Address>} All historic LP KEEP-ETH / KEEP-TBTC token stakers
   * */
  async findStakers(pairName, pairObj) {
    const lpRewardsContractAddress = pairObj.lpRewardsContract.options.address

    logger.info(
      `looking for Transfer events emitted from ${lpRewardsContractAddress} ` +
        `to ${pairName} pair ${pairObj.lpTokenContract.options.address} ` +
        `between blocks ${pairObj.lpCreationBlock} and ${this.targetBlock}`
    )

    const events = await getPastEvents(
      this.context.web3,
      pairObj.lpTokenContract,
      "Transfer",
      pairObj.lpCreationBlock,
      this.targetBlock
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
   * @param {Object} pairObj LP pair object
   *
   * @return {Map<Address,BN>} LP Balances by lp stakers
   */
  async getLpTokenStakersBalances(lpStakers, pairObj) {
    const lpBalanceByStaker = new Map()
    let expectedTotalSupply = new BN(0)

    for (const staker of lpStakers) {
      const lpBalance = new BN(
        await callWithRetry(
          pairObj.lpRewardsContract.methods.balanceOf(staker),
          undefined,
          undefined,
          this.targetBlock
        )
      )
      if (!lpBalance.isZero()) {
        lpBalanceByStaker.set(staker, lpBalance)
        expectedTotalSupply = expectedTotalSupply.add(lpBalance)
      }
    }
    const actualTotalSupply = new BN(
      await callWithRetry(
        pairObj.lpRewardsContract.methods.totalSupply(),
        undefined,
        undefined,
        this.targetBlock
      )
    )

    if (!expectedTotalSupply.eq(actualTotalSupply)) {
      logger.error(
        `sum of LP staker balances ${expectedTotalSupply} does not match the total supply ${actualTotalSupply}`
      )
    }

    logger.info(`total supply of LP Token: ${expectedTotalSupply.toString()}`)

    return lpBalanceByStaker
  }

  /**
   * Calculates KEEP for all LP KEEP-ETH / KEEP-TBTC stakers.
   *
   * @param {Map<Address, BN>} stakersBalances LP KEEP-ETH / KEEP-TBTC Token amounts by stakers
   * @param {String} pairName LP pair name
   * @param {Object} pairObj LP pair object
   *
   * @return {Map<Address,BN>} KEEP Tokens in LP KEEP-ETH / KEEP-TBTC at the target block
   */
  async calcKeepInStakersBalances(stakersBalances, pairName, pairObj) {
    logger.info(`check token stakers at block ${this.targetBlock}`)

    const totalSupply = await callWithRetry(
      pairObj.lpTokenContract.methods.totalSupply(),
      undefined,
      undefined,
      this.targetBlock
    )
    logger.info(
      `total supply of LP ${pairName} token at block ${this.targetBlock} is: ${totalSupply}`
    )

    const lpReserves = await callWithRetry(
      pairObj.lpTokenContract.methods.getReserves(),
      undefined,
      undefined,
      this.targetBlock
    )
    logger.info(
      `KEEP Liquidity pool at block ${this.targetBlock} is: ${lpReserves._reserve0}`
    )

    const lpPairData = {
      keepLiquidityPool: lpReserves._reserve0,
      lpTotalSupply: totalSupply,
    }

    const keepInLpByStakers = new Map()

    for (const [stakerAddress, lpStakerBalance] of stakersBalances.entries()) {
      const keepInLPToken = await this.calcKeepTokenfromLPToken(
        lpStakerBalance,
        lpPairData
      )
      keepInLpByStakers.set(stakerAddress, keepInLPToken)

      logger.info(
        `staker: ${stakerAddress} - LP Balance: ${lpStakerBalance} - KEEP in LP: ${keepInLPToken}`
      )
    }

    logger.info(
      `found ${keepInLpByStakers.size} stakers at block ${this.targetBlock}`
    )

    dumpDataToFile(keepInLpByStakers, pairObj.keepInLpTokenFilePath)

    return keepInLpByStakers
  }

  /**
   * Calculates amount of KEEP token which makes a KEEP-ETH / KEEP-TBTC Uniswap pair.
   * A Uniswap LP pair is a bookeeping tool to keep track of how much the liquidity
   * stakers are owed. They store two assets of equivalent value of each, ex. KEEP-ETH.
   * This means that the value of KEEP owned is dependent on the ratio of staked
   * LP tokens and the total LP supply. Ratio between LP tokens and KEEP tokens
   * should be equal:
   * LP_staker_balance / LP_total_supply_pool == KEEP_staker_owed / KEEP_total_liquidity_pool
   * Now, the number of KEEP tokens which makes a KEEP-ETH pair can be calculated
   * using the following equation:
   * KEEP_staker_owed = (LP_staker_balance * KEEP_total_liquidity_pool) / LP_total_supply_pool
   * where:
   * LP_staker_balance is retrieved from LPRewardsContract
   * KEEP_total_liquidity_pool is fetched from Uniswap LP Token - lpToken.getReserves()._reserve0
   * LP_total_supply_pool is fetched from Uniswap LP Token - lpToken.totalSupply()
   *
   * References:
   * Returns in Uniswap: https://uniswap.org/docs/v2/advanced-topics/understanding-returns/
   *
   * @param {BN} lpStakerBalance LP amount staked by a staker in a LPRewardsContract
   * @param {PairData} lpPairData Pair data fetched from LP Token Contract
   *
   * @return {BN} KEEP token amounts in LP token balance
   */
  async calcKeepTokenfromLPToken(lpStakerBalance, lpPairData) {
    return lpStakerBalance
      .mul(new BN(lpPairData.keepLiquidityPool))
      .div(new BN(lpPairData.lpTotalSupply))
  }

  /**
   * @return {Map<Address,BN>} KEEP token amounts staked by stakers at the target block
   */
  async getTokenHoldingsAtTargetBlock() {
    await this.initialize()
    const keepInLPsByStakers = new Map()

    for (const [pairName, pairObj] of Object.entries(
      this.liquidityStakingObjects
    )) {
      const lpStakers = await this.findStakers(pairName, pairObj)
      const stakersBalances = await this.getLpTokenStakersBalances(
        lpStakers,
        pairObj
      )
      const keepInLpByStakers = await this.calcKeepInStakersBalances(
        stakersBalances,
        pairName,
        pairObj
      )

      keepInLpByStakers.forEach((balance, staker) => {
        if (keepInLPsByStakers.has(staker)) {
          keepInLPsByStakers.get(staker).iadd(balance)
        } else {
          keepInLPsByStakers.set(staker, new BN(balance))
        }
      })
    }

    return keepInLPsByStakers
  }
}
