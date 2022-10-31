const { MerkleTree } = require("merkletreejs")
const BigNumber = require("bignumber.js")
const keccak256 = require("keccak256")

/**
 * Generate a Merkle distribution from a Threshold rewards input
 * @param {Object} merkleInput      Merkle input generated from rewards
 * @return {Object}                 Merkle distribution
 */
exports.genMerkleDist = function (merkleInput) {
  const stakingProviders = Object.keys(merkleInput)
  const data = Object.values(merkleInput)

  const elements = stakingProviders.map(
    (stakingProvider, i) =>
      stakingProvider +
      data[i].beneficiary.substr(2) +
      BigNumber(data[i].amount).toString(16).padStart(64, "0")
  )

  const tree = new MerkleTree(elements, keccak256, {
    hashLeaves: true,
    sort: true,
  })

  const root = tree.getHexRoot()
  const leaves = tree.getHexLeaves()
  const proofs = leaves.map(tree.getHexProof, tree)

  const totalAmount = data
    .map((claim) => BigNumber(claim.amount))
    .reduce((a, b) => a.plus(b))
    .toFixed()

  const claims = Object.entries(merkleInput).map(([stakingProvider, data]) => {
    const leaf = MerkleTree.bufferToHex(
      keccak256(
        stakingProvider +
          data.beneficiary.substr(2) +
          BigNumber(data.amount).toString(16).padStart(64, "0")
      )
    )
    return {
      stakingProvider: stakingProvider,
      beneficiary: data.beneficiary,
      amount: data.amount,
      proof: proofs[leaves.indexOf(leaf)],
    }
  })

  const dist = {
    totalAmount: totalAmount,
    merkleRoot: root,
    claims: claims.reduce(
      (a, { stakingProvider, beneficiary, amount, proof }) => {
        a[stakingProvider] = { beneficiary, amount, proof }
        return a
      },
      {}
    ),
  }

  return dist
}

/**
 * Combine two Threshold rewards inputs, adding the amounts and taking the
 * beneficiary of the second input
 * @param {Object} baseMerkleInput  Merkle input used as base
 * @param {Object} addedMerkleInput Merkle input to be added
 * @return {Object}                 Combination of two Merkle inputs
 */
exports.combineMerkleInputs = function (baseMerkleInput, addedMerkleInput) {
  const combined = JSON.parse(JSON.stringify(baseMerkleInput))
  Object.keys(addedMerkleInput).map((stakingProvider) => {
    const combinedClaim = combined[stakingProvider]
    const addedClaim = addedMerkleInput[stakingProvider]
    if (combinedClaim) {
      combinedClaim.beneficiary = addedClaim.beneficiary
      combinedClaim.amount = BigNumber(combinedClaim.amount)
        .plus(BigNumber(addedClaim.amount))
        .toFixed()
    } else {
      combined[stakingProvider] = {
        beneficiary: addedClaim.beneficiary,
        amount: addedClaim.amount,
      }
    }
  })
  return combined
}
