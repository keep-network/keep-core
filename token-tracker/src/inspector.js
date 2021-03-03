/** @typedef { import("./lib/context").default } Context */
/** @typedef { import("./truth-sources/truth-source").ITruthSource } ITruthSource */
/** @typedef { import("./lib/ethereum-helper").Address } Address */

import BN from "bn.js"

import Web3 from "web3"
const { toChecksumAddress } = Web3.utils

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
    for (const TruthSource of this.#truthSources) {
      /** @type {ITruthSource} */
      const truthSourceInstance = new TruthSource(this.context, targetBlock)

      const newHoldings = await truthSourceInstance.getTokenHoldingsAtTargetBlock()

      this.addTokenHoldings(newHoldings)
    }
    return this.tokenHoldings
  }

  /**
   * Adds balances of new tokens holdings to the existing map of holdings. 
   * @param {Map<Address,BN>} newHoldings Map of token holders to add.
   */
  addTokenHoldings(newHoldings) {
    newHoldings.forEach((value, holder) => {
      holder = toChecksumAddress(holder)
      value = new BN(value)

      if (!this.tokenHoldings.has(holder)) {
        this.tokenHoldings.set(holder, new BN(0))
      }

      this.tokenHoldings.get(holder).iadd(value)
    })
  }
}
