// Script that verifies if a claim proof of a specified account is valid
// Use: node src/scripts/verify_merkle_dist.js <YYYY-MM-DD>

const { MerkleTree } = require("merkletreejs")
const BigNumber = require("bignumber.js")
const keccak256 = require("keccak256")
const fs = require("fs")

function verifyProof(wallet, beneficiary, amount, proof, root) {
  amount = BigNumber(amount)
  const tree = new MerkleTree([], keccak256, { sortPairs: true })
  const element =
    wallet + beneficiary.substr(2) + amount.toString(16).padStart(64, "0")
  const node = MerkleTree.bufferToHex(keccak256(element))
  return tree.verify(proof, node, root)
}

function main() {
  const args = process.argv.slice(2)
  const distDate = args[0]
  const proofPath = `distributions/${distDate}/MerkleDist.json`

  let proofVerification = true

  const json = JSON.parse(fs.readFileSync(proofPath, { encoding: "utf8" }))
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
