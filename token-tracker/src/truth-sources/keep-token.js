/** @typedef { import("../lib/context").Context } Context */
/** @typedef { import("../lib/ethereum-helper").Address } Address */

import { ITruthSource } from "./index.js"
import { getPastEvents, callWithRetry } from "../lib/ethereum-helper.js"
import { writeFileSync, readFileSync } from "fs"
import BN from "bn.js"
import { mapToObject } from "../lib/map-helper.js"

const KEEP_TOKEN_HISTORIC_HOLDERS_DUMP_PATH = "./tmp/keep-token-holders.json"
const KEEP_TOKEN_BALANCES_DUMP_PATH = "./tmp/keep-token-balances.json"

/**
 * TODO: Write docs
 * Short description.
 * @typedef {ITruthSource} KeepTokenTruthSource
 */

export class KeepTokenTruthSource extends ITruthSource {
  /**
   * @param {Context} context
   * @param {Number} finalBlock
   */
  constructor(context, finalBlock) {
    super(context, finalBlock)
  }

  async initialize() {
    this.keepToken = await this.context.contracts.KeepToken.deployed()
  }

  /**
   * @returns {Array<Address>} All historic token holders.
   */
  async findHistoricHolders() {
    console.info(
      `looking for Transfer events emitted from ${this.keepToken.options.address} ` +
        `between blocks ${this.context.contracts.deploymentBlock} and ${this.finalBlock}`
    )

    const events = await getPastEvents(
      this.context.web3,
      this.keepToken,
      "Transfer",
      this.context.contracts.deploymentBlock,
      this.finalBlock
    )
    console.info(`found ${events.length} token transfer events`)

    const allTokenHoldersSet = new Set()
    events.forEach((event) => allTokenHoldersSet.add(event.returnValues.to))

    const allTokenHolders = Array.from(allTokenHoldersSet)
    console.info(`found ${allTokenHolders.length} unique historic holders`)

    console.info(
      `dump all historic token holders to a file: ${KEEP_TOKEN_HISTORIC_HOLDERS_DUMP_PATH}`
    )
    writeFileSync(
      KEEP_TOKEN_HISTORIC_HOLDERS_DUMP_PATH,
      JSON.stringify(allTokenHolders, null, 2)
    )

    return allTokenHolders
  }

  /**
   * @param {Array<Address>} tokenHolders Token holders to check.
   * @returns {Map<Address,BN} Token holdings at the final blocks.
   */
  async checkFinalHoldersBalances(tokenHolders) {
    console.info(`check token holding at block ${this.finalBlock}`)

    /** @type {Map<Address,BN>} */
    const holdersBalances = new Map()

    for (const holderAddress of tokenHolders) {
      const finalBalance = new BN(
        await callWithRetry(
          this.keepToken.methods.balanceOf(holderAddress),
          undefined,
          undefined,
          this.finalBlock
        )
      )

      console.debug(`${holderAddress}: ${finalBalance}`)

      if (finalBalance.gtn(0)) {
        holdersBalances.set(holderAddress, finalBalance)
      }
    }

    console.info(
      `found ${holdersBalances.length} holders at block ${this.finalBlock}`
    )

    console.info(
      `dump token holdings to a file: ${KEEP_TOKEN_BALANCES_DUMP_PATH}`
    )
    writeFileSync(
      KEEP_TOKEN_BALANCES_DUMP_PATH,
      JSON.stringify(mapToObject(holdersBalances), null, 2)
    )

    return holdersBalances
  }

  /**
   * @returns {Map<Address,BN>} Token holdings at the final blocks.
   */
  async getTokenHoldingsAtFinalBlock() {
    await this.initialize()

    const allTokenHolders = await this.findHistoricHolders()

    const holdersBalances = await this.checkFinalHoldersBalances(
      allTokenHolders
    )

    return holdersBalances
  }
}
