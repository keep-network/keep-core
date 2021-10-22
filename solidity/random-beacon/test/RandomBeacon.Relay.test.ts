import { ethers, waffle, helpers } from "hardhat"
import { expect } from "chai"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import blsData from "./data/bls"
import { to1e18 } from "./functions"
import { randomBeaconDeployment } from "./fixtures"
import type { RandomBeacon, TestToken, RelayStub } from "../typechain"

const { time } = helpers
const { mineBlocks } = time

describe("RandomBeacon - Relay", () => {
  const relayRequestFee = to1e18(100)

  let requester: SignerWithAddress
  let submitter: SignerWithAddress
  let other: SignerWithAddress
  let invalidEntrySubmitter: SignerWithAddress

  let randomBeacon: RandomBeacon
  let testToken: TestToken
  let relayStub: RelayStub

  // prettier-ignore
  before(async () => {
    [requester, submitter, other, invalidEntrySubmitter] = await ethers.getSigners()
  })

  beforeEach("load test fixture", async () => {
    const contracts = await waffle.loadFixture(randomBeaconDeployment)

    randomBeacon = contracts.randomBeacon as RandomBeacon
    testToken = contracts.testToken as TestToken
    relayStub = contracts.relayStub as RelayStub

    await randomBeacon.updateRelayEntryParameters(to1e18(100), 10, 5760, 0)
  })

  describe("requestRelayEntry", () => {
    context("when groups exist", () => {
      beforeEach(async () => {
        // TODO: Currently `selectGroup` returns a hardcoded group. Once
        //       proper implementation is ready, add the group manually here.
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
            tx = await randomBeacon.connect(requester).requestRelayEntry()
          })

          it("should deposit relay request fee to the maintenance pool", async () => {
            const actualMaintenancePoolBalance = await testToken.balanceOf(
              randomBeacon.address
            )
            expect(
              actualMaintenancePoolBalance.sub(previousMaintenancePoolBalance)
            ).to.be.equal(relayRequestFee)
          })

          it("should emit RelayEntryRequested event", async () => {
            await expect(tx)
              .to.emit(randomBeacon, "RelayEntryRequested")
              .withArgs(1, 1, blsData.previousEntry)
          })
        })

        context("when the requester doesn't pay the relay request fee", () => {
          it("should revert", async () => {
            await expect(
              randomBeacon.connect(requester).requestRelayEntry()
            ).to.be.revertedWith("Transfer amount exceeds allowance")
          })
        })
      })

      context("when there is an other relay entry in progress", () => {
        beforeEach(async () => {
          await approveTestToken()
          await randomBeacon.connect(requester).requestRelayEntry()
        })

        it("should revert", async () => {
          await expect(
            randomBeacon.connect(requester).requestRelayEntry()
          ).to.be.revertedWith("Another relay request in progress")
        })
      })
    })

    context("when no groups exist", () => {
      it("should revert", async () => {
        // TODO: Implement once proper `selectGroup` is ready.
      })
    })
  })

  describe("submitRelayEntry", () => {
    context("when relay request is in progress", () => {
      beforeEach(async () => {
        await approveTestToken()
        await randomBeacon.connect(requester).requestRelayEntry()
      })

      context("when relay entry is not timed out", () => {
        context("when submitter index is valid", () => {
          context("when submitter is eligible", () => {
            context("when entry is valid", () => {
              it("should emit RelayEntrySubmitted event", async () => {
                await expect(
                  randomBeacon
                    .connect(submitter)
                    .submitRelayEntry(16, blsData.groupSignature)
                )
                  .to.emit(randomBeacon, "RelayEntrySubmitted")
                  .withArgs(1, blsData.groupSignature)
              })
            })

            context("when entry is not valid", () => {
              it("should revert", async () => {
                // In that case nextGroupSignature % 64 gives 3 so that
                // member needs to submit the wrong relay entry.
                await expect(
                  randomBeacon
                    .connect(invalidEntrySubmitter)
                    .submitRelayEntry(3, blsData.nextGroupSignature)
                ).to.be.revertedWith("Invalid entry")
              })
            })
          })

          context("when submitter is not eligible", () => {
            it("should revert", async () => {
              await expect(
                randomBeacon
                  .connect(other)
                  .submitRelayEntry(17, blsData.groupSignature)
              ).to.be.revertedWith("Submitter is not eligible")
            })
          })
        })

        context("when submitter index is beyond valid range", () => {
          it("should revert", async () => {
            await expect(
              randomBeacon
                .connect(submitter)
                .submitRelayEntry(0, blsData.nextGroupSignature)
            ).to.be.revertedWith("Invalid submitter index")

            await expect(
              randomBeacon
                .connect(submitter)
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
                  .connect(submitter)
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
              .connect(submitter)
              .submitRelayEntry(16, blsData.nextGroupSignature)
          ).to.be.revertedWith("Relay request timed out")
        })
      })
    })

    context("when relay request is not in progress", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(submitter)
            .submitRelayEntry(16, blsData.nextGroupSignature)
        ).to.be.revertedWith("No relay request in progress")
      })
    })
  })

  describe("isEligible", () => {
    it("should correctly manage the eligibility queue", async () => {
      await relayStub.setCurrentRequestStartBlock()

      // At the beginning only member 8 is eligible because
      // (blsData.groupSignature % groupSize) + 1 = 8.
      await assertMembersEligible([8])
      await assertMembersNotEligible([1, 2, 3, 4, 5, 6, 7])

      await mineBlocks(10)

      await assertMembersEligible([8, 1])
      await assertMembersNotEligible([2, 3, 4, 5, 6, 7])

      await mineBlocks(10)

      await assertMembersEligible([8, 1, 2])
      await assertMembersNotEligible([3, 4, 5, 6, 7])

      await mineBlocks(10)

      await assertMembersEligible([8, 1, 2, 3])
      await assertMembersNotEligible([4, 5, 6, 7])

      await mineBlocks(10)

      await assertMembersEligible([8, 1, 2, 3, 4])
      await assertMembersNotEligible([5, 6, 7])

      await mineBlocks(10)

      await assertMembersEligible([8, 1, 2, 3, 4, 5])
      await assertMembersNotEligible([6, 7])

      await mineBlocks(10)

      await assertMembersEligible([8, 1, 2, 3, 4, 5, 6])
      await assertMembersNotEligible([7])

      await mineBlocks(10)

      await assertMembersEligible([8, 1, 2, 3, 4, 5, 6, 7])
    })
  })

  async function approveTestToken() {
    await testToken.mint(requester.address, relayRequestFee)
    await testToken
      .connect(requester)
      .approve(randomBeacon.address, relayRequestFee)
  }

  async function assertMembersEligible(members: number[]) {
    for (let i = 0; i < members.length; i++) {
      // eslint-disable-next-line no-await-in-loop,@typescript-eslint/no-unused-expressions
      expect(await relayStub.isEligible(members[i], blsData.groupSignature)).to
        .be.true
    }
  }

  async function assertMembersNotEligible(members: number[]) {
    for (let i = 0; i < members.length; i++) {
      // eslint-disable-next-line no-await-in-loop,@typescript-eslint/no-unused-expressions
      expect(await relayStub.isEligible(members[i], blsData.groupSignature)).to
        .be.false
    }
  }
})
