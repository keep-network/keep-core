/** @typedef { import("../lib/context").default } Context */
/** @typedef { import("../lib/ethereum-helper").Address } Address */
/** @typedef { import("bn.js") } BN */

/**
 * @typedef {Object} ITruthSource
 * @prop {Context} context Description.
 * @prop {Number} finalBlock Description.
 */

import Web3 from "web3"
const { toChecksumAddress } = Web3.utils

/**
 * This is an interface that should be implemented by all sources of truth for
 * token ownership inspection.
 * @property {Context} context Context instance.
 * @property {Number} finalBlock Block at which ownership should be inspected.
 * @property {IgnoredContractsResolver} ignoredContractsResolver Resolver for
 * filtering owners that are contracts. If the owner is one of the ignored contracts
 * it should be removed as it should be already handled by a truth source implementation,
 * to get an actual owner.
 */
export class ITruthSource {
  constructor(
    /** @type {Context} */ context,
    /** @type {Number} */ finalBlock
  ) {
    this.context = context
    this.finalBlock = finalBlock

    this.addressesResolver = new AddressesResolver(context.web3)
  }

  /**
   * Returns a map of addresses with their token holdings discovered using the
   * specific source of truth.
   * @return {Map<Address,BN>} Token holdings at the final blocks.
   */
  async getTokenHoldingsAtFinalBlock() {
    throw new Error("getTokenHoldingsAtFinalBlock not implemented")
  }
}
