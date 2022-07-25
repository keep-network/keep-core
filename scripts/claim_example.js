// Script shown as example about how to claim tokens

const { ethers } = require("hardhat")
const { BN } = require("@openzeppelin/test-helpers")
const fs = require("fs")

// input JSON file location with merkle root and claims data
const PROOF_PATH = "scripts/examples/example_proof_generated_10.json"
const DIST_PATH = "scripts/examples/example_dist_generated_10.json"

async function main() {
  const [owner, rewardsHolder] = await ethers.getSigners()

  // Read Merkle distribution data
  const proofJson = JSON.parse(
    fs.readFileSync(PROOF_PATH, { encoding: "utf8" })
  )
  if (typeof proofJson !== "object") throw new Error("Invalid JSON")
  const distJson = JSON.parse(fs.readFileSync(DIST_PATH, { encoding: "utf8" }))
  if (typeof distJson !== "object") throw new Error("Invalid JSON")

  const merkleRoot = proofJson.merkleRoot
  const totalAmount = new BN(distJson.totalAmount.slice(2), 16)
  const account = Object.keys(proofJson.claims)[0]
  const beneficiary = proofJson.claims[account].beneficiary
  const amount = proofJson.claims[account].amount
  const merkleProof = proofJson.claims[account].proof

  // Deploy token and cumulative merkle contracts
  const Token = await ethers.getContractFactory("TokenMock")
  const token = await Token.deploy()

  await token.mint(rewardsHolder.address, totalAmount.toString())
  console.log(
    "Mint:",
    (await token.balanceOf(rewardsHolder.address)).toString(),
    "tokens minted"
  )

  // Deploy cumulative merkle contract
  const CumulativeMerkle = await ethers.getContractFactory(
    "CumulativeMerkleDrop"
  )
  const cumulativeMerkle = await CumulativeMerkle.deploy(
    token.address,
    rewardsHolder.address,
    owner.address
  )

  // Set Merkle Root on contract
  await cumulativeMerkle.setMerkleRoot(merkleRoot)
  merkleRootContract = await cumulativeMerkle.merkleRoot()
  console.log("Merkle Root in contract:", merkleRootContract)

  await token
    .connect(rewardsHolder)
    .approve(cumulativeMerkle.address, totalAmount.toNumber())

  console.log(
    "Balance of beneficiary before claim:",
    (await token.balanceOf(beneficiary)).toString()
  )

  await cumulativeMerkle.claim(
    account,
    beneficiary,
    amount,
    merkleRoot,
    merkleProof
  )
  console.log(
    "Balance of account after claim:",
    (await token.balanceOf(account)).toString()
  )
}

main()
