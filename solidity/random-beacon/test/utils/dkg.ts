/* eslint-disable no-await-in-loop */

import { ethers } from "hardhat"
import type { BigNumber, ContractTransaction } from "ethers"
import blsData from "../data/bls"
import type { RandomBeacon } from "../../typechain"
import { Operator } from "./sortitionpool"

export interface DkgResult {
  submitterMemberIndex: number
  groupPubKey: string
  misbehavedMembersIndices: number[]
  signatures: string
  signingMembersIndices: number[]
  members: number[]
}

export const noMisbehaved = []

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
  signers: Operator[],
  startBlock: number,
  submitterIndex = 1
): Promise<{
  transaction: ContractTransaction
  dkgResult: DkgResult
  dkgResultHash: string
  members: number[]
}> {
  const { members, signingMembersIndices, signaturesBytes } =
    await signDkgResult(signers, groupPublicKey, noMisbehaved, startBlock)

  const dkgResult: DkgResult = {
    submitterMemberIndex: submitterIndex,
    groupPubKey: blsData.groupPubKey,
    misbehavedMembersIndices: noMisbehaved,
    signatures: signaturesBytes,
    signingMembersIndices,
    members,
  }

  const dkgResultHash = ethers.utils.keccak256(
    ethers.utils.defaultAbiCoder.encode(
      [
        "(uint256 submitterMemberIndex, bytes groupPubKey, uint8[] misbehavedMembersIndices, bytes signatures, uint256[] signingMembersIndices, uint32[] members)",
      ],
      [dkgResult]
    )
  )

  const transaction = await randomBeacon
    .connect(await ethers.getSigner(signers[submitterIndex - 1].address))
    .submitDkgResult(dkgResult)

  return { transaction, dkgResult, dkgResultHash, members }
}

async function signDkgResult(
  signers: Operator[],
  groupPublicKey: string,
  misbehaved: number[],
  startBlock: number
): Promise<{
  members: number[]
  signingMembersIndices: number[]
  signaturesBytes: string
}> {
  const resultHash = ethers.utils.solidityKeccak256(
    ["bytes", "uint8[]", "uint256"],
    [groupPublicKey, misbehaved, startBlock]
  )

  const members: number[] = []
  const signingMembersIndices: number[] = []
  const signatures: string[] = []

  for (let i = 0; i < signers.length; i++) {
    const { id, address } = signers[i]
    const signerIndex: number = i + 1

    members.push(id)
    signingMembersIndices.push(signerIndex)

    const ethersSigner = await ethers.getSigner(address)
    const signature = await ethersSigner.signMessage(
      ethers.utils.arrayify(resultHash)
    )

    signatures.push(signature)
  }

  const signaturesBytes: string = ethers.utils.hexConcat(signatures)

  return { members, signingMembersIndices, signaturesBytes }
}
