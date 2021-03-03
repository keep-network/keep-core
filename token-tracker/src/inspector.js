/** @typedef { import("./lib/context").default } Context */
/** @typedef { import("./truth-sources/truth-source").ITruthSource } ITruthSource */
/** @typedef { import("./lib/ethereum-helper").Address } Address */

import { logger } from "./lib/winston.js"
import { getDeploymentBlockNumber } from "./lib/contract-helper.js"
import { Contract } from "./lib/contract-helper.js"
import { addTokenHoldings } from "./lib/map-helper.js"

import KeepTokenJson from "@keep-network/keep-core/artifacts/KeepToken.json"
import TokenStakingJSON from "@keep-network/keep-core/artifacts/TokenStaking.json"
import OldTokenStakingJSON from "@keep-network/keep-core/artifacts/OldTokenStaking.json"

export class Inspector {
  /** @type {Set<Function>} */
  #truthSources

  constructor(/** @type {Context}*/ context) {
    this.context = context

    this.#truthSources = new Set()

    /** @type {Map<Address,BN>} */
    this.tokenHoldings = new Map()
  }

  /**
   * Registers new ITruthSource class implementation.
   * @param {Function} TruthSource
   */
  registerTruthSource(TruthSource) {
    if (this.#truthSources.has(TruthSource)) {
      throw new Error(`truth source already registered: ${TruthSource.name}`)
    }

    this.#truthSources.add(TruthSource)
  }

  /**
   * Gets token ownership at the given block for all registered sources of truth.
   * @param {Number} targetBlock Block to check token holdings at.
   * @return {Promise<Map<Address,BN>>} Map of owners to their balances.
   */
  async getOwnershipsAtBlock(targetBlock) {
    // FIXME: We can get rid of global deployment block tracking and switch to
    // particular contracts deployment blocks if we use tbtc.js like functions
    // to get past events with a contract instance defining deployment block.
    this.context.deploymentBlock = await getDeploymentBlockNumber(
      KeepTokenJson,
      this.context.web3
    )

    logger.debug(`deployment block: ${this.context.deploymentBlock}`)

    // Initialize common contracts used by multiple sources of truth.
    this.context.addContract(
      "TokenStaking",
      new Contract(TokenStakingJSON, this.context.web3)
    )
    this.context.addContract(
      "OldTokenStaking",
      new Contract(OldTokenStakingJSON, this.context.web3)
    )

    for (const TruthSource of this.#truthSources) {
      /** @type {ITruthSource} */
      const truthSourceInstance = new TruthSource(this.context, targetBlock)

      const newHoldings = await truthSourceInstance.getTokenHoldingsAtTargetBlock()

      this.tokenHoldings = addTokenHoldings(this.tokenHoldings, newHoldings)
    }
    return this.tokenHoldings
  }
}
