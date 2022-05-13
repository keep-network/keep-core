import { ethers, helpers, deployments } from "hardhat"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { Contract } from "ethers"
import type {
  SortitionPool,
  BeaconDkgValidator as DKGValidator,
  RandomBeaconStub,
  TokenStaking,
  RandomBeaconGovernance,
  T,
} from "../../typechain"

const { to1e18 } = helpers.number

export const constants = {
  groupSize: 64,
  groupThreshold: 33,
  offchainDkgTime: 72, // 5 * (1 + 5) + 2 * (1 + 10) + 20
  poolWeightDivisor: to1e18(1),
  tokenStakingNotificationReward: to1e18(10000), // 10k T
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
  relayEntrySoftTimeout: 35,
  relayEntryHardTimeout: 100,
  callbackGasLimit: 200000,
  groupCreationFrequency: 10,
  groupLifeTime: 1000,
  dkgResultChallengePeriodLength: 100,
  dkgResultSubmissionTimeout: 30,
  dkgSubmitterPrecedencePeriodLength: 5,
  sortitionPoolRewardsBanDuration: 1209600, // 2 weeks
  relayEntrySubmissionFailureSlashingAmount: to1e18(1000),
  maliciousDkgResultSlashingAmount: to1e18(50000),
  relayEntryTimeoutNotificationRewardMultiplier: 40,
  unauthorizedSigningNotificationRewardMultiplier: 50,
  dkgMaliciousResultNotificationRewardMultiplier: 100,
  unauthorizedSigningSlashingAmount: to1e18(100000),
  minimumAuthorization: to1e18(200000),
  authorizationDecreaseDelay: 403200,
  authorizationDecreaseChangePeriod: 403200,
  reimbursementPoolStaticGas: 40800,
  reimbursementPoolMaxGasPrice: ethers.utils.parseUnits("500", "gwei"),
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

export async function reimbursmentPoolDeployment(): Promise<DeployedContracts> {
  const ReimbursementPool = await ethers.getContractFactory("ReimbursementPool")
  const reimbursementPool = await ReimbursementPool.deploy(
    params.reimbursementPoolStaticGas,
    params.reimbursementPoolMaxGasPrice
  )
  await reimbursementPool.deployed()

  const contracts: DeployedContracts = { reimbursementPool }

  return contracts
}

export async function randomBeaconDeployment(): Promise<DeployedContracts> {
  await deployments.fixture(["TokenStaking"])
  const t: T = await helpers.contracts.getContract<T>("T")

  const staking: TokenStaking =
    await helpers.contracts.getContract<TokenStaking>("TokenStaking")

  // TODO: Implement Hardhat deployment scripts and load deployed contracts, same
  // as it's done above for T and TokenStaking.
  const { deployer } = await helpers.signers.getNamedSigners()

  const SortitionPool = await ethers.getContractFactory("SortitionPool")
  const sortitionPool = (await SortitionPool.deploy(
    t.address,
    constants.poolWeightDivisor
  )) as SortitionPool

  const Authorization = await ethers.getContractFactory("BeaconAuthorization")
  const authorization = await Authorization.deploy()
  await authorization.deployed()

  const BeaconDkg = await ethers.getContractFactory("BeaconDkg")
  const dkg = await BeaconDkg.deploy()
  await dkg.deployed()

  const BeaconInactivity = await ethers.getContractFactory("BeaconInactivity")
  const inactivity = await BeaconInactivity.deploy()
  await inactivity.deployed()

  const BeaconDkgValidator = await ethers.getContractFactory(
    "BeaconDkgValidator"
  )
  const dkgValidator = (await BeaconDkgValidator.deploy(
    sortitionPool.address
  )) as DKGValidator
  await dkgValidator.deployed()

  const ReimbursementPool = await ethers.getContractFactory("ReimbursementPool")
  const reimbursementPool = await ReimbursementPool.deploy(
    params.reimbursementPoolStaticGas,
    params.reimbursementPoolMaxGasPrice
  )
  await reimbursementPool.deployed()

  await deployer.sendTransaction({
    to: reimbursementPool.address,
    value: ethers.utils.parseEther("100.0"), // Send 100.0 ETH
  })

  const RandomBeacon = await ethers.getContractFactory("RandomBeaconStub", {
    libraries: {
      BLS: (await blsDeployment()).bls.address,
      BeaconAuthorization: authorization.address,
      BeaconDkg: dkg.address,
      BeaconInactivity: inactivity.address,
    },
  })

  const randomBeacon: RandomBeaconStub = await RandomBeacon.deploy(
    sortitionPool.address,
    t.address,
    staking.address,
    dkgValidator.address,
    reimbursementPool.address
  )
  await randomBeacon.deployed()

  await staking.connect(deployer).approveApplication(randomBeacon.address)

  await sortitionPool.connect(deployer).transferOwnership(randomBeacon.address)
  await reimbursementPool.connect(deployer).authorize(randomBeacon.address)

  await updateTokenStakingParams(t, staking, deployer)
  await setFixtureParameters(randomBeacon)

  const contracts: DeployedContracts = {
    sortitionPool,
    staking,
    randomBeacon,
    t,
    reimbursementPool,
  }

  return contracts
}

export async function testDeployment(): Promise<DeployedContracts> {
  const { deployer, governance } = await helpers.signers.getNamedSigners()

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
  await contracts.randomBeacon
    .connect(deployer)
    .transferGovernance(randomBeaconGovernance.address)
  await randomBeaconGovernance
    .connect(deployer)
    .transferOwnership(governance.address)

  const newContracts = { randomBeaconGovernance }

  return { ...contracts, ...newContracts }
}

async function updateTokenStakingParams(
  t: T,
  staking: TokenStaking,
  deployer: SignerWithAddress
) {
  // initialNotifierTreasury should be configured high enough to execute all the
  // slashing in test suites.
  const initialNotifierTreasury = to1e18(9_000_000) // 9MM T
  await t.connect(deployer).approve(staking.address, initialNotifierTreasury)
  await staking
    .connect(deployer)
    .pushNotificationReward(initialNotifierTreasury)
  await staking
    .connect(deployer)
    .setNotificationReward(constants.tokenStakingNotificationReward)
}

async function setFixtureParameters(randomBeacon: RandomBeaconStub) {
  await randomBeacon.updateAuthorizationParameters(
    params.minimumAuthorization,
    params.authorizationDecreaseDelay,
    params.authorizationDecreaseChangePeriod
  )

  await randomBeacon.updateRelayEntryParameters(
    params.relayEntrySoftTimeout,
    params.relayEntryHardTimeout,
    params.callbackGasLimit
  )

  await randomBeacon.updateRewardParameters(
    params.sortitionPoolRewardsBanDuration,
    params.relayEntryTimeoutNotificationRewardMultiplier,
    params.unauthorizedSigningNotificationRewardMultiplier,
    params.dkgMaliciousResultNotificationRewardMultiplier
  )

  await randomBeacon.updateGroupCreationParameters(
    params.groupCreationFrequency,
    params.groupLifeTime,
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
