import { helpers } from "hardhat"

import { constants, dkgParams } from "../fixtures"
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
    await walletRegistry.relayEntry(),
    startBlock
  )

  return { tx, startBlock, dkgSeed }
}

export async function createNewWallet(
  walletRegistry: WalletRegistry,
  walletOwner: SignerWithAddress,
  publicKey?: string
): Promise<{
  members: Operator[]
  walletID: string
}> {
  const { dkgSeed, startBlock } = await requestNewWallet(
    walletRegistry,
    walletOwner
  )

  await mineBlocks(constants.offchainDkgTime)

  const {
    dkgResult,
    submitter,
    signers: members,
  } = await signAndSubmitCorrectDkgResult(
    walletRegistry,
    publicKey || ecdsaData.group1.publicKey,
    dkgSeed,
    startBlock,
    noMisbehaved
  )

  await mineBlocks(dkgParams.dkgResultChallengePeriodLength)

  const approveDkgResultTx = await walletRegistry
    .connect(submitter)
    .approveDkgResult(dkgResult)

  const walletID: string = await getWalletID(approveDkgResultTx)

  return { members, walletID }
}

export async function getWalletID(
  approveDkgResultTx: ContractTransaction
): Promise<string> {
  const { walletID } = (await approveDkgResultTx.wait()).events.find(
    (e) => e.event === "WalletCreated"
  ).args

  return walletID
}
