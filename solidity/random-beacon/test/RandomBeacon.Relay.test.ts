/* eslint-disable @typescript-eslint/no-unused-expressions, no-await-in-loop, @typescript-eslint/no-extra-semi */

import {
  ethers,
  waffle,
  helpers,
  getUnnamedAccounts,
  getNamedAccounts,
} from "hardhat"
import { expect } from "chai"
import { BigNumber, ContractTransaction } from "ethers"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { Address } from "hardhat-deploy/types"
import blsData from "./data/bls"
import {
  constants,
  dkgState,
  params,
  randomBeaconDeployment,
  blsDeployment,
} from "./fixtures"
import { createGroup, hashUint32Array } from "./utils/groups"
import { signHeartbeatFailureClaim } from "./utils/heartbeat"
import type {
  RandomBeacon,
  RandomBeaconStub,
  TestToken,
  RelayStub,
  SortitionPool,
  StakingStub,
  BLS,
} from "../typechain"
import { registerOperators, Operator, OperatorID } from "./utils/operators"

const { mineBlocks, mineBlocksTo } = helpers.time
const { to1e18 } = helpers.number
const ZERO_ADDRESS = ethers.constants.AddressZero
const { createSnapshot, restoreSnapshot } = helpers.snapshot

async function fixture() {
  const deployment = await randomBeaconDeployment()

  // Additional contracts needed by this test suite.
  const bls = (await blsDeployment()).bls as BLS
  const relayStub = (await (
    await ethers.getContractFactory("RelayStub", {
      libraries: {
        BLS: (await blsDeployment()).bls.address,
      },
    })
  ).deploy()) as RelayStub

  // Register operators in the sortition pool to make group creation
  // possible.
  const operators = await registerOperators(
    deployment.randomBeacon as RandomBeacon,
    (await getUnnamedAccounts()).slice(0, constants.groupSize)
  )

  return {
    randomBeacon: deployment.randomBeacon as RandomBeacon,
    sortitionPool: deployment.sortitionPool as SortitionPool,
    testToken: deployment.testToken as TestToken,
    staking: deployment.stakingStub as StakingStub,
    relayStub,
    bls,
    operators,
  }
}

