/* eslint-disable @typescript-eslint/no-unused-expressions, no-await-in-loop */

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
import { createGroup } from "./utils/groups"
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

const fixture = async () => {
  const deployment = await randomBeaconDeployment()

  const operators = await registerOperators(
    deployment.randomBeacon as RandomBeacon,
    (await getUnnamedAccounts()).slice(0, constants.groupSize)
  )

  const bls = await blsDeployment()

  return {
    randomBeacon: deployment.randomBeacon as RandomBeacon,
    sortitionPool: deployment.sortitionPool as SortitionPool,
    testToken: deployment.testToken as TestToken,
    staking: deployment.stakingStub as StakingStub,
    relayStub: (await (
      await ethers.getContractFactory("RelayStub")
    ).deploy()) as RelayStub,
    operators,
    bls: bls.bls as BLS,
  }
}

describe("RandomBeacon - Relay", () => {
  const relayRequestFee = to1e18(100)
  const ineligibleOperatorNotifierReward = to1e18(200)

  let deployer: SignerWithAddress
  let requester: SignerWithAddress
  let notifier: SignerWithAddress
  let submitter: SignerWithAddress
  let members: Operator[]
  let membersIDs: OperatorID[]
  let membersAddresses: Address[]
  let bls: BLS

  let randomBeacon: RandomBeacon
  let sortitionPool: SortitionPool
  let testToken: TestToken
  let staking: StakingStub
  let relayStub: RelayStub

  before(async () => {
    deployer = await ethers.getSigner((await getNamedAccounts()).deployer)
    requester = await ethers.getSigner((await getUnnamedAccounts())[1])
    notifier = await ethers.getSigner((await getUnnamedAccounts())[2])
    submitter = await ethers.getSigner((await getUnnamedAccounts())[3])
  })

  beforeEach("load test fixture", async () => {
    let operators
      // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;({
      randomBeacon,
      sortitionPool,
      testToken,
      staking,
      relayStub,
      operators,
      bls,
    } = await waffle.loadFixture(fixture))

    members = operators // All operators will be members of the group used in tests.
    membersIDs = members.map((member) => member.id)
    membersAddresses = members.map((member) => member.address)

    await randomBeacon.updateRelayEntryParameters(to1e18(100), 10, 5760, 0)
    // groupLifetime: 64 * 10 * 5760 + 100
    // if the group lifetime is less than 6400 it will be expired on selection.
    await randomBeacon.updateGroupCreationParameters(100, 6500)
    await randomBeacon.updateRewardParameters(
      0,
      0,
      ineligibleOperatorNotifierReward,
      1209600,
      5,
      5,
      100
    )
  })

  describe("requestRelayEntry", () => {
    context("when groups exist", () => {
      beforeEach(async () => {
        await createGroup(randomBeacon as RandomBeaconStub, members)
      })

      context("when there is no other relay entry in progress", () => {
        context("when the requester pays the relay request fee", () => {
          let tx: ContractTransaction
          let previousDkgRewardsPoolBalance: BigNumber
          let previousRandomBeaconBalance: BigNumber

          beforeEach(async () => {
            previousDkgRewardsPoolBalance = await randomBeacon.dkgRewardsPool()
            previousRandomBeaconBalance = await testToken.balanceOf(
              randomBeacon.address
            )
            await approveTestToken()
          })

          context(
            "when relay request does not hit group creation frequency threshold",
            () => {
              beforeEach(async () => {
                tx = await randomBeacon
                  .connect(requester)
                  .requestRelayEntry(ZERO_ADDRESS)
              })

              it("should deposit relay request fee to the DKG rewards pool", async () => {
                // Assert correct pool bookkeeping.
                const currentDkgRewardsPoolBalance =
                  await randomBeacon.dkgRewardsPool()
                expect(
                  currentDkgRewardsPoolBalance.sub(
                    previousDkgRewardsPoolBalance
                  )
                ).to.be.equal(relayRequestFee)

                // Assert actual transfer took place.
                const currentRandomBeaconBalance = await testToken.balanceOf(
                  randomBeacon.address
                )
                expect(
                  currentRandomBeaconBalance.sub(previousRandomBeaconBalance)
                ).to.be.equal(relayRequestFee)
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
              beforeEach(async () => {
                // Force group creation on each relay entry.
                await randomBeacon
                  .connect(deployer)
                  .updateGroupCreationParameters(1, params.groupLifeTime)

                tx = await randomBeacon
                  .connect(requester)
                  .requestRelayEntry(ZERO_ADDRESS)
              })

              it("should deposit relay request fee to the DKG rewards pool", async () => {
                // Assert correct pool bookkeeping.
                const currentDkgRewardsPoolBalance =
                  await randomBeacon.dkgRewardsPool()
                expect(
                  currentDkgRewardsPoolBalance.sub(
                    previousDkgRewardsPoolBalance
                  )
                ).to.be.equal(relayRequestFee)

                // Assert actual transfer took place.
                const currentRandomBeaconBalance = await testToken.balanceOf(
                  randomBeacon.address
                )
                expect(
                  currentRandomBeaconBalance.sub(previousRandomBeaconBalance)
                ).to.be.equal(relayRequestFee)
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
        beforeEach(async () => {
          await approveTestToken()
          await randomBeacon.connect(requester).requestRelayEntry(ZERO_ADDRESS)
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

  describe("submitRelayEntry", () => {
    beforeEach(async () => {
      await createGroup(randomBeacon as RandomBeaconStub, members)
    })

    context("when relay request is in progress", () => {
      beforeEach(async () => {
        await approveTestToken()
        await randomBeacon.connect(requester).requestRelayEntry(ZERO_ADDRESS)
      })

      context("when relay entry is not timed out", () => {
        context("when entry is valid", () => {
          context("when result is submitted before the soft timeout", () => {
            let tx: ContractTransaction

            beforeEach(async () => {
              tx = await randomBeacon
                .connect(submitter)
                .submitRelayEntry(blsData.groupSignature)
            })

            it("should not slash any members", async () => {
              await expect(tx).to.not.emit(staking, "Slashed")
            })

            it("should emit RelayEntrySubmitted event", async () => {
              await expect(tx)
                .to.emit(randomBeacon, "RelayEntrySubmitted")
                .withArgs(1, 0, submitter.address, blsData.groupSignature)
            })

            it("should terminate the relay request", async () => {
              expect(await randomBeacon.isRelayRequestInProgress()).to.be.false
            })
          })

          context("when result is submitted after the soft timeout", () => {
            let tx: ContractTransaction

            beforeEach(async () => {
              // Let's assume we want to submit the relay entry after 75%
              // of the soft timeout period elapses. If so we need to
              // mine the following number of blocks:
              // `groupSize * relayEntrySubmissionEligibilityDelay +
              // (0.75 * relayEntryHardTimeout)`. However, we need to
              // subtract one block because the relay entry submission
              // transaction will move the blockchain ahead by one block
              // due to the Hardhat auto-mine feature.
              await mineBlocks(64 * 10 + 0.75 * 5760 - 1)

              tx = await randomBeacon
                .connect(submitter)
                .submitRelayEntry(blsData.groupSignature)
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
                .withArgs(1, 0, submitter.address, blsData.groupSignature)
            })

            it("should terminate the relay request", async () => {
              expect(await randomBeacon.isRelayRequestInProgress()).to.be.false
            })
          })

          context("when DKG is awaiting a seed", () => {
            let tx: ContractTransaction

            beforeEach(async () => {
              // Simulate DKG is awaiting a seed.
              await (randomBeacon as RandomBeaconStub).publicDkgLockState()

              tx = await randomBeacon
                .connect(submitter)
                .submitRelayEntry(blsData.groupSignature)
            })

            it("should emit DkgStarted event", async () => {
              await expect(tx)
                .to.emit(randomBeacon, "DkgStarted")
                .withArgs(blsData.groupSignatureUint256)
            })
          })
        })

        context("when entry is not valid", () => {
          it("should revert", async () => {
            await expect(
              randomBeacon
                .connect(submitter)
                .submitRelayEntry(blsData.nextGroupSignature)
            ).to.be.revertedWith("Invalid entry")
          })
        })
      })

      context("when relay entry is timed out", () => {
        it("should revert", async () => {
          // groupSize * relayEntrySubmissionEligibilityDelay + relayEntryHardTimeout
          await mineBlocks(64 * 10 + 5760)

          await expect(
            randomBeacon
              .connect(submitter)
              .submitRelayEntry(blsData.nextGroupSignature)
          ).to.be.revertedWith("Relay request timed out")
        })
      })
    })

    context("when relay request is not in progress", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(submitter)
            .submitRelayEntry(blsData.nextGroupSignature)
        ).to.be.revertedWith("No relay request in progress")
      })
    })
  })

  describe("reportRelayEntryTimeout", () => {
    beforeEach(async () => {
      await createGroup(randomBeacon, members)

      await approveTestToken()
      await randomBeacon.connect(requester).requestRelayEntry(ZERO_ADDRESS)
    })

    context("when relay entry timed out", () => {
      context(
        "when other active groups exist after timeout is reported",
        () => {
          let tx: ContractTransaction

          beforeEach(async () => {
            // `groupSize * relayEntrySubmissionEligibilityDelay +
            // relayEntryHardTimeout`.
            await mineBlocks(64 * 10 + 5760)

            await (randomBeacon as RandomBeaconStub).roughlyAddGroup(
              "0x01",
              membersIDs
            )

            tx = await randomBeacon.connect(notifier).reportRelayEntryTimeout()
          })

          it("should slash the full slashing amount for all group members", async () => {
            await expect(tx)
              .to.emit(staking, "Seized")
              .withArgs(to1e18(1000), 5, notifier.address, membersAddresses)

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

          beforeEach(async () => {
            // Create another group in a rough way just to have an active group
            // once the one handling the timed out request gets terminated.
            // This makes the request retry possible. That group will not
            // perform any signing so their public key can be arbitrary bytes.
            // Also, that group is created just after the relay request is
            // made to ensure it is not selected for signing the original request.
            await (randomBeacon as RandomBeaconStub).roughlyAddGroup(
              "0x01",
              membersIDs
            )
          })

          it("should clean up current relay request data", async () => {
            // `groupSize * relayEntrySubmissionEligibilityDelay +
            // relayEntryHardTimeout`. This times out the relay entry
            await mineBlocks(64 * 10 + 5760)

            const registry = await randomBeacon.getGroupsRegistry()
            const secondGroupLifetime = await (
              randomBeacon as RandomBeaconStub
            ).groupLifetimeOf(registry[1])

            // Expire second group
            await mineBlocksTo(Number(secondGroupLifetime) + 1)

            tx = await randomBeacon.reportRelayEntryTimeout()

            await expect(tx).to.not.emit(randomBeacon, "RelayEntryRequested")
            expect(await randomBeacon.isRelayRequestInProgress()).to.be.false
          })
        }
      )

      context("when no active groups exist after timeout is reported", () => {
        let tx: ContractTransaction

        beforeEach(async () => {
          // `groupSize * relayEntrySubmissionEligibilityDelay +
          // relayEntryHardTimeout`.
          await mineBlocks(64 * 10 + 5760)

          tx = await randomBeacon.connect(notifier).reportRelayEntryTimeout()
        })

        it("should slash the full slashing amount for all group members", async () => {
          await expect(tx)
            .to.emit(staking, "Seized")
            .withArgs(to1e18(1000), 5, notifier.address, membersAddresses)

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
          beforeEach(async () => {
            // `groupSize * relayEntrySubmissionEligibilityDelay +
            // relayEntryHardTimeout`.
            await mineBlocks(64 * 10 + 5760)

            // Simulate DKG is awaiting a seed.
            await (randomBeacon as RandomBeaconStub).publicDkgLockState()

            await randomBeacon.connect(notifier).reportRelayEntryTimeout()
          })

          it("should notify DKG seed timed out", async () => {
            // TODO: Uncomment those assertions once termination is implemented.
            // expect(await randomBeacon.getGroupCreationState()).to.be.equal(
            //   dkgState.IDLE
            // )
            // expect(await sortitionPool.isLocked()).to.be.false
          })

          it("should emit DkgSeedTimedOut event", async () => {
            // TODO: Uncomment those assertions once termination is implemented.
            // await expect(tx).to.emit(randomBeacon, "DkgSeedTimedOut")
          })
        }
      )
    })

    context("when relay entry did not time out", () => {
      it("should revert", async () => {
        await expect(randomBeacon.reportRelayEntryTimeout()).to.be.revertedWith(
          "Relay request did not time out"
        )
      })
    })
  })

  describe("reportUnauthorizedSigning", () => {
    beforeEach(async () => {
      await createGroup(randomBeacon as RandomBeaconStub, members)

      await approveTestToken()
      await randomBeacon.connect(requester).requestRelayEntry(ZERO_ADDRESS)
    })

    context("when a group is active", () => {
      context("when provided signature is valid", () => {
        let tx
        beforeEach(async () => {
          const notifierSignature = await bls.sign(
            notifier.address,
            blsData.secretKey
          )
          tx = await randomBeacon
            .connect(notifier)
            .reportUnauthorizedSigning(notifierSignature, 0)
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
            .withArgs(to1e18(100000), 5, notifier.address, membersAddresses)
        })

        it("should emit unauthorized signing slashing event", async () => {
          await expect(tx)
            .to.emit(randomBeacon, "UnauthorizedSigningSlashed")
            .withArgs(0, to1e18(100000), membersAddresses)
        })
      })
    })

    context("when group is terminated", () => {
      it("should revert", async () => {
        // `groupSize * relayEntrySubmissionEligibilityDelay +
        // relayEntryHardTimeout`.
        await mineBlocks(64 * 10 + 5760)
        await (randomBeacon as RandomBeaconStub).roughlyTerminateGroup(0)

        const notifierSignature = await bls.sign(
          notifier.address,
          blsData.secretKey
        )
        await expect(
          randomBeacon
            .connect(notifier)
            .reportUnauthorizedSigning(notifierSignature, 0)
        ).to.be.revertedWith("Group cannot be terminated")
      })
    })

    context("when provided signature is not valid", () => {
      it("should revert", async () => {
        // the valid key is 123 instead of 42
        const notifierSignature = await bls.sign(notifier.address, 42)
        await expect(
          randomBeacon
            .connect(notifier)
            .reportUnauthorizedSigning(notifierSignature, 0)
        ).to.be.revertedWith("Invalid signature")
      })
    })
  })

  describe("getSlashingFactor", () => {
    const testGroupSize = 8

    beforeEach(async () => {
      await relayStub.setCurrentRequestStartBlock()
    })

    context("when soft timeout has not been exceeded yet", () => {
      it("should return a slashing factor equal to zero", async () => {
        // `groupSize * relayEntrySubmissionEligibilityDelay`
        await mineBlocks(8 * 10)

        expect(await relayStub.getSlashingFactor(testGroupSize)).to.be.equal(0)
      })
    })

    context("when soft timeout has been exceeded by one block", () => {
      it("should return a correct slashing factor", async () => {
        // `groupSize * relayEntrySubmissionEligibilityDelay + 1 block`
        await mineBlocks(8 * 10 + 1)

        // We are exceeded the soft timeout by `1` block so this is the
        // `submissionDelay` factor. If so we can calculate the slashing factor
        // as `(submissionDelay * 1e18) / relayEntryHardTimeout` which
        // gives `1 * 1e18 / 5760 = 173611111111111` (0.017%).
        expect(await relayStub.getSlashingFactor(testGroupSize)).to.be.equal(
          BigNumber.from("173611111111111")
        )
      })
    })

    context(
      "when soft timeout has been exceeded by the number of blocks equal to the hard timeout",
      () => {
        it("should return a correct slashing factor", async () => {
          // `groupSize * relayEntrySubmissionEligibilityDelay + relayEntryHardTimeout`
          await mineBlocks(8 * 10 + 5760)

          // We are exceeded the soft timeout by `5760` blocks so this is the
          // `submissionDelay` factor. If so we can calculate the slashing
          // factor as `(submissionDelay * 1e18) / relayEntryHardTimeout` which
          // gives `5760 * 1e18 / 5760 = 1000000000000000000` (100%).
          expect(await relayStub.getSlashingFactor(testGroupSize)).to.be.equal(
            BigNumber.from("1000000000000000000")
          )
        })
      }
    )

    context(
      "when soft timeout has been exceeded by the number of blocks bigger than the hard timeout",
      () => {
        it("should return a correct slashing factor", async () => {
          // `groupSize * relayEntrySubmissionEligibilityDelay +
          // relayEntryHardTimeout + 1 block`.
          await mineBlocks(8 * 10 + 5760 + 1)

          // We are exceeded the soft timeout by a value bigger than the
          // hard timeout. In that case the maximum value (100%) of the slashing
          // factor should be returned.
          expect(await relayStub.getSlashingFactor(testGroupSize)).to.be.equal(
            BigNumber.from("1000000000000000000")
          )
        })
      }
    )
  })

  describe("fundHeartbeatNotifierRewardsPool", () => {
    const amount = to1e18(1000)

    let previousHeartbeatNotifierRewardsPoolBalance: BigNumber
    let previousRandomBeaconBalance: BigNumber

    beforeEach(async () => {
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
    const validFailedMembersIndices = Array.from(Array(31), (_, i) => i + 1)
    const groupThreshold = 33

    let group

    beforeEach(async () => {
      await createGroup(randomBeacon as RandomBeaconStub, members)
      group = await randomBeacon["getGroup(uint64)"](groupId)
    })

    context("when group is active and non-terminated", () => {
      context("when failed members indices are correct", () => {
        context("when signatures array is correct", () => {
          context("when signing members indices are correct", () => {
            context("when all signatures are correct", () => {
              let tx: ContractTransaction
              let nonce: BigNumber
              let initialNotifierBalance: BigNumber
              let initialHeartbeatNotifierRewardsPoolBalance: BigNumber

              beforeEach(async () => {
                await donateHeartbeatNotifierRewardsPool(
                  ineligibleOperatorNotifierReward.mul(
                    validFailedMembersIndices.length
                  )
                )

                nonce = await randomBeacon.failedHeartbeatNonce(groupId)

                initialNotifierBalance = await testToken.balanceOf(
                  notifier.address
                )

                initialHeartbeatNotifierRewardsPoolBalance =
                  await randomBeacon.heartbeatNotifierRewardsPool()

                const { signatures, signingMembersIndices } =
                  await signHeartbeatFailureClaim(
                    members,
                    nonce.toNumber(),
                    group.groupPubKey,
                    validFailedMembersIndices,
                    groupThreshold
                  )

                tx = await randomBeacon
                  .connect(notifier)
                  .notifyFailedHeartbeat({
                    groupId,
                    failedMembersIndices: validFailedMembersIndices,
                    signatures,
                    signingMembersIndices,
                  })
              })

              it("should increment failed heartbeat nonce for the group", async () => {
                expect(
                  await randomBeacon.failedHeartbeatNonce(groupId)
                ).to.be.equal(nonce.add(1))
              })

              it("should emit FailedHeartbeatNotified event", async () => {
                await expect(tx)
                  .to.emit(randomBeacon, "FailedHeartbeatNotified")
                  .withArgs(
                    groupId,
                    nonce.toNumber(),
                    membersIDs.slice(0, 31),
                    notifier.address
                  )
              })

              it("should ban sortition pool rewards for ineligible operators", async () => {
                const now = await helpers.time.lastBlockTime()
                const expectedUntil = now + 1209600 // 2 weeks

                await expect(tx)
                  .to.emit(sortitionPool, "IneligibleForRewards")
                  .withArgs(membersIDs.slice(0, 31), expectedUntil)
              })

              it("should pay notifier reward from heartbeat notifier rewards pool", async () => {
                const expectedReward = ineligibleOperatorNotifierReward.mul(
                  validFailedMembersIndices.length
                )

                const currentNotifierBalance = await testToken.balanceOf(
                  notifier.address
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
                    validFailedMembersIndices,
                    groupThreshold - 1
                  )

                await expect(
                  randomBeacon.notifyFailedHeartbeat({
                    groupId,
                    failedMembersIndices: validFailedMembersIndices,
                    // Slice removes `0x` prefix from wrong signature.
                    signatures: signatures + invalidSignature.slice(2),
                    signingMembersIndices: [...signingMembersIndices, 33],
                  })
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
                        validFailedMembersIndices,
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
                        validFailedMembersIndices,
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
                      validFailedMembersIndices,
                      groupThreshold
                    )

                  await expect(
                    randomBeacon.notifyFailedHeartbeat({
                      groupId,
                      failedMembersIndices: validFailedMembersIndices,
                      signatures,
                      // Remove the first signing member index
                      signingMembersIndices: signingMembersIndices.slice(1),
                    })
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
                    validFailedMembersIndices,
                    groupThreshold
                  )

                signingMembersIndices[0] = 0

                await expect(
                  randomBeacon.notifyFailedHeartbeat({
                    groupId,
                    failedMembersIndices: validFailedMembersIndices,
                    signatures,
                    signingMembersIndices,
                  })
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
                      validFailedMembersIndices,
                      groupThreshold
                    )

                  signingMembersIndices[signingMembersIndices.length - 1] = 65

                  await expect(
                    randomBeacon.notifyFailedHeartbeat({
                      groupId,
                      failedMembersIndices: validFailedMembersIndices,
                      signatures,
                      signingMembersIndices,
                    })
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
                      validFailedMembersIndices,
                      groupThreshold
                    )

                  // eslint-disable-next-line prefer-destructuring
                  signingMembersIndices[10] = signingMembersIndices[11]

                  await expect(
                    randomBeacon.notifyFailedHeartbeat({
                      groupId,
                      failedMembersIndices: validFailedMembersIndices,
                      signatures,
                      signingMembersIndices,
                    })
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
                randomBeacon.notifyFailedHeartbeat({
                  groupId,
                  failedMembersIndices: validFailedMembersIndices,
                  signatures,
                  signingMembersIndices: stubMembersIndices,
                })
              ).to.be.revertedWith("No signatures provided")
            })
          })

          context(
            "when signatures count is not divisible by signature byte size",
            () => {
              it("should revert", async () => {
                const signatures = "0x010203"

                await expect(
                  randomBeacon.notifyFailedHeartbeat({
                    groupId,
                    failedMembersIndices: validFailedMembersIndices,
                    signatures,
                    signingMembersIndices: stubMembersIndices,
                  })
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
                    validFailedMembersIndices,
                    groupThreshold
                  )

                await expect(
                  randomBeacon.notifyFailedHeartbeat({
                    groupId,
                    failedMembersIndices: validFailedMembersIndices,
                    // Remove the first signature to cause a mismatch with
                    // the signing members count.
                    signatures: `0x${signatures.slice(132)}`,
                    signingMembersIndices,
                  })
                ).to.be.revertedWith("Unexpected signatures count")
              })
            }
          )

          context("when signatures count is less than group threshold", () => {
            it("should revert", async () => {
              const { signatures, signingMembersIndices } =
                await signHeartbeatFailureClaim(
                  members,
                  0,
                  group.groupPubKey,
                  validFailedMembersIndices,
                  // Provide one signature too few.
                  groupThreshold - 1
                )

              await expect(
                randomBeacon.notifyFailedHeartbeat({
                  groupId,
                  failedMembersIndices: validFailedMembersIndices,
                  signatures,
                  signingMembersIndices,
                })
              ).to.be.revertedWith("Too few signatures")
            })
          })

          context("when signatures count is bigger than group size", () => {
            it("should revert", async () => {
              const { signatures, signingMembersIndices } =
                await signHeartbeatFailureClaim(
                  members,
                  0,
                  group.groupPubKey,
                  validFailedMembersIndices,
                  // All group signs.
                  members.length
                )

              await expect(
                randomBeacon.notifyFailedHeartbeat({
                  groupId,
                  failedMembersIndices: validFailedMembersIndices,
                  // Provide one signature too much.
                  signatures: signatures + signatures.slice(2, 132),
                  signingMembersIndices: [
                    ...signingMembersIndices,
                    signingMembersIndices[0],
                  ],
                })
              ).to.be.revertedWith("Too many signatures")
            })
          })
        })
      })

      context("when failed members indices are incorrect", () => {
        const assertFailedMembersIndicesCorrupted = async (
          failedMembersIndices
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
            randomBeacon.notifyFailedHeartbeat({
              groupId,
              failedMembersIndices,
              signatures,
              signingMembersIndices,
            })
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
            const failedMembersIndices = Array.from(Array(64), (_, i) => i + 1)
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

    context("when group active but terminated", () => {
      beforeEach(async () => {
        // Simulate group was terminated.
        await (randomBeacon as RandomBeaconStub).roughlyTerminateGroup(groupId)
      })

      it("should revert", async () => {
        await expect(
          randomBeacon.notifyFailedHeartbeat({
            groupId,
            failedMembersIndices: stubMembersIndices,
            signatures: stubSignatures,
            signingMembersIndices: stubMembersIndices,
          })
        ).to.be.revertedWith("Group must be active and non-terminated")
      })
    })

    context("when group is expired", () => {
      beforeEach(async () => {
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

      it("should revert", async () => {
        await expect(
          randomBeacon.notifyFailedHeartbeat({
            groupId,
            failedMembersIndices: stubMembersIndices,
            signatures: stubSignatures,
            signingMembersIndices: stubMembersIndices,
          })
        ).to.be.revertedWith("Group must be active and non-terminated")
      })
    })
  })

  async function approveTestToken() {
    await testToken.mint(requester.address, relayRequestFee)
    await testToken
      .connect(requester)
      .approve(randomBeacon.address, relayRequestFee)
  }

  async function donateHeartbeatNotifierRewardsPool(donate: BigNumber) {
    await testToken.mint(deployer.address, donate)
    await testToken.connect(deployer).approve(randomBeacon.address, donate)

    await randomBeacon.fundHeartbeatNotifierRewardsPool(
      deployer.address,
      donate
    )
  }
})
