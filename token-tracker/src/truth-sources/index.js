/** @typedef { import("../lib/context").default } Context */
/** @typedef { import("../lib/ethereum-helper").Address } Address */
/** @typedef { import("web3").default } Web3 */
/** @typedef { import("bn.js") } BN */

/**
 * @typedef {Object} ITruthSource
 * @prop {Context} context Description.
 * @prop {Number} finalBlock Description.
 */

export class ITruthSource {
  /**
   * @param {Context} context
   * @param {Number} finalBlock
   */
  constructor(context, finalBlock) {
    this.context = context
    this.finalBlock = finalBlock
  }

  /**
   * @return {Map<Address,BN>} Token holdings at the final blocks.
   */
  async getTokenHoldingsAtFinalBlock() {
    throw new Error("getTokenHoldingsAtFinalBlock not implemented")
  }
}

export function TruthSources() {
  return [KeepTokenTruthSource]
}
