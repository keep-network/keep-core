/* eslint-disable @typescript-eslint/no-unused-expressions, no-await-in-loop */

import {
  ethers,
  waffle,
  helpers,
  getUnnamedAccounts,
  getNamedAccounts,
} from "hardhat"
import { expect } from "chai"
import { BigNumber, ContractReceipt, ContractTransaction } from "ethers"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { Address } from "hardhat-deploy/types"
import { describe } from "mocha"
import blsData from "./data/bls"
import { constants, dkgState, params, randomBeaconDeployment } from "./fixtures"
import { createGroup } from "./utils/groups"
import type {
  RandomBeacon,
  RandomBeaconStub,
  TestToken,
  RelayStub,
  SortitionPool,
  StakingStub,
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

  return {
    randomBeacon: deployment.randomBeacon as RandomBeacon,
    sortitionPool: deployment.sortitionPool as SortitionPool,
    testToken: deployment.testToken as TestToken,
    staking: deployment.stakingStub as StakingStub,
    relayStub: (await (
      await ethers.getContractFactory("RelayStub")
    ).deploy()) as RelayStub,
    operators,
  }
}

describe("RandomBeacon - Relay", () => {
  const relayRequestFee = to1e18(100)

  // When determining the eligibility queue, the
  // `(blsData.groupSignature % 64) + 1` equation points member`16` as the first
  // eligible one. This is why we use that index as `submitRelayEntry` parameter.
  // The `submitter` signer represents that member too.
  const firstEligibleMemberIndex = 16
  // In the invalid entry scenario `(blsData.nextGroupSignature % 64) + 1`
  // gives 3 so that  member needs to submit the wrong relay entry.
  const invalidEntryFirstEligibleMemberIndex = 3
  const ineligibleOperatorNotifierReward = to1e18(200)

  let deployer: SignerWithAddress
  let requester: SignerWithAddress
  let notifier: SignerWithAddress
  let member3: SignerWithAddress
  let member16: SignerWithAddress
  let member18: SignerWithAddress
  let members: Operator[]
  let membersIDs: OperatorID[]
  let membersAddresses: Address[]

  let randomBeacon: RandomBeacon
  let sortitionPool: SortitionPool
  let testToken: TestToken
  let staking: StakingStub
  let relayStub: RelayStub

  before(async () => {
    deployer = await ethers.getSigner((await getNamedAccounts()).deployer)
    requester = await ethers.getSigner((await getUnnamedAccounts())[1])
    notifier = await ethers.getSigner((await getUnnamedAccounts())[2])
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
    } = await waffle.loadFixture(fixture))

    members = operators // All operators will be members of the group used in tests.
    membersIDs = members.map((member) => member.id)
    membersAddresses = members.map((member) => member.address)

    member3 = await ethers.getSigner(
      members[invalidEntryFirstEligibleMemberIndex - 1].address
    )
    member16 = await ethers.getSigner(
      members[firstEligibleMemberIndex - 1].address
    )
    member18 = await ethers.getSigner(
      members[firstEligibleMemberIndex + 2 - 1].address
    )

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
      5
    )
  })

  describe("requestRelayEntry", () => {
    context("when groups exist", () => {
      beforeEach(async () => {
        await createGroup(randomBeacon as RandomBeaconStub, members)
      })

      context("when there is no other relay entry in progress", () => {
        context("when the requester pays the relay request fee", () => {
          let tx
          let previousMaintenancePoolBalance

          beforeEach(async () => {
            previousMaintenancePoolBalance = await testToken.balanceOf(
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

              it("should deposit relay request fee to the maintenance pool", async () => {
                const currentMaintenancePoolBalance = await testToken.balanceOf(
                  randomBeacon.address
                )
                expect(
                  currentMaintenancePoolBalance.sub(
                    previousMaintenancePoolBalance
                  )
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

              it("should deposit relay request fee to the maintenance pool", async () => {
                const currentMaintenancePoolBalance = await testToken.balanceOf(
                  randomBeacon.address
                )
                expect(
                  currentMaintenancePoolBalance.sub(
                    previousMaintenancePoolBalance
                  )
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
    let requestStartBlock: number

    beforeEach(async () => {
      await createGroup(randomBeacon as RandomBeaconStub, members)
    })

    context("when relay request is in progress", () => {
      beforeEach(async () => {
        await approveTestToken()
        const tx = await randomBeacon
          .connect(requester)
          .requestRelayEntry(ZERO_ADDRESS)
        const receipt = await tx.wait()
        requestStartBlock = receipt.blockNumber
      })

      context("when relay entry is not timed out", () => {
        context("when submitter index is valid", () => {
          context("when entry is valid", () => {
            context("when result is submitted before the soft timeout", () => {
              let tx: ContractTransaction
              let receipt: ContractReceipt

              beforeEach(async () => {
                tx = await randomBeacon
                  .connect(member16)
                  .submitRelayEntry(
                    firstEligibleMemberIndex,
                    blsData.groupSignature
                  )

                receipt = await tx.wait()
              })

              it("should not slash any members", async () => {
                await expect(tx).to.not.emit(staking, "Slashed")
              })

              it("should emit RelayEntrySubmitted event", async () => {
                await expect(tx)
                  .to.emit(randomBeacon, "RelayEntrySubmitted")
                  .withArgs(1, blsData.groupSignature)
              })

              it("should terminate the relay request", async () => {
                expect(await randomBeacon.isRelayRequestInProgress()).to.be
                  .false
              })
            })

            context("when result is submitted after the soft timeout", () => {
              let tx: ContractTransaction
              let receipt: ContractReceipt

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
                  .connect(member16)
                  .submitRelayEntry(
                    firstEligibleMemberIndex,
                    blsData.groupSignature
                  )

                receipt = await tx.wait()
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
                  .withArgs(1, blsData.groupSignature)
              })

              it("should terminate the relay request", async () => {
                expect(await randomBeacon.isRelayRequestInProgress()).to.be
                  .false
              })
            })

            context(
              "when result is submitted by a member who is not yet eligible",
              () => {
                let tx: ContractTransaction

                beforeEach(async () => {
                  // Member 18 is not yet eligible as the first eligible member
                  // is member 16.
                  tx = await randomBeacon
                    .connect(member18)
                    .submitRelayEntry(
                      firstEligibleMemberIndex + 2,
                      blsData.groupSignature
                    )
                })

                // We don't repeat all checks made in the above scenarios and
                // just assert the member was able to submit the result.
                it("should succeed", async () => {
                  await expect(tx).to.emit(randomBeacon, "RelayEntrySubmitted")
                })
              }
            )

            context("when DKG is awaiting a seed", () => {
              let tx: ContractTransaction

              beforeEach(async () => {
                // Simulate DKG is awaiting a seed.
                await (randomBeacon as RandomBeaconStub).publicDkgLockState()

                tx = await randomBeacon
                  .connect(member16)
                  .submitRelayEntry(
                    firstEligibleMemberIndex,
                    blsData.groupSignature
                  )
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
                  .connect(member3)
                  .submitRelayEntry(
                    invalidEntryFirstEligibleMemberIndex,
                    blsData.nextGroupSignature
                  )
              ).to.be.revertedWith("Invalid entry")
            })
          })
        })

        context("when submitter index is beyond valid range", () => {
          it("should revert", async () => {
            await expect(
              randomBeacon
                .connect(member16)
                .submitRelayEntry(0, blsData.nextGroupSignature)
            ).to.be.revertedWith("Invalid submitter index")

            await expect(
              randomBeacon
                .connect(member16)
                .submitRelayEntry(65, blsData.nextGroupSignature)
            ).to.be.revertedWith("Invalid submitter index")
          })
        })

        context(
          "when submitter index does not correspond to sender address",
          () => {
            it("should revert", async () => {
              await expect(
                randomBeacon
                  .connect(member16)
                  .submitRelayEntry(17, blsData.nextGroupSignature)
              ).to.be.revertedWith("Unexpected submitter index")
            })
          }
        )
      })

      context("when relay entry is timed out", () => {
        it("should revert", async () => {
          // groupSize * relayEntrySubmissionEligibilityDelay + relayEntryHardTimeout
          await mineBlocks(64 * 10 + 5760)

          await expect(
            randomBeacon
              .connect(member16)
              .submitRelayEntry(
                firstEligibleMemberIndex,
                blsData.nextGroupSignature
              )
          ).to.be.revertedWith("Relay request timed out")
        })
      })
    })

    context("when relay request is not in progress", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(member16)
            .submitRelayEntry(
              firstEligibleMemberIndex,
              blsData.nextGroupSignature
            )
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
            // TODO: Implementation once `Groups` library is ready.
          })

          it("should emit RelayEntryTimedOut event", async () => {
            await expect(tx)
              .to.emit(randomBeacon, "RelayEntryTimedOut")
              .withArgs(1, 0)
          })

          it("should retry current relay request", async () => {
            // We expect the same request ID because this is a retry.
            // Group ID is still `0` because there is only one group
            // after termination was performed.
            await expect(tx)
              .to.emit(randomBeacon, "RelayEntryRequested")
              .withArgs(1, 0, blsData.previousEntry)

            expect(await randomBeacon.isRelayRequestInProgress()).to.be.true
          })
        }
      )

      context(
        "when group is terminated that was supposed to submit a relay request and another group expires",
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

            await (randomBeacon as RandomBeaconStub).roughlyTerminateGroup(0)

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
          // TODO: Implementation once `Groups` library is ready.
        })

        it("should emit RelayEntryTimedOut event", async () => {
          await expect(tx)
            .to.emit(randomBeacon, "RelayEntryTimedOut")
            .withArgs(1, 0)
        })

        it("should clean up current relay request data", async () => {
          // TODO: Uncomment those assertions once termination is implemented.
          // await expect(tx).to.not.emit(randomBeacon, "RelayEntryRequested")
          // expect(await randomBeacon.isRelayRequestInProgress()).to.be.false
        })
      })

      context(
        "when no active groups exist after timeout is reported and DKG is awaiting seed",
        () => {
          let tx: ContractTransaction

          beforeEach(async () => {
            // `groupSize * relayEntrySubmissionEligibilityDelay +
            // relayEntryHardTimeout`.
            await mineBlocks(64 * 10 + 5760)

            // Simulate DKG is awaiting a seed.
            await (randomBeacon as RandomBeaconStub).publicDkgLockState()

            tx = await randomBeacon.connect(notifier).reportRelayEntryTimeout()
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

  async function approveTestToken() {
    await testToken.mint(requester.address, relayRequestFee)
    await testToken
      .connect(requester)
      .approve(randomBeacon.address, relayRequestFee)
  }
})
