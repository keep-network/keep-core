/** @typedef { import("./lib/context").Context } Context */
/** @typedef { import("./truth-sources/index").ITruthSource } ITruthSource */
/** @typedef { import("./lib/ethereum-helper").Address } Address */

import BN from "bn.js"

import web3utils from "web3-utils"
const { toChecksumAddress } = web3utils

export class Inspector {
  /**
   * @param {Context} context
   */
  constructor(context) {
    this.context = context

    /** @type {Set<Function>} */
    this.truthSources = new Set()

    /** @type {Map<Address,BN>} */
    this.tokenHoldings = new Map()
  }

  /**
   *
   * @param {Function} TruthSource
   */
  registerTruthSource(TruthSource) {
    if (this.truthSources.has(TruthSource)) {
      throw new Error(`truth source already registered: ${TruthSource.name}`)
    }

    this.truthSources.add(TruthSource)
  }

  /**
   *
   * @param {Number} finalBlock
   * @returns {Map<Address,BN>}
   */
  async getOwnershipsAtBlock(finalBlock) {
    for (const TruthSource of this.truthSources) {
      /** @type {ITruthSource} */
      const truthSourceInstance = new TruthSource(this.context, finalBlock)

      const newHoldings = await truthSourceInstance.getTokenHoldingsAtFinalBlock()

      this.addTokenHoldings(newHoldings)
    }
    return this.tokenHoldings
  }

  /**
   *
   * @param {Map<Address,BN>} newHoldings
   */
  addTokenHoldings(newHoldings) {
    newHoldings.forEach((value, holder) => {
      holder = toChecksumAddress(holder)
      value = new BN(value)

      if (!this.tokenHoldings[holder]) {
        this.tokenHoldings[holder] = value
      } else {
        this.tokenHoldings[holder].iadd(value)
      }
    })
  }
}