describe("RandomBeacon - Relay", () => {
  let deployer: SignerWithAddress
  let requester: SignerWithAddress
  let notifier: SignerWithAddress
  let submitter: SignerWithAddress
  let members: Operator[]
  let membersIDs: OperatorID[]
  let membersAddresses: Address[]

  let randomBeacon: RandomBeacon
  let sortitionPool: SortitionPool
  let testToken: TestToken
  let staking: StakingStub
  let relayStub: RelayStub
  let bls: BLS

  before(async () => {
    deployer = await ethers.getSigner((await getNamedAccounts()).deployer)
    requester = await ethers.getSigner((await getUnnamedAccounts())[1])
    notifier = await ethers.getSigner((await getUnnamedAccounts())[2])
    submitter = await ethers.getSigner((await getUnnamedAccounts())[3])
    ;({
      randomBeacon,
      sortitionPool,
      testToken,
      staking,
      relayStub,
      bls,
      operators: members,
    } = await waffle.loadFixture(fixture))

    membersIDs = members.map((member) => member.id)
    membersAddresses = members.map((member) => member.address)
  })

  describe("requestRelayEntry", () => {
    context("when groups exist", () => {
      before(async () => {
        await createSnapshot()

        await createGroup(randomBeacon as RandomBeaconStub, members)
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("when there is no other relay entry in progress", () => {
        context("when the requester pays the relay request fee", () => {
          let tx: ContractTransaction
          let previousDkgRewardsPoolBalance: BigNumber
          let previousRandomBeaconBalance: BigNumber

          before(async () => {
            await createSnapshot()

            previousDkgRewardsPoolBalance = await randomBeacon.dkgRewardsPool()
            previousRandomBeaconBalance = await testToken.balanceOf(
              randomBeacon.address
            )
            await approveTestToken()
          })

          after(async () => {
            await restoreSnapshot()
          })

          context(
            "when relay request does not hit group creation frequency threshold",
            () => {
              before(async () => {
                await createSnapshot()

                tx = await randomBeacon
                  .connect(requester)
                  .requestRelayEntry(ZERO_ADDRESS)
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should deposit relay request fee to the DKG rewards pool", async () => {
                // Assert correct pool bookkeeping.
                const currentDkgRewardsPoolBalance =
                  await randomBeacon.dkgRewardsPool()
                expect(
                  currentDkgRewardsPoolBalance.sub(
                    previousDkgRewardsPoolBalance
                  )
                ).to.be.equal(params.relayRequestFee)

                // Assert actual transfer took place.
                const currentRandomBeaconBalance = await testToken.balanceOf(
                  randomBeacon.address
                )
                expect(
                  currentRandomBeaconBalance.sub(previousRandomBeaconBalance)
                ).to.be.equal(params.relayRequestFee)
              })

              it("should emit RelayEntryRequested event", async () => {
                await expect(tx)
                  .to.emit(randomBeacon, "RelayEntryRequested")
                  .withArgs(1, 0, blsData.previousEntry)
              })

              it("should not lock DKG state", async () => {
                expect(await randomBeacon.getGroupCreationState()).to.be.equal(
                  dkgState.IDLE
                )
                expect(await sortitionPool.isLocked()).to.be.false
              })
            }
          )

          context(
            "when relay request hits group creation frequency threshold",
            () => {
              before(async () => {
                await createSnapshot()

                // Force group creation on each relay entry.
                await randomBeacon
                  .connect(deployer)
                  .updateGroupCreationParameters(1, params.groupLifeTime)

                tx = await randomBeacon
                  .connect(requester)
                  .requestRelayEntry(ZERO_ADDRESS)
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should deposit relay request fee to the DKG rewards pool", async () => {
                // Assert correct pool bookkeeping.
                const currentDkgRewardsPoolBalance =
                  await randomBeacon.dkgRewardsPool()
                expect(
                  currentDkgRewardsPoolBalance.sub(
                    previousDkgRewardsPoolBalance
                  )
                ).to.be.equal(params.relayRequestFee)

                // Assert actual transfer took place.
                const currentRandomBeaconBalance = await testToken.balanceOf(
                  randomBeacon.address
                )
                expect(
                  currentRandomBeaconBalance.sub(previousRandomBeaconBalance)
                ).to.be.equal(params.relayRequestFee)
              })

              it("should emit RelayEntryRequested event", async () => {
                await expect(tx)
                  .to.emit(randomBeacon, "RelayEntryRequested")
                  .withArgs(1, 0, blsData.previousEntry)
              })

              it("should lock DKG state", async () => {
                expect(await randomBeacon.getGroupCreationState()).to.be.equal(
                  dkgState.AWAITING_SEED
                )
                expect(await sortitionPool.isLocked()).to.be.true
              })

              it("should emit DkgStateLocked event", async () => {
                await expect(tx).to.emit(randomBeacon, "DkgStateLocked")
              })
            }
          )
        })

        context("when the requester doesn't pay the relay request fee", () => {
          it("should revert", async () => {
            await expect(
              randomBeacon.connect(requester).requestRelayEntry(ZERO_ADDRESS)
            ).to.be.revertedWith("Transfer amount exceeds allowance")
          })
        })
      })

      context("when there is an other relay entry in progress", () => {
        before(async () => {
          await createSnapshot()

          await approveTestToken()
          await randomBeacon.connect(requester).requestRelayEntry(ZERO_ADDRESS)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should revert", async () => {
          await expect(
            randomBeacon.connect(requester).requestRelayEntry(ZERO_ADDRESS)
          ).to.be.revertedWith("Another relay request in progress")
        })
      })
    })

    context("when no groups exist", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon.connect(requester).requestRelayEntry(ZERO_ADDRESS)
        ).to.be.revertedWith("No active groups")
      })
    })
  })

  describe("submitRelayEntry happy path", () => {
    before(async () => {
      await createSnapshot()

      await createGroup(randomBeacon as RandomBeaconStub, members)
    })

    after(async () => {
      await restoreSnapshot()
    })

    context("when relay request is in progress", () => {
      before(async () => {
        await createSnapshot()

        await approveTestToken()
        await randomBeacon.connect(requester).requestRelayEntry(ZERO_ADDRESS)
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("when relay entry has not timed out", () => {
        context("when entry is valid", () => {
          context("when result is submitted before the soft timeout", () => {
            let tx: ContractTransaction
            before(async () => {
              await createSnapshot()

              tx = await randomBeacon
                .connect(submitter)
                ["submitRelayEntry(bytes)"](blsData.groupSignature)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should not slash any members", async () => {
              await expect(tx).to.not.emit(staking, "Slashed")
            })

            it("should emit RelayEntrySubmitted event", async () => {
              await expect(tx)
                .to.emit(randomBeacon, "RelayEntrySubmitted")
                .withArgs(1, submitter.address, blsData.groupSignature)
            })

            it("should terminate the relay request", async () => {
              expect(await randomBeacon.isRelayRequestInProgress()).to.be.false
            })
          })

          context("when result is submitted after the soft timeout", () => {
            before(async () => {
              await createSnapshot()
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should revert", async () => {
              await mineBlocks(
                constants.groupSize *
                  params.relayEntrySubmissionEligibilityDelay +
                  1
              )
              await expect(
                randomBeacon
                  .connect(submitter)
                  ["submitRelayEntry(bytes)"](blsData.groupSignature)
              ).to.be.revertedWith("Relay submission passed a soft timeout")
            })
          })

          context("when DKG is awaiting a seed", () => {
            let tx: ContractTransaction

            before(async () => {
              await createSnapshot()

              // Simulate DKG is awaiting a seed.
              await (randomBeacon as RandomBeaconStub).publicDkgLockState()

              tx = await randomBeacon
                .connect(submitter)
                ["submitRelayEntry(bytes)"](blsData.groupSignature)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should emit DkgStarted event", async () => {
              await expect(tx)
                .to.emit(randomBeacon, "DkgStarted")
                .withArgs(blsData.groupSignatureUint256)
            })
          })
        })

        context("when entry is invalid", () => {
          it("should revert", async () => {
            await expect(
              randomBeacon
                .connect(submitter)
                ["submitRelayEntry(bytes)"](blsData.nextGroupSignature)
            ).to.be.revertedWith("Invalid entry")
          })
        })
      })
    })
  })

  describe("submitRelayEntry after the soft timeout", () => {
    before(async () => {
      await createSnapshot()

      await createGroup(randomBeacon as RandomBeaconStub, members)
    })

    after(async () => {
      await restoreSnapshot()
    })

    context("when relay request is in progress", () => {
      before(async () => {
        await createSnapshot()

        await approveTestToken()
        await randomBeacon.connect(requester).requestRelayEntry(ZERO_ADDRESS)
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("when the input params are valid", () => {
        context("when result is submitted before the soft timeout", () => {
          let tx: ContractTransaction

          before(async () => {
            await createSnapshot()
            tx = await randomBeacon
              .connect(submitter)
              ["submitRelayEntry(bytes,uint32[])"](
                blsData.groupSignature,
                membersIDs
              )
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should not slash members ", async () => {
            await expect(tx).to.not.emit(staking, "Slashed")

            await expect(tx).to.not.emit(randomBeacon, "RelayEntryDelaySlashed")
          })

          it("should emit RelayEntrySubmitted event", async () => {
            await expect(tx)
              .to.emit(randomBeacon, "RelayEntrySubmitted")
              .withArgs(1, submitter.address, blsData.groupSignature)
          })

          it("should terminate the relay request", async () => {
            expect(await randomBeacon.isRelayRequestInProgress()).to.be.false
          })
        })

        context("when result is submitted after the soft timeout", () => {
          let tx: ContractTransaction

          before(async () => {
            await createSnapshot()

            // Let's assume we want to submit the relay entry after 75%
            // of the soft timeout period elapses. If so we need to
            // mine the following number of blocks:
            // `groupSize * relayEntrySubmissionEligibilityDelay +
            // (0.75 * relayEntryHardTimeout)`. However, we need to
            // subtract one block because the relay entry submission
            // transaction will move the blockchain ahead by one block
            // due to the Hardhat auto-mine feature.
            await mineBlocks(
              constants.groupSize *
                params.relayEntrySubmissionEligibilityDelay +
                0.75 * params.relayEntryHardTimeout -
                1
            )
            tx = await randomBeacon
              .connect(submitter)
              ["submitRelayEntry(bytes,uint32[])"](
                blsData.groupSignature,
                membersIDs
              )
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should slash a correct portion of the slashing amount for all members ", async () => {
            // `relayEntrySubmissionFailureSlashingAmount = 1000e18`.
            // 75% of the soft timeout period elapsed so we expect
            // `750e18` to be slashed.
            await expect(tx)
              .to.emit(staking, "Slashed")
              .withArgs(to1e18(750), membersAddresses)

            await expect(tx)
              .to.emit(randomBeacon, "RelayEntryDelaySlashed")
              .withArgs(1, to1e18(750), membersAddresses)
          })

          it("should emit RelayEntrySubmitted event", async () => {
            await expect(tx)
              .to.emit(randomBeacon, "RelayEntrySubmitted")
              .withArgs(1, submitter.address, blsData.groupSignature)
          })

          it("should terminate the relay request", async () => {
            expect(await randomBeacon.isRelayRequestInProgress()).to.be.false
          })
        })
      })

      context("when the input params are invalid", () => {
        before(async () => {
          await createSnapshot()
          await mineBlocks(
            constants.groupSize * params.relayEntrySubmissionEligibilityDelay +
              1
          )
        })

        after(async () => {
          await restoreSnapshot()
        })

        context("when entry is not valid", () => {
          it("should revert", async () => {
            await expect(
              randomBeacon
                .connect(submitter)
                ["submitRelayEntry(bytes,uint32[])"](
                  blsData.nextGroupSignature,
                  membersIDs
                )
            ).to.be.revertedWith("Invalid entry")
          })
        })

        context("when group members are invalid", () => {
          it("should revert", async () => {
            const invalidMembersId = [0, 1, 42]
            await expect(
              randomBeacon
                .connect(submitter)
                ["submitRelayEntry(bytes,uint32[])"](
                  blsData.nextGroupSignature,
                  invalidMembersId
                )
            ).to.be.revertedWith("Invalid group members")
          })
        })
      })

      context("when a relay entry has timed out", () => {
        it("should revert", async () => {
          await mineBlocks(
            constants.groupSize * params.relayEntrySubmissionEligibilityDelay +
              params.relayEntryHardTimeout
          )

          await expect(
            randomBeacon
              .connect(submitter)
              ["submitRelayEntry(bytes,uint32[])"](
                blsData.nextGroupSignature,
                membersIDs
              )
          ).to.be.revertedWith("Relay request timed out")
        })
      })
    })
  })

  describe("submitEntry", () => {
    context("when relay request is not in progress", () => {
      it("should revert", async () => {
        await expect(
          relayStub.callStatic.submitEntry(
            blsData.groupSignature,
            blsData.groupPubKey
          )
        ).to.be.revertedWith("No relay request in progress")
      })
    })
  })

  describe("reportRelayEntryTimeout", () => {
    before(async () => {
      await createSnapshot()

      await createGroup(randomBeacon, members)
      await approveTestToken()
      await randomBeacon.connect(requester).requestRelayEntry(ZERO_ADDRESS)
    })

    after(async () => {
      await restoreSnapshot()
    })

    context("when relay entry timed out", () => {
      context(
        "when other active groups exist after timeout is reported",
        () => {
          let tx: ContractTransaction

          before(async () => {
            await createSnapshot()

            await mineBlocks(
              constants.groupSize *
                params.relayEntrySubmissionEligibilityDelay +
                params.relayEntryHardTimeout
            )

            await (randomBeacon as RandomBeaconStub).roughlyAddGroup(
              "0x01",
              hashUint32Array(membersIDs)
            )

            tx = await randomBeacon
              .connect(notifier)
              .reportRelayEntryTimeout(membersIDs)
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should slash the full slashing amount for all group members", async () => {
            await expect(tx)
              .to.emit(staking, "Seized")
              .withArgs(
                to1e18(1000),
                params.relayEntryTimeoutNotificationRewardMultiplier,
                notifier.address,
                membersAddresses
              )

            await expect(tx)
              .to.emit(randomBeacon, "RelayEntryTimeoutSlashed")
              .withArgs(1, to1e18(1000), membersAddresses)
          })

          it("should terminate the group", async () => {
            const isGroupTeminated = await (
              randomBeacon as RandomBeaconStub
            ).isGroupTerminated(0)
            expect(isGroupTeminated).to.be.equal(true)
          })

          it("should emit RelayEntryTimedOut event", async () => {
            await expect(tx)
              .to.emit(randomBeacon, "RelayEntryTimedOut")
              .withArgs(1, 0)
          })

          it("should retry current relay request", async () => {
            // We expect the same request ID because this is a retry.
            // Group ID is `1` because we take an active group from `groupsRegistry`
            // array. Group with an index `0` was terminated.
            await expect(tx)
              .to.emit(randomBeacon, "RelayEntryRequested")
              .withArgs(1, 1, blsData.previousEntry)

            expect(await randomBeacon.isRelayRequestInProgress()).to.be.true
          })
        }
      )

      context(
        "when a group that was supposed to submit a relay request is terminated and another group expires",
        () => {
          let tx: ContractTransaction

          before(async () => {
            await createSnapshot()

            // Create another group in a rough way just to have an active group
            // once the one handling the timed out request gets terminated.
            // This makes the request retry possible. That group will not
            // perform any signing so their public key can be arbitrary bytes.
            // Also, that group is created just after the relay request is
            // made to ensure it is not selected for signing the original request.
            await (randomBeacon as RandomBeaconStub).roughlyAddGroup(
              "0x01",
              hashUint32Array(membersIDs)
            )

            await mineBlocks(
              constants.groupSize *
                params.relayEntrySubmissionEligibilityDelay +
                params.relayEntryHardTimeout
            )

            const registry = await randomBeacon.getGroupsRegistry()
            const secondGroupLifetime = await (
              randomBeacon as RandomBeaconStub
            ).groupLifetimeOf(registry[1])

            // Expire second group
            await mineBlocksTo(Number(secondGroupLifetime) + 1)

            tx = await randomBeacon.reportRelayEntryTimeout(membersIDs)
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should clean up current relay request data", async () => {
            await expect(tx).to.not.emit(randomBeacon, "RelayEntryRequested")
            expect(await randomBeacon.isRelayRequestInProgress()).to.be.false
          })
        }
      )

      context("when no active groups exist after timeout is reported", () => {
        let tx: ContractTransaction

        before(async () => {
          await createSnapshot()

          await mineBlocks(
            constants.groupSize * params.relayEntrySubmissionEligibilityDelay +
              params.relayEntryHardTimeout
          )

          tx = await randomBeacon
            .connect(notifier)
            .reportRelayEntryTimeout(membersIDs)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should slash the full slashing amount for all group members", async () => {
          await expect(tx)
            .to.emit(staking, "Seized")
            .withArgs(
              to1e18(1000),
              params.relayEntryTimeoutNotificationRewardMultiplier,
              notifier.address,
              membersAddresses
            )

          await expect(tx)
            .to.emit(randomBeacon, "RelayEntryTimeoutSlashed")
            .withArgs(1, to1e18(1000), membersAddresses)
        })

        it("should terminate the group", async () => {
          const isGroupTeminated = await (
            randomBeacon as RandomBeaconStub
          ).isGroupTerminated(0)
          expect(isGroupTeminated).to.be.equal(true)
        })

        it("should emit RelayEntryTimedOut event", async () => {
          await expect(tx)
            .to.emit(randomBeacon, "RelayEntryTimedOut")
            .withArgs(1, 0)
        })

        it("should clean up current relay request data", async () => {
          await expect(tx).to.not.emit(randomBeacon, "RelayEntryRequested")
          expect(await randomBeacon.isRelayRequestInProgress()).to.be.false
        })
      })

      context(
        "when no active groups exist after timeout is reported and DKG is awaiting seed",
        () => {
          let tx: ContractTransaction

          before(async () => {
            await createSnapshot()

            await mineBlocks(
              constants.groupSize *
                params.relayEntrySubmissionEligibilityDelay +
                params.relayEntryHardTimeout
            )

            // Simulate DKG is awaiting a seed.
            await (randomBeacon as RandomBeaconStub).publicDkgLockState()

            tx = await randomBeacon
              .connect(notifier)
              .reportRelayEntryTimeout(membersIDs)
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should notify DKG seed timed out", async () => {
            expect(await randomBeacon.getGroupCreationState()).to.be.equal(
              dkgState.IDLE
            )
            expect(await sortitionPool.isLocked()).to.be.false
          })

          it("should emit DkgSeedTimedOut event", async () => {
            await expect(tx).to.emit(randomBeacon, "DkgSeedTimedOut")
          })
        }
      )
    })

    context("when relay entry did not time out", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon.reportRelayEntryTimeout(membersIDs)
        ).to.be.revertedWith("Relay request did not time out")
      })
    })

    context("when group members are invalid", () => {
      it("should revert", async () => {
        const invalidMembersId = [0, 1, 42]
        await expect(
          randomBeacon.reportRelayEntryTimeout(invalidMembersId)
        ).to.be.revertedWith("Invalid group members")
      })
    })
  })

  describe("reportUnauthorizedSigning", () => {
    before(async () => {
      await createSnapshot()

      await createGroup(randomBeacon as RandomBeaconStub, members)
      await approveTestToken()
      await randomBeacon.connect(requester).requestRelayEntry(ZERO_ADDRESS)
    })

    after(async () => {
      await restoreSnapshot()
    })

    context("when a group is active", () => {
      context("when provided signature is valid", () => {
        let tx

        before(async () => {
          await createSnapshot()

          const notifierSignature = await bls.sign(
            notifier.address,
            blsData.secretKey
          )
          tx = await randomBeacon
            .connect(notifier)
            .reportUnauthorizedSigning(notifierSignature, 0, membersIDs)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should terminate the group", async () => {
          const isGroupTeminated = await (
            randomBeacon as RandomBeaconStub
          ).isGroupTerminated(0)
          expect(isGroupTeminated).to.be.equal(true)
        })

        it("should call staking contract to seize the min stake", async () => {
          await expect(tx)
            .to.emit(staking, "Seized")
            .withArgs(
              to1e18(100000),
              params.unauthorizedSigningNotificationRewardMultiplier,
              notifier.address,
              membersAddresses
            )
        })

        it("should emit unauthorized signing slashing event", async () => {
          await expect(tx)
            .to.emit(randomBeacon, "UnauthorizedSigningSlashed")
            .withArgs(0, to1e18(100000), membersAddresses)
        })
      })
    })

    context("when group is terminated", () => {
      before(async () => {
        await createSnapshot()

        await mineBlocks(
          constants.groupSize * params.relayEntrySubmissionEligibilityDelay +
            params.relayEntryHardTimeout
        )

        await (randomBeacon as RandomBeaconStub).roughlyTerminateGroup(0)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should revert", async () => {
        const notifierSignature = await bls.sign(
          notifier.address,
          blsData.secretKey
        )

        await expect(
          randomBeacon
            .connect(notifier)
            .reportUnauthorizedSigning(notifierSignature, 0, membersIDs)
        ).to.be.revertedWith("Group cannot be terminated")
      })
    })

    context("when provided signature is invalid", () => {
      it("should revert", async () => {
        // the valid key is 123 instead of 42
        const notifierSignature = await bls.sign(notifier.address, 42)
        await expect(
          randomBeacon
            .connect(notifier)
            .reportUnauthorizedSigning(notifierSignature, 0, membersIDs)
        ).to.be.revertedWith("Invalid signature")
      })
    })

    context("when group members are invalid", () => {
      it("should revert", async () => {
        const notifierSignature = await bls.sign(notifier.address, 42)
        const invalidMembersId = [0, 1, 42]
        await expect(
          randomBeacon
            .connect(notifier)
            .reportUnauthorizedSigning(notifierSignature, 0, invalidMembersId)
        ).to.be.revertedWith("Invalid group members")
      })
    })
  })

  describe("calculateSlashingAmount", () => {
    const testGroupSize = 64

    before(async () => {
      await relayStub.setCurrentRequestStartBlock()
    })

    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    context("when a soft timeout has been exceeded by one block", () => {
      it("should return a correct slashing amount", async () => {
        await mineBlocks(
          testGroupSize * params.relayEntrySubmissionEligibilityDelay + 1
        )

        // We exceeded the soft timeout by `1`
        // slashing amount: 1 * 1000e18 / 100 = 10e18
        expect(
          await relayStub.callStatic.calculateSlashingAmount()
        ).to.be.equal(BigNumber.from("10000000000000000000"))
      })
    })

    context(
      "when soft timeout has been exceeded by the number of blocks equal to the hard timeout",
      () => {
        it("should return a correct slashing amount", async () => {
          await mineBlocks(
            testGroupSize * params.relayEntrySubmissionEligibilityDelay +
              params.relayEntryHardTimeout
          )

          // We exceeded the soft timeout by `100`
          // slashing amount: 100 * 1000e18 / 100 = 1000e18
          expect(
            await relayStub.callStatic.calculateSlashingAmount()
          ).to.be.equal(BigNumber.from("1000000000000000000000"))
        })
      }
    )

    context(
      "when soft timeout has been exceeded by the number of blocks bigger than the hard timeout",
      () => {
        it("should return a correct slashing factor", async () => {
          await mineBlocks(
            testGroupSize * params.relayEntrySubmissionEligibilityDelay +
              params.relayEntryHardTimeout +
              1
          )

          // We are exceeded the soft timeout by a value bigger than the
          // hard timeout. In that case the maximum value (100%) of the slashing
          // amount should be returned.
          expect(
            await relayStub.callStatic.calculateSlashingAmount()
          ).to.be.equal(BigNumber.from("1000000000000000000000"))
        })
      }
    )
  })

  describe("fundHeartbeatNotifierRewardsPool", () => {
    const amount = to1e18(1000)

    let previousHeartbeatNotifierRewardsPoolBalance: BigNumber
    let previousRandomBeaconBalance: BigNumber

    before(async () => {
      await createSnapshot()

      previousHeartbeatNotifierRewardsPoolBalance =
        await randomBeacon.heartbeatNotifierRewardsPool()
      previousRandomBeaconBalance = await testToken.balanceOf(
        randomBeacon.address
      )

      await testToken.mint(deployer.address, amount)
      await testToken.connect(deployer).approve(randomBeacon.address, amount)

      await randomBeacon.fundHeartbeatNotifierRewardsPool(
        deployer.address,
        amount
      )
    })

    after(async () => {
      await restoreSnapshot()
    })

    it("should increase the heartbeat notifier rewards pool balance", async () => {
      const currentHeartbeatNotifierRewardsPoolBalance =
        await randomBeacon.heartbeatNotifierRewardsPool()
      expect(
        currentHeartbeatNotifierRewardsPoolBalance.sub(
          previousHeartbeatNotifierRewardsPoolBalance
        )
      ).to.be.equal(amount)
    })

    it("should transfer tokens to the random beacon contract", async () => {
      const currentRandomBeaconBalance = await testToken.balanceOf(
        randomBeacon.address
      )
      expect(
        currentRandomBeaconBalance.sub(previousRandomBeaconBalance)
      ).to.be.equal(amount)
    })
  })

  describe("notifyFailedHeartbeat", () => {
    const groupId = 0
    const stubSignatures = "0x00"
    const stubMembersIndices = []
    // Use 31 element `failedMembersIndices` array to simulate the most gas
    // expensive real-world case. If group size is 64, the required threshold
    // is 33 so we assume 31 operators at most will be marked as ineligible
    // during a single `notifyFailedHeartbeat` call.
    const subsequentFailedMembersIndices = Array.from(
      Array(31),
      (_, i) => i + 1
    )
    const nonSubsequentFailedMembersIndices = [2, 5, 7, 23, 56]
    const groupThreshold = 33

    let group

    before(async () => {
      await createSnapshot()

      await createGroup(randomBeacon as RandomBeaconStub, members)
      group = await randomBeacon["getGroup(uint64)"](groupId)
    })

    after(async () => {
      await restoreSnapshot()
    })

    context("when passed nonce is valid", () => {
      context("when group is active and non-terminated", () => {
        context("when failed members indices are correct", () => {
          context("when signatures array is correct", () => {
            context("when signing members indices are correct", () => {
              context("when all signatures are correct", () => {
                context("when claim sender signed the claim", () => {
                  const assertNotifyFailedHeartbeatSucceed = async (
                    failedMembersIndices: number[],
                    signaturesCount: number,
                    modifySignatures: (signatures: string) => string,
                    modifySigningMemberIndices: (
                      signingMemberIndices: number[]
                    ) => number[]
                  ) => {
                    let tx: ContractTransaction
                    let initialNonce: BigNumber
                    let initialNotifierBalance: BigNumber
                    let initialHeartbeatNotifierRewardsPoolBalance: BigNumber
                    let claimSender: SignerWithAddress

                    before(async () => {
                      await createSnapshot()

                      // Assume claim sender is the first signing member.
                      claimSender = await ethers.getSigner(members[0].address)

                      await fundHeartbeatNotifierRewardsPool(
                        params.ineligibleOperatorNotifierReward.mul(
                          failedMembersIndices.length
                        )
                      )

                      initialNonce = await randomBeacon.failedHeartbeatNonce(
                        groupId
                      )

                      initialNotifierBalance = await testToken.balanceOf(
                        claimSender.address
                      )

                      initialHeartbeatNotifierRewardsPoolBalance =
                        await randomBeacon.heartbeatNotifierRewardsPool()

                      const { signatures, signingMembersIndices } =
                        await signHeartbeatFailureClaim(
                          members,
                          0,
                          group.groupPubKey,
                          failedMembersIndices,
                          signaturesCount
                        )

                      tx = await randomBeacon
                        .connect(claimSender)
                        .notifyFailedHeartbeat(
                          {
                            groupId,
                            failedMembersIndices,
                            signatures: modifySignatures(signatures),
                            signingMembersIndices: modifySigningMemberIndices(
                              signingMembersIndices
                            ),
                          },
                          0,
                          membersIDs
                        )
                    })

                    after(async () => {
                      await restoreSnapshot()
                    })

                    it("should increment failed heartbeat nonce for the group", async () => {
                      expect(
                        await randomBeacon.failedHeartbeatNonce(groupId)
                      ).to.be.equal(initialNonce.add(1))
                    })

                    it("should emit HeartbeatFailed event", async () => {
                      await expect(tx)
                        .to.emit(randomBeacon, "HeartbeatFailed")
                        .withArgs(
                          groupId,
                          initialNonce.toNumber(),
                          claimSender.address
                        )
                    })

                    it("should ban sortition pool rewards for ineligible operators", async () => {
                      const now = await helpers.time.lastBlockTime()
                      const expectedUntil =
                        now + params.sortitionPoolRewardsBanDuration

                      const expectedIneligibleMembersIDs =
                        failedMembersIndices.map((i) => membersIDs[i - 1])

                      await expect(tx)
                        .to.emit(sortitionPool, "IneligibleForRewards")
                        .withArgs(expectedIneligibleMembersIDs, expectedUntil)
                    })

                    it("should pay notifier reward from heartbeat notifier rewards pool", async () => {
                      const expectedReward =
                        params.ineligibleOperatorNotifierReward.mul(
                          failedMembersIndices.length
                        )

                      const currentNotifierBalance = await testToken.balanceOf(
                        claimSender.address
                      )
                      expect(
                        currentNotifierBalance.sub(initialNotifierBalance)
                      ).to.be.equal(expectedReward)

                      const currentHeartbeatNotifierRewardsPoolBalance =
                        await randomBeacon.heartbeatNotifierRewardsPool()
                      expect(
                        initialHeartbeatNotifierRewardsPoolBalance.sub(
                          currentHeartbeatNotifierRewardsPoolBalance
                        )
                      ).to.be.equal(expectedReward)
                    })
                  }

                  context(
                    "when there are multiple subsequent failed members indices",
                    async () => {
                      await assertNotifyFailedHeartbeatSucceed(
                        subsequentFailedMembersIndices,
                        groupThreshold,
                        (signatures) => signatures,
                        (signingMembersIndices) => signingMembersIndices
                      )
                    }
                  )

                  context(
                    "when there is only one failed members index",
                    async () => {
                      await assertNotifyFailedHeartbeatSucceed(
                        [32],
                        groupThreshold,
                        (signatures) => signatures,
                        (signingMembersIndices) => signingMembersIndices
                      )
                    }
                  )

                  context(
                    "when there are multiple non-subsequent failed members indices",
                    async () => {
                      await assertNotifyFailedHeartbeatSucceed(
                        nonSubsequentFailedMembersIndices,
                        groupThreshold,
                        (signatures) => signatures,
                        (signingMembersIndices) => signingMembersIndices
                      )
                    }
                  )

                  context(
                    "when there are multiple non-subsequent signing members indices",
                    async () => {
                      const newSigningMembersIndices = [
                        1, 5, 8, 11, 14, 15, 18, 20, 22, 24, 25, 27, 29, 30, 31,
                        33, 38, 39, 41, 42, 44, 47, 48, 49, 51, 53, 55, 56, 57,
                        59, 61, 62, 64,
                      ]

                      const getSignature = (signatures, index) =>
                        signatures
                          .slice(2)
                          .slice(130 * index, 130 * index + 130)

                      const modifySignatures = (signatures) => {
                        let newSignatures = "0x"

                        for (
                          let i = 0;
                          i < newSigningMembersIndices.length;
                          i++
                        ) {
                          const newSigningMemberIndex =
                            newSigningMembersIndices[i]
                          newSignatures += getSignature(
                            signatures,
                            newSigningMemberIndex - 1
                          )
                        }

                        return newSignatures
                      }

                      await assertNotifyFailedHeartbeatSucceed(
                        subsequentFailedMembersIndices,
                        // Make more signatures than needed to allow picking up
                        // arbitrary signatures.
                        64,
                        modifySignatures,
                        (_) => newSigningMembersIndices
                      )
                    }
                  )
                })

                context(
                  "when claim sender did not sign the claim",
                  async () => {
                    it("should revert", async () => {
                      const { signatures, signingMembersIndices } =
                        await signHeartbeatFailureClaim(
                          members,
                          0,
                          group.groupPubKey,
                          subsequentFailedMembersIndices,
                          groupThreshold
                        )

                      // Assume claim sender is member `34` - the first member
                      // who did not sign the claim. We take index `33` since
                      // `members` array is zero-based.
                      const claimSender = await ethers.getSigner(
                        members[33].address
                      )

                      await expect(
                        randomBeacon.connect(claimSender).notifyFailedHeartbeat(
                          {
                            groupId,
                            failedMembersIndices:
                              subsequentFailedMembersIndices,
                            signatures,
                            signingMembersIndices,
                          },
                          0,
                          membersIDs
                        )
                      ).to.be.revertedWith("Sender must be claim signer")
                    })
                  }
                )
              })

              context("when one of the signatures is incorrect", () => {
                const assertInvalidSignature = async (invalidSignature) => {
                  // The 32 signers sign correct parameters. Invalid signature
                  // is expected to be provided by signer 33.
                  const { signatures, signingMembersIndices } =
                    await signHeartbeatFailureClaim(
                      members,
                      0,
                      group.groupPubKey,
                      subsequentFailedMembersIndices,
                      groupThreshold - 1
                    )

                  await expect(
                    randomBeacon.notifyFailedHeartbeat(
                      {
                        groupId,
                        failedMembersIndices: subsequentFailedMembersIndices,
                        // Slice removes `0x` prefix from wrong signature.
                        signatures: signatures + invalidSignature.slice(2),
                        signingMembersIndices: [...signingMembersIndices, 33],
                      },
                      0,
                      membersIDs
                    )
                  ).to.be.revertedWith("Invalid signature")
                }

                context(
                  "when one of the signatures signed the wrong nonce",
                  () => {
                    it("should revert", async () => {
                      // Signer 33 signs wrong nonce.
                      const invalidSignature = (
                        await signHeartbeatFailureClaim(
                          [members[32]],
                          1,
                          group.groupPubKey,
                          subsequentFailedMembersIndices,
                          1
                        )
                      ).signatures

                      await assertInvalidSignature(invalidSignature)
                    })
                  }
                )

                context(
                  "when one of the signatures signed the wrong group public key",
                  () => {
                    it("should revert", async () => {
                      // Signer 33 signs wrong group public key.
                      const invalidSignature = (
                        await signHeartbeatFailureClaim(
                          [members[32]],
                          0,
                          "0x010203",
                          subsequentFailedMembersIndices,
                          1
                        )
                      ).signatures

                      await assertInvalidSignature(invalidSignature)
                    })
                  }
                )

                context(
                  "when one of the signatures signed the wrong failed group members indices",
                  () => {
                    it("should revert", async () => {
                      // Signer 33 signs wrong failed group members indices.
                      const invalidSignature = (
                        await signHeartbeatFailureClaim(
                          [members[32]],
                          0,
                          group.groupPubKey,
                          [1, 2, 3, 4, 5, 6, 7, 8],
                          1
                        )
                      ).signatures

                      await assertInvalidSignature(invalidSignature)
                    })
                  }
                )
              })
            })

            context("when signing members indices are incorrect", () => {
              context(
                "when signing members indices count is different than signatures count",
                () => {
                  it("should revert", async () => {
                    const { signatures, signingMembersIndices } =
                      await signHeartbeatFailureClaim(
                        members,
                        0,
                        group.groupPubKey,
                        subsequentFailedMembersIndices,
                        groupThreshold
                      )

                    await expect(
                      randomBeacon.notifyFailedHeartbeat(
                        {
                          groupId,
                          failedMembersIndices: subsequentFailedMembersIndices,
                          signatures,
                          // Remove the first signing member index
                          signingMembersIndices: signingMembersIndices.slice(1),
                        },
                        0,
                        membersIDs
                      )
                    ).to.be.revertedWith("Unexpected signatures count")
                  })
                }
              )

              context("when first signing member index is zero", () => {
                it("should revert", async () => {
                  const { signatures, signingMembersIndices } =
                    await signHeartbeatFailureClaim(
                      members,
                      0,
                      group.groupPubKey,
                      subsequentFailedMembersIndices,
                      groupThreshold
                    )

                  signingMembersIndices[0] = 0

                  await expect(
                    randomBeacon.notifyFailedHeartbeat(
                      {
                        groupId,
                        failedMembersIndices: subsequentFailedMembersIndices,
                        signatures,
                        signingMembersIndices,
                      },
                      0,
                      membersIDs
                    )
                  ).to.be.revertedWith("Corrupted members indices")
                })
              })

              context(
                "when last signing member index is bigger than group size",
                () => {
                  it("should revert", async () => {
                    const { signatures, signingMembersIndices } =
                      await signHeartbeatFailureClaim(
                        members,
                        0,
                        group.groupPubKey,
                        subsequentFailedMembersIndices,
                        groupThreshold
                      )

                    signingMembersIndices[signingMembersIndices.length - 1] = 65

                    await expect(
                      randomBeacon.notifyFailedHeartbeat(
                        {
                          groupId,
                          failedMembersIndices: subsequentFailedMembersIndices,
                          signatures,
                          signingMembersIndices,
                        },
                        0,
                        membersIDs
                      )
                    ).to.be.revertedWith("Corrupted members indices")
                  })
                }
              )

              context(
                "when signing members indices are not ordered in ascending order",
                () => {
                  it("should revert", async () => {
                    const { signatures, signingMembersIndices } =
                      await signHeartbeatFailureClaim(
                        members,
                        0,
                        group.groupPubKey,
                        subsequentFailedMembersIndices,
                        groupThreshold
                      )

                    // eslint-disable-next-line prefer-destructuring
                    signingMembersIndices[10] = signingMembersIndices[11]

                    await expect(
                      randomBeacon.notifyFailedHeartbeat(
                        {
                          groupId,
                          failedMembersIndices: subsequentFailedMembersIndices,
                          signatures,
                          signingMembersIndices,
                        },
                        0,
                        membersIDs
                      )
                    ).to.be.revertedWith("Corrupted members indices")
                  })
                }
              )
            })
          })

          context("when signatures array is incorrect", () => {
            context("when signatures count is zero", () => {
              it("should revert", async () => {
                const signatures = "0x"

                await expect(
                  randomBeacon.notifyFailedHeartbeat(
                    {
                      groupId,
                      failedMembersIndices: subsequentFailedMembersIndices,
                      signatures,
                      signingMembersIndices: stubMembersIndices,
                    },
                    0,
                    membersIDs
                  )
                ).to.be.revertedWith("No signatures provided")
              })
            })

            context(
              "when signatures count is not divisible by signature byte size",
              () => {
                it("should revert", async () => {
                  const signatures = "0x010203"

                  await expect(
                    randomBeacon.notifyFailedHeartbeat(
                      {
                        groupId,
                        failedMembersIndices: subsequentFailedMembersIndices,
                        signatures,
                        signingMembersIndices: stubMembersIndices,
                      },
                      0,
                      membersIDs
                    )
                  ).to.be.revertedWith("Malformed signatures array")
                })
              }
            )

            context(
              "when signatures count is different than signing members count",
              () => {
                it("should revert", async () => {
                  const { signatures, signingMembersIndices } =
                    await signHeartbeatFailureClaim(
                      members,
                      0,
                      group.groupPubKey,
                      subsequentFailedMembersIndices,
                      groupThreshold
                    )

                  await expect(
                    randomBeacon.notifyFailedHeartbeat(
                      {
                        groupId,
                        failedMembersIndices: subsequentFailedMembersIndices,
                        // Remove the first signature to cause a mismatch with
                        // the signing members count.
                        signatures: `0x${signatures.slice(132)}`,
                        signingMembersIndices,
                      },
                      0,
                      membersIDs
                    )
                  ).to.be.revertedWith("Unexpected signatures count")
                })
              }
            )

            context(
              "when signatures count is less than group threshold",
              () => {
                it("should revert", async () => {
                  const { signatures, signingMembersIndices } =
                    await signHeartbeatFailureClaim(
                      members,
                      0,
                      group.groupPubKey,
                      subsequentFailedMembersIndices,
                      // Provide one signature too few.
                      groupThreshold - 1
                    )

                  await expect(
                    randomBeacon.notifyFailedHeartbeat(
                      {
                        groupId,
                        failedMembersIndices: subsequentFailedMembersIndices,
                        signatures,
                        signingMembersIndices,
                      },
                      0,
                      membersIDs
                    )
                  ).to.be.revertedWith("Too few signatures")
                })
              }
            )

            context("when signatures count is bigger than group size", () => {
              it("should revert", async () => {
                const { signatures, signingMembersIndices } =
                  await signHeartbeatFailureClaim(
                    members,
                    0,
                    group.groupPubKey,
                    subsequentFailedMembersIndices,
                    // All group signs.
                    members.length
                  )

                await expect(
                  randomBeacon.notifyFailedHeartbeat(
                    {
                      groupId,
                      failedMembersIndices: subsequentFailedMembersIndices,
                      // Provide one signature too much.
                      signatures: signatures + signatures.slice(2, 132),
                      signingMembersIndices: [
                        ...signingMembersIndices,
                        signingMembersIndices[0],
                      ],
                    },
                    0,
                    membersIDs
                  )
                ).to.be.revertedWith("Too many signatures")
              })
            })
          })
        })

        context("when failed members indices are incorrect", () => {
          const assertFailedMembersIndicesCorrupted = async (
            failedMembersIndices: number[]
          ) => {
            const { signatures, signingMembersIndices } =
              await signHeartbeatFailureClaim(
                members,
                0,
                group.groupPubKey,
                failedMembersIndices,
                groupThreshold
              )

            await expect(
              randomBeacon.notifyFailedHeartbeat(
                {
                  groupId,
                  failedMembersIndices,
                  signatures,
                  signingMembersIndices,
                },
                0,
                membersIDs
              )
            ).to.be.revertedWith("Corrupted members indices")
          }

          context("when failed members indices count is zero", () => {
            it("should revert", async () => {
              const failedMembersIndices = []

              await assertFailedMembersIndicesCorrupted(failedMembersIndices)
            })
          })

          context(
            "when failed members indices count is bigger than group size",
            () => {
              it("should revert", async () => {
                const failedMembersIndices = Array.from(
                  Array(65),
                  (_, i) => i + 1
                )

                await assertFailedMembersIndicesCorrupted(failedMembersIndices)
              })
            }
          )

          context("when first failed member index is zero", () => {
            it("should revert", async () => {
              const failedMembersIndices = Array.from(
                Array(64),
                (_, i) => i + 1
              )
              failedMembersIndices[0] = 0

              await assertFailedMembersIndicesCorrupted(failedMembersIndices)
            })
          })

          context(
            "when last failed member index is bigger than group size",
            () => {
              it("should revert", async () => {
                const failedMembersIndices = Array.from(
                  Array(64),
                  (_, i) => i + 1
                )
                failedMembersIndices[failedMembersIndices.length - 1] = 65

                await assertFailedMembersIndicesCorrupted(failedMembersIndices)
              })
            }
          )

          context(
            "when failed members indices are not ordered in ascending order",
            () => {
              it("should revert", async () => {
                const failedMembersIndices = Array.from(
                  Array(64),
                  (_, i) => i + 1
                )
                // eslint-disable-next-line prefer-destructuring
                failedMembersIndices[10] = failedMembersIndices[11]

                await assertFailedMembersIndicesCorrupted(failedMembersIndices)
              })
            }
          )
        })
      })

      context("when group is active but terminated", () => {
        before(async () => {
          await createSnapshot()

          // Simulate group was terminated.
          await (randomBeacon as RandomBeaconStub).roughlyTerminateGroup(
            groupId
          )
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should revert", async () => {
          await expect(
            randomBeacon.notifyFailedHeartbeat(
              {
                groupId,
                failedMembersIndices: stubMembersIndices,
                signatures: stubSignatures,
                signingMembersIndices: stubMembersIndices,
              },
              0,
              membersIDs
            )
          ).to.be.revertedWith("Group must be active and non-terminated")
        })
      })

      context("when group is expired", () => {
        before(async () => {
          await createSnapshot()

          // Set a short value of group lifetime to avoid long test execution
          // if original value is used.
          const newGroupLifetime = 10
          await randomBeacon.updateGroupCreationParameters(
            params.groupCreationFrequency,
            newGroupLifetime
          )
          // Simulate group was expired.
          await mineBlocks(newGroupLifetime)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should revert", async () => {
          await expect(
            randomBeacon.notifyFailedHeartbeat(
              {
                groupId,
                failedMembersIndices: stubMembersIndices,
                signatures: stubSignatures,
                signingMembersIndices: stubMembersIndices,
              },
              0,
              membersIDs
            )
          ).to.be.revertedWith("Group must be active and non-terminated")
        })
      })
    })

    context("when passed nonce is invalid", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon.notifyFailedHeartbeat(
            {
              groupId,
              failedMembersIndices: stubMembersIndices,
              signatures: stubSignatures,
              signingMembersIndices: stubMembersIndices,
            },
            1,
            membersIDs
          ) // Initial nonce is `0`.
        ).to.be.revertedWith("Invalid nonce")
      })
    })

    context("when group members are invalid", () => {
      it("should revert", async () => {
        const invalidMembersId = [0, 1, 42]
        await expect(
          randomBeacon.notifyFailedHeartbeat(
            {
              groupId,
              failedMembersIndices: stubMembersIndices,
              signatures: stubSignatures,
              signingMembersIndices: stubMembersIndices,
            },
            0,
            invalidMembersId
          )
        ).to.be.revertedWith("Invalid group members")
      })
    })
  })

  async function approveTestToken() {
    await testToken.mint(requester.address, params.relayRequestFee)
    await testToken
      .connect(requester)
      .approve(randomBeacon.address, params.relayRequestFee)
  }

  async function fundHeartbeatNotifierRewardsPool(donateAmount: BigNumber) {
    await testToken.mint(deployer.address, donateAmount)
    await testToken
      .connect(deployer)
      .approve(randomBeacon.address, donateAmount)

    await randomBeacon.fundHeartbeatNotifierRewardsPool(
      deployer.address,
      donateAmount
    )
  }
})
