import { deployments, ethers, helpers, getUnnamedAccounts } from "hardhat"
import { smock } from "@defi-wonderland/smock"

// eslint-disable-next-line import/no-cycle
import { registerOperators } from "../utils/operators"

import type { IWalletOwner } from "../../typechain/IWalletOwner"
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
import type { FakeContract } from "@defi-wonderland/smock"

const { to1e18 } = helpers.number

export const constants = {
  groupSize: 100,
  groupThreshold: 51,
  poolWeightDivisor: to1e18(1),
  governanceDelayCritical: 1209600, // 2 weeks
}

export const dkgState = {
  IDLE: 0,
  AWAITING_SEED: 1,
  AWAITING_RESULT: 2,
  CHALLENGE: 3,
}

export const params = {
  minimumAuthorization: to1e18(400000),
  authorizationDecreaseDelay: 5184000,
  dkgSeedTimeout: 8,
  dkgResultChallengePeriodLength: 10,
  dkgResultSubmissionTimeout: 30,
  dkgSubmitterPrecedencePeriodLength: 5,
}

export const walletRegistryFixture = deployments.createFixture(
  async (): Promise<{
    walletRegistry: WalletRegistryStub & WalletRegistry
    walletRegistryGovernance: WalletRegistryGovernance
    sortitionPool: SortitionPool
    staking: StakingStub
    walletOwner: FakeContract<IWalletOwner>
    deployer: SignerWithAddress
    governance: SignerWithAddress
    thirdParty: SignerWithAddress
    operators: Operator[]
  }> => {
    await deployments.fixture(["WalletRegistry"])

    const walletRegistry: WalletRegistryStub & WalletRegistry =
      await ethers.getContract("WalletRegistry")
    const walletRegistryGovernance: WalletRegistryGovernance =
      await ethers.getContract("WalletRegistryGovernance")
    const sortitionPool: SortitionPool = await ethers.getContract(
      "SortitionPool"
    )
    const tToken: T = await ethers.getContract("T")
    const staking: StakingStub = await ethers.getContract("StakingStub")

    const deployer: SignerWithAddress = await ethers.getNamedSigner("deployer")
    const governance: SignerWithAddress = await ethers.getNamedSigner(
      "governance"
    )

    const thirdParty: SignerWithAddress = await ethers.getSigner(
      (
        await getUnnamedAccounts()
      )[0]
    )

    // Accounts offset provided to slice getUnnamedAccounts have to include number
    // of unnamed accounts that were already used.
    const unnamedAccountsOffset = 1
    const operators: Operator[] = await registerOperators(
      walletRegistry,
      tToken,
      (
        await getUnnamedAccounts()
      ).slice(
        unnamedAccountsOffset,
        unnamedAccountsOffset + constants.groupSize
      )
    )

    // Set parameters with tweaked values to reduce test execution time.
    await updateWalletRegistryParams(walletRegistryGovernance, governance)

    // Mock Wallet Owner contract.
    const walletOwner: FakeContract<IWalletOwner> = await initializeWalletOwner(
      walletRegistryGovernance,
      governance
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
)

export async function updateWalletRegistryParams(
  walletRegistryGovernance: WalletRegistryGovernance,
  governance: SignerWithAddress
) {
  await walletRegistryGovernance
    .connect(governance)
    .beginMinimumAuthorizationUpdate(params.minimumAuthorization)

  await walletRegistryGovernance
    .connect(governance)
    .beginAuthorizationDecreaseDelayUpdate(params.authorizationDecreaseDelay)

  await walletRegistryGovernance
    .connect(governance)
    .beginDkgSeedTimeoutUpdate(params.dkgSeedTimeout)

  await walletRegistryGovernance
    .connect(governance)
    .beginDkgResultChallengePeriodLengthUpdate(
      params.dkgResultChallengePeriodLength
    )

  await walletRegistryGovernance
    .connect(governance)
    .beginDkgResultSubmissionTimeoutUpdate(params.dkgResultSubmissionTimeout)

  await walletRegistryGovernance
    .connect(governance)
    .beginDkgSubmitterPrecedencePeriodLengthUpdate(
      params.dkgSubmitterPrecedencePeriodLength
    )

  await helpers.time.increaseTime(constants.governanceDelayCritical)

  await walletRegistryGovernance
    .connect(governance)
    .finalizeMinimumAuthorizationUpdate()

  await walletRegistryGovernance
    .connect(governance)
    .finalizeAuthorizationDecreaseDelayUpdate()

  await walletRegistryGovernance
    .connect(governance)
    .finalizeDkgSeedTimeoutUpdate()

  await walletRegistryGovernance
    .connect(governance)
    .finalizeDkgResultChallengePeriodLengthUpdate()

  await walletRegistryGovernance
    .connect(governance)
    .finalizeDkgResultSubmissionTimeoutUpdate()

  await walletRegistryGovernance
    .connect(governance)
    .finalizeDkgSubmitterPrecedencePeriodLengthUpdate()
}

async function initializeWalletOwner(
  walletRegistryGovernance: WalletRegistryGovernance,
  governance: SignerWithAddress
): Promise<FakeContract<IWalletOwner>> {
  const deployer: SignerWithAddress = await ethers.getNamedSigner("deployer")

  const walletOwner: FakeContract<IWalletOwner> =
    await smock.fake<IWalletOwner>("IWalletOwner")

  await deployer.sendTransaction({
    to: walletOwner.address,
    value: ethers.utils.parseEther("1"),
  })

  await walletRegistryGovernance
    .connect(governance)
    .initializeWalletOwner(walletOwner.address)

  return walletOwner
}
