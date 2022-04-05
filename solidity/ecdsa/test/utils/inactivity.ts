import { ethers } from "hardhat"

import type { Operator } from "./operators"

// eslint-disable-next-line import/prefer-default-export
export async function signOperatorInactivityClaim(
  signers: Operator[],
  nonce: number,
  groupPubKey: string,
  failedHeartbeat: boolean,
  inactiveMembersIndices: number[],
  numberOfSignatures: number
): Promise<{
  signatures: string
  signingMembersIndices: number[]
}> {
  const messageHash = ethers.utils.solidityKeccak256(
    ["uint256", "bytes", "uint8[]", "bool"],
    [nonce, groupPubKey, inactiveMembersIndices, failedHeartbeat]
  )

  const signingMembersIndices: number[] = []
  const signatures: string[] = []

  for (let i = 0; i < signers.length; i++) {
    if (signatures.length === numberOfSignatures) {
      // eslint-disable-next-line no-continue
      continue
    }

    const signerIndex: number = i + 1

    signingMembersIndices.push(signerIndex)

    // eslint-disable-next-line no-await-in-loop
    const signature = await signers[i].signer.signMessage(
      ethers.utils.arrayify(messageHash)
    )

    signatures.push(signature)
  }

  return {
    signatures: ethers.utils.hexConcat(signatures),
    signingMembersIndices,
  }
}
