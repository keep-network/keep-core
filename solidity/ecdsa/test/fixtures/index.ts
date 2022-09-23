import { deployments, ethers, helpers } from "hardhat"
import { smock } from "@defi-wonderland/smock"

// eslint-disable-next-line import/no-cycle
import { registerOperators } from "../utils/operators"
import { fakeRandomBeacon } from "../utils/randomBeacon"

import type { IWalletOwner } from "../../typechain/IWalletOwner"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { Operator } from "../utils/operators"
import type {
  SortitionPool,
  ReimbursementPool,
  WalletRegistry,
  WalletRegistryStub,
  WalletRegistryGovernance,
  TokenStaking,
  T,
  IRandomBeacon,
} from "../../typechain"
import type { FakeContract } from "@defi-wonderland/smock"

const { to1e18 } = helpers.number

export const constants = {
  groupSize: 100,
  groupThreshold: 51,
  poolWeightDivisor: to1e18(1),
  tokenStakingNotificationReward: to1e18(10_000), // 10k T
  governanceDelay: 604_800, // 1 week
}

export const dkgState = {
  IDLE: 0,
  AWAITING_SEED: 1,
  AWAITING_RESULT: 2,
  CHALLENGE: 3,
}

export const params = {
  minimumAuthorization: to1e18(40_000),
  authorizationDecreaseDelay: 3_888_000,
  authorizationDecreaseChangePeriod: 3_888_000,
  dkgSeedTimeout: 8,
  dkgResultChallengePeriodLength: 10,
  dkgResultChallengeExtraGas: 50_000,
  dkgResultSubmissionTimeout: 30,
  dkgSubmitterPrecedencePeriodLength: 5,
  sortitionPoolRewardsBanDuration: 1_209_600, // 14 days
}

export const walletRegistryFixture = deployments.createFixture(
  async (): Promise<{
    tToken: T
    walletRegistry: WalletRegistryStub & WalletRegistry
    walletRegistryGovernance: WalletRegistryGovernance
    sortitionPool: SortitionPool
    reimbursementPool: ReimbursementPool
    staking: TokenStaking
    randomBeacon: FakeContract<IRandomBeacon>
    walletOwner: FakeContract<IWalletOwner>
    deployer: SignerWithAddress
    governance: SignerWithAddress
    thirdParty: SignerWithAddress
    operators: Operator[]
  }> => {
    // Due to a [bug] in hardhat-gas-reporter plugin we avoid using `--deploy-fixture`
    // flag of `hardhat-deploy` plugin. This requires us to load a global fixture
    // (`deployments.fixture()`) instead of loading a specific tag (`deployments.fixture(<tag>)`.
    // bug: https://github.com/cgewecke/hardhat-gas-reporter/issues/86
    await deployments.fixture()

    const walletRegistry: WalletRegistryStub & WalletRegistry =
      await helpers.contracts.getContract("WalletRegistry")
    const walletRegistryGovernance: WalletRegistryGovernance =
      await helpers.contracts.getContract("WalletRegistryGovernance")
    const sortitionPool: SortitionPool = await helpers.contracts.getContract(
      "EcdsaSortitionPool"
    )
    const tToken: T = await helpers.contracts.getContract("T")
    const staking: TokenStaking = await helpers.contracts.getContract(
      "TokenStaking"
    )

    const reimbursementPool: ReimbursementPool =
      await helpers.contracts.getContract("ReimbursementPool")

    const randomBeacon: FakeContract<IRandomBeacon> = await fakeRandomBeacon(
      walletRegistry
    )

    const { deployer, governance, chaosnetOwner } =
      await helpers.signers.getNamedSigners()

    await sortitionPool.connect(chaosnetOwner).deactivateChaosnet()

    const [thirdParty] = await helpers.signers.getUnnamedSigners()

    // Accounts offset provided to slice getUnnamedAccounts have to include number
    // of unnamed accounts that were already used.
    const unnamedAccountsOffset = 1
    const operators: Operator[] = await registerOperators(
      walletRegistry,
      tToken,
      constants.groupSize,
      unnamedAccountsOffset
    )

    // Set up TokenStaking parameters
    await updateTokenStakingParams(tToken, staking, deployer)

    // Set parameters with tweaked values to reduce test execution time.
    await updateWalletRegistryParams(walletRegistryGovernance, governance)

    await fundReimbursementPool(deployer, reimbursementPool)

    // Mock Wallet Owner contract.
    const walletOwner: FakeContract<IWalletOwner> = await initializeWalletOwner(
      walletRegistryGovernance,
      governance
    )

    return {
      tToken,
      walletRegistry,
      sortitionPool,
      reimbursementPool,
      randomBeacon,
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

async function updateTokenStakingParams(
  tToken: T,
  staking: TokenStaking,
  deployer: SignerWithAddress
) {
  const initialNotifierTreasury = constants.tokenStakingNotificationReward.mul(
    constants.groupSize
  )
  await tToken
    .connect(deployer)
    .approve(staking.address, initialNotifierTreasury)
  await staking
    .connect(deployer)
    .pushNotificationReward(initialNotifierTreasury)
  await staking
    .connect(deployer)
    .setNotificationReward(constants.tokenStakingNotificationReward)
}

export async function updateWalletRegistryParams(
  walletRegistryGovernance: WalletRegistryGovernance,
  governance: SignerWithAddress
): Promise<void> {
  await walletRegistryGovernance
    .connect(governance)
    .beginMinimumAuthorizationUpdate(params.minimumAuthorization)

  await walletRegistryGovernance
    .connect(governance)
    .beginAuthorizationDecreaseDelayUpdate(params.authorizationDecreaseDelay)

  await walletRegistryGovernance
    .connect(governance)
    .beginAuthorizationDecreaseChangePeriodUpdate(
      params.authorizationDecreaseChangePeriod
    )

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

  await walletRegistryGovernance
    .connect(governance)
    .beginSortitionPoolRewardsBanDurationUpdate(
      params.sortitionPoolRewardsBanDuration
    )

  await helpers.time.increaseTime(constants.governanceDelay)

  await walletRegistryGovernance
    .connect(governance)
    .finalizeMinimumAuthorizationUpdate()

  await walletRegistryGovernance
    .connect(governance)
    .finalizeAuthorizationDecreaseDelayUpdate()

  await walletRegistryGovernance
    .connect(governance)
    .finalizeAuthorizationDecreaseChangePeriodUpdate()

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

  await walletRegistryGovernance
    .connect(governance)
    .finalizeSortitionPoolRewardsBanDurationUpdate()
}

export async function initializeWalletOwner(
  walletRegistryGovernance: WalletRegistryGovernance,
  governance: SignerWithAddress
): Promise<FakeContract<IWalletOwner>> {
  const { deployer } = await helpers.signers.getNamedSigners()

  const walletOwner: FakeContract<IWalletOwner> =
    await smock.fake<IWalletOwner>("IWalletOwner")

  await deployer.sendTransaction({
    to: walletOwner.address,
    value: ethers.utils.parseEther("1000"),
  })

  await walletRegistryGovernance
    .connect(governance)
    .initializeWalletOwner(walletOwner.address)

  return walletOwner
}

async function fundReimbursementPool(
  deployer: SignerWithAddress,
  reimbursementPool: ReimbursementPool
) {
  await deployer.sendTransaction({
    to: reimbursementPool.address,
    value: ethers.utils.parseEther("100.0"), // Send 100.0 ETH
  })
}
