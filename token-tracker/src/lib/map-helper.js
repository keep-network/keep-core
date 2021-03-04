/** @typedef { import("./ethereum-helper").Address } Address */

import BN from "bn.js"

import web3utils from "web3-utils"
const { toChecksumAddress } = web3utils

/**
 * Adds balances of new tokens holdings to the existing map of holdings.
 * @param {Map<Address,BN>} oldHoldings Map of token holders to add.
 * @param {Map<Address,BN>} newHoldings Map of token holders to add.
 * @return {Map<Address,BN>}
 */
export function addTokenHoldings(oldHoldings, newHoldings) {
  const result = oldHoldings

  newHoldings.forEach((value, holder) => {
    holder = toChecksumAddress(holder)
    value = new BN(value)

    if (result.has(holder)) {
      result.get(holder).iadd(value)
    } else {
      result.set(holder, value)
    }
  })

  return result
}
