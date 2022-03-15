/* eslint-disable @typescript-eslint/no-unused-expressions */
import { deployments, ethers, getUnnamedAccounts, helpers } from "hardhat"
import { smock } from "@defi-wonderland/smock"
import { expect } from "chai"
import { to1e18 } from "@keep-network/hardhat-helpers/dist/src/number"

import { constants, params, updateWalletRegistryParams } from "./fixtures"

import type { FakeContract } from "@defi-wonderland/smock"
import type { ContractTransaction } from "ethers"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type {
  WalletRegistry,
  SortitionPool,
  TokenStaking,
  T,
  IApplication,
  WalletRegistryGovernance,
} from "../typechain"

const { mineBlocks } = helpers.time

const { createSnapshot, restoreSnapshot } = helpers.snapshot

const ZERO_ADDRESS = ethers.constants.AddressZero
const MAX_UINT64 = ethers.BigNumber.from("18446744073709551615") // 2^64 - 1

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
  let walletOwner: SignerWithAddress
  let slasher: FakeContract<IApplication>

  const stakedAmount = to1e18(1000000) // 1M T
  let minimumAuthorization

  before("load test fixture", async () => {
    await deployments.fixture(["WalletRegistry"])

    t = await ethers.getContract("T")
    walletRegistry = await ethers.getContract("WalletRegistry")
    sortitionPool = await ethers.getContract("SortitionPool")
    staking = await ethers.getContract("TokenStaking")

    deployer = await ethers.getNamedSigner("deployer")
    walletOwner = await ethers.getNamedSigner("walletOwner")

    const accounts = await getUnnamedAccounts()
    owner = await ethers.getSigner(accounts[1])
    stakingProvider = await ethers.getSigner(accounts[2])
    operator = await ethers.getSigner(accounts[3])
    authorizer = await ethers.getSigner(accounts[4])
    beneficiary = await ethers.getSigner(accounts[5])
    thirdParty = await ethers.getSigner(accounts[6])

    const governanceContract: WalletRegistryGovernance =
      await ethers.getContract("WalletRegistryGovernance")
    const governance = await ethers.getNamedSigner("governance")
    await updateWalletRegistryParams(governanceContract, governance)

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

    minimumAuthorization = await walletRegistry.minimumAuthorization()

    // Initialize slasher - fake application capable of slashing the
    // staking provider.
    slasher = await smock.fake<IApplication>("IApplication")
    await staking.connect(deployer).approveApplication(slasher.address)
    await staking
      .connect(authorizer)
      .increaseAuthorization(
        stakingProvider.address,
        slasher.address,
        stakedAmount
      )

    // Fund slasher so that it can call T TokenStaking functions
    await (
      await ethers.getSigners()
    )[0].sendTransaction({
      to: slasher.address,
      value: ethers.utils.parseEther("1"),
    })
  })

  describe("registerOperator", () => {
    context("when called with zero-address operator", () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(stakingProvider).registerOperator(ZERO_ADDRESS)
        ).to.be.revertedWith("Operator can not be zero address")
      })
    })

    // It is not possible to update operator address. Once the operator is
    // registered for the given staking provider, it must remain the same.
    // Staking provider address is unique for each stake delegation - see T
    // TokenStaking contract.
    context(
      "when operator has been already registered for the staking provider",
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

    // Some other staking provider is using the given operator address.
    // Should not happen in practice but we should protect against it.
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

    // This is the normal, happy path. Stake owner delegated their stake to
    // the staking provider, and the staking provider is registering operator
    // for ECDSA application.
    context("when staking provider is registering new operator", () => {
      let tx: ContractTransaction

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
    })

    // It is possible to approve authorization decrease request immediately
    // in case the operator was not yet registered by the staking provider.
    // It makes sense because non-registered operator could not be in the
    // sortition pool, so there is no state that could be not in sync.
    // However, we need to ensure this is not exploited by malicious stakers.
    // We do not want to let operators with a pending authorization decrease
    // request that can be immediately approved to join the sortition pool.
    // If there is a pending authorization decrease for the staking provider,
    // it must be first approved before operator for that staking provider is
    // registered.
    context("when there is a pending authorization decrease request", () => {
      before(async () => {
        await createSnapshot()

        await staking
          .connect(authorizer)
          .increaseAuthorization(
            stakingProvider.address,
            walletRegistry.address,
            stakedAmount
          )

        const deauthorizingBy = to1e18(1)

        await staking
          .connect(authorizer)
          ["requestAuthorizationDecrease(address,address,uint96)"](
            stakingProvider.address,
            walletRegistry.address,
            deauthorizingBy
          )
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should revert", async () => {
        await expect(
          walletRegistry
            .connect(stakingProvider)
            .registerOperator(operator.address)
        ).to.be.revertedWith(
          "There is a pending authorization decrease request"
        )
      })
    })

    // This is a continuation of the previous test case - in case there is
    // a staking provider who has not yet registered the operator and there is
    // an authorization decrease requested for that staking provider, upon
    // approving that authorization decrease request, staking provider can
    // register an operator.
    context("when authorization decrease request was approved", () => {
      let tx: ContractTransaction

      before(async () => {
        await createSnapshot()

        await staking
          .connect(authorizer)
          .increaseAuthorization(
            stakingProvider.address,
            walletRegistry.address,
            stakedAmount
          )

        const deauthorizingBy = to1e18(1)

        await staking
          .connect(authorizer)
          ["requestAuthorizationDecrease(address,address,uint96)"](
            stakingProvider.address,
            walletRegistry.address,
            deauthorizingBy
          )

        await walletRegistry.approveAuthorizationDecrease(
          stakingProvider.address
        )

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
    })
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

    context("when authorization is below the minimum", () => {
      it("should revert", async () => {
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

    // This is normal, happy path for a new delegation. Stake owner delegated
    // their stake to the staking provider, and while still being in the
    // dashboard (assuming staker is the authorizer), increased authorization
    // for ECDSA application. Staking provider has not registered operator yet.
    // This will happen later.
    context("when the operator is unknown", () => {
      // Minimum possible authorization - the minimum authorized amount for
      // ECDSA as set in `minimumAuthorization` parameter.
      context("when increasing to the minimum possible value", () => {
        let tx: ContractTransaction

        before(async () => {
          await createSnapshot()
          tx = await staking
            .connect(authorizer)
            .increaseAuthorization(
              stakingProvider.address,
              walletRegistry.address,
              minimumAuthorization
            )
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should emit AuthorizationIncreaseRequested", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "AuthorizationIncreaseRequested")
            .withArgs(
              stakingProvider.address,
              ZERO_ADDRESS,
              0,
              minimumAuthorization
            )
        })
      })

      // Maximum possible authorization - the entire stake delegated to the
      // staking provider.
      context("when increasing to the maximum possible value", () => {
        let tx: ContractTransaction

        before(async () => {
          await createSnapshot()
          tx = await staking
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

        it("should emit AuthorizationIncreaseRequested", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "AuthorizationIncreaseRequested")
            .withArgs(stakingProvider.address, ZERO_ADDRESS, 0, stakedAmount)
        })
      })
    })

    // This is normal, happy path for staking provider acting before the
    // authorizer, most probably because authorizer is someone else than the
    // stake owner. Stake owner delegated their stake to the staking provider,
    // staking provider registered operator for ECDSA, and after that, the
    // authorizer increased the authorization for the staking provider.
    context("when the operator is registered", () => {
      before(async () => {
        await createSnapshot()

        await walletRegistry
          .connect(stakingProvider)
          .registerOperator(operator.address)
      })

      after(async () => {
        await restoreSnapshot()
      })

      // Minimum possible authorization - the minimum authorized amount for
      // ECDSA as set in `minimumAuthorization` parameter.
      context("when increasing to the minimum possible value", () => {
        let tx: ContractTransaction

        before(async () => {
          await createSnapshot()

          tx = await staking
            .connect(authorizer)
            .increaseAuthorization(
              stakingProvider.address,
              walletRegistry.address,
              minimumAuthorization
            )
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should emit AuthorizationIncreaseRequested", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "AuthorizationIncreaseRequested")
            .withArgs(
              stakingProvider.address,
              operator.address,
              0,
              minimumAuthorization
            )
        })
      })

      // Maximum possible authorization - the entire stake delegated to the
      // staking provider.
      context("when increasing to the maximum possible value", () => {
        let tx: ContractTransaction

        before(async () => {
          await createSnapshot()

          tx = await staking
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

        it("should emit AuthorizationIncreaseRequested", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "AuthorizationIncreaseRequested")
            .withArgs(
              stakingProvider.address,
              operator.address,
              0,
              stakedAmount
            )
        })
      })
    })
  })

  describe("authorizationDecreaseRequested", () => {
    context("when called not by the staking contract", () => {
      it("should revert", async () => {
        await expect(
          walletRegistry
            .connect(thirdParty)
            .authorizationDecreaseRequested(stakingProvider.address, 100, 99)
        ).to.be.revertedWith("Caller is not the staking contract")
      })
    })

    // This is normal happy path in case the stake owner wants to cancel the
    // authorization before staking provider started their set up procedure.
    // Given the operator was not registered yet by the staking provider, we
    // can allow the authorization decrease to be processed immediately if it
    // is valid.
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

      // This is not valid authorization decrease request - one most decrease
      // to 0 or to some value above the minimum.
      context("when decreasing to a non-zero value below the minimum", () => {
        it("should revert", async () => {
          const deauthorizingTo = minimumAuthorization.sub(1)
          const deauthorizingBy = stakedAmount.sub(deauthorizingTo)

          await expect(
            staking
              .connect(authorizer)
              ["requestAuthorizationDecrease(address,address,uint96)"](
                stakingProvider.address,
                walletRegistry.address,
                deauthorizingBy
              )
          ).to.be.revertedWith(
            "Authorization amount should be 0 or above the minimum"
          )
        })
      })

      // Decreasing to zero when operator was not set up yet - authorization
      // decrease request is valid and can be approved
      context("when decreasing to zero", () => {
        let tx: ContractTransaction
        const decreasingTo = 0

        before(async () => {
          await createSnapshot()

          const decreasingBy = stakedAmount.sub(decreasingTo)
          tx = await staking
            .connect(authorizer)
            ["requestAuthorizationDecrease(address,address,uint96)"](
              stakingProvider.address,
              walletRegistry.address,
              decreasingBy
            )
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should require no time delay before approving", async () => {
          expect(
            await walletRegistry.remainingAuthorizationDecreaseDelay(
              stakingProvider.address
            )
          ).to.equal(0)
        })

        it("should emit AuthorizationDecreaseRequested event", async () => {
          const now = await helpers.time.lastBlockTime()
          await expect(tx)
            .to.emit(walletRegistry, "AuthorizationDecreaseRequested")
            .withArgs(
              stakingProvider.address,
              ZERO_ADDRESS,
              stakedAmount,
              decreasingTo,
              now
            )
        })
      })

      context("when decreasing to the minimum", () => {
        let tx: ContractTransaction
        let decreasingTo

        before(async () => {
          await createSnapshot()

          decreasingTo = minimumAuthorization
          const decreasingBy = stakedAmount.sub(decreasingTo)
          tx = await staking
            .connect(authorizer)
            ["requestAuthorizationDecrease(address,address,uint96)"](
              stakingProvider.address,
              walletRegistry.address,
              decreasingBy
            )
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should require no time delay before approving", async () => {
          expect(
            await walletRegistry.remainingAuthorizationDecreaseDelay(
              stakingProvider.address
            )
          ).to.equal(0)
        })

        it("should emit AuthorizationDecreaseRequested event", async () => {
          const now = await helpers.time.lastBlockTime()
          await expect(tx)
            .to.emit(walletRegistry, "AuthorizationDecreaseRequested")
            .withArgs(
              stakingProvider.address,
              ZERO_ADDRESS,
              stakedAmount,
              decreasingTo,
              now
            )
        })
      })

      context("when decreasing to a value above the minimum", () => {
        let tx: ContractTransaction
        let decreasingTo

        before(async () => {
          await createSnapshot()

          decreasingTo = minimumAuthorization.add(1)
          const decreasingBy = stakedAmount.sub(decreasingTo)
          tx = await staking
            .connect(authorizer)
            ["requestAuthorizationDecrease(address,address,uint96)"](
              stakingProvider.address,
              walletRegistry.address,
              decreasingBy
            )
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should require no time delay before approving", async () => {
          expect(
            await walletRegistry.remainingAuthorizationDecreaseDelay(
              stakingProvider.address
            )
          ).to.equal(0)
        })

        it("should emit AuthorizationDecreaseRequested event", async () => {
          const now = await helpers.time.lastBlockTime()
          await expect(tx)
            .to.emit(walletRegistry, "AuthorizationDecreaseRequested")
            .withArgs(
              stakingProvider.address,
              ZERO_ADDRESS,
              stakedAmount,
              decreasingTo,
              now
            )
        })
      })
    })

    // The most popular scenario - operator is registered, it has an
    // authorization and that authorization is decreased after some time.
    context("when the operator is registered", () => {
      before(async () => {
        await createSnapshot()
        await staking
          .connect(authorizer)
          .increaseAuthorization(
            stakingProvider.address,
            walletRegistry.address,
            stakedAmount
          )
        await walletRegistry
          .connect(stakingProvider)
          .registerOperator(operator.address)
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("when decreasing to a non-zero value below the minimum", () => {
        it("should revert", async () => {
          const deauthorizingTo = minimumAuthorization.sub(1)
          const deauthorizingBy = stakedAmount.sub(deauthorizingTo)

          await expect(
            staking
              .connect(authorizer)
              ["requestAuthorizationDecrease(address,address,uint96)"](
                stakingProvider.address,
                walletRegistry.address,
                deauthorizingBy
              )
          ).to.be.revertedWith(
            "Authorization amount should be 0 or above the minimum"
          )
        })
      })

      context("when decreasing to zero", () => {
        let tx: ContractTransaction
        const decreasingTo = 0

        before(async () => {
          await createSnapshot()

          const decreasingBy = stakedAmount.sub(decreasingTo)
          tx = await staking
            .connect(authorizer)
            ["requestAuthorizationDecrease(address,address,uint96)"](
              stakingProvider.address,
              walletRegistry.address,
              decreasingBy
            )
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should require updating the pool before approving", async () => {
          expect(
            await walletRegistry.remainingAuthorizationDecreaseDelay(
              stakingProvider.address
            )
          ).to.equal(MAX_UINT64)
        })

        it("should emit AuthorizationDecreaseRequested event", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "AuthorizationDecreaseRequested")
            .withArgs(
              stakingProvider.address,
              operator.address,
              stakedAmount,
              decreasingTo,
              MAX_UINT64
            )
        })
      })

      context("when decreasing to the minimum", () => {
        let tx: ContractTransaction
        let decreasingTo

        before(async () => {
          await createSnapshot()

          decreasingTo = minimumAuthorization
          const decreasingBy = stakedAmount.sub(decreasingTo)
          tx = await staking
            .connect(authorizer)
            ["requestAuthorizationDecrease(address,address,uint96)"](
              stakingProvider.address,
              walletRegistry.address,
              decreasingBy
            )
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should require updating the pool before approving", async () => {
          expect(
            await walletRegistry.remainingAuthorizationDecreaseDelay(
              stakingProvider.address
            )
          ).to.equal(MAX_UINT64)
        })

        it("should emit AuthorizationDecreaseRequested event", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "AuthorizationDecreaseRequested")
            .withArgs(
              stakingProvider.address,
              operator.address,
              stakedAmount,
              decreasingTo,
              MAX_UINT64
            )
        })
      })

      context("when decreasing to a value above the minimum", () => {
        let tx: ContractTransaction
        let decreasingTo

        before(async () => {
          await createSnapshot()

          decreasingTo = minimumAuthorization.add(1)
          const decreasingBy = stakedAmount.sub(decreasingTo)
          tx = await staking
            .connect(authorizer)
            ["requestAuthorizationDecrease(address,address,uint96)"](
              stakingProvider.address,
              walletRegistry.address,
              decreasingBy
            )
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should require updating the pool before approving", async () => {
          expect(
            await walletRegistry.remainingAuthorizationDecreaseDelay(
              stakingProvider.address
            )
          ).to.equal(MAX_UINT64)
        })

        it("should emit AuthorizationDecreaseRequested event", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "AuthorizationDecreaseRequested")
            .withArgs(
              stakingProvider.address,
              operator.address,
              stakedAmount,
              decreasingTo,
              MAX_UINT64
            )
        })
      })
    })
  })

  describe("approveAuthorizationDecrease", () => {
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

    context("when decrease was not requested", () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.approveAuthorizationDecrease(stakingProvider.address)
        ).to.be.revertedWith("Authorization decrease not requested")
      })
    })

    context("when the operator is unknown", () => {
      context("when the decrease was requested", () => {
        before(async () => {
          await createSnapshot()

          const deauthorizingBy = stakedAmount

          staking
            .connect(authorizer)
            ["requestAuthorizationDecrease(address,address,uint96)"](
              stakingProvider.address,
              walletRegistry.address,
              deauthorizingBy
            )
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should let to approve immediately", async () => {
          const tx = await walletRegistry.approveAuthorizationDecrease(
            stakingProvider.address
          )
          // ok, did not revert
          await expect(tx)
            .to.emit(walletRegistry, "AuthorizationDecreaseApproved")
            .withArgs(stakingProvider.address)
        })
      })
    })

    context("when the operator is registered", () => {
      before(async () => {
        await createSnapshot()

        await walletRegistry
          .connect(stakingProvider)
          .registerOperator(operator.address)

        const deauthorizingBy = stakedAmount
        staking
          .connect(authorizer)
          ["requestAuthorizationDecrease(address,address,uint96)"](
            stakingProvider.address,
            walletRegistry.address,
            deauthorizingBy
          )
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("when the pool was not updated", () => {
        before(async () => {
          await createSnapshot()

          // even if we wait for the entire delay, it should not help
          await helpers.time.increaseTime(params.authorizationDecreaseDelay)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should revert", async () => {
          await expect(
            walletRegistry.approveAuthorizationDecrease(stakingProvider.address)
          ).to.be.revertedWith("Authorization decrease request not activated")
        })
      })

      context("when the pool was updated but the delay did not pass", () => {
        before(async () => {
          await createSnapshot()

          await walletRegistry.updateOperatorStatus(operator.address)
          await helpers.time.increaseTime(params.authorizationDecreaseDelay - 1)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should revert", async () => {
          await expect(
            walletRegistry.approveAuthorizationDecrease(stakingProvider.address)
          ).to.be.revertedWith("Authorization decrease delay not passed")
        })
      })

      context("when the pool was updated and the delay passed", () => {
        let tx: ContractTransaction

        before(async () => {
          await createSnapshot()

          await walletRegistry.updateOperatorStatus(operator.address)
          await helpers.time.increaseTime(params.authorizationDecreaseDelay)

          tx = await walletRegistry.approveAuthorizationDecrease(
            stakingProvider.address
          )
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should reduce authorized stake amount", async () => {
          expect(await sortitionPool.getPoolWeight(operator.address)).to.equal(
            0
          )
        })

        it("should emit AuthorizationDecreaseApproved event", async () => {
          expect(tx)
            .to.emit(walletRegistry, "AuthorizationDecreaseApproved")
            .withArgs(stakingProvider.address)
        })
      })
    })
  })

  describe("involuntaryAuthorizationDecrease", () => {
    before(async () => {
      await createSnapshot()
      await staking
        .connect(authorizer)
        .increaseAuthorization(
          stakingProvider.address,
          walletRegistry.address,
          stakedAmount
        )

      await walletRegistry
        .connect(stakingProvider)
        .registerOperator(operator.address)
      await walletRegistry.connect(operator).joinSortitionPool()
    })

    after(async () => {
      await restoreSnapshot()
    })

    context("when called not by the staking contract", () => {
      it("should revert", async () => {
        await expect(
          walletRegistry
            .connect(thirdParty)
            .involuntaryAuthorizationDecrease(stakingProvider.address, 100, 99)
        ).to.be.revertedWith("Caller is not the staking contract")
      })
    })

    context("when the operator is in the sortition pool", () => {
      context("when the sortition pool is locked", () => {
        before(async () => {
          await createSnapshot()

          // lock the pool for DKG
          await walletRegistry.connect(walletOwner).requestNewWallet()

          // and slash!
          await staking
            .connect(slasher.wallet)
            .slash(to1e18(100), [stakingProvider.address])
          await staking.connect(thirdParty).processSlashing(1)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should ignore the update", async () => {
          expect(await walletRegistry.isOperatorUpToDate(operator.address)).to
            .be.false
        })
      })

      context("when the sortition pool is not locked", () => {
        before(async () => {
          await createSnapshot()

          // slash!
          await staking
            .connect(slasher.wallet)
            .slash(to1e18(100), [stakingProvider.address])
          await staking.connect(thirdParty).processSlashing(1)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update operator status", async () => {
          expect(await walletRegistry.isOperatorUpToDate(operator.address)).to
            .be.true
        })
      })
    })
  })

  describe("joinSortitionPool", () => {
    context("when the operator is unknown", () => {
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
      let tx: ContractTransaction

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
            minimumAuthorization
          )

        tx = await walletRegistry.connect(operator).joinSortitionPool()
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

      it("should emit OperatorJoinedSortitionPool", async () => {
        await expect(tx)
          .to.emit(walletRegistry, "OperatorJoinedSortitionPool")
          .withArgs(stakingProvider.address, operator.address)
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

    context("when operator is in the process of deauthorizing", () => {
      let deauthorizingTo

      before(async () => {
        await createSnapshot()

        await walletRegistry
          .connect(stakingProvider)
          .registerOperator(operator.address)

        const authorizedStake = stakedAmount

        await staking
          .connect(authorizer)
          .increaseAuthorization(
            stakingProvider.address,
            walletRegistry.address,
            authorizedStake
          )

        deauthorizingTo = minimumAuthorization.add(to1e18(1337))
        const deauthorizingBy = authorizedStake.sub(deauthorizingTo)

        await staking
          .connect(authorizer)
          ["requestAuthorizationDecrease(address,address,uint96)"](
            stakingProvider.address,
            walletRegistry.address,
            deauthorizingBy
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
          deauthorizingTo.div(constants.poolWeightDivisor)
        )
      })

      it("should activate authorization decrease delay", async () => {
        expect(
          await walletRegistry.remainingAuthorizationDecreaseDelay(
            stakingProvider.address
          )
        ).to.equal(params.authorizationDecreaseDelay)
      })
    })

    context(
      "when operator is in the process of deauthorizing but also increased authorization in the meantime",
      () => {
        let expectedAuthorizedStake

        before(async () => {
          await createSnapshot()

          await walletRegistry
            .connect(stakingProvider)
            .registerOperator(operator.address)

          const authorizedStake = minimumAuthorization.add(to1e18(100))

          await staking
            .connect(authorizer)
            .increaseAuthorization(
              stakingProvider.address,
              walletRegistry.address,
              authorizedStake
            )

          const deauthorizingTo = minimumAuthorization.add(to1e18(50))
          const deauthorizingBy = authorizedStake.sub(deauthorizingTo)

          await staking
            .connect(authorizer)
            ["requestAuthorizationDecrease(address,address,uint96)"](
              stakingProvider.address,
              walletRegistry.address,
              deauthorizingBy
            )

          const increasingBy = to1e18(5000)
          await staking
            .connect(authorizer)
            .increaseAuthorization(
              stakingProvider.address,
              walletRegistry.address,
              increasingBy
            )

          expectedAuthorizedStake = deauthorizingTo.add(increasingBy)

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
            expectedAuthorizedStake.div(constants.poolWeightDivisor)
          )
        })

        it("should activate authorization decrease delay", async () => {
          expect(
            await walletRegistry.remainingAuthorizationDecreaseDelay(
              stakingProvider.address
            )
          ).to.equal(params.authorizationDecreaseDelay)
        })
      }
    )
  })

  describe("updateOperatorStatus", () => {
    context("when the operator is unknown", () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.updateOperatorStatus(thirdParty.address)
        ).to.be.revertedWith("Unknown operator")
      })
    })

    context("when operator is not in the sortition pool", () => {
      before(async () => {
        await createSnapshot()

        await walletRegistry
          .connect(stakingProvider)
          .registerOperator(operator.address)
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("when the authorization increased", () => {
        let tx: ContractTransaction

        before(async () => {
          await createSnapshot()

          await staking
            .connect(authorizer)
            .increaseAuthorization(
              stakingProvider.address,
              walletRegistry.address,
              minimumAuthorization
            )

          tx = await walletRegistry
            .connect(thirdParty)
            .updateOperatorStatus(operator.address)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should not insert operator into the pool", async () => {
          expect(await sortitionPool.isOperatorInPool(operator.address)).to.be
            .false
        })

        it("should emit OperatorStatusUpdated", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "OperatorStatusUpdated")
            .withArgs(stakingProvider.address, operator.address)
        })
      })

      context("when there was an authorization decrease request", () => {
        let tx: ContractTransaction //

        before(async () => {
          await createSnapshot()

          await staking
            .connect(authorizer)
            .increaseAuthorization(
              stakingProvider.address,
              walletRegistry.address,
              stakedAmount
            )

          const deauthorizingBy = to1e18(100)
          await staking
            .connect(authorizer)
            ["requestAuthorizationDecrease(address,address,uint96)"](
              stakingProvider.address,
              walletRegistry.address,
              deauthorizingBy
            )

          tx = await walletRegistry
            .connect(thirdParty)
            .updateOperatorStatus(operator.address)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should not insert operator into the pool", async () => {
          expect(await sortitionPool.isOperatorInPool(operator.address)).to.be
            .false
        })

        it("should activate authorization decrease delay", async () => {
          expect(
            await walletRegistry.remainingAuthorizationDecreaseDelay(
              stakingProvider.address
            )
          ).to.equal(params.authorizationDecreaseDelay)
        })

        it("should emit OperatorStatusUpdated", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "OperatorStatusUpdated")
            .withArgs(stakingProvider.address, operator.address)
        })
      })
    })

    context("when operator is in the sortition pool", () => {
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
            minimumAuthorization.mul(2)
          )

        await walletRegistry.connect(operator).joinSortitionPool()
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("when the authorization increased", () => {
        let tx: ContractTransaction
        let expectedWeight

        before(async () => {
          await createSnapshot()

          const topUp = to1e18(1337)
          await staking
            .connect(authorizer)
            .increaseAuthorization(
              stakingProvider.address,
              walletRegistry.address,
              topUp
            )

          // initial authorization was 2 x minimum
          // it was increased by 1337 tokens
          // so the final authorization should be 2 x minimum + 1337
          expectedWeight = minimumAuthorization
            .mul(2)
            .add(topUp)
            .div(constants.poolWeightDivisor)

          tx = await walletRegistry
            .connect(thirdParty)
            .updateOperatorStatus(operator.address)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the pool", async () => {
          expect(await sortitionPool.getPoolWeight(operator.address)).to.equal(
            expectedWeight
          )
        })

        it("should emit OperatorStatusUpdated", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "OperatorStatusUpdated")
            .withArgs(stakingProvider.address, operator.address)
        })
      })

      context("when there was an authorization decrease request", () => {
        let tx: ContractTransaction
        let expectedWeight

        before(async () => {
          await createSnapshot()

          // initial authorization was 2 x minimum
          // we want to decrease to minimum + 1337
          const deauthorizingTo = minimumAuthorization.add(to1e18(1337))
          const deauthorizingBy = minimumAuthorization
            .mul(2)
            .sub(deauthorizingTo)
          expectedWeight = deauthorizingTo.div(constants.poolWeightDivisor)

          await staking
            .connect(authorizer)
            ["requestAuthorizationDecrease(address,address,uint96)"](
              stakingProvider.address,
              walletRegistry.address,
              deauthorizingBy
            )

          tx = await walletRegistry
            .connect(thirdParty)
            .updateOperatorStatus(operator.address)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the pool", async () => {
          expect(await sortitionPool.getPoolWeight(operator.address)).to.equal(
            expectedWeight
          )
        })

        it("should activate authorization decrease delay", async () => {
          expect(
            await walletRegistry.remainingAuthorizationDecreaseDelay(
              stakingProvider.address
            )
          ).to.equal(params.authorizationDecreaseDelay)
        })

        it("should emit OperatorStatusUpdated", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "OperatorStatusUpdated")
            .withArgs(stakingProvider.address, operator.address)
        })
      })

      context(
        "when operator is in the process of deauthorizing but also increased authorization in the meantime",
        () => {
          let tx: ContractTransaction
          let expectedWeight

          before(async () => {
            await createSnapshot()

            // initial authorization was 2 x minimum
            // we want to decrease to minimum + 1337
            // and then decrease by 7331
            const deauthorizingTo = minimumAuthorization.add(to1e18(1337))
            const deauthorizingBy = minimumAuthorization
              .mul(2)
              .sub(deauthorizingTo)
            const increasingBy = to1e18(7331)
            const increasingTo = deauthorizingTo.add(increasingBy)
            expectedWeight = increasingTo.div(constants.poolWeightDivisor)

            await staking
              .connect(authorizer)
              ["requestAuthorizationDecrease(address,address,uint96)"](
                stakingProvider.address,
                walletRegistry.address,
                deauthorizingBy
              )

            await staking
              .connect(authorizer)
              .increaseAuthorization(
                stakingProvider.address,
                walletRegistry.address,
                increasingBy
              )

            tx = await walletRegistry
              .connect(thirdParty)
              .updateOperatorStatus(operator.address)
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should update the pool", async () => {
            expect(
              await sortitionPool.getPoolWeight(operator.address)
            ).to.equal(expectedWeight)
          })

          it("should activate authorization decrease delay", async () => {
            expect(
              await walletRegistry.remainingAuthorizationDecreaseDelay(
                stakingProvider.address
              )
            ).to.equal(params.authorizationDecreaseDelay)
          })

          it("should emit OperatorStatusUpdated", async () => {
            await expect(tx)
              .to.emit(walletRegistry, "OperatorStatusUpdated")
              .withArgs(stakingProvider.address, operator.address)
          })
        }
      )
    })
  })

  describe("isOperatorUpToDate", () => {
    context("when the operator is unknown", () => {
      it("should revert", async () => {
        it("should revert", async () => {
          await expect(
            walletRegistry.isOperatorUpToDate(thirdParty.address)
          ).to.be.revertedWith("Unknown operator")
        })
      })
    })

    context("when the operator is not in the sortition pool", () => {
      before(async () => {
        await createSnapshot()

        await walletRegistry
          .connect(stakingProvider)
          .registerOperator(operator.address)
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("when the operator has no authorized stake", () => {
        it("should return true", async () => {
          expect(await walletRegistry.isOperatorUpToDate(operator.address)).to
            .be.true
        })
      })

      context("when the operator has authorized stake", () => {
        before(async () => {
          await createSnapshot()

          await staking
            .connect(authorizer)
            .increaseAuthorization(
              stakingProvider.address,
              walletRegistry.address,
              minimumAuthorization
            )
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should return false", async () => {
          expect(await walletRegistry.isOperatorUpToDate(operator.address)).to
            .be.false
        })
      })
    })

    context("when the operator is in the sortition pool", () => {
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
            minimumAuthorization.mul(2)
          )

        await walletRegistry.connect(operator).joinSortitionPool()
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("when the operator just joined the pool", () => {
        it("should return true", async () => {
          expect(await walletRegistry.isOperatorUpToDate(operator.address)).to
            .be.true
        })
      })

      context("when authorization was increased", () => {
        before(async () => {
          await createSnapshot()

          await staking
            .connect(authorizer)
            .increaseAuthorization(
              stakingProvider.address,
              walletRegistry.address,
              to1e18(1337)
            )
        })

        after(async () => {
          await restoreSnapshot()
        })

        context("when sortition pool was not yet updated", () => {
          it("should return false", async () => {
            expect(await walletRegistry.isOperatorUpToDate(operator.address)).to
              .be.false
          })
        })

        context("when the sortition pool was updated", () => {
          it("should return true", async () => {
            await walletRegistry.updateOperatorStatus(operator.address)
            expect(await walletRegistry.isOperatorUpToDate(operator.address)).to
              .be.true
          })
        })
      })

      context("when authorization decrease was requested", () => {
        before(async () => {
          await createSnapshot()

          const deauthorizingBy = to1e18(1)
          await staking
            .connect(authorizer)
            ["requestAuthorizationDecrease(address,address,uint96)"](
              stakingProvider.address,
              walletRegistry.address,
              deauthorizingBy
            )
        })

        after(async () => {
          await restoreSnapshot()
        })

        context("when sortition pool was not yet updated", () => {
          it("should return false", async () => {
            expect(await walletRegistry.isOperatorUpToDate(operator.address)).to
              .be.false
          })
        })

        context("when the sortition pool was updated", () => {
          it("should return true", async () => {
            await walletRegistry.updateOperatorStatus(operator.address)
            expect(await walletRegistry.isOperatorUpToDate(operator.address)).to
              .be.true
          })
        })
      })

      context("when operator was slashed when the pool was locked", () => {
        before(async () => {
          await createSnapshot()

          // Increase authorization to the maximum possible value and update
          // sortition pool. This way, any slashing from `slasher` application
          // will affect authorized stake amount for WalletRegistry.
          const authorized = await staking.authorizedStake(
            stakingProvider.address,
            walletRegistry.address
          )
          const increaseBy = stakedAmount.sub(authorized)
          await staking
            .connect(authorizer)
            .increaseAuthorization(
              stakingProvider.address,
              walletRegistry.address,
              increaseBy
            )
          await walletRegistry.updateOperatorStatus(operator.address)

          // lock the pool for DKG
          await walletRegistry.connect(walletOwner).requestNewWallet()

          // and slash!
          await staking
            .connect(slasher.wallet)
            .slash(to1e18(100), [stakingProvider.address])
          await staking.connect(thirdParty).processSlashing(1)

          // unlock the pool by stopping DKG
          await mineBlocks(params.dkgSeedTimeout)
          await walletRegistry.notifySeedTimeout()
        })

        after(async () => {
          await restoreSnapshot()
        })

        context("when sortition pool was not yet updated", () => {
          it("should return false", async () => {
            expect(await walletRegistry.isOperatorUpToDate(operator.address)).to
              .be.false
          })
        })

        context("when the sortition pool was updated", () => {
          it("should return true", async () => {
            await walletRegistry.updateOperatorStatus(operator.address)
            expect(await walletRegistry.isOperatorUpToDate(operator.address)).to
              .be.true
          })
        })
      })
    })
  })
})
