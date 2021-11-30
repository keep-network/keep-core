import { helpers, ethers } from "hardhat"
import type { BigNumber } from "ethers"
// eslint-disable-next-line import/no-cycle
import { noMisbehaved, signAndSubmitArbitraryDkgResult } from "./dkg"
import { constants, params } from "../fixtures"
import blsData from "../data/bls"
import type { Operator } from "./operators"
import type { RandomBeacon, SortitionPool } from "../../typechain"

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
  randomBeacon: RandomBeacon,
  seed: BigNumber
): Promise<Operator[]> {
  const sortitionPool = (await ethers.getContractAt(
    "SortitionPool",
    await randomBeacon.sortitionPool()
  )) as SortitionPool

  const identifiers = await randomBeacon.selectGroup(seed.toHexString())
  const addresses = await sortitionPool.getIDOperators(identifiers)

  return identifiers.map((identifier, i) => ({
    id: identifier,
    address: addresses[i],
  }))
}
