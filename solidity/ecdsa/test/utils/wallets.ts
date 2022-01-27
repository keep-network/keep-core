import { helpers } from "hardhat"
import { keccak256 } from "ethers/lib/utils"

import { params } from "../fixtures"
import ecdsaData from "../data/ecdsa"

import {
  calculateDkgSeed,
  noMisbehaved,
  signAndSubmitCorrectDkgResult,
} from "./dkg"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { BigNumber, ContractTransaction } from "ethers"
import type { WalletRegistry } from "../../typechain"
import type { Operator } from "./operators"

const { mineBlocks } = helpers.time

export async function requestNewWallet(
  walletRegistry: WalletRegistry,
  walletOwner: SignerWithAddress
): Promise<{
  tx: ContractTransaction
  startBlock: number
  dkgSeed: BigNumber
}> {
  const tx: ContractTransaction = await walletRegistry
    .connect(walletOwner)
    .requestNewWallet()

  const startBlock: number = tx.blockNumber

  const dkgSeed: BigNumber = calculateDkgSeed(
    await walletRegistry.randomRelayEntry(),
    startBlock
  )

  return { tx, startBlock, dkgSeed }
}

export async function createNewWallet(
  walletRegistry: WalletRegistry,
  walletOwner: SignerWithAddress,
  publicKey = ecdsaData.group1.publicKey
): Promise<{
  members: Operator[]
  publicKeyHash: string
}> {
  const { dkgSeed, startBlock } = await requestNewWallet(
    walletRegistry,
    walletOwner
  )

  const {
    dkgResult,
    submitter,
    signers: members,
  } = await signAndSubmitCorrectDkgResult(
    walletRegistry,
    publicKey,
    dkgSeed,
    startBlock,
    noMisbehaved
  )

  await mineBlocks(params.dkgResultChallengePeriodLength)

  await walletRegistry.connect(submitter).approveDkgResult(dkgResult)

  return { members, publicKeyHash: keccak256(publicKey) }
}
