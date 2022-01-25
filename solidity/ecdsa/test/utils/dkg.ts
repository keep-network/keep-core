/* eslint-disable no-await-in-loop */

import { ethers } from "hardhat"
import { expect } from "chai"
import { BigNumber } from "ethers"

import { constants } from "../fixtures"

import { selectGroup } from "./groups"
import { firstEligibleIndex } from "./submission"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { BigNumberish, ContractTransaction } from "ethers"
import type { SortitionPool, WalletRegistry } from "../../typechain"
import type { Operator } from "./operators"
import type { DkgResultSubmittedEvent, ResultStruct } from "../../typechain/DKG"

export interface DkgResult {
  submitterMemberIndex: number
  groupPubKey: string
  misbehavedMembersIndices: number[]
  signatures: string
  signingMembersIndices: number[]
  members: number[]
  membersHash: string
}

export const noMisbehaved: number[] = []

export function calculateDkgSeed(
  relayEntry: BigNumberish,
  blockNumber: BigNumberish
): BigNumber {
  return ethers.BigNumber.from(
    ethers.utils.keccak256(
      ethers.utils.solidityPack(
        ["uint256", "uint256"],
        [ethers.BigNumber.from(relayEntry), ethers.BigNumber.from(blockNumber)]
      )
    )
  )
}

// Sign and submit a correct DKG result which cannot be challenged because used
// signers belong to an actual group selected by the sortition pool for given
// seed.
export async function signAndSubmitCorrectDkgResult(
  walletRegistry: WalletRegistry,
  groupPublicKey: string,
  seed: BigNumber,
  startBlock: number,
  misbehavedIndices: number[],
  submitterIndex?: number,
  numberOfSignatures = 51
): Promise<{
  signers: Operator[]
  dkgResult: DkgResult
  dkgResultHash: string
  members: number[]
  submitter: SignerWithAddress
  transaction: ContractTransaction
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
    await walletRegistry.sortitionPool()
  )) as SortitionPool

  const signers = await selectGroup(sortitionPool, seed)

  return {
    signers,
    ...(await signAndSubmitArbitraryDkgResult(
      walletRegistry,
      groupPublicKey,
      signers,
      startBlock,
      misbehavedIndices,
      submitterIndex,
      numberOfSignatures
    )),
  }
}

const DKG_RESULT_PARAMS_SIGNATURE =
  "(uint256 submitterMemberIndex, bytes groupPubKey, uint8[] misbehavedMembersIndices, bytes signatures, uint256[] signingMembersIndices, uint32[] members, bytes32 membersHash)"

// Sign and submit an arbitrary DKG result using given signers. Signers don't
// need to be part of the actual sortition pool group. This function is useful
// for preparing invalid or malicious results for testing purposes.
export async function signAndSubmitArbitraryDkgResult(
  walletRegistry: WalletRegistry,
  groupPublicKey: string,
  signers: Operator[],
  startBlock: number,
  misbehavedIndices: number[],
  submitterIndex?: number,
  numberOfSignatures = 51
): Promise<{
  dkgResult: DkgResult
  dkgResultHash: string
  members: number[]
  submitter: SignerWithAddress
  transaction: ContractTransaction
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

  const dkgResult: DkgResult = {
    submitterMemberIndex: submitterIndex,
    groupPubKey: groupPublicKey,
    misbehavedMembersIndices: misbehavedIndices,
    signatures: signaturesBytes,
    signingMembersIndices,
    members,
    membersHash: hashDKGMembers(members, misbehavedIndices),
  }

  const dkgResultHash = ethers.utils.keccak256(
    ethers.utils.defaultAbiCoder.encode(
      [DKG_RESULT_PARAMS_SIGNATURE],
      [dkgResult]
    )
  )

  const submitter = await ethers.getSigner(signers[submitterIndex - 1].address)

  return {
    dkgResult,
    dkgResultHash,
    members,
    submitter,
    ...(await submitDkgResult(walletRegistry, dkgResult, submitter)),
  }
}

