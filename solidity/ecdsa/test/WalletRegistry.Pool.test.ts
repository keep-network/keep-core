import { deployments, ethers, getUnnamedAccounts, helpers } from "hardhat"
import { expect } from "chai"

import { params } from "./fixtures"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { WalletRegistry, SortitionPool, StakingStub } from "../typechain"

const { createSnapshot, restoreSnapshot } = helpers.snapshot

describe("WalletRegistry - Pool", () => {
  let walletRegistry: WalletRegistry
  let sortitionPool: SortitionPool
  let staking: StakingStub
  let stakingProvider: SignerWithAddress
  let operator: SignerWithAddress
  let thirdParty: SignerWithAddress

  before("load test fixture", async () => {
    await deployments.fixture(["WalletRegistry"])

    walletRegistry = await ethers.getContract("WalletRegistry")
    sortitionPool = await ethers.getContract("SortitionPool")
    staking = await ethers.getContract("StakingStub")

    stakingProvider = await ethers.getSigner((await getUnnamedAccounts())[0])
    operator = await ethers.getSigner((await getUnnamedAccounts())[1])
    thirdParty = await ethers.getSigner((await getUnnamedAccounts())[2])

    // FIXME: Remove this assignment once Token Staking integration is implemented.
    stakingProvider = operator
  })

  describe("registerOperator", () => {
    context("when the staking provider has stake", () => {
      before(async () => {
        await createSnapshot()

        await staking.increaseAuthorization(
          stakingProvider.address,
          walletRegistry.address,
          params.minimumAuthorization
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
          params.minimumAuthorization
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
