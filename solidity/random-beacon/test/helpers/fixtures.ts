import { Contract } from "ethers"
import { ethers, getNamedAccounts } from "hardhat"

import type { RandomBeacon, RandomBeaconGovernance } from "../../typechain"

export const constants = {
  groupSize: 64,
  signatureThreshold: 33,
  timeDKG: 5 * (1 + 5) + 2 * (1 + 10) + 20,
}

export const params = {
  relayRequestFee: 0,
  relayEntrySubmissionEligibilityDelay: 10,
  relayEntryHardTimeout: 5760,
  callbackGasLimit: 200000,
  groupCreationFrequency: 10,
  groupLifeTime: 60 * 60 * 24 * 14,
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

interface DeployedContracts {
  [key: string]: Contract
}

// TODO: Contract deployemnts should be replaced by hardhat-deploy plugin.

export async function randomBeaconDeployment(): Promise<DeployedContracts> {
  const deployer = await ethers.getSigner((await getNamedAccounts()).deployer)
  const governance = await ethers.getSigner(
    (await getNamedAccounts()).governance
  )

  const DKG = await ethers.getContractFactory("DKG")
  const dkg = await DKG.deploy()

  const Groups = await ethers.getContractFactory("Groups")
  const groups = await Groups.deploy()

  const { testToken } = await testTokenDeployment()
  const { maintenancePool } = await maintenancePoolDeployment(testToken)

  const RandomBeacon = await ethers.getContractFactory("RandomBeacon", {
    libraries: {
      DKG: dkg.address,
      Groups: groups.address,
      BLS: (await blsDeployment()).bls.address
    },
  })

  const randomBeacon: RandomBeacon = await RandomBeacon.connect(
    deployer
  ).deploy(
    testToken.address,
    maintenancePool.address,
    constants.groupSize,
    constants.signatureThreshold,
    constants.timeDKG
  ) as RandomBeacon

  await randomBeacon.transferOwnership(await governance.getAddress())

  const contracts: DeployedContracts = {
    dkg,
    groups,
    randomBeacon,
    testToken,
    maintenancePool
  }

  return contracts
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

export async function testTokenDeployment(): Promise<DeployedContracts> {
  const TestToken = await ethers.getContractFactory("TestToken")
  const testToken = await TestToken.deploy()
  await testToken.deployed()

  const contracts: DeployedContracts = { testToken }

  return contracts
}

export async function maintenancePoolDeployment(tToken: Contract): Promise<DeployedContracts> {
  const MaintenancePool = await ethers.getContractFactory("MaintenancePool")
  const maintenancePool = await MaintenancePool.deploy(tToken.address)
  await maintenancePool.deployed()

  const contracts: DeployedContracts = { maintenancePool }

  return contracts
}

export async function testDeployment(): Promise<DeployedContracts> {
  const deployer = await ethers.getSigner((await getNamedAccounts()).deployer)
  const governance = await ethers.getSigner(
    (await getNamedAccounts()).governance
  )

  const contracts = await randomBeaconDeployment()

  const RandomBeaconGovernance = await ethers.getContractFactory(
    "RandomBeaconGovernance"
  )
  const randomBeaconGovernance: RandomBeaconGovernance = await RandomBeaconGovernance.connect(
    deployer
  ).deploy(contracts.randomBeacon.address) as RandomBeaconGovernance

  await contracts.randomBeacon
    .connect(governance)
    .transferOwnership(randomBeaconGovernance.address)

  const newContracts = { randomBeaconGovernance }

  return { ...contracts, ...newContracts }
}
