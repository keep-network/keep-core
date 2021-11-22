import { helpers, ethers } from "hardhat"
import { noMisbehaved, signAndSubmitDkgResult } from "./dkg"
import { constants, params } from "../fixtures"
import blsData from "../data/bls"
import { Operator } from "./sortitionpool"
import type { RandomBeacon } from "../../typechain"

const { mineBlocks } = helpers.time

// eslint-disable-next-line import/prefer-default-export
export async function createGroup(
  randomBeacon: RandomBeacon,
  signers: Operator[]
): Promise<void> {
  const { blockNumber: startBlock } = await randomBeacon.genesis()
  const submitterIndex = 1
  await mineBlocks(constants.offchainDkgTime)
  await signAndSubmitDkgResult(
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
    .approveDkgResult()
}
