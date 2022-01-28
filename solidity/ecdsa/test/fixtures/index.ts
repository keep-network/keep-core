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
} from "../../typechain"

const { to1e18 } = helpers.number

export const constants = {
  groupSize: 100,
  groupThreshold: 51,
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

export const dkgParams = {
  dkgResultChallengePeriodLength: 10,
  dkgResultSubmissionEligibilityDelay: 5,
}

export async function walletRegistryFixture(): Promise<{
  walletRegistry: WalletRegistryStub & WalletRegistry
  sortitionPool: SortitionPool
  walletOwner: SignerWithAddress
  deployer: SignerWithAddress
  thirdParty: SignerWithAddress
  operators: Operator[]
  staking: StakingStub
}> {
  await deployments.fixture(["WalletRegistry"])

  const walletRegistry: WalletRegistryStub & WalletRegistry =
    await ethers.getContract("WalletRegistry")
  const sortitionPool: SortitionPool = await ethers.getContract("SortitionPool")
  const staking: StakingStub = await ethers.getContract("StakingStub")

  const deployer: SignerWithAddress = await ethers.getNamedSigner("deployer")
  const walletOwner: SignerWithAddress = await ethers.getNamedSigner(
    "walletOwner"
  )

  const thirdParty = await ethers.getSigner((await getUnnamedAccounts())[0])

  // Accounts offset provided to slice getUnnamedAccounts have to include number
  // of unnamed accounts that were already used.
  const operators = await registerOperators(
    walletRegistry,
    (await getUnnamedAccounts()).slice(1, 1 + constants.groupSize)
  )

  await walletRegistry.updateDkgParameters(
    dkgParams.dkgResultChallengePeriodLength,
    dkgParams.dkgResultSubmissionEligibilityDelay
  )

  return {
    walletRegistry,
    sortitionPool,
    walletOwner,
    deployer,
    thirdParty,
    operators,
    staking,
  }
}
