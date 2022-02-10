import {
  deployments,
  ethers,
  getUnnamedAccounts,
  waffle,
  helpers,
} from "hardhat"
import { expect } from "chai"

import { constants } from "./fixtures"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { WalletRegistry, SortitionPool, StakingStub } from "../typechain"

const { createSnapshot, restoreSnapshot } = helpers.snapshot

async function fixture(): Promise<{
  walletRegistry: WalletRegistry
  sortitionPool: SortitionPool
  staking: StakingStub
  stakingProvider: SignerWithAddress
  operator: SignerWithAddress
  thirdParty: SignerWithAddress
}> {
  await deployments.fixture(["WalletRegistry"])

  const walletRegistry: WalletRegistry = await ethers.getContract(
    "WalletRegistry"
  )
  const sortitionPool: SortitionPool = await ethers.getContract("SortitionPool")
  const staking: StakingStub = await ethers.getContract("StakingStub")

  const stakingProvider = await ethers.getSigner(
    (
      await getUnnamedAccounts()
    )[0]
  )
  const operator = await ethers.getSigner((await getUnnamedAccounts())[1])
  const thirdParty = await ethers.getSigner((await getUnnamedAccounts())[2])

  return {
    walletRegistry,
    sortitionPool,
    staking,
    stakingProvider,
    operator,
    thirdParty,
  }
}

describe("WalletRegistry - Pool", () => {
  let walletRegistry: WalletRegistry
  let sortitionPool: SortitionPool
  let staking: StakingStub
  let stakingProvider: SignerWithAddress
  let operator: SignerWithAddress
  let thirdParty: SignerWithAddress

  before("load test fixture", async () => {
    await createSnapshot()

    // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;({
      walletRegistry,
      sortitionPool,
      staking,
      stakingProvider,
      operator,
      thirdParty,
    } = await waffle.loadFixture(fixture))

    // FIXME: Remove this assignment once Token Staking integration is implemented.
    stakingProvider = operator
  })

  after(async () => {
    await restoreSnapshot()
  })

  describe("registerOperator", () => {
    context("when the staking provider has stake", () => {
      before(async () => {
        await createSnapshot()

        await staking.increaseAuthorization(
          stakingProvider.address,
          walletRegistry.address,
          constants.minimumStake
        )
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("when the operator is not registered yet", () => {
        before(async () => {
          await createSnapshot()

          await walletRegistry.connect(operator).registerOperator()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should register the operator in the sortition pool", async () => {
          await expect(await sortitionPool.isOperatorInPool(operator.address))
            .to.be.true
        })
      })

      context("when the operator is already registered", () => {
        before(async () => {
          await createSnapshot()

          await walletRegistry.connect(operator).registerOperator()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should revert", async () => {
          await expect(
            walletRegistry.connect(operator).registerOperator()
          ).to.be.revertedWith("Operator is already registered")
        })
      })
    })
  })

  describe("updateOperatorStatus", () => {
    context("when operator is not registered", () => {
      it("should revert", async () => {
        await expect(
          walletRegistry
            .connect(thirdParty)
            .updateOperatorStatus(operator.address)
        ).to.be.revertedWith("Operator is not registered in the pool")
      })
    })

    context("when operator is registered", () => {
      before(async () => {
        await createSnapshot()

        await staking.increaseAuthorization(
          stakingProvider.address,
          ethers.constants.AddressZero,
          constants.minimumStake
        )
        await walletRegistry.connect(operator).registerOperator()
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("when status update removes operator from sortition pool", () => {
        before(async () => {
          await createSnapshot()

          // Simulate the operator became ineligible.
          await staking.requestAuthorizationDecrease(operator.address)

          await walletRegistry
            .connect(thirdParty)
            .updateOperatorStatus(operator.address)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should remove operator from the pool", async () => {
          await expect(await sortitionPool.isOperatorInPool(operator.address))
            .to.be.false
        })
      })
    })
  })
})
