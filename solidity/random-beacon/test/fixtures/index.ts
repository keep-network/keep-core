import { Contract } from "ethers"
import { ethers } from "hardhat"
import type {
  SortitionPoolStub,
  RandomBeaconStub,
  RandomBeaconGovernance,
} from "../../typechain"

export const constants = {
  groupSize: 64,
  groupThreshold: 33,
  signatureThreshold: 48, // groupThreshold + (groupSize - groupThreshold) / 2
  offchainDkgTime: 72, // 5 * (1 + 5) + 2 * (1 + 10) + 20
}

export const params = {
  relayRequestFee: 0,
  relayEntrySubmissionEligibilityDelay: 10,
  relayEntryHardTimeout: 5760,
  callbackGasLimit: 200000,
  groupCreationFrequency: 10,
  groupLifeTime: 60 * 60 * 24 * 14, // 2 weeks
  dkgResultChallengePeriodLength: 1440,
  dkgResultSubmissionEligibilityDelay: 10,
  dkgResultSubmissionReward: 0,
  sortitionPoolUnlockingReward: 0,
  relayEntrySubmissionFailureSlashingAmount: ethers.BigNumber.from(10)
    .pow(18)
    .mul(1000),
  maliciousDkgResultSlashingAmount: ethers.BigNumber.from(10)
    .pow(18)
    .mul(50000),
}

// TODO: We should consider using hardhat-deploy plugin for contracts deployment.

export interface DeployedContracts {
  [key: string]: Contract
}

export async function blsDeployment(): Promise<DeployedContracts> {
  const BLS = await ethers.getContractFactory("BLS")
  const bls = await BLS.deploy()
  await bls.deployed()

  const contracts: DeployedContracts = { bls }

  return contracts
}

export async function testTokenDeployment(): Promise<DeployedContracts> {
  const TestToken = await ethers.getContractFactory("TestToken")
  const testToken = await TestToken.deploy()
  await testToken.deployed()

  const contracts: DeployedContracts = { testToken }

  return contracts
}

export async function randomBeaconDeployment(): Promise<DeployedContracts> {
  const SortitionPoolStub = await ethers.getContractFactory("SortitionPoolStub")
  const sortitionPoolStub: SortitionPoolStub = await SortitionPoolStub.deploy()

  const { testToken } = await testTokenDeployment()

  const RandomBeacon = await ethers.getContractFactory("RandomBeaconStub", {
    libraries: {
      BLS: (await blsDeployment()).bls.address,
    },
  })

  const randomBeacon: RandomBeaconStub = await RandomBeacon.deploy(
    sortitionPoolStub.address,
    testToken.address
  )
  await randomBeacon.deployed()

  const contracts: DeployedContracts = {
    sortitionPoolStub,
    randomBeacon,
    testToken,
  }

  return contracts
}

export async function testDeployment(): Promise<DeployedContracts> {
  const contracts = await randomBeaconDeployment()

  const RandomBeaconGovernance = await ethers.getContractFactory(
    "RandomBeaconGovernance"
  )
  const randomBeaconGovernance: RandomBeaconGovernance =
    await RandomBeaconGovernance.deploy(contracts.randomBeacon.address)
  await randomBeaconGovernance.deployed()
  await contracts.randomBeacon.transferOwnership(randomBeaconGovernance.address)

  const newContracts = { randomBeaconGovernance }

  return { ...contracts, ...newContracts }
}
