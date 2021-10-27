import { ethers, waffle, helpers, getUnnamedAccounts } from "hardhat"
import { expect } from "chai"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import blsData from "./data/bls"
import { getDkgGroupSigners, signAndSubmitDkgResult } from "./utils/dkg"
import { to1e18, ZERO_ADDRESS } from "./functions"
import { constants, params, randomBeaconDeployment } from "./fixtures"
import type { RandomBeacon, TestToken, RelayStub } from "../typechain"
import type { DkgGroupSigners } from "./utils/dkg"

const { time } = helpers
const { mineBlocks } = time

describe("RandomBeacon - Relay", () => {
  const relayRequestFee = to1e18(100)

  // When determining the eligibility queue, the `(blsData.groupSignature % 64) + 1`
  // equation points member`16` as the first eligible one. This is why we use that
  // index as `submitRelayEntry` parameter. The `submitter` signer represents that
  // member too.
  const submitterMemberIndex = 16
  // In that case `(blsData.nextGroupSignature % 64) + 1` gives 3 so that  member needs
  // to submit the wrong relay entry.
  const invalidSubmitterMemberIndex = 3

  let requester: SignerWithAddress
  let submitter: SignerWithAddress
  let other: SignerWithAddress
  let invalidEntrySubmitter: SignerWithAddress
  let signers: DkgGroupSigners

  let randomBeacon: RandomBeacon
  let testToken: TestToken
  let relayStub: RelayStub

  const fixture = async () => {
    const deployment = await randomBeaconDeployment()

    return {
      randomBeacon: deployment.randomBeacon,
      testToken: deployment.testToken,
      relayStub: await (await ethers.getContractFactory("RelayStub")).deploy(),
    }
  }

  before(async () => {
    requester = await ethers.getSigner((await getUnnamedAccounts())[1])

    signers = await getDkgGroupSigners(constants.groupSize, 1)

    submitter = await ethers.getSigner(signers.get(submitterMemberIndex))
    invalidEntrySubmitter = await ethers.getSigner(
      signers.get(invalidSubmitterMemberIndex)
    )
    other = await ethers.getSigner(signers.get(submitterMemberIndex + 1))
  })

  beforeEach("load test fixture", async () => {
    const contracts = await waffle.loadFixture(fixture)

    randomBeacon = contracts.randomBeacon as RandomBeacon
    testToken = contracts.testToken as TestToken
    relayStub = contracts.relayStub as RelayStub

    await randomBeacon.updateRelayEntryParameters(to1e18(100), 10, 5760, 0)
  })

  describe("requestRelayEntry", () => {
    context("when groups exist", () => {
      beforeEach(async () => {
        await createGroup()
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
            tx = await randomBeacon
              .connect(requester)
              .requestRelayEntry(ZERO_ADDRESS)
          })

          it("should deposit relay request fee to the maintenance pool", async () => {
            const currentMaintenancePoolBalance = await testToken.balanceOf(
              randomBeacon.address
            )
            expect(
              currentMaintenancePoolBalance.sub(previousMaintenancePoolBalance)
            ).to.be.equal(relayRequestFee)
          })

          it("should emit RelayEntryRequested event", async () => {
            await expect(tx)
              .to.emit(randomBeacon, "RelayEntryRequested")
              .withArgs(1, 0, blsData.previousEntry)
          })
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
        // TODO: Implement once proper `selectGroup` is ready.
        await expect(
          randomBeacon.connect(requester).requestRelayEntry(ZERO_ADDRESS)
        ).to.be.revertedWith(
          "reverted with panic code 0x11 (Arithmetic operation underflowed or overflowed outside of an unchecked block)"
        )
      })
    })
  })

  describe("submitRelayEntry", () => {
    beforeEach(async () => {
      await createGroup()
    })

    context("when relay request is in progress", () => {
      beforeEach(async () => {
        await approveTestToken()
        await randomBeacon.connect(requester).requestRelayEntry(ZERO_ADDRESS)
      })

      context("when relay entry is not timed out", () => {
        context("when submitter index is valid", () => {
          context("when submitter is eligible", () => {
            context("when entry is valid", () => {
              it("should emit RelayEntrySubmitted event", async () => {
                await expect(
                  randomBeacon
                    .connect(submitter)
                    .submitRelayEntry(
                      submitterMemberIndex,
                      blsData.groupSignature
                    )
                )
                  .to.emit(randomBeacon, "RelayEntrySubmitted")
                  .withArgs(1, blsData.groupSignature)
              })
            })

            context("when entry is not valid", () => {
              it("should revert", async () => {
                await expect(
                  randomBeacon
                    .connect(invalidEntrySubmitter)
                    .submitRelayEntry(
                      invalidSubmitterMemberIndex,
                      blsData.nextGroupSignature
                    )
                ).to.be.revertedWith("Invalid entry")
              })
            })
          })

          context("when submitter is not eligible", () => {
            it("should revert", async () => {
              await expect(
                randomBeacon
                  .connect(other)
                  .submitRelayEntry(
                    submitterMemberIndex + 1,
                    blsData.groupSignature
                  )
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
              .submitRelayEntry(
                submitterMemberIndex,
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
            .connect(submitter)
            .submitRelayEntry(submitterMemberIndex, blsData.nextGroupSignature)
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

  async function createGroup() {
    const { blockNumber: startBlock } = await randomBeacon.genesis()
    await mineBlocks(constants.offchainDkgTime)
    await signAndSubmitDkgResult(
      randomBeacon,
      blsData.groupPubKey,
      signers,
      startBlock
    )
    await mineBlocks(params.dkgResultChallengePeriodLength)
    await randomBeacon.approveDkgResult()
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
