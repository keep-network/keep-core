import { deployments, ethers, helpers, getUnnamedAccounts } from "hardhat"

// eslint-disable-next-line import/no-cycle
import { registerOperators } from "../utils/operators"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { Operator } from "../utils/operators"
import type {
  SortitionPool,
  WalletRegistry,
  WalletRegistryStub,
  StakingStub,
  WalletRegistryGovernance,
  T,
} from "../../typechain"

const { to1e18 } = helpers.number

export const constants = {
  groupSize: 100,
  groupThreshold: 51,
  minimumStake: to1e18(100000),
  poolWeightDivisor: to1e18(1),
}

export const dkgState = {
  IDLE: 0,
  AWAITING_SEED: 1,
  AWAITING_RESULT: 2,
  CHALLENGE: 3,
}

export const params = {
  dkgResultChallengePeriodLength: 10,
  governanceDelay: 43200, // 12 hours
  dkgResultSubmissionTimeout: 30,
  dkgSubmitterPrecedencePeriodLength: 5,
}

export async function walletRegistryFixture(): Promise<{
  walletRegistry: WalletRegistryStub & WalletRegistry
  walletRegistryGovernance: WalletRegistryGovernance
  sortitionPool: SortitionPool
  walletOwner: SignerWithAddress
  deployer: SignerWithAddress
  governance: SignerWithAddress
  thirdParty: SignerWithAddress
  operators: Operator[]
  staking: StakingStub
}> {
  await deployments.fixture(["WalletRegistry"])

  const walletRegistry: WalletRegistryStub & WalletRegistry =
    await ethers.getContract("WalletRegistry")
  const walletRegistryGovernance: WalletRegistryGovernance =
    await ethers.getContract("WalletRegistryGovernance")
  const sortitionPool: SortitionPool = await ethers.getContract("SortitionPool")
  const tToken: T = await ethers.getContract("T")
  const staking: StakingStub = await ethers.getContract("StakingStub")

  const deployer: SignerWithAddress = await ethers.getNamedSigner("deployer")
  const governance: SignerWithAddress = await ethers.getNamedSigner(
    "governance"
  )
  const walletOwner: SignerWithAddress = await ethers.getNamedSigner(
    "walletOwner"
  )

  const thirdParty = await ethers.getSigner((await getUnnamedAccounts())[0])

  // Accounts offset provided to slice getUnnamedAccounts have to include number
  // of unnamed accounts that were already used.
  const operators = await registerOperators(
    walletRegistry,
    tToken,
    (await getUnnamedAccounts()).slice(1, 1 + constants.groupSize)
  )

  await walletRegistryGovernance
    .connect(governance)
    .beginDkgResultChallengePeriodLengthUpdate(
      params.dkgResultChallengePeriodLength
    )
  await walletRegistryGovernance
    .connect(governance)
    .beginDkgResultSubmissionEligibilityDelayUpdate(
      params.dkgResultSubmissionEligibilityDelay
    )

  await helpers.time.increaseTime(params.governanceDelay)

  await walletRegistryGovernance
    .connect(governance)
    .finalizeDkgResultChallengePeriodLengthUpdate()

  await walletRegistryGovernance
    .connect(governance)
    .finalizeDkgResultSubmissionEligibilityDelayUpdate()

  // TODO: remove and replace by the above
  await walletRegistry.updateDkgParams(
    params.dkgResultChallengePeriodLength,
    params.dkgResultSubmissionTimeout,
    params.dkgSubmitterPrecedencePeriodLength
  )

  return {
    walletRegistry,
    sortitionPool,
    walletOwner,
    deployer,
    governance,
    thirdParty,
    operators,
    staking,
    walletRegistryGovernance,
  }
}
