import { helpers, ethers } from "hardhat"

import { params } from "../fixtures"
import ecdsaData from "../data/ecdsa"

import { noMisbehaved, signAndSubmitCorrectDkgResult } from "./dkg"
import { fakeRandomBeacon } from "./randomBeacon"

import type { WalletRegistry } from "../../typechain"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { Operator } from "./operators"

const { mineBlocks } = helpers.time
const { keccak256 } = ethers.utils

// eslint-disable-next-line import/prefer-default-export
export async function createNewWallet(
  walletRegistry: WalletRegistry,
  walletOwner: SignerWithAddress,
  publicKey = ecdsaData.group1.publicKey
): Promise<{
  members: Operator[]
  publicKeyHash: string
}> {
  const tx = await walletRegistry.connect(walletOwner).requestNewWallet()

  const randomBeacon = await fakeRandomBeacon(walletRegistry)

  const relayEntry = ethers.utils.randomBytes(32)

  const dkgSeed = ethers.BigNumber.from(keccak256(relayEntry))

  // eslint-disable-next-line no-underscore-dangle
  await walletRegistry
    .connect(randomBeacon.wallet)
    .__beaconCallback(relayEntry, 0)

  const {
    dkgResult,
    submitter,
    signers: members,
  } = await signAndSubmitCorrectDkgResult(
    walletRegistry,
    publicKey,
    dkgSeed,
    tx.blockNumber,
    noMisbehaved
  )

  await mineBlocks(params.dkgResultChallengePeriodLength)

  await walletRegistry.connect(submitter).approveDkgResult(dkgResult)

  return { members, publicKeyHash: keccak256(publicKey) }
}
