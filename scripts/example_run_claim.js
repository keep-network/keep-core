// Script shown as example about how to claim tokens

const { ethers } = require('hardhat')
const { BN } = require('@openzeppelin/test-helpers')
const fs = require('fs')

// input JSON file location with merkle root and claims data
const PROOF_PATH = 'scripts/examples/example_proof_generated.json'
const DIST_PATH = 'scripts/examples/example_dist_generated.json'

async function main () {
  // Deploy token and cumulative merkle contracts
  const Token = await ethers.getContractFactory('TokenMock')
  const token = await Token.deploy('token', 't')
  const CumulativeMerkle = await ethers.getContractFactory('CumulativeMerkleDrop')
  const cumulativeMerkle = await CumulativeMerkle.deploy(token.address)

  const proofJson = JSON.parse(fs.readFileSync(PROOF_PATH, { encoding: 'utf8' }))
  if (typeof proofJson !== 'object') throw new Error('Invalid JSON')
  const distJson = JSON.parse(fs.readFileSync(DIST_PATH, { encoding: 'utf8' }))
  if (typeof distJson !== 'object') throw new Error('Invalid JSON')

  const merkleRoot = proofJson.merkleRoot
  const totalAmount = new BN(distJson.totalAmount.slice(2), 16)
  const account = Object.keys(proofJson.claims)[0]
  const amount = proofJson.claims[account].amount
  const merkleProof = proofJson.claims[account].proof

  // Set Merkle Root on contract
  await cumulativeMerkle.setMerkleRoot(merkleRoot)

  merkleRootContract = await cumulativeMerkle.merkleRoot()
  console.log('Merkle Root in contract:', merkleRootContract)

  await token.mint(cumulativeMerkle.address, totalAmount.toString())
  console.log('Mint:', (await token.balanceOf(cumulativeMerkle.address)).toString(), 'tokens minted')

  console.log('Balance of account before claim:', (await token.balanceOf(account)).toString())
  await cumulativeMerkle.claim(account, amount, merkleRoot, merkleProof)
  console.log('Balance of account after claim:', (await token.balanceOf(account)).toString())
}

main()
