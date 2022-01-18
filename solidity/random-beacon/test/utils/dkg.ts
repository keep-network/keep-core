/* eslint-disable no-await-in-loop */

import { ethers } from "hardhat"
import type { BigNumber, ContractTransaction } from "ethers"
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { RandomBeacon, SortitionPool } from "../../typechain"
import { Operator } from "./operators"
// eslint-disable-next-line import/no-cycle
import { selectGroup } from "./groups"
import { firstEligibleIndex } from "./submission"
import { constants } from "../fixtures"

export interface DkgResult {
  submitterMemberIndex: number
  groupPubKey: string
  misbehavedMembersIndices: number[]
  signatures: string
  signingMembersIndices: number[]
  members: number[]
  membersHash: string
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

// Sign and submit a correct DKG result which cannot be challenged because used
// signers belong to an actual group selected by the sortition pool for given
// seed.
export async function signAndSubmitCorrectDkgResult(
  randomBeacon: RandomBeacon,
  groupPublicKey: string,
  seed: BigNumber,
  startBlock: number,
  misbehavedIndices: number[],
  submitterIndex?: number,
  membersHash?: string,
  numberOfSignatures = 33
): Promise<{
  transaction: ContractTransaction
  dkgResult: DkgResult
  dkgResultHash: string
  members: number[]
  submitter: SignerWithAddress
}> {
  if (!submitterIndex) {
    // eslint-disable-next-line no-param-reassign
    submitterIndex = firstEligibleIndex(
      ethers.utils.keccak256(groupPublicKey),
      constants.groupSize
    )
  }

  const sortitionPool = (await ethers.getContractAt(
    "SortitionPool",
    await randomBeacon.sortitionPool()
  )) as SortitionPool

  return signAndSubmitArbitraryDkgResult(
    randomBeacon,
    groupPublicKey,
    await selectGroup(sortitionPool, seed),
    startBlock,
    misbehavedIndices,
    submitterIndex,
    membersHash,
    numberOfSignatures
  )
}

// Sign and submit an arbitrary DKG result using given signers. Signers don't
// need to be part of the actual sortition pool group. This function is useful
// for preparing invalid or malicious results for testing purposes.
export async function signAndSubmitArbitraryDkgResult(
  randomBeacon: RandomBeacon,
  groupPublicKey: string,
  signers: Operator[],
  startBlock: number,
  misbehavedIndices: number[],
  submitterIndex?: number,
  groupMembersHash?: string,
  numberOfSignatures = 33
): Promise<{
  transaction: ContractTransaction
  dkgResult: DkgResult
  dkgResultHash: string
  members: number[]
  submitter: SignerWithAddress
}> {
  const { members, signingMembersIndices, signaturesBytes } =
    await signDkgResult(
      signers,
      groupPublicKey,
      misbehavedIndices,
      startBlock,
      numberOfSignatures
    )

  if (!submitterIndex) {
    // eslint-disable-next-line no-param-reassign
    submitterIndex = firstEligibleIndex(ethers.utils.keccak256(groupPublicKey))
  }
  let membersHash = groupMembersHash
  if (!membersHash) {
    membersHash = hashDKGMembers(members, misbehavedIndices)
  }

  const dkgResult: DkgResult = {
    submitterMemberIndex: submitterIndex,
    groupPubKey: groupPublicKey,
    misbehavedMembersIndices: misbehavedIndices,
    signatures: signaturesBytes,
    signingMembersIndices,
    members,
    membersHash,
  }

  const dkgResultHash = ethers.utils.keccak256(
    ethers.utils.defaultAbiCoder.encode(
      [
        "(uint256 submitterMemberIndex, bytes groupPubKey, uint8[] misbehavedMembersIndices, bytes signatures, uint256[] signingMembersIndices, uint32[] members, bytes32 membersHash)",
      ],
      [dkgResult]
    )
  )

  const submitter = await ethers.getSigner(signers[submitterIndex - 1].address)

  const transaction = await randomBeacon
    .connect(submitter)
    .submitDkgResult(dkgResult)

  return {
    transaction,
    dkgResult,
    dkgResultHash,
    members,
    submitter,
  }
}

// Signs and submits a DKG result containing signatures with random bytes.
// Attempting to recover addresses from such signatures causes a revert. It is
// useful for preparing malicious DKG results.
export async function signAndSubmitUnrecoverableDkgResult(
  randomBeacon: RandomBeacon,
  groupPublicKey: string,
  signers: Operator[],
  startBlock: number,
  misbehavedIndices: number[],
  submitterIndex?: number,
  numberOfSignatures = 33
): Promise<{
  transaction: ContractTransaction
  dkgResult: DkgResult
  dkgResultHash: string
  members: number[]
  submitter: SignerWithAddress
}> {
  const { members, signingMembersIndices } = await signDkgResult(
    signers,
    groupPublicKey,
    misbehavedIndices,
    startBlock,
    numberOfSignatures
  )

  if (!submitterIndex) {
    // eslint-disable-next-line no-param-reassign
    submitterIndex = firstEligibleIndex(ethers.utils.keccak256(groupPublicKey))
  }

  const signatureHexStrLength = 2 * 65
  const unrecoverableSignatures = `0x${"a".repeat(
    signatureHexStrLength * numberOfSignatures
  )}`

  const membersHash = hashDKGMembers(members, misbehavedIndices)

  const dkgResult: DkgResult = {
    submitterMemberIndex: submitterIndex,
    groupPubKey: groupPublicKey,
    misbehavedMembersIndices: misbehavedIndices,
    signatures: unrecoverableSignatures,
    signingMembersIndices,
    members,
    membersHash,
  }

  const dkgResultHash = ethers.utils.keccak256(
    ethers.utils.defaultAbiCoder.encode(
      [
        "(uint256 submitterMemberIndex, bytes groupPubKey, uint8[] misbehavedMembersIndices, bytes signatures, uint256[] signingMembersIndices, uint32[] members, bytes32 membersHash)",
      ],
      [dkgResult]
    )
  )

  const submitter = await ethers.getSigner(signers[submitterIndex - 1].address)

  const transaction = await randomBeacon
    .connect(submitter)
    .submitDkgResult(dkgResult)

  return { transaction, dkgResult, dkgResultHash, members, submitter }
}

export async function signDkgResult(
  signers: Operator[],
  groupPublicKey: string,
  misbehavedMembersIndices: number[],
  startBlock: number,
  numberOfSignatures: number
): Promise<{
  members: number[]
  signingMembersIndices: number[]
  signaturesBytes: string
}> {
  const resultHash = ethers.utils.solidityKeccak256(
    ["bytes", "uint8[]", "uint256"],
    [groupPublicKey, misbehavedMembersIndices, startBlock]
  )

  const members: number[] = []
  const signingMembersIndices: number[] = []
  const signatures: string[] = []
  for (let i = 0; i < signers.length; i++) {
    const { id, address } = signers[i]
    members.push(id)

    if (signatures.length === numberOfSignatures) {
      // eslint-disable-next-line no-continue
      continue
    }

    const signerIndex: number = i + 1

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

// Creates a members hash that actively participated in dkg
export function hashDKGMembers(
  members: number[],
  misbehavedMembersIndices: number[]
): string {
  if (misbehavedMembersIndices.length > 0) {
    const activeDkgMembers = [...members]
    for (let i = 0; i < misbehavedMembersIndices.length; i++) {
      activeDkgMembers.splice(misbehavedMembersIndices[i] - i - 1, 1)
    }

    return ethers.utils.keccak256(
      ethers.utils.defaultAbiCoder.encode(["uint32[]"], [activeDkgMembers])
    )
  }

  return ethers.utils.keccak256(
    ethers.utils.defaultAbiCoder.encode(["uint32[]"], [members])
  )
}
