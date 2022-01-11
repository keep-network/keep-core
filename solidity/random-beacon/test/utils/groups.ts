import { helpers, ethers } from "hardhat"
import type { BigNumber } from "ethers"
// eslint-disable-next-line import/no-cycle
import { noMisbehaved, signAndSubmitArbitraryDkgResult } from "./dkg"
import { constants, params } from "../fixtures"
import blsData from "../data/bls"
import type { Operator } from "./operators"
import type { RandomBeacon, SortitionPool } from "../../typechain"

const { keccak256, defaultAbiCoder } = ethers.utils
const { mineBlocks } = helpers.time

export async function createGroup(
  randomBeacon: RandomBeacon,
  signers: Operator[]
): Promise<void> {
  const { blockNumber: startBlock } = await randomBeacon.genesis()

  await mineBlocks(constants.offchainDkgTime)

  const { dkgResult, submitter } = await signAndSubmitArbitraryDkgResult(
    randomBeacon,
    blsData.groupPubKey,
    signers,
    startBlock,
    noMisbehaved
  )

  await mineBlocks(params.dkgResultChallengePeriodLength)

  await randomBeacon.connect(submitter).approveDkgResult(dkgResult)
}

export async function selectGroup(
  sortitionPool: SortitionPool,
  seed: BigNumber
): Promise<Operator[]> {
  const identifiers = await sortitionPool.selectGroup(
    constants.groupSize,
    seed.toHexString()
  )
  const addresses = await sortitionPool.getIDOperators(identifiers)

  return identifiers.map((identifier, i) => ({
    id: identifier,
    address: addresses[i],
  }))
}

export function hashUint32Array(arrayToHash: number[]) {
  return keccak256(defaultAbiCoder.encode(["uint32[]"], [arrayToHash]))
}
