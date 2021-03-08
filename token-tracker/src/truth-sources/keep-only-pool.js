/** @typedef { import("../lib/context").Context } Context */
/** @typedef { import("../lib/ethereum-helper").Address } Address */

import BN from "bn.js"

import { ITruthSource } from "./truth-source.js"
import { getDeploymentBlockNumber } from "../lib/contract-helper.js"
import { getPastEvents, getChainID } from "../lib/ethereum-helper.js"
import { dumpDataToFile } from "../lib/file-helper.js"
import { logger } from "../lib/winston.js"
import { EthereumHelpers } from "@keep-network/tbtc.js"
const { callWithRetry } = EthereumHelpers

import KeepVaultJSON from "@keep-network/keep-core/artifacts/KeepVault.json"

const KEEP_STAKER_BALANCES_IN_KEEP_ONLY_POOL_PATH =
  "./tmp/keep-staker-balances-in-keep-only-pool.json"

export class KeepOnlyPoolTruthSource extends ITruthSource {
  /**
   * @param {Context} context
   * @param {Number} targetBlock
   */
  constructor(context, targetBlock) {
    super(context, targetBlock)
  }

  async initialize() {
    const chainID = await getChainID(this.context.web3)

    this.keepVault = EthereumHelpers.buildContract(
      this.context.web3,
      KeepVaultJSON.abi,
      KeepVaultJSON.networks[chainID].address
    )
  }

  /**
   * Finds all historic KEEP-only pool stakers in KeepVault contract based on
   * "Staked" events.
   *
   * @return {Set<Address>} All historic KEEP-only pool stakers.
   */
  async findHistoricKeepOnlyPoolStakers() {
    const keepVaultDeploymentBlock = await getDeploymentBlockNumber(
      KeepVaultJSON,
      this.context.web3
    )

    logger.info(
      `looking for "Staked" events emitted from ${this.keepVault.options.address} ` +
        `between blocks ${keepVaultDeploymentBlock} and ${this.targetBlock}`
    )

    const events = await getPastEvents(
      this.context.web3,
      this.keepVault,
      "Staked",
      keepVaultDeploymentBlock,
      this.targetBlock
    )
    logger.info(`found ${events.length} token "Staked" events`)

    const tokenStakers = new Set()
    events.forEach((event) => tokenStakers.add(event.returnValues.user))

    logger.info(
      `found ${tokenStakers.size} unique historic KEEP-only pool stakers`
    )

    return tokenStakers
  }

  /**
   * Retrieves balances of KEEP tokens locked in KeepVault contract for stakers.
   *
   * @param {Array<Address>} stakers KEEP-only pool stakers addresses.
   *
   * @return {Map<Address,BN>} Staked KEEP balances by stakers addresses.
   */
  async getTokenStakersBalances(stakers) {
    const balancesByStakerAddress = new Map()
    let expectedTotalStaking = new BN(0)

    for (const staker of stakers) {
      const stakerBalance = new BN(
        await callWithRetry(
          this.keepVault.methods.totalStakedFor(staker),
          undefined,
          undefined,
          this.targetBlock
        )
      )
      if (!stakerBalance.isZero()) {
        balancesByStakerAddress.set(staker, stakerBalance)
        expectedTotalStaking = expectedTotalStaking.add(stakerBalance)

        logger.info(`staker: ${staker} - KEEP balance: ${stakerBalance}`)
      }
    }

    const actualTotalStaking = new BN(
      await callWithRetry(
        this.keepVault.methods.totalStakingShares(),
        undefined,
        undefined,
        this.targetBlock
      )
    )

    if (!expectedTotalStaking.eq(actualTotalStaking)) {
      logger.error(
        `sum of KEEP staker balances ${expectedTotalStaking} does not match
        the total staking supply ${actualTotalStaking} of KEEP tokens
        in KeepVault contract`
      )
    }

    logger.info(
      `total supply of staked KEEP tokens locked in KeepVault is: ${actualTotalStaking.toString()}`
    )

    dumpDataToFile(
      balancesByStakerAddress,
      KEEP_STAKER_BALANCES_IN_KEEP_ONLY_POOL_PATH
    )

    return balancesByStakerAddress
  }

  /**
   * Initializes KeepVault contract and retrieves all the stakers with their
   * KEEP balances at the target block.
   *
   * @return {Map<Address,BN>} KEEP token amounts staked by stakers at the target block.
   */
  async getTokenHoldingsAtTargetBlock() {
    await this.initialize()

    const stakers = await this.findHistoricKeepOnlyPoolStakers()

    return await this.getTokenStakersBalances(stakers)
  }
}
