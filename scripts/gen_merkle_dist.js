// Script that generates a new Merkle Distribution and outputs the data to a JSON file

const { MerkleTree } = require("merkletreejs")
const { BN } = require("@openzeppelin/test-helpers")
const keccak256 = require("keccak256")
const fs = require("fs")

// Input JSON file location with account addresses and balances
const ACC_PATH = "scripts/examples/example_stake_list_10.json"
// Output JSON file location of generated merkle distribution data
const GEN_DIST_PATH = "scripts/examples/example_dist_generated_10.json"
// Generated merkle proof data: first element of claims
const GEN_PROOF_PATH = "scripts/examples/example_proof_generated_10.json"

function toBN(x) {
  return new BN(x)
}

function makeDist(stakingProviders, data) {
  const elements = stakingProviders.map(
    (w, i) => w + data[i].beneficiary.substr(2) + toBN(data[i].amount).toString(16, 64)
  )
  const tree = new MerkleTree(elements, keccak256, {
    hashLeaves: true,
    sort: true,
  })
  const root = tree.getHexRoot()
  const leaves = tree.getHexLeaves()
  const proofs = leaves.map(tree.getHexProof, tree)

  return { leaves, root, proofs }
}

function main() {
  const json = JSON.parse(fs.readFileSync(ACC_PATH, { encoding: "utf8" }))
  if (typeof json !== "object") throw new Error("Invalid JSON")

  const dist = makeDist(Object.keys(json), Object.values(json))

  const totalAmount = Object.values(json)
    .map((data) => toBN(data.amount))
    .reduce((a, b) => a.add(b))
    .toString()

  const claims = Object.entries(json).map(([stakingProvider, data]) => {
    leaf = MerkleTree.bufferToHex(
      keccak256(stakingProvider + data.beneficiary.substr(2) + toBN(data.amount).toString(16, 64))
    )
    return {
      stakingProvider: stakingProvider,
      beneficiary: data.beneficiary,
      amount: data.amount,
      proof: dist.proofs[dist.leaves.indexOf(leaf)],
    }
  })

  const dist_json = {
    totalAmount: totalAmount,
    merkleRoot: dist.root,
    claims: claims.reduce((a, { stakingProvider, beneficiary, amount, proof }) => {
      a[stakingProvider] = { beneficiary, amount, proof }
      return a
    }, {}),
  }

  // Take the first element of generated distribution
  const claimAccount = Object.keys(dist_json.claims)[0]
  const proof_json = {
    merkleRoot: dist_json.merkleRoot,
    claims: {
      [claimAccount]: dist_json.claims[claimAccount],
    },
  }

  fs.writeFileSync(GEN_DIST_PATH, JSON.stringify(dist_json, null, 2))
  fs.writeFileSync(GEN_PROOF_PATH, JSON.stringify(proof_json, null, 2))
}

main()
