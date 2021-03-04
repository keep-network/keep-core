// Code copied from https://github.com/keep-network/tbtc.js/blob/4d28dc4d6c2db2d5430ad03ccfd891b1856a5f4c/bin/owner-lookup.js

/** @typedef {import ('web3').default} Web3 */

import EthereumHelpers from "@keep-network/tbtc.js/src/EthereumHelpers.js"
import ManagedGrantJSON from "@keep-network/keep-core/artifacts/ManagedGrant.json"

const ManagedGrantABI = ManagedGrantJSON.abi

export async function resolveGrantee(
  /** @type {Web3} */ web3,
  /** @type {string} */ grantee
) {
  if ((await web3.eth.getStorageAt(grantee, 0)) === "0x") {
    return grantee // grantee is already a user-owned account
  } else {
    try {
      const grant = EthereumHelpers.buildContract(
        web3,
        // @ts-ignore Oh but this is an AbiItem[]
        ManagedGrantABI,
        grantee
      )

      return await grant.methods.grantee().call()
    } catch (_) {
      // If we threw, assume this isn't a ManagedGrant and the
      // grantee is just an unknown contract---e.g. Gnosis Safe.
      return grantee
    }
  }
}
