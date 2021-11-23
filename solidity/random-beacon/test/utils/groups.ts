import { helpers, ethers } from "hardhat"
import { BigNumber } from "ethers"
// eslint-disable-next-line import/no-cycle
import { noMisbehaved, signAndSubmitArbitraryDkgResult } from "./dkg"
import { constants, params } from "../fixtures"
import blsData from "../data/bls"
import { Operator } from "./operators"
import type { RandomBeacon, SortitionPool } from "../../typechain"

const { mineBlocks } = helpers.time

// eslint-disable-next-line import/prefer-default-export
export async function createGroup(
  randomBeacon: RandomBeacon,
  signers: Operator[]
): Promise<void> {
  const { blockNumber: startBlock } = await randomBeacon.genesis()
  const submitterIndex = 1
  await mineBlocks(constants.offchainDkgTime)
  const { dkgResult } = await signAndSubmitArbitraryDkgResult(
    randomBeacon,
    blsData.groupPubKey,
    signers,
    startBlock,
    noMisbehaved,
    submitterIndex
  )

  await mineBlocks(params.dkgResultChallengePeriodLength)
  await randomBeacon
    .connect(await ethers.getSigner(signers[submitterIndex - 1].address))
    .approveDkgResult(dkgResult)
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
