import { Contract } from "ethers"
import { ethers } from "hardhat"

interface DeployedContracts {
  [key: string]: Contract
}

export async function testModUtilsDeployment(): Promise<DeployedContracts> {
  const TestModUtils = await ethers.getContractFactory("TestModUtils")
  const testModUtils = await TestModUtils.deploy()
  await testModUtils.deployed()

  const contracts: DeployedContracts = { testModUtils }

  return contracts
}

export async function testAltBn128Deployment(): Promise<DeployedContracts> {
  const TestAltBn128 = await ethers.getContractFactory("TestAltBn128")
  const testAltBn128 = await TestAltBn128.deploy()
  await testAltBn128.deployed()

  const contracts: DeployedContracts = { testAltBn128 }

  return contracts
}

export async function blsDeployment(): Promise<DeployedContracts> {
  const BLS = await ethers.getContractFactory("BLS")
  const bls = await BLS.deploy()
  await bls.deployed()

  const contracts: DeployedContracts = { bls }

  return contracts
}
