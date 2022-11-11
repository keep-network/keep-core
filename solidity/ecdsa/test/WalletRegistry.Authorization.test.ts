/* eslint-disable @typescript-eslint/no-unused-expressions */
import { deployments, ethers, getUnnamedAccounts, helpers } from "hardhat"
import { smock } from "@defi-wonderland/smock"
import { expect } from "chai"

import {
  constants,
  params,
  initializeWalletOwner,
  updateWalletRegistryParams,
} from "./fixtures"

import type { IWalletOwner } from "../typechain/IWalletOwner"
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
const { to1e18 } = helpers.number

const { createSnapshot, restoreSnapshot } = helpers.snapshot

const ZERO_ADDRESS = ethers.constants.AddressZero
const MAX_UINT64 = ethers.BigNumber.from("18446744073709551615") // 2^64 - 1

describe("WalletRegistry - Authorization", () => {
  let t: T
  let walletRegistry: WalletRegistry
  let walletRegistryGovernance: WalletRegistryGovernance
  let sortitionPool: SortitionPool
  let staking: TokenStaking

  let deployer: SignerWithAddress
  let governance: SignerWithAddress

  let owner: SignerWithAddress
  let stakingProvider: SignerWithAddress
  let operator: SignerWithAddress
  let authorizer: SignerWithAddress
  let beneficiary: SignerWithAddress
  let thirdParty: SignerWithAddress
  let walletOwner: FakeContract<IWalletOwner>
  let slasher: FakeContract<IApplication>

  const stakedAmount = to1e18(1000000) // 1M T
  let minimumAuthorization

  before("load test fixture", async () => {
    await deployments.fixture(["WalletRegistry"])

    t = await helpers.contracts.getContract("T")
    walletRegistry = await helpers.contracts.getContract("WalletRegistry")
    walletRegistryGovernance = await helpers.contracts.getContract(
      "WalletRegistryGovernance"
    )
    sortitionPool = await helpers.contracts.getContract("EcdsaSortitionPool")
    staking = await helpers.contracts.getContract("TokenStaking")

    const accounts = await getUnnamedAccounts()
    owner = await ethers.getSigner(accounts[1])
    stakingProvider = await ethers.getSigner(accounts[2])
    operator = await ethers.getSigner(accounts[3])
    authorizer = await ethers.getSigner(accounts[4])
    beneficiary = await ethers.getSigner(accounts[5])
    thirdParty = await ethers.getSigner(accounts[6])
    ;({ deployer, governance } = await helpers.signers.getNamedSigners())

    const { chaosnetOwner } = await helpers.signers.getNamedSigners()
    await sortitionPool.connect(chaosnetOwner).deactivateChaosnet()

    walletOwner = await initializeWalletOwner(
      walletRegistryGovernance,
      governance
    )

    await updateWalletRegistryParams(walletRegistryGovernance, governance)

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
      value: ethers.utils.parseEther("100"),
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

          // should revert even if it's another operator than the one previously
          // registered for the staking provider
          await expect(
            walletRegistry
              .connect(stakingProvider)
              .registerOperator(thirdParty.address)
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

      it("should not register operator in the pool", async () => {
        expect(await walletRegistry.isOperatorInPool(operator.address)).to.be
          .false
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

        it("should emit AuthorizationIncreased", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "AuthorizationIncreased")
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

        it("should emit AuthorizationIncreased", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "AuthorizationIncreased")
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

        it("should emit AuthorizationIncreased", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "AuthorizationIncreased")
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

        it("should emit AuthorizationIncreased", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "AuthorizationIncreased")
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
        let decreasingBy

        before(async () => {
          await createSnapshot()

          decreasingBy = stakedAmount.sub(decreasingTo)
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

        it("should capture deauthorizing amount", async () => {
          expect(
            await walletRegistry.pendingAuthorizationDecrease(
              stakingProvider.address
            )
          ).to.equal(decreasingBy)
        })
      })

      context("when decreasing to the minimum", () => {
        let tx: ContractTransaction
        let decreasingTo
        let decreasingBy

        before(async () => {
          await createSnapshot()

          decreasingTo = minimumAuthorization
          decreasingBy = stakedAmount.sub(decreasingTo)
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

        it("should capture deauthorizing amount", async () => {
          expect(
            await walletRegistry.pendingAuthorizationDecrease(
              stakingProvider.address
            )
          ).to.equal(decreasingBy)
        })
      })

      context("when decreasing to a value above the minimum", () => {
        let tx: ContractTransaction
        let decreasingTo
        let decreasingBy

        before(async () => {
          await createSnapshot()

          decreasingTo = minimumAuthorization.add(1)
          decreasingBy = stakedAmount.sub(decreasingTo)
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

        it("should capture deauthorizing amount", async () => {
          expect(
            await walletRegistry.pendingAuthorizationDecrease(
              stakingProvider.address
            )
          ).to.equal(decreasingBy)
        })
      })

      context("when called one more time", () => {
        const deauthorizingFirst = to1e18(10)
        const deauthorizingSecond = to1e18(20)

        before(async () => {
          await createSnapshot()

          await staking
            .connect(authorizer)
            ["requestAuthorizationDecrease(address,address,uint96)"](
              stakingProvider.address,
              walletRegistry.address,
              deauthorizingFirst
            )
        })

        after(async () => {
          await restoreSnapshot()
        })

        context("when change period is equal delay", () => {
          before(async () => {
            // this should be the default situation from the fixture setup so we
            // just confirm it here
            const {
              authorizationDecreaseDelay,
              authorizationDecreaseChangePeriod,
            } = await walletRegistry.authorizationParameters()
            expect(authorizationDecreaseDelay).to.equal(
              authorizationDecreaseChangePeriod
            )
          })

          context("when delay passed", () => {
            before(async () => {
              await createSnapshot()
              await helpers.time.increaseTime(params.authorizationDecreaseDelay)

              await staking
                .connect(authorizer)
                ["requestAuthorizationDecrease(address,address,uint96)"](
                  stakingProvider.address,
                  walletRegistry.address,
                  deauthorizingSecond
                )
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should overwrite the previous request", async () => {
              expect(
                await walletRegistry.pendingAuthorizationDecrease(
                  stakingProvider.address
                )
              ).to.be.equal(deauthorizingSecond)
            })
          })

          context("when delay did not pass", () => {
            before(async () => {
              await createSnapshot()
              await helpers.time.increaseTime(
                params.authorizationDecreaseDelay - 60 // -1min
              )

              await staking
                .connect(authorizer)
                ["requestAuthorizationDecrease(address,address,uint96)"](
                  stakingProvider.address,
                  walletRegistry.address,
                  deauthorizingSecond
                )
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should overwrite the previous request", async () => {
              expect(
                await walletRegistry.pendingAuthorizationDecrease(
                  stakingProvider.address
                )
              ).to.be.equal(deauthorizingSecond)
            })
          })
        })

        context("when change period is zero", () => {
          before(async () => {
            await createSnapshot()

            await walletRegistryGovernance
              .connect(governance)
              .beginAuthorizationDecreaseChangePeriodUpdate(0)
            await helpers.time.increaseTime(constants.governanceDelay)
            await walletRegistryGovernance
              .connect(governance)
              .finalizeAuthorizationDecreaseChangePeriodUpdate()

            await staking
              .connect(authorizer)
              ["requestAuthorizationDecrease(address,address,uint96)"](
                stakingProvider.address,
                walletRegistry.address,
                deauthorizingSecond
              )
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should overwrite the previous request", async () => {
            expect(
              await walletRegistry.pendingAuthorizationDecrease(
                stakingProvider.address
              )
            ).to.be.equal(deauthorizingSecond)
          })
        })

        context("when change period is not equal delay and is non-zero", () => {
          const newChangePeriod = 3600 // 1h before delay end

          before(async () => {
            await createSnapshot()

            await walletRegistryGovernance
              .connect(governance)
              .beginAuthorizationDecreaseChangePeriodUpdate(newChangePeriod)
            await helpers.time.increaseTime(constants.governanceDelay)
            await walletRegistryGovernance
              .connect(governance)
              .finalizeAuthorizationDecreaseChangePeriodUpdate()
          })

          after(async () => {
            await restoreSnapshot()
          })

          context("when delay passed", () => {
            before(async () => {
              await createSnapshot()
              await helpers.time.increaseTime(params.authorizationDecreaseDelay)

              await staking
                .connect(authorizer)
                ["requestAuthorizationDecrease(address,address,uint96)"](
                  stakingProvider.address,
                  walletRegistry.address,
                  deauthorizingSecond
                )
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should overwrite the previous request", async () => {
              expect(
                await walletRegistry.pendingAuthorizationDecrease(
                  stakingProvider.address
                )
              ).to.be.equal(deauthorizingSecond)
            })
          })

          context("when change period activated", () => {
            before(async () => {
              await createSnapshot()
              await helpers.time.increaseTime(
                params.authorizationDecreaseDelay - newChangePeriod + 60
              ) // +1min

              await staking
                .connect(authorizer)
                ["requestAuthorizationDecrease(address,address,uint96)"](
                  stakingProvider.address,
                  walletRegistry.address,
                  deauthorizingSecond
                )
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should overwrite the previous request", async () => {
              expect(
                await walletRegistry.pendingAuthorizationDecrease(
                  stakingProvider.address
                )
              ).to.be.equal(deauthorizingSecond)
            })
          })

          context("when change period did not activate", () => {
            before(async () => {
              await createSnapshot()
              await helpers.time.increaseTime(
                params.authorizationDecreaseDelay - newChangePeriod - 60 // -1min
              )

              await staking
                .connect(authorizer)
                ["requestAuthorizationDecrease(address,address,uint96)"](
                  stakingProvider.address,
                  walletRegistry.address,
                  deauthorizingSecond
                )
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should overwrite the previous request", async () => {
              expect(
                await walletRegistry.pendingAuthorizationDecrease(
                  stakingProvider.address
                )
              ).to.be.equal(deauthorizingSecond)
            })
          })
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
        let decreasingBy

        before(async () => {
          await createSnapshot()

          decreasingBy = stakedAmount.sub(decreasingTo)
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

        it("should capture deauthorizing amount", async () => {
          expect(
            await walletRegistry.pendingAuthorizationDecrease(
              stakingProvider.address
            )
          ).to.equal(decreasingBy)
        })
      })

      context("when decreasing to the minimum", () => {
        let tx: ContractTransaction
        let decreasingTo
        let decreasingBy

        before(async () => {
          await createSnapshot()

          decreasingTo = minimumAuthorization
          decreasingBy = stakedAmount.sub(decreasingTo)
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

        it("should capture deauthorizing amount", async () => {
          expect(
            await walletRegistry.pendingAuthorizationDecrease(
              stakingProvider.address
            )
          ).to.equal(decreasingBy)
        })
      })

      context("when decreasing to a value above the minimum", () => {
        let tx: ContractTransaction
        let decreasingTo
        let decreasingBy

        before(async () => {
          await createSnapshot()

          decreasingTo = minimumAuthorization.add(1)
          decreasingBy = stakedAmount.sub(decreasingTo)
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

        it("should capture deauthorizing amount", async () => {
          expect(
            await walletRegistry.pendingAuthorizationDecrease(
              stakingProvider.address
            )
          ).to.equal(decreasingBy)
        })
      })

      context("when called one more time", () => {
        const deauthorizingFirst = to1e18(11)
        const deauthorizingSecond = to1e18(21)

        before(async () => {
          await createSnapshot()

          await walletRegistry.connect(operator).joinSortitionPool()

          await staking
            .connect(authorizer)
            ["requestAuthorizationDecrease(address,address,uint96)"](
              stakingProvider.address,
              walletRegistry.address,
              deauthorizingFirst
            )
        })

        after(async () => {
          await restoreSnapshot()
        })

        context("when change period is equal delay", () => {
          before(async () => {
            // this should be the default situation from the fixture setup so we
            // just confirm it here
            const {
              authorizationDecreaseDelay,
              authorizationDecreaseChangePeriod,
            } = await walletRegistry.authorizationParameters()
            expect(authorizationDecreaseDelay).to.equal(
              authorizationDecreaseChangePeriod
            )
          })

          context("when called before sortition pool was updated", () => {
            before(async () => {
              await createSnapshot()

              await staking
                .connect(authorizer)
                ["requestAuthorizationDecrease(address,address,uint96)"](
                  stakingProvider.address,
                  walletRegistry.address,
                  deauthorizingSecond
                )
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should overwrite the previous request", async () => {
              expect(
                await walletRegistry.pendingAuthorizationDecrease(
                  stakingProvider.address
                )
              ).to.be.equal(deauthorizingSecond)
            })

            it("should require updating the pool before approving", async () => {
              expect(
                await walletRegistry.remainingAuthorizationDecreaseDelay(
                  stakingProvider.address
                )
              ).to.equal(MAX_UINT64)
            })
          })

          context("when called after sortition pool was updated", () => {
            before(async () => {
              await createSnapshot()
              await walletRegistry.updateOperatorStatus(operator.address)
            })

            after(async () => {
              await restoreSnapshot()
            })

            context("when delay passed", () => {
              before(async () => {
                await createSnapshot()
                await helpers.time.increaseTime(
                  params.authorizationDecreaseDelay
                )
              })

              after(async () => {
                await restoreSnapshot()
              })

              before(async () => {
                await createSnapshot()

                await staking
                  .connect(authorizer)
                  ["requestAuthorizationDecrease(address,address,uint96)"](
                    stakingProvider.address,
                    walletRegistry.address,
                    deauthorizingSecond
                  )
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should overwrite the previous request", async () => {
                expect(
                  await walletRegistry.pendingAuthorizationDecrease(
                    stakingProvider.address
                  )
                ).to.be.equal(deauthorizingSecond)
              })

              it("should require updating the pool before approving", async () => {
                expect(
                  await walletRegistry.remainingAuthorizationDecreaseDelay(
                    stakingProvider.address
                  )
                ).to.equal(MAX_UINT64)
              })
            })

            context("when delay did not pass", () => {
              before(async () => {
                await createSnapshot()

                await helpers.time.increaseTime(
                  params.authorizationDecreaseDelay - 60 // -1min
                )

                await staking
                  .connect(authorizer)
                  ["requestAuthorizationDecrease(address,address,uint96)"](
                    stakingProvider.address,
                    walletRegistry.address,
                    deauthorizingSecond
                  )
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should overwrite the previous request", async () => {
                expect(
                  await walletRegistry.pendingAuthorizationDecrease(
                    stakingProvider.address
                  )
                ).to.be.equal(deauthorizingSecond)
              })

              it("should require updating the pool before approving", async () => {
                expect(
                  await walletRegistry.remainingAuthorizationDecreaseDelay(
                    stakingProvider.address
                  )
                ).to.equal(MAX_UINT64)
              })
            })
          })
        })

        context("when change period is zero", () => {
          before(async () => {
            await createSnapshot()

            await walletRegistryGovernance
              .connect(governance)
              .beginAuthorizationDecreaseChangePeriodUpdate(0)
            await helpers.time.increaseTime(constants.governanceDelay)
            await walletRegistryGovernance
              .connect(governance)
              .finalizeAuthorizationDecreaseChangePeriodUpdate()
          })

          after(async () => {
            await restoreSnapshot()
          })

          context("when called before sortition pool was updated", () => {
            before(async () => {
              await createSnapshot()

              await staking
                .connect(authorizer)
                ["requestAuthorizationDecrease(address,address,uint96)"](
                  stakingProvider.address,
                  walletRegistry.address,
                  deauthorizingSecond
                )
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should overwrite the previous request", async () => {
              expect(
                await walletRegistry.pendingAuthorizationDecrease(
                  stakingProvider.address
                )
              ).to.be.equal(deauthorizingSecond)
            })

            it("should require updating the pool before approving", async () => {
              expect(
                await walletRegistry.remainingAuthorizationDecreaseDelay(
                  stakingProvider.address
                )
              ).to.equal(MAX_UINT64)
            })
          })

          context("when called after sortition pool was updated", () => {
            before(async () => {
              await createSnapshot()

              await walletRegistry.updateOperatorStatus(operator.address)
            })

            after(async () => {
              await restoreSnapshot()
            })

            context("when called before delay passed", () => {
              it("should revert", async () => {
                await expect(
                  staking
                    .connect(authorizer)
                    ["requestAuthorizationDecrease(address,address,uint96)"](
                      stakingProvider.address,
                      walletRegistry.address,
                      deauthorizingSecond
                    )
                ).to.be.revertedWith(
                  "Not enough time passed since the original request"
                )
              })
            })

            context("when called after delay passed", () => {
              before(async () => {
                await createSnapshot()
                await helpers.time.increaseTime(
                  params.authorizationDecreaseDelay
                )

                await staking
                  .connect(authorizer)
                  ["requestAuthorizationDecrease(address,address,uint96)"](
                    stakingProvider.address,
                    walletRegistry.address,
                    deauthorizingSecond
                  )
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should overwrite the previous request", async () => {
                expect(
                  await walletRegistry.pendingAuthorizationDecrease(
                    stakingProvider.address
                  )
                ).to.be.equal(deauthorizingSecond)
              })

              it("should require updating the pool before approving", async () => {
                expect(
                  await walletRegistry.remainingAuthorizationDecreaseDelay(
                    stakingProvider.address
                  )
                ).to.equal(MAX_UINT64)
              })
            })
          })
        })

        context("when change period is not equal delay and is non-zero", () => {
          const newChangePeriod = 3600 // 1h before delay end

          before(async () => {
            await createSnapshot()

            await walletRegistryGovernance
              .connect(governance)
              .beginAuthorizationDecreaseChangePeriodUpdate(newChangePeriod)
            await helpers.time.increaseTime(constants.governanceDelay)
            await walletRegistryGovernance
              .connect(governance)
              .finalizeAuthorizationDecreaseChangePeriodUpdate()
          })

          after(async () => {
            await restoreSnapshot()
          })

          context("when called before sortition pool was updated", () => {
            before(async () => {
              await createSnapshot()

              await staking
                .connect(authorizer)
                ["requestAuthorizationDecrease(address,address,uint96)"](
                  stakingProvider.address,
                  walletRegistry.address,
                  deauthorizingSecond
                )
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should overwrite the previous request", async () => {
              expect(
                await walletRegistry.pendingAuthorizationDecrease(
                  stakingProvider.address
                )
              ).to.be.equal(deauthorizingSecond)
            })

            it("should require updating the pool before approving", async () => {
              expect(
                await walletRegistry.remainingAuthorizationDecreaseDelay(
                  stakingProvider.address
                )
              ).to.equal(MAX_UINT64)
            })
          })

          context("when called after sortition pool was updated", () => {
            before(async () => {
              await createSnapshot()

              await walletRegistry.updateOperatorStatus(operator.address)
            })

            after(async () => {
              await restoreSnapshot()
            })

            context("when change period did not activate", () => {
              before(async () => {
                await createSnapshot()
                await helpers.time.increaseTime(
                  params.authorizationDecreaseDelay - newChangePeriod - 60 // -1min
                )
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should revert", async () => {
                await expect(
                  staking
                    .connect(authorizer)
                    ["requestAuthorizationDecrease(address,address,uint96)"](
                      stakingProvider.address,
                      walletRegistry.address,
                      deauthorizingSecond
                    )
                ).to.be.revertedWith(
                  "Not enough time passed since the original request"
                )
              })
            })

            context("when change period did activate", () => {
              before(async () => {
                await createSnapshot()
                await helpers.time.increaseTime(
                  params.authorizationDecreaseDelay - newChangePeriod + 60 // +1min
                )

                await staking
                  .connect(authorizer)
                  ["requestAuthorizationDecrease(address,address,uint96)"](
                    stakingProvider.address,
                    walletRegistry.address,
                    deauthorizingSecond
                  )
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should overwrite the previous request", async () => {
                expect(
                  await walletRegistry.pendingAuthorizationDecrease(
                    stakingProvider.address
                  )
                ).to.be.equal(deauthorizingSecond)
              })

              it("should require updating the pool before approving", async () => {
                expect(
                  await walletRegistry.remainingAuthorizationDecreaseDelay(
                    stakingProvider.address
                  )
                ).to.equal(MAX_UINT64)
              })
            })

            context("when delay passed", () => {
              before(async () => {
                await createSnapshot()
                await helpers.time.increaseTime(
                  params.authorizationDecreaseDelay
                )

                await staking
                  .connect(authorizer)
                  ["requestAuthorizationDecrease(address,address,uint96)"](
                    stakingProvider.address,
                    walletRegistry.address,
                    deauthorizingSecond
                  )
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should overwrite the previous request", async () => {
                expect(
                  await walletRegistry.pendingAuthorizationDecrease(
                    stakingProvider.address
                  )
                ).to.be.equal(deauthorizingSecond)
              })

              it("should require updating the pool before approving", async () => {
                expect(
                  await walletRegistry.remainingAuthorizationDecreaseDelay(
                    stakingProvider.address
                  )
                ).to.equal(MAX_UINT64)
              })
            })
          })
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
          await helpers.time.increaseTime(
            params.authorizationDecreaseDelay - 60 // -1min
          )
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

        it("should clear pending authorization decrease", async () => {
          expect(
            await walletRegistry.pendingAuthorizationDecrease(
              stakingProvider.address
            )
          ).to.equal(0)
        })
      })
    })
  })

  describe("involuntaryAuthorizationDecrease", () => {
    context("when called not by the staking contract", () => {
      it("should revert", async () => {
        await expect(
          walletRegistry
            .connect(thirdParty)
            .involuntaryAuthorizationDecrease(stakingProvider.address, 100, 99)
        ).to.be.revertedWith("Caller is not the staking contract")
      })
    })

    context("when the operator is unknown", () => {
      const slashedAmount = to1e18(100)
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

        // lock the pool for DKG
        // we lock the pool to ensure that the update is ignored for the
        // operator and that involuntaryAuthorizationDecrease logic in this
        // case is basically a pass-through
        await walletRegistry.connect(walletOwner.wallet).requestNewWallet()

        // slash!
        await staking
          .connect(slasher.wallet)
          .slash(slashedAmount, [stakingProvider.address])
        tx = await staking.connect(thirdParty).processSlashing(1)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should ignore the update", async () => {
        await expect(tx).to.not.emit(
          walletRegistry,
          "InvoluntaryAuthorizationDecreaseFailed"
        )
      })
    })

    context("when the operator is known", () => {
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

      context("when the operator is not in the sortition pool", () => {
        const slashedAmount = to1e18(100)
        let tx: ContractTransaction

        before(async () => {
          await createSnapshot()

          // lock the pool for DKG
          // we lock the pool to ensure that the update is ignored for the
          // operator and that involuntaryAuthorizationDecrease logic in this
          // case is basically a pass-through
          await walletRegistry.connect(walletOwner.wallet).requestNewWallet()

          // slash!
          await staking
            .connect(slasher.wallet)
            .slash(slashedAmount, [stakingProvider.address])
          tx = await staking.connect(thirdParty).processSlashing(1)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should ignore the update", async () => {
          await expect(tx).to.not.emit(
            walletRegistry,
            "InvoluntaryAuthorizationDecreaseFailed"
          )
        })
      })

      context("when the operator is in the sortition pool", () => {
        before(async () => {
          await createSnapshot()
          await walletRegistry.connect(operator).joinSortitionPool()
        })

        after(async () => {
          await restoreSnapshot()
        })

        context("when the sortition pool is locked", () => {
          const slashedAmount = to1e18(100)
          let tx: ContractTransaction

          before(async () => {
            await createSnapshot()

            // lock the pool for DKG
            await walletRegistry.connect(walletOwner.wallet).requestNewWallet()

            // and slash!
            await staking
              .connect(slasher.wallet)
              .slash(slashedAmount, [stakingProvider.address])
            tx = await staking.connect(thirdParty).processSlashing(1)
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should not update the pool", async () => {
            expect(await walletRegistry.isOperatorUpToDate(operator.address)).to
              .be.false
          })

          it("should emit InvoluntaryAuthorizationDecreaseFailed event", async () => {
            await expect(tx)
              .to.emit(walletRegistry, "InvoluntaryAuthorizationDecreaseFailed")
              .withArgs(
                stakingProvider.address,
                operator.address,
                stakedAmount,
                stakedAmount.sub(slashedAmount)
              )
          })
        })

        context("when the sortition pool is not locked", () => {
          context("when the authorization drops to above the minimum", () => {
            const slashedAmount = to1e18(100)
            let tx: ContractTransaction

            before(async () => {
              await createSnapshot()

              // slash!
              await staking
                .connect(slasher.wallet)
                .slash(slashedAmount, [stakingProvider.address])
              tx = await staking.connect(thirdParty).processSlashing(1)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should update operator status", async () => {
              expect(await walletRegistry.isOperatorUpToDate(operator.address))
                .to.be.true
            })

            it("should not emit InvoluntaryAuthorizationDecreaseFailed event", async () => {
              await expect(tx).to.not.emit(
                walletRegistry,
                "InvoluntaryAuthorizationDecreaseFailed"
              )
            })
          })

          context(
            "when the authorized amount drops to below the minimum",
            () => {
              before(async () => {
                const slashingTo = minimumAuthorization.sub(1)
                const slashingBy = stakedAmount.sub(slashingTo)

                await createSnapshot()

                // slash!
                await staking
                  .connect(slasher.wallet)
                  .slash(slashingBy, [stakingProvider.address])

                await staking.connect(thirdParty).processSlashing(1)
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should remove operator from the sortition pool", async () => {
                expect(await walletRegistry.isOperatorInPool(operator.address))
                  .to.be.false
              })
            }
          )
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

    // The only option for it to happen is when there was a slashing.
    context(
      "when the authorization dropped below the minimum but is still non-zero",
      () => {
        before(async () => {
          await createSnapshot()

          await walletRegistry
            .connect(stakingProvider)
            .registerOperator(operator.address)

          const authorizedAmount = minimumAuthorization
          await staking
            .connect(authorizer)
            .increaseAuthorization(
              stakingProvider.address,
              walletRegistry.address,
              authorizedAmount
            )

          const slashingTo = minimumAuthorization.sub(1)
          const slashedAmount = authorizedAmount.sub(slashingTo)

          await staking
            .connect(slasher.wallet)
            .slash(slashedAmount, [stakingProvider.address])
          await staking.connect(thirdParty).processSlashing(1)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should revert", async () => {
          await expect(
            walletRegistry.connect(operator).joinSortitionPool()
          ).to.be.revertedWith("Authorization below the minimum")
        })
      }
    )

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
        expect(await walletRegistry.isOperatorInPool(operator.address)).to.be
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
          expect(await walletRegistry.isOperatorInPool(operator.address)).to.be
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
        expect(await walletRegistry.isOperatorInPool(operator.address)).to.be
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
          expect(await walletRegistry.isOperatorInPool(operator.address)).to.be
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
          expect(await walletRegistry.isOperatorInPool(operator.address)).to.be
            .false
        })

        it("should emit OperatorStatusUpdated", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "OperatorStatusUpdated")
            .withArgs(stakingProvider.address, operator.address)
        })
      })

      context("when there was an authorization decrease request", () => {
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
          expect(await walletRegistry.isOperatorInPool(operator.address)).to.be
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

      context(
        "when there was an authorization decrease request to non-zero",
        () => {
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

      context(
        "when there was an authorization decrease request to zero",
        () => {
          let tx: ContractTransaction

          before(async () => {
            await createSnapshot()

            // initial authorization was 2 x minimum
            // we want to decrease to zero
            const deauthorizingBy = minimumAuthorization.mul(2)

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

          it("should remove operator from the sortition pool", async () => {
            expect(await walletRegistry.isOperatorInPool(operator.address)).to
              .be.false
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

  describe("eligibleStake", () => {
    context("when staking provider has no stake authorized", () => {
      it("should return zero", async () => {
        expect(
          await walletRegistry.eligibleStake(stakingProvider.address)
        ).to.equal(0)
      })
    })

    context("when staking provider has stake authorized", () => {
      let authorizedAmount

      before(async () => {
        await createSnapshot()

        authorizedAmount = minimumAuthorization
        await staking
          .connect(authorizer)
          .increaseAuthorization(
            stakingProvider.address,
            walletRegistry.address,
            authorizedAmount
          )
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should return authorized amount", async () => {
        expect(
          await walletRegistry.eligibleStake(stakingProvider.address)
        ).to.equal(authorizedAmount)
      })
    })

    context(
      "when staking provider has some part of the stake deauthorizing",
      () => {
        let authorizedAmount
        let deauthorizingAmount

        before(async () => {
          await createSnapshot()

          authorizedAmount = minimumAuthorization.add(to1e18(2000))

          await staking
            .connect(authorizer)
            .increaseAuthorization(
              stakingProvider.address,
              walletRegistry.address,
              authorizedAmount
            )

          deauthorizingAmount = to1e18(1337)
          await staking
            .connect(authorizer)
            ["requestAuthorizationDecrease(address,address,uint96)"](
              stakingProvider.address,
              walletRegistry.address,
              deauthorizingAmount
            )
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should return authorized amount minus deauthorizing amount", async () => {
          expect(
            await walletRegistry.eligibleStake(stakingProvider.address)
          ).to.equal(authorizedAmount.sub(deauthorizingAmount))
        })
      }
    )

    context("when staking provider has all of the stake deauthorizing", () => {
      before(async () => {
        await createSnapshot()

        const authorizedAmount = minimumAuthorization
        await staking
          .connect(authorizer)
          .increaseAuthorization(
            stakingProvider.address,
            walletRegistry.address,
            authorizedAmount
          )

        await staking
          .connect(authorizer)
          ["requestAuthorizationDecrease(address,address,uint96)"](
            stakingProvider.address,
            walletRegistry.address,
            authorizedAmount
          )
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should return zero", async () => {
        expect(
          await walletRegistry.eligibleStake(stakingProvider.address)
        ).to.equal(0)
      })
    })

    context("when staking provider has all of the stake deauthorized", () => {
      before(async () => {
        await createSnapshot()

        const authorizedAmount = minimumAuthorization.add(1200)
        await staking
          .connect(authorizer)
          .increaseAuthorization(
            stakingProvider.address,
            walletRegistry.address,
            authorizedAmount
          )

        await staking
          .connect(authorizer)
          ["requestAuthorizationDecrease(address,address,uint96)"](
            stakingProvider.address,
            walletRegistry.address,
            authorizedAmount
          )

        await walletRegistry.approveAuthorizationDecrease(
          stakingProvider.address
        )
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should return zero", async () => {
        expect(
          await walletRegistry.eligibleStake(stakingProvider.address)
        ).to.equal(0)
      })
    })

    // The only option for it to happen is when there was a slashing.
    context(
      "when the authorization dropped below the minimum but is still non-zero",
      () => {
        before(async () => {
          await createSnapshot()

          const authorizedAmount = minimumAuthorization
          await staking
            .connect(authorizer)
            .increaseAuthorization(
              stakingProvider.address,
              walletRegistry.address,
              authorizedAmount
            )

          const slashingTo = minimumAuthorization.sub(1)
          const slashedAmount = authorizedAmount.sub(slashingTo)

          await staking
            .connect(slasher.wallet)
            .slash(slashedAmount, [stakingProvider.address])
          await staking.connect(thirdParty).processSlashing(1)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should return zero", async () => {
          expect(
            await walletRegistry.eligibleStake(stakingProvider.address)
          ).to.equal(0)
        })
      }
    )
  })

  describe("remainingAuthorizationDecreaseDelay", () => {
    before(async () => {
      await createSnapshot()

      const authorizedAmount = minimumAuthorization.add(1200)
      await staking
        .connect(authorizer)
        .increaseAuthorization(
          stakingProvider.address,
          walletRegistry.address,
          authorizedAmount
        )

      await walletRegistry
        .connect(stakingProvider)
        .registerOperator(operator.address)
      await walletRegistry.connect(operator).joinSortitionPool()

      await staking
        .connect(authorizer)
        ["requestAuthorizationDecrease(address,address,uint96)"](
          stakingProvider.address,
          walletRegistry.address,
          authorizedAmount
        )
    })

    after(async () => {
      await restoreSnapshot()
    })

    // These tests cover only basic cases. More scenarios such as operator not
    // registered for the staking provider has been covered in tests for other
    // functions.

    it("should not activate before sortition pool is updated", async () => {
      expect(
        await walletRegistry.remainingAuthorizationDecreaseDelay(
          stakingProvider.address
        )
      ).to.equal(MAX_UINT64)
    })

    it("should activate after updating sortition pool", async () => {
      await walletRegistry.updateOperatorStatus(operator.address)
      expect(
        await walletRegistry.remainingAuthorizationDecreaseDelay(
          stakingProvider.address
        )
      ).to.equal(params.authorizationDecreaseDelay)
    })

    it("should reduce over time", async () => {
      await walletRegistry.updateOperatorStatus(operator.address)
      await helpers.time.increaseTime(params.authorizationDecreaseDelay / 2)
      expect(
        await walletRegistry.remainingAuthorizationDecreaseDelay(
          stakingProvider.address
        )
      ).to.be.closeTo(
        ethers.BigNumber.from(params.authorizationDecreaseDelay / 2),
        5 // +- 5sec
      )
    })

    it("should eventually go to zero", async () => {
      await walletRegistry.updateOperatorStatus(operator.address)
      await helpers.time.increaseTime(params.authorizationDecreaseDelay)
      expect(
        await walletRegistry.remainingAuthorizationDecreaseDelay(
          stakingProvider.address
        )
      ).to.equal(0)

      // ...and should remain zero
      await helpers.time.increaseTime(3600) // +1h
      expect(
        await walletRegistry.remainingAuthorizationDecreaseDelay(
          stakingProvider.address
        )
      ).to.equal(0)
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
          await walletRegistry.connect(walletOwner.wallet).requestNewWallet()

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

  // Testing final states for scenarios when functions are invoked one after
  // another. Operator is known and registered in the sortition pool.
  context("mixed interactions", () => {
    let initialIncrease

    before(async () => {
      await createSnapshot()

      await walletRegistry
        .connect(stakingProvider)
        .registerOperator(operator.address)

      // Authorized almost the entire staked amount but leave some margin for
      // authorization increase.
      initialIncrease = stakedAmount.sub(to1e18(20000))
      await staking
        .connect(authorizer)
        .increaseAuthorization(
          stakingProvider.address,
          walletRegistry.address,
          initialIncrease
        )
      await walletRegistry.connect(operator).joinSortitionPool()
    })

    after(async () => {
      await restoreSnapshot()
    })

    // Invoke `increaseAuthorization` after `increaseAuthorization`.
    describe("authorizationIncreased -> authorizationIncreased", () => {
      let secondIncrease

      before(async () => {
        await createSnapshot()

        secondIncrease = to1e18(11111)
        await staking
          .connect(authorizer)
          .increaseAuthorization(
            stakingProvider.address,
            walletRegistry.address,
            secondIncrease
          )

        await walletRegistry
          .connect(operator)
          .updateOperatorStatus(operator.address)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should have correct eligible stake", async () => {
        expect(
          await walletRegistry.eligibleStake(stakingProvider.address)
        ).to.equal(initialIncrease.add(secondIncrease))
      })

      it("should have operator status updated", async () => {
        expect(await walletRegistry.isOperatorUpToDate(operator.address)).to.be
          .true
      })
    })

    // Invoke `increaseAuthorization` after `authorizationDecreaseRequested`.
    // The decrease is not yet approved when `increaseAuthorization` is called.
    describe("authorizationDecreaseRequested -> authorizationIncreased", () => {
      let firstDecrease
      let secondIncrease

      before(async () => {
        await createSnapshot()

        firstDecrease = to1e18(111)
        await staking
          .connect(authorizer)
          ["requestAuthorizationDecrease(address,address,uint96)"](
            stakingProvider.address,
            walletRegistry.address,
            firstDecrease
          )
        await walletRegistry
          .connect(operator)
          .updateOperatorStatus(operator.address)

        secondIncrease = to1e18(11111)
        await staking
          .connect(authorizer)
          .increaseAuthorization(
            stakingProvider.address,
            walletRegistry.address,
            secondIncrease
          )
        await walletRegistry
          .connect(operator)
          .updateOperatorStatus(operator.address)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should have correct eligible stake", async () => {
        expect(
          await walletRegistry.eligibleStake(stakingProvider.address)
        ).to.equal(initialIncrease.sub(firstDecrease).add(secondIncrease))
      })

      it("should have operator status updated", async () => {
        expect(await walletRegistry.isOperatorUpToDate(operator.address)).to.be
          .true
      })
    })

    // Invoke `increaseAuthorization` after `approveAuthorizationDecrease`.
    // The decrease is approved when `increaseAuthorization` is called.
    describe("non-zero approveAuthorizationDecrease -> authorizationIncreased", () => {
      let firstDecrease
      let secondIncrease

      before(async () => {
        await createSnapshot()

        firstDecrease = to1e18(222)
        await staking
          .connect(authorizer)
          ["requestAuthorizationDecrease(address,address,uint96)"](
            stakingProvider.address,
            walletRegistry.address,
            firstDecrease
          )
        await walletRegistry
          .connect(operator)
          .updateOperatorStatus(operator.address)

        await helpers.time.increaseTime(params.authorizationDecreaseDelay)
        await walletRegistry.approveAuthorizationDecrease(
          stakingProvider.address
        )

        secondIncrease = to1e18(7311)
        await staking
          .connect(authorizer)
          .increaseAuthorization(
            stakingProvider.address,
            walletRegistry.address,
            secondIncrease
          )
        await walletRegistry
          .connect(operator)
          .updateOperatorStatus(operator.address)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should have correct eligible stake", async () => {
        expect(
          await walletRegistry.eligibleStake(stakingProvider.address)
        ).to.equal(initialIncrease.sub(firstDecrease).add(secondIncrease))
      })

      it("should have operator status updated", async () => {
        expect(await walletRegistry.isOperatorUpToDate(operator.address)).to.be
          .true
      })
    })

    // Invoke `increaseAuthorization` after the authorization was decreased to 0.
    describe("to-zero approveAuthorizationDecrease -> authorizationIncreased", () => {
      let secondIncrease

      before(async () => {
        await createSnapshot()

        await staking
          .connect(authorizer)
          ["requestAuthorizationDecrease(address,address,uint96)"](
            stakingProvider.address,
            walletRegistry.address,
            initialIncrease
          )
        await walletRegistry
          .connect(operator)
          .updateOperatorStatus(operator.address)

        await helpers.time.increaseTime(params.authorizationDecreaseDelay)
        await walletRegistry.approveAuthorizationDecrease(
          stakingProvider.address
        )

        secondIncrease = minimumAuthorization.add(to1e18(21))
        await staking
          .connect(authorizer)
          .increaseAuthorization(
            stakingProvider.address,
            walletRegistry.address,
            secondIncrease
          )
        await walletRegistry.connect(operator).joinSortitionPool()
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should have correct eligible stake", async () => {
        expect(
          await walletRegistry.eligibleStake(stakingProvider.address)
        ).to.equal(secondIncrease)
      })

      it("should have operator status updated", async () => {
        expect(await walletRegistry.isOperatorUpToDate(operator.address)).to.be
          .true
      })
    })

    // Invoke `increaseAuthorization` after `involuntaryAuthorizationDecrease`
    // when the authorization amount dropped below the minimum authorization.
    describe("below-minimum involuntaryAuthorizationDecrease -> authorizationIncreased", () => {
      let slashingTo
      let secondIncrease

      before(async () => {
        await createSnapshot()

        slashingTo = minimumAuthorization.sub(1)
        const slashedAmount = initialIncrease.sub(slashingTo)

        await staking
          .connect(slasher.wallet)
          .slash(slashedAmount, [stakingProvider.address])
        await staking.connect(thirdParty).processSlashing(1)

        // Give the stake owner some more T and let them top-up the stake before
        // they increase the authorization again.
        secondIncrease = to1e18(10000)
        await t.connect(deployer).mint(owner.address, secondIncrease)
        await t.connect(owner).approve(staking.address, secondIncrease)
        await staking
          .connect(owner)
          .topUp(stakingProvider.address, secondIncrease)

        // And finally increase!
        await staking
          .connect(authorizer)
          .increaseAuthorization(
            stakingProvider.address,
            walletRegistry.address,
            secondIncrease
          )
        await walletRegistry.connect(operator).joinSortitionPool()
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should have correct eligible stake", async () => {
        expect(
          await walletRegistry.eligibleStake(stakingProvider.address)
        ).to.equal(slashingTo.add(secondIncrease))
      })

      it("should have operator status updated", async () => {
        expect(await walletRegistry.isOperatorUpToDate(operator.address)).to.be
          .true
      })
    })

    describe("authorizationDecreaseRequested -> involuntaryAuthorizationDecrease", () => {
      let decreasedAmount
      let slashingTo

      before(async () => {
        await createSnapshot()

        decreasedAmount = to1e18(20000)
        await staking
          .connect(authorizer)
          ["requestAuthorizationDecrease(address,address,uint96)"](
            stakingProvider.address,
            walletRegistry.address,
            decreasedAmount
          )
        await walletRegistry
          .connect(operator)
          .updateOperatorStatus(operator.address)

        slashingTo = initialIncrease.sub(to1e18(100))
        const slashedAmount = initialIncrease.sub(slashingTo)

        await staking
          .connect(slasher.wallet)
          .slash(slashedAmount, [stakingProvider.address])
        await staking.connect(thirdParty).processSlashing(1)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should have correct eligible stake", async () => {
        expect(
          await walletRegistry.eligibleStake(stakingProvider.address)
        ).to.equal(slashingTo.sub(decreasedAmount))
      })

      it("should have operator status updated", async () => {
        expect(await walletRegistry.isOperatorUpToDate(operator.address)).to.be
          .true
      })
    })

    describe("authorizationDecreaseRequested -> involuntaryAuthorizationDecrease -> approveAuthorizationDecrease", () => {
      let decreasedAmount
      let slashingTo

      before(async () => {
        await createSnapshot()

        decreasedAmount = to1e18(20000)
        await staking
          .connect(authorizer)
          ["requestAuthorizationDecrease(address,address,uint96)"](
            stakingProvider.address,
            walletRegistry.address,
            decreasedAmount
          )
        await walletRegistry
          .connect(operator)
          .updateOperatorStatus(operator.address)

        slashingTo = initialIncrease.sub(to1e18(100))
        const slashedAmount = initialIncrease.sub(slashingTo)

        await staking
          .connect(slasher.wallet)
          .slash(slashedAmount, [stakingProvider.address])
        await staking.connect(thirdParty).processSlashing(1)

        await helpers.time.increaseTime(params.authorizationDecreaseDelay)
        await walletRegistry.approveAuthorizationDecrease(
          stakingProvider.address
        )
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should have correct eligible stake", async () => {
        expect(
          await walletRegistry.eligibleStake(stakingProvider.address)
        ).to.equal(slashingTo.sub(decreasedAmount))
      })

      it("should have operator status updated", async () => {
        expect(await walletRegistry.isOperatorUpToDate(operator.address)).to.be
          .true
      })
    })

    describe("approveAuthorizationDecrease -> involuntaryAuthorizationDecrease", () => {
      let decreasedAmount
      let slashingTo

      before(async () => {
        await createSnapshot()

        decreasedAmount = to1e18(1000)
        await staking
          .connect(authorizer)
          ["requestAuthorizationDecrease(address,address,uint96)"](
            stakingProvider.address,
            walletRegistry.address,
            decreasedAmount
          )
        await walletRegistry
          .connect(operator)
          .updateOperatorStatus(operator.address)

        await helpers.time.increaseTime(params.authorizationDecreaseDelay)
        await walletRegistry.approveAuthorizationDecrease(
          stakingProvider.address
        )

        slashingTo = initialIncrease.sub(to1e18(2500))
        const slashedAmount = initialIncrease
          .sub(decreasedAmount)
          .sub(slashingTo)

        await staking
          .connect(slasher.wallet)
          .slash(slashedAmount, [stakingProvider.address])
        await staking.connect(thirdParty).processSlashing(1)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should have correct eligible stake", async () => {
        expect(
          await walletRegistry.eligibleStake(stakingProvider.address)
        ).to.equal(slashingTo)
      })

      it("should have operator status updated", async () => {
        expect(await walletRegistry.isOperatorUpToDate(operator.address)).to.be
          .true
      })
    })
  })
})
