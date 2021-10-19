/* eslint-disable no-await-in-loop */

import { ethers, getUnnamedAccounts } from "hardhat"
import { expect } from "chai"
import type { Address } from "hardhat-deploy/types"
import type { BigNumber, ContractTransaction } from "ethers"
import blsData from "../data/bls"
import type { RandomBeacon } from "../../typechain"

export type DkgGroupSigners = Map<number, Address>

export interface DkgResult {
  submitterMemberIndex: number
  groupPubKey: string
  misbehaved: string
  signatures: string
  signingMemberIndices: number[]
  members: string[]
}

export async function getDkgGroupSigners(
  groupSize: number,
  startAccountsOffset: number
): Promise<DkgGroupSigners> {
  const signers = new Map<number, Address>()

  for (let i = 1; i <= groupSize; i++) {
    const signer = (await getUnnamedAccounts())[startAccountsOffset + i]

    await expect(
      signer,
      `signer [${i}] is not defined; check hardhat network configuration`
    ).is.not.empty

    signers.set(i, signer)
  }

  return signers
}

export async function genesis(
  randomBeacon: RandomBeacon
): Promise<[ContractTransaction, BigNumber]> {
  const tx = await randomBeacon.genesis()

  const expectedSeed = ethers.BigNumber.from(
    ethers.utils.keccak256(
      ethers.utils.solidityPack(
        ["uint256", "uint256"],
        [await randomBeacon.genesisSeed(), tx.blockNumber]
      )
    )
  )

  return [tx, expectedSeed]
}

export async function signAndSubmitDkgResult(
  randomBeacon: RandomBeacon,
  groupPublicKey: string,
  signers: DkgGroupSigners,
  startBlock: number,
  submitterIndex = 1
): Promise<{
  transaction: ContractTransaction
  dkgResult: DkgResult
  dkgResultHash: string
  members: string[]
}> {
  const noMisbehaved = "0x"

  const { members, signingMemberIndices, signaturesBytes } =
    await signDkgResult(signers, groupPublicKey, noMisbehaved, startBlock)

  const dkgResult = {
    submitterMemberIndex: submitterIndex,
    groupPubKey: blsData.groupPubKey,
    misbehaved: noMisbehaved,
    signatures: signaturesBytes,
    signingMemberIndices,
    members,
  }

  const dkgResultHash = ethers.utils.keccak256(
    ethers.utils.defaultAbiCoder.encode(
      [
        "(uint256 submitterMemberIndex, bytes groupPubKey, bytes misbehaved, bytes signatures, uint256[] signingMemberIndices, address[] members)",
      ],
      [dkgResult]
    )
  )

  const transaction = await randomBeacon
    .connect(await ethers.getSigner(signers.get(submitterIndex)))
    .submitDkgResult(dkgResult)

  return { transaction, dkgResult, dkgResultHash, members }
}

async function signDkgResult(
  signers: DkgGroupSigners,
  groupPublicKey: string,
  misbehaved: string,
  startBlock: number
) {
  const resultHash = ethers.utils.solidityKeccak256(
    ["bytes", "bytes", "uint256"],
    [groupPublicKey, misbehaved, startBlock]
  )

  const members: string[] = []
  const signingMemberIndices: number[] = []
  const signatures: string[] = []

  // eslint-disable-next-line no-restricted-syntax
  for (const [memberIndex, signer] of signers) {
    members.push(signer)

    signingMemberIndices.push(memberIndex)

    const ethersSigner = await ethers.getSigner(signer)

    const signature = await ethersSigner.signMessage(
      ethers.utils.arrayify(resultHash)
    )

    signatures.push(signature)
  }

  const signaturesBytes: string = ethers.utils.hexConcat(signatures)

  return { members, signingMemberIndices, signaturesBytes }
}
