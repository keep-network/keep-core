import { helpers } from "hardhat"

import { constants, params } from "../fixtures"
import ecdsaData from "../data/ecdsa"

import {
  calculateDkgSeed,
  noMisbehaved,
  signAndSubmitCorrectDkgResult,
} from "./dkg"

import type { BigNumber, ContractTransaction } from "ethers"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { WalletFactory } from "../../typechain"
import type { Operator } from "./operators"

const { mineBlocks } = helpers.time

export async function requestNewWallet(
  walletFactory: WalletFactory,
  walletManager: SignerWithAddress
): Promise<{
  tx: ContractTransaction
  startBlock: number
  dkgSeed: BigNumber
}> {
  const tx: ContractTransaction = await walletFactory
    .connect(walletManager)
    .requestNewWallet()

  const startBlock: number = tx.blockNumber

  const dkgSeed: BigNumber = calculateDkgSeed(
    await walletFactory.relayEntry(),
    startBlock
  )

  return { tx, startBlock, dkgSeed }
}

export async function createNewWallet(
  walletFactory: WalletFactory,
  walletManager: SignerWithAddress
): Promise<{ members: Operator[] }> {
  const { dkgSeed, startBlock } = await requestNewWallet(
    walletFactory,
    walletManager
  )

  await mineBlocks(constants.offchainDkgTime)

  const {
    dkgResult,
    submitter,
    signers: members,
  } = await signAndSubmitCorrectDkgResult(
    walletFactory,
    ecdsaData.group1.publicKey,
    dkgSeed,
    startBlock,
    noMisbehaved
  )

  await mineBlocks(params.dkgResultChallengePeriodLength)

  await walletFactory.connect(submitter).approveDkgResult(dkgResult)

  return { members }
}
