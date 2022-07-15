// Script that verifies if a claim proof of a specified account is valid

const { MerkleTree } = require("merkletreejs")
const { BN } = require("@openzeppelin/test-helpers")
const keccak256 = require("keccak256")
const fs = require("fs")

// Merkle Distribution JSON file location
const PROOF_PATH = "scripts/examples/MerkleDist.json"

function verifyProof(wallet, beneficiary, amount, proof, root) {
  amount = new BN(amount)
  const tree = new MerkleTree([], keccak256, { sortPairs: true })
  const element = wallet + beneficiary.substr(2) + amount.toString(16, 64)
  const node = MerkleTree.bufferToHex(keccak256(element))
  return tree.verify(proof, node, root)
}

function main() {
  const json = JSON.parse(fs.readFileSync(PROOF_PATH, { encoding: "utf8" }))
  if (typeof json !== "object") throw new Error("Invalid JSON")

  const merkleRoot = json.merkleRoot
  const claims = json.claims
  Object.keys(claims).forEach( stakingProvider => {
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
    console.log("Proof result: ", proofResult)
  })

}

main()
