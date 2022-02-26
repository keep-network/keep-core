/* eslint-disable @typescript-eslint/no-unused-expressions */
import { Sign } from "crypto"

import { deployments, ethers, getUnnamedAccounts, helpers } from "hardhat"
import { expect } from "chai"
import { to1e18 } from "@keep-network/hardhat-helpers/dist/src/number"

import { constants, params } from "./fixtures"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type {
  WalletRegistry,
  SortitionPool,
  TokenStaking,
  T,
} from "../typechain"

const { createSnapshot, restoreSnapshot } = helpers.snapshot

describe("WalletRegistry - Pool", () => {
  let t: T
  let walletRegistry: WalletRegistry
  let sortitionPool: SortitionPool
  let staking: TokenStaking

  let deployer: SignerWithAddress

  let owner: SignerWithAddress
  let stakingProvider: SignerWithAddress
  let operator: SignerWithAddress
  let authorizer: SignerWithAddress
  let beneficiary: SignerWithAddress
  let thirdParty: SignerWithAddress

  const stakedAmount = to1e18(1000000) // 1M T

  before("load test fixture", async () => {
    await deployments.fixture(["WalletRegistry"])

    t = await ethers.getContract("T")
    walletRegistry = await ethers.getContract("WalletRegistry")
    sortitionPool = await ethers.getContract("SortitionPool")
    staking = await ethers.getContract("TokenStaking")

    deployer = await ethers.getNamedSigner("deployer")

    const accounts = await getUnnamedAccounts()
    owner = await ethers.getSigner(accounts[1])
    stakingProvider = await ethers.getSigner(accounts[2])
    operator = await ethers.getSigner(accounts[3])
    authorizer = await ethers.getSigner(accounts[4])
    beneficiary = await ethers.getSigner(accounts[5])
    thirdParty = await ethers.getSigner(accounts[6])

    await t.connect(deployer).mint(owner.address, stakedAmount)
    await t.connect(owner).approve(staking.address, stakedAmount)
    await staking
      .connect(owner)
      .stake(
        stakingProvider.address,
        beneficiary.address,
        authorizer.address,
        stakedAmount
      )
  })

  describe("registerOperator", () => {
    context(
      "when operator has been already set for the staking provider",
      () => {
        before(async () => {
          await createSnapshot()
          await walletRegistry
            .connect(stakingProvider)
            .registerOperator(operator.address)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should revert", async () => {
          await expect(
            walletRegistry
              .connect(stakingProvider)
              .registerOperator(operator.address)
          ).to.be.revertedWith("Operator already set for the staking provider")
        })
      }
    )

    context("when the operator is already in use", () => {
      before(async () => {
        await createSnapshot()
        await walletRegistry
          .connect(thirdParty)
          .registerOperator(operator.address)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should revert", async () => {
        await expect(
          walletRegistry
            .connect(stakingProvider)
            .registerOperator(operator.address)
        ).to.be.revertedWith("Operator address already in use")
      })
    })

    context(
      "when staking provider is setting operator for the first time",
      () => {
        let tx

        before(async () => {
          await createSnapshot()
          tx = await walletRegistry
            .connect(stakingProvider)
            .registerOperator(operator.address)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should set staking provider -> operator mapping", async () => {
          expect(
            await walletRegistry.stakingProviderToOperator(
              stakingProvider.address
            )
          ).to.equal(operator.address)
        })

        it("should set operator -> staking provider mapping", async () => {
          expect(
            await walletRegistry.operatorToStakingProvider(operator.address)
          ).to.equal(stakingProvider.address)
        })

        it("should emit OperatorRegistered event", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "OperatorRegistered")
            .withArgs(stakingProvider.address, operator.address)
        })
      }
    )
  })

  describe("authorizationIncreased", () => {
    context("when called not by the staking contract", () => {
      it("should revert", async () => {
        await expect(
          walletRegistry
            .connect(thirdParty)
            .authorizationIncreased(stakingProvider.address, 0, stakedAmount)
        ).to.be.revertedWith("Caller is not the staking contract")
      })
    })

    context("when operator does not have the minimum authorization", () => {
      it("should revert", async () => {
        const minimumAuthorization = await walletRegistry.minimumAuthorization()

        await expect(
          staking
            .connect(authorizer)
            .increaseAuthorization(
              stakingProvider.address,
              walletRegistry.address,
              minimumAuthorization.sub(1)
            )
        ).to.be.revertedWith("Authorization below the minimum")
      })
    })

    context("when operator has just the minimum authorization", () => {
      before(async () => {
        await createSnapshot()
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should let the update pass", async () => {
        const minimumAuthorization = await walletRegistry.minimumAuthorization()

        await staking
          .connect(authorizer)
          .increaseAuthorization(
            stakingProvider.address,
            walletRegistry.address,
            minimumAuthorization
          )
        // ok, did not revert
      })
    })

    context("when the operator is unknown", () => {
      before(async () => {
        await createSnapshot()
        await staking
          .connect(authorizer)
          .increaseAuthorization(
            stakingProvider.address,
            walletRegistry.address,
            stakedAmount
          )
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should ignore the update", async () => {
        expect(await sortitionPool.isOperatorInPool(stakingProvider.address)).to
          .be.false
        expect(await sortitionPool.isOperatorInPool(operator.address)).to.be
          .false
      })
    })

    context("when the operator is not in the sortition pool", () => {
      before(async () => {
        await createSnapshot()
        await walletRegistry
          .connect(stakingProvider)
          .registerOperator(operator.address)
        await staking
          .connect(authorizer)
          .increaseAuthorization(
            stakingProvider.address,
            walletRegistry.address,
            stakedAmount
          )
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should ignore the update", async () => {
        expect(await sortitionPool.isOperatorInPool(stakingProvider.address)).to
          .be.false
        expect(await sortitionPool.isOperatorInPool(operator.address)).to.be
          .false
      })
    })

    context("when the operator is in the sortition pool", () => {
      const initialAuthorization = to1e18(500000)
      const authorizationIncrease = to1e18(30000)

      before(async () => {
        await createSnapshot()

        await walletRegistry
          .connect(stakingProvider)
          .registerOperator(operator.address)

        await staking
          .connect(authorizer)
          .increaseAuthorization(
            stakingProvider.address,
            walletRegistry.address,
            initialAuthorization
          )
        await walletRegistry.connect(operator).joinSortitionPool()

        await staking
          .connect(authorizer)
          .increaseAuthorization(
            stakingProvider.address,
            walletRegistry.address,
            authorizationIncrease
          )
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should update the operator status in the pool", async () => {
        expect(await sortitionPool.isOperatorInPool(operator.address)).to.be
          .true
        expect(await sortitionPool.getPoolWeight(operator.address)).to.equal(
          initialAuthorization
            .add(authorizationIncrease)
            .div(constants.poolWeightDivisor)
        )
      })
    })
  })

  describe("joinSortitionPool", () => {
    context("when the operator is unknown", () => {
      before(async () => {
        await createSnapshot()
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should revert", async () => {
        await expect(
          walletRegistry.connect(thirdParty).joinSortitionPool()
        ).to.be.revertedWith("Unknown operator")
      })
    })

    context("when the operator has no stake authorized", () => {
      before(async () => {
        await createSnapshot()

        await walletRegistry
          .connect(stakingProvider)
          .registerOperator(operator.address)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should revert", async () => {
        await expect(
          walletRegistry.connect(operator).joinSortitionPool()
        ).to.be.revertedWith("Authorization below the minimum")
      })
    })

    context("when the operator has the minimum stake authorized", () => {
      let minimumAuthorization

      before(async () => {
        await createSnapshot()

        await walletRegistry
          .connect(stakingProvider)
          .registerOperator(operator.address)

        minimumAuthorization = await walletRegistry.minimumAuthorization()

        await staking
          .connect(authorizer)
          .increaseAuthorization(
            stakingProvider.address,
            walletRegistry.address,
            minimumAuthorization
          )

        await walletRegistry.connect(operator).joinSortitionPool()
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should insert operator into the pool", async () => {
        expect(await sortitionPool.isOperatorInPool(operator.address)).to.be
          .true
      })

      it("should use a correct stake weight", async () => {
        expect(await sortitionPool.getPoolWeight(operator.address)).to.equal(
          minimumAuthorization.div(constants.poolWeightDivisor)
        )
      })
    })

    context(
      "when the operator has more than the minimum stake authorized",
      () => {
        let authorizedStake

        before(async () => {
          await createSnapshot()

          await walletRegistry
            .connect(stakingProvider)
            .registerOperator(operator.address)

          const minimumAuthorization =
            await walletRegistry.minimumAuthorization()
          authorizedStake = minimumAuthorization.mul(2)

          await staking
            .connect(authorizer)
            .increaseAuthorization(
              stakingProvider.address,
              walletRegistry.address,
              authorizedStake
            )

          await walletRegistry.connect(operator).joinSortitionPool()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should insert operator into the pool", async () => {
          expect(await sortitionPool.isOperatorInPool(operator.address)).to.be
            .true
        })

        it("should use a correct stake weight", async () => {
          expect(await sortitionPool.getPoolWeight(operator.address)).to.equal(
            authorizedStake.div(constants.poolWeightDivisor)
          )
        })
      }
    )
  })
})
