// Script that verifies if a claim proof of a specified account is valid

const { MerkleTree } = require("merkletreejs")
const BigNumber = require("bignumber.js")
const keccak256 = require("keccak256")
const fs = require("fs")

// Merkle Distribution JSON file location
const PROOF_PATH = "distributions/2022-07-15/MerkleDist.json"

function verifyProof(wallet, beneficiary, amount, proof, root) {
  amount = BigNumber(amount)
  const tree = new MerkleTree([], keccak256, { sortPairs: true })
  const element =
    wallet + beneficiary.substr(2) + amount.toString(16).padStart(64, "0")
  const node = MerkleTree.bufferToHex(keccak256(element))
  return tree.verify(proof, node, root)
}

function main() {
  let proofVerification = true

  const json = JSON.parse(fs.readFileSync(PROOF_PATH, { encoding: "utf8" }))
  if (typeof json !== "object") throw new Error("Invalid JSON")

  const merkleRoot = json.merkleRoot
  const claims = json.claims
  Object.keys(claims).forEach((stakingProvider) => {
    const beneficiary = claims[stakingProvider].beneficiary
    const amount = claims[stakingProvider].amount
    const proof = claims[stakingProvider].proof
    const proofResult = verifyProof(
      stakingProvider,
      beneficiary,
      amount,
      proof,
      merkleRoot
    )

    if (!proofResult) {
      proofVerification = false
    }
  })

  console.log("Proof result: ", proofVerification)
}

main()
