import { ethers } from "hardhat"
import { Operator } from "./operators"

// eslint-disable-next-line import/prefer-default-export
export async function signHeartbeatFailureClaim(
  signers: Operator[],
  nonce: number,
  groupPubKey: string,
  failedMembersIndices: number[],
  numberOfSignatures: number
): Promise<{
  signatures: string
  signingMembersIndices: number[]
}> {
  const messageHash = ethers.utils.solidityKeccak256(
    ["uint256", "bytes", "uint8[]"],
    [nonce, groupPubKey, failedMembersIndices]
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
    const ethersSigner = await ethers.getSigner(signers[i].address)
    // eslint-disable-next-line no-await-in-loop
    const signature = await ethersSigner.signMessage(
      ethers.utils.arrayify(messageHash)
    )

    signatures.push(signature)
  }

  return {
    signatures: ethers.utils.hexConcat(signatures),
    signingMembersIndices,
  }
}
