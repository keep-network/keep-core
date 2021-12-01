import { Contract } from "ethers"
import { ethers, helpers, getNamedAccounts } from "hardhat"
import type {
  SortitionPool,
  DKGValidator,
  RandomBeaconStub,
  RandomBeaconGovernance,
  StakingStub,
} from "../../typechain"

const { to1e18 } = helpers.number

export const constants = {
  groupSize: 64,
  groupThreshold: 33,
  offchainDkgTime: 72, // 5 * (1 + 5) + 2 * (1 + 10) + 20
  minimumStake: to1e18(100000),
  poolWeightDivisor: to1e18(1),
}

export const dkgState = {
  IDLE: 0,
  AWAITING_SEED: 1,
  KEY_GENERATION: 2,
  AWAITING_RESULT: 3,
  CHALLENGE: 4,
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
  const deployer = await ethers.getSigner((await getNamedAccounts()).deployer)

  const { testToken } = await testTokenDeployment()

  const StakingStub = await ethers.getContractFactory("StakingStub")
  const stakingStub: StakingStub = await StakingStub.deploy()

  const SortitionPool = await ethers.getContractFactory("SortitionPool")
  const sortitionPool = (await SortitionPool.deploy(
    stakingStub.address,
    testToken.address,
    constants.poolWeightDivisor
  )) as SortitionPool

  const DKG = await ethers.getContractFactory("DKG")
  const dkg = await DKG.deploy()
  await dkg.deployed()

  const DKGValidator = await ethers.getContractFactory("DKGValidator")
  const dkgValidator = (await DKGValidator.deploy(
    sortitionPool.address
  )) as DKGValidator
  await dkgValidator.deployed()

  const RandomBeacon = await ethers.getContractFactory("RandomBeaconStub", {
    libraries: {
      BLS: (await blsDeployment()).bls.address,
      DKG: dkg.address,
    },
  })
  const randomBeacon: RandomBeaconStub = await RandomBeacon.deploy(
    sortitionPool.address,
    testToken.address,
    stakingStub.address,
    dkgValidator.address
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
