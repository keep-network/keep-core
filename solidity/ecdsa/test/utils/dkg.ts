/* eslint-disable no-await-in-loop */

import { ethers } from "hardhat"

import { constants } from "../fixtures"

import { selectGroup } from "./groups"
import { firstEligibleIndex } from "./submission"

import type { BigNumber, BigNumberish, ContractTransaction } from "ethers"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { SortitionPool, WalletFactory } from "../../typechain"
import type { Operator } from "./operators"

export interface DkgResult {
  submitterMemberIndex: number
  groupPubKey: string
  misbehavedMembersIndices: number[]
  signatures: string
  signingMembersIndices: number[]
  members: number[]
}

export const noMisbehaved = []

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
  walletFactory: WalletFactory,
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
  walletAddress: string
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
    await walletFactory.sortitionPool()
  )) as SortitionPool

  const signers = await selectGroup(sortitionPool, seed)

  return {
    signers,
    ...(await signAndSubmitArbitraryDkgResult(
      walletFactory,
      groupPublicKey,
      signers,
      startBlock,
      misbehavedIndices,
      submitterIndex,
      numberOfSignatures
    )),
  }
}

// Sign and submit an arbitrary DKG result using given signers. Signers don't
// need to be part of the actual sortition pool group. This function is useful
// for preparing invalid or malicious results for testing purposes.
export async function signAndSubmitArbitraryDkgResult(
  walletFactory: WalletFactory,
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
  walletAddress: string
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
  }

  const dkgResultHash = ethers.utils.keccak256(
    ethers.utils.defaultAbiCoder.encode(
      [
        "(uint256 submitterMemberIndex, bytes groupPubKey, uint8[] misbehavedMembersIndices, bytes signatures, uint256[] signingMembersIndices, uint32[] members)",
      ],
      [dkgResult]
    )
  )

  const submitter = await ethers.getSigner(signers[submitterIndex - 1].address)

  return {
    dkgResult,
    dkgResultHash,
    members,
    submitter,
    ...(await submitDkgResult(walletFactory, dkgResult, submitter)),
  }
}

// Signs and submits a DKG result containing signatures with random bytes.
// Attempting to recover addresses from such signatures causes a revert. It is
// useful for preparing malicious DKG results.
export async function signAndSubmitUnrecoverableDkgResult(
  walletFactory: WalletFactory,
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
  walletAddress: string
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
  }

  const dkgResultHash = ethers.utils.keccak256(
    ethers.utils.defaultAbiCoder.encode(
      [
        "(uint256 submitterMemberIndex, bytes groupPubKey, uint8[] misbehavedMembersIndices, bytes signatures, uint256[] signingMembersIndices, uint32[] members)",
      ],
      [dkgResult]
    )
  )

  const submitter = await ethers.getSigner(signers[submitterIndex - 1].address)

  return {
    dkgResult,
    dkgResultHash,
    members,
    submitter,
    ...(await submitDkgResult(walletFactory, dkgResult, submitter)),
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
  walletFactory: WalletFactory,
  dkgResult: DkgResult,
  submitter: SignerWithAddress
): Promise<{
  transaction: ContractTransaction
  walletAddress: string
}> {
  const transaction = await walletFactory
    .connect(submitter)
    .submitDkgResult(dkgResult)

  const { walletAddress } = (await transaction.wait()).events.find(
    (e) => e.event === "WalletCreated"
  ).args

  return { transaction, walletAddress }
}
