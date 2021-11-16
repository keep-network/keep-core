import { Contract } from "ethers"
import { ethers, getNamedAccounts } from "hardhat"
import type {
  SortitionPool,
  SortitionPoolStub,
  RandomBeaconStub,
  RandomBeaconGovernance,
  StakingStub,
} from "../../typechain"
import { to1e18 } from "../functions"

export const constants = {
  groupSize: 64,
  groupThreshold: 33,
  signatureThreshold: 48, // groupThreshold + (groupSize - groupThreshold) / 2
  offchainDkgTime: 72, // 5 * (1 + 5) + 2 * (1 + 10) + 20
  minimumStake: to1e18(100000),
  poolWeightDivisor: 2000,
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

export async function randomBeaconDeployment(
  sortitionPoolStub?: SortitionPoolStub
): Promise<DeployedContracts> {
  const deployer = await ethers.getSigner((await getNamedAccounts()).deployer)

  const StakingStub = await ethers.getContractFactory("StakingStub")
  const stakingStub: StakingStub = await StakingStub.deploy()

  // Use the sortition pool stub if it's passed or the real sortition
  // pool otherwise.
  let sortitionPool: SortitionPool | SortitionPoolStub
  if (typeof sortitionPoolStub !== "undefined") {
    sortitionPool = sortitionPoolStub
  } else {
    const SortitionPool = await ethers.getContractFactory("SortitionPool")
    sortitionPool = (await SortitionPool.deploy(
      stakingStub.address,
      constants.minimumStake,
      constants.poolWeightDivisor
    )) as SortitionPool
  }

  const { testToken } = await testTokenDeployment()

  const RandomBeacon = await ethers.getContractFactory("RandomBeaconStub", {
    libraries: {
      BLS: (await blsDeployment()).bls.address,
    },
  })
  const randomBeacon: RandomBeaconStub = await RandomBeacon.deploy(
    sortitionPool.address,
    testToken.address,
    stakingStub.address
  )
  await randomBeacon.deployed()

  await sortitionPool.connect(deployer).transferOwnership(randomBeacon.address)

  const contracts: DeployedContracts = {
    sortitionPool,
    stakingStub,
    randomBeacon,
    testToken,
  }

  return contracts
}

export async function testDeployment(): Promise<DeployedContracts> {
  const SortitionPoolStub = await ethers.getContractFactory("SortitionPoolStub")
  const sortitionPoolStub: SortitionPoolStub = await SortitionPoolStub.deploy()
  const contracts = await randomBeaconDeployment(sortitionPoolStub)

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
