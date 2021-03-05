/** @typedef { import("../lib/context").Context } Context */
/** @typedef { import("../lib/ethereum-helper").Address } Address */

import BN from "bn.js"

import { ITruthSource } from "./truth-source.js"
import { Contract, getDeploymentBlockNumber } from "../lib/contract-helper.js"
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

  async findHistoricKeepOnlyPoolStakers() {
    const keepVaultDeploymentBlock = await getDeploymentBlockNumber(
      KeepVaultJSON,
      this.context.web3
    )

    console.log("keepVaultDeploymentBlock: ", keepVaultDeploymentBlock)

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

    logger.info(`found ${tokenStakers.size} unique historic stakers`)

    return tokenStakers
  }

  async getTokenStakersBalances(stakers) {
    const balancesByStaker = new Map()

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
        balancesByStaker.set(staker, stakerBalance)
      }
    }

    console.log("balancesByStaker: ", balancesByStaker)

    return balancesByStaker
  }

  async getTokenHoldingsAtTargetBlock() {
    await this.initialize()

    const stakers = await this.findHistoricKeepOnlyPoolStakers()

    await this.getTokenStakersBalances(stakers)

    return {}
  }
}
