import { Contract } from "ethers"
import { ethers } from "hardhat"
import type {
  SortitionPoolStub,
  RandomBeacon,
  RandomBeaconGovernance
} from "../../typechain"

export const constants = {
  groupSize: 3,
  signatureThreshold: 2,
  timeDKG: 13,
  dkgSubmissionEligibilityDelay: 10
}

// TODO: We should consider using hardhat-deploy plugin for contracts deployment.

interface DeployedContracts {
  [key: string]: Contract
}

export async function blsDeployment(): Promise<DeployedContracts> {
  const BLS = await ethers.getContractFactory("BLS")
  const bls = await BLS.deploy()
  await bls.deployed()

  const contracts: DeployedContracts = { bls }

  return contracts
}

export async function randomBeaconDeployment(): Promise<DeployedContracts> {
  const SortitionPoolStub = await ethers.getContractFactory("SortitionPoolStub")
  const sortitionPoolStub: SortitionPoolStub = await SortitionPoolStub.deploy()
  await sortitionPoolStub.deployed()

  const RandomBeacon = await ethers.getContractFactory("RandomBeacon")
  const randomBeacon: RandomBeacon = await RandomBeacon.deploy(
    sortitionPoolStub.address
  )
  await randomBeacon.deployed()

  const contracts: DeployedContracts = { sortitionPoolStub, randomBeacon }

  return contracts
}

export async function testDeployment(): Promise<DeployedContracts> {
  const contracts = await randomBeaconDeployment()

  const RandomBeaconGovernance = await ethers.getContractFactory(
    "RandomBeaconGovernance"
  )
  const randomBeaconGovernance: RandomBeaconGovernance = await RandomBeaconGovernance.deploy(
    contracts.randomBeacon.address
  )
  await randomBeaconGovernance.deployed()
  await contracts.randomBeacon.transferOwnership(randomBeaconGovernance.address)

  const newContracts = { randomBeaconGovernance }

  return { ...contracts, ...newContracts }
}