// Signs and submits a DKG result containing signatures with random bytes.
// Attempting to recover addresses from such signatures causes a revert. It is
// useful for preparing malicious DKG results.
export async function signAndSubmitUnrecoverableDkgResult(
  walletRegistry: WalletRegistry,
  groupPublicKey: string,
  signers: Operator[],
  startBlock: number,
  misbehavedIndices: number[],
  submitterIndex?: number,
  numberOfSignatures = 51
): Promise<{
  dkgResult: DkgResult
  dkgResultHash: string
  members: number[]
  submitter: SignerWithAddress
  transaction: ContractTransaction
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

  const dkgResult: DkgResult = {
    submitterMemberIndex: submitterIndex,
    groupPubKey: groupPublicKey,
    misbehavedMembersIndices: misbehavedIndices,
    signatures: unrecoverableSignatures,
    signingMembersIndices,
    members,
    membersHash: hashDKGMembers(members, misbehavedIndices),
  }

  const dkgResultHash = ethers.utils.keccak256(
    ethers.utils.defaultAbiCoder.encode(
      [DKG_RESULT_PARAMS_SIGNATURE],
      [dkgResult]
    )
  )

  const submitter = await ethers.getSigner(signers[submitterIndex - 1].address)

  return {
    dkgResult,
    dkgResultHash,
    members,
    submitter,
    ...(await submitDkgResult(walletRegistry, dkgResult, submitter)),
  }
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

export async function submitDkgResult(
  walletRegistry: WalletRegistry,
  dkgResult: DkgResult,
  submitter: SignerWithAddress
): Promise<{
  transaction: ContractTransaction
}> {
  const transaction = await walletRegistry
    .connect(submitter)
    .submitDkgResult(dkgResult)

  return { transaction }
}

// Creates a members hash that actively participated in dkg
export function hashDKGMembers(
  members: number[],
  misbehavedMembersIndices?: number[]
): string {
  if (misbehavedMembersIndices && misbehavedMembersIndices.length > 0) {
    const activeDkgMembers = [...members]
    for (let i = 0; i < misbehavedMembersIndices.length; i++) {
      if (misbehavedMembersIndices[i] !== 0) {
        activeDkgMembers.splice(misbehavedMembersIndices[i] - i - 1, 1)
      }
    }

    return ethers.utils.keccak256(
      ethers.utils.defaultAbiCoder.encode(["uint32[]"], [activeDkgMembers])
    )
  }

  return ethers.utils.keccak256(
    ethers.utils.defaultAbiCoder.encode(["uint32[]"], [members])
  )
}

interface DkgResultSubmittedEventArgs {
  resultHash: string
  seed: BigNumber
  result: ResultStruct
}

// This is a workaround for a bug in ethereum-waffle library that doesn't let
// verify events that have an array nested in a struct.
// See: https://github.com/EthWorks/Waffle/issues/245
export async function expectDkgResultSubmittedEvent(
  tx: ContractTransaction,
  expectedArgs: DkgResultSubmittedEventArgs
): Promise<void> {
  const eventName = "DkgResultSubmitted"

  const event: DkgResultSubmittedEvent = (await tx.wait()).events.find(
    (e) => e.event === eventName
  ) as unknown as DkgResultSubmittedEvent

  await expect(event, `Event ${eventName} not emitted`).to.be.not.null

  const actualArgs = event.args

  await expect(actualArgs.length, "invalid event args length").to.be.equal(
    Object.keys(expectedArgs).length
  )

  await expect(
    actualArgs.result.length,
    "invalid result args length"
  ).to.be.equal(Object.keys(expectedArgs.result).length)

  await expect(actualArgs.resultHash, "invalid resultHash").to.be.equal(
    expectedArgs.resultHash
  )

  await expect(actualArgs.resultHash, "invalid resultHash").to.be.equal(
    expectedArgs.resultHash
  )

  await expect(actualArgs.seed, "invalid seed").to.be.equal(expectedArgs.seed)

  await expect(
    actualArgs.result.submitterMemberIndex,
    "invalid submitterMemberIndex"
  ).to.be.equal(expectedArgs.result.submitterMemberIndex)

  await expect(
    actualArgs.result.groupPubKey,
    "invalid groupPubKey"
  ).to.be.equal(expectedArgs.result.groupPubKey)

  await expect(
    actualArgs.result.misbehavedMembersIndices,
    "invalid misbehavedMembersIndices"
  ).to.be.deep.equal(expectedArgs.result.misbehavedMembersIndices)

  await expect(actualArgs.result.signatures, "invalid signatures").to.be.equal(
    expectedArgs.result.signatures
  )

  await expect(
    actualArgs.result.signingMembersIndices,
    "invalid signingMembersIndices"
  ).to.be.deep.equal(
    expectedArgs.result.signingMembersIndices.map(BigNumber.from)
  )

  await expect(actualArgs.result.members, "invalid members").to.be.deep.equal(
    expectedArgs.result.members
  )

  await expect(
    actualArgs.result.membersHash,
    "invalid membersHash"
  ).to.be.equal(expectedArgs.result.membersHash)
}
