// Snapshots are a feature of some EVM implementations for improved dev UX.
// They allow us to snapshot the entire state of the chain, and restore it at a later point.
// https://github.com/trufflesuite/ganache-core/blob/master/README.md#custom-methods
const { web3 } = require("@openzeppelin/test-environment")

const snapshotIdsStack = []

/**
 * Snapshot the state of the blockchain at the current block
 */
async function createSnapshot() {
  return await new Promise((res, rej) => {
    web3.currentProvider.send(
      {
        jsonrpc: "2.0",
        method: "evm_snapshot",
        params: [],
      },
      function (err, result) {
        if (err) rej(err)
        const snapshotId = result.result
        snapshotIdsStack.push(snapshotId)
        res()
      }
    )
  })
}

/**
 * Restores the chain to a latest snapshot
 */
async function restoreSnapshot() {
  const snapshotId = snapshotIdsStack.pop()
  return await new Promise((res, rej) => {
    web3.currentProvider.send(
      {
        jsonrpc: "2.0",
        method: "evm_revert",
        params: [snapshotId],
      },
      function (err, result) {
        if (err) rej(err)
        else res()
      }
    )
  })
}

module.exports.createSnapshot = createSnapshot
module.exports.restoreSnapshot = restoreSnapshot
