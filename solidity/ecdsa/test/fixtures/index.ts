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
  dkgSeedTimeout: 8,
  dkgResultChallengePeriodLength: 10,
  dkgResultSubmissionTimeout: 30,
  dkgSubmitterPrecedencePeriodLength: 5,
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
  const tToken: T = await ethers.getContract("T")
  const staking: StakingStub = await ethers.getContract("StakingStub")

  const deployer: SignerWithAddress = await ethers.getNamedSigner("deployer")
  const walletOwner: SignerWithAddress = await ethers.getNamedSigner(
    "walletOwner"
  )

  const thirdParty = await ethers.getSigner((await getUnnamedAccounts())[0])

  // Accounts offset provided to slice getUnnamedAccounts have to include number
  // of unnamed accounts that were already used.
  const unnamedAccountsOffset = 1
  const operators = await registerOperators(
    walletRegistry,
    tToken,
    (
      await getUnnamedAccounts()
    ).slice(unnamedAccountsOffset, unnamedAccountsOffset + constants.groupSize)
  )

  await walletRegistry.updateDkgParams(
    params.dkgSeedTimeout,
    params.dkgResultChallengePeriodLength,
    params.dkgResultSubmissionTimeout,
    params.dkgSubmitterPrecedencePeriodLength
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
