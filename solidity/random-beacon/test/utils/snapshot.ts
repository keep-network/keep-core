import { ethers } from "hardhat"

const snapshotIdsStack = []

// Snapshot the state of the blockchain at the current block.
export async function createSnapshot() {
  const snapshotId = await ethers.provider.send("evm_snapshot", [])
  snapshotIdsStack.push(snapshotId)
}

// Restores the chain to a latest snapshot.
export async function restoreSnapshot() {
  const snapshotId = snapshotIdsStack.pop()
  await ethers.provider.send("evm_revert", [snapshotId])
}
