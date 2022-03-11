import { ethers, helpers, getNamedAccounts } from "hardhat"

import type { Contract } from "ethers"
import type {
  SortitionPool,
  BeaconDkgValidator as DKGValidator,
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
  governanceDelay: 604800, // 1 week
  relayRequestFee: to1e18(100),
  relayEntrySoftTimeout: 35,
  relayEntryHardTimeout: 100,
  callbackGasLimit: 200000,
  groupCreationFrequency: 10,
  groupLifeTime: 1000,
  dkgResultChallengePeriodLength: 100,
  dkgResultSubmissionTimeout: 30,
  dkgSubmitterPrecedencePeriodLength: 5,
  dkgResultSubmissionReward: to1e18(5),
  sortitionPoolUnlockingReward: to1e18(10),
  sortitionPoolRewardsBanDuration: 1209600, // 2 weeks
  relayEntrySubmissionFailureSlashingAmount: to1e18(1000),
  maliciousDkgResultSlashingAmount: to1e18(50000),
  relayEntryTimeoutNotificationRewardMultiplier: 40,
  unauthorizedSigningNotificationRewardMultiplier: 50,
  dkgMaliciousResultNotificationRewardMultiplier: 100,
  ineligibleOperatorNotifierReward: to1e18(200),
  unauthorizedSigningSlashingAmount: to1e18(100000),
  minimumAuthorization: to1e18(100000),
  authorizationDecreaseDelay: 0,
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

  const BeaconDkg = await ethers.getContractFactory("BeaconDkg")
  const dkg = await BeaconDkg.deploy()
  await dkg.deployed()

  const Heartbeat = await ethers.getContractFactory("Heartbeat")
  const heartbeat = await Heartbeat.deploy()
  await heartbeat.deployed()

  const BeaconDkgValidator = await ethers.getContractFactory(
    "BeaconDkgValidator"
  )
  const dkgValidator = (await BeaconDkgValidator.deploy(
    sortitionPool.address
  )) as DKGValidator
  await dkgValidator.deployed()

  const RandomBeacon = await ethers.getContractFactory("RandomBeaconStub", {
    libraries: {
      BLS: (await blsDeployment()).bls.address,
      BeaconDkg: dkg.address,
      Heartbeat: heartbeat.address,
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

  await setFixtureParameters(randomBeacon)

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
    await RandomBeaconGovernance.deploy(
      contracts.randomBeacon.address,
      params.governanceDelay
    )
  await randomBeaconGovernance.deployed()
  await contracts.randomBeacon.transferOwnership(randomBeaconGovernance.address)

  const newContracts = { randomBeaconGovernance }

  return { ...contracts, ...newContracts }
}

async function setFixtureParameters(randomBeacon: RandomBeaconStub) {
  await randomBeacon.updateAuthorizationParameters(
    params.minimumAuthorization,
    params.authorizationDecreaseDelay
  )

  await randomBeacon.updateRelayEntryParameters(
    params.relayRequestFee,
    params.relayEntrySoftTimeout,
    params.relayEntryHardTimeout,
    params.callbackGasLimit
  )

  await randomBeacon.updateRewardParameters(
    params.dkgResultSubmissionReward,
    params.sortitionPoolUnlockingReward,
    params.ineligibleOperatorNotifierReward,
    params.sortitionPoolRewardsBanDuration,
    params.relayEntryTimeoutNotificationRewardMultiplier,
    params.unauthorizedSigningNotificationRewardMultiplier,
    params.dkgMaliciousResultNotificationRewardMultiplier
  )

  await randomBeacon.updateGroupCreationParameters(
    params.groupCreationFrequency,
    params.groupLifeTime
  )

  await randomBeacon.updateDkgParameters(
    params.dkgResultChallengePeriodLength,
    params.dkgResultSubmissionTimeout,
    params.dkgSubmitterPrecedencePeriodLength
  )

  await randomBeacon.updateSlashingParameters(
    params.relayEntrySubmissionFailureSlashingAmount,
    params.maliciousDkgResultSlashingAmount,
    params.unauthorizedSigningSlashingAmount
  )
}
