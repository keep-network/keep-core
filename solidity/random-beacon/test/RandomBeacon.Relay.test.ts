import { ethers, waffle, helpers } from "hardhat"
import { expect } from "chai"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import { ContractReceipt, ContractTransaction } from "ethers"
import blsData from "./data/bls"
import group from "./data/group"
import { to1e18 } from "./functions"
import { randomBeaconDeployment } from "./fixtures"
import type {
  RandomBeacon,
  SortitionPoolStub,
  TestToken,
  StakingStub,
  RelayStub,
} from "../typechain"

const { time } = helpers
const { mineBlocks, mineBlocksTo } = time

describe("RandomBeacon - Relay", () => {
  const relayRequestFee = to1e18(100)

  let requester: SignerWithAddress
  let member3: SignerWithAddress
  let member16: SignerWithAddress
  let member17: SignerWithAddress
  let member18: SignerWithAddress

  let randomBeacon: RandomBeacon
  let sortitionPool: SortitionPoolStub
  let testToken: TestToken
  let staking: StakingStub
  let relayStub: RelayStub

  const fixture = async () => {
    const deployment = await randomBeaconDeployment()

    return {
      randomBeacon: deployment.randomBeacon,
      sortitionPoolStub: deployment.sortitionPoolStub,
      testToken: deployment.testToken,
      stakingStub: deployment.stakingStub,
      relayStub: await (await ethers.getContractFactory("RelayStub")).deploy(),
    }
  }

  // prettier-ignore
  before(async () => {
    [requester, member3, member16, member17, member18] = await ethers.getSigners()
  })

  beforeEach("load test fixture", async () => {
    const contracts = await waffle.loadFixture(fixture)

    randomBeacon = contracts.randomBeacon as RandomBeacon
    sortitionPool = contracts.sortitionPoolStub as SortitionPoolStub
    testToken = contracts.testToken as TestToken
    staking = contracts.stakingStub as StakingStub
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
              context(
                "when first eligible member submits before the soft timeout",
                () => {
                  let tx: ContractTransaction

                  beforeEach(async () => {
                    // When determining the eligibility queue, the
                    // `(groupSignature % 64) + 1` equation points member `16`
                    // as the first eligible one. This is why we use that
                    // index as `submitRelayEntry` parameter. The `submitter`
                    // signer represents that member too.
                    tx = await randomBeacon
                      .connect(member16)
                      .submitRelayEntry(16, blsData.groupSignature)
                  })

                  it("should not remove any members from the sortition pool", async () => {
                    await expect(tx).to.not.emit(
                      sortitionPool,
                      "OperatorsRemoved"
                    )
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
                    // eslint-disable-next-line @typescript-eslint/no-unused-expressions
                    expect(await randomBeacon.isRelayRequestInProgress()).to.be
                      .false
                  })
                }
              )

              context(
                "when other than first eligible member submits before the soft timeout",
                () => {
                  let tx: ContractTransaction

                  beforeEach(async () => {
                    // When determining the eligibility queue, the
                    // `(groupSignature % 64) + 1` equation points member `16`
                    // as the first eligible one. However, we wait 20 blocks
                    // to make two more members eligible. The member `18`
                    // submits the result.
                    await mineBlocks(20)

                    tx = await randomBeacon
                      .connect(member18)
                      .submitRelayEntry(18, blsData.groupSignature)
                  })

                  it("should remove members who did not submit from the sortition pool", async () => {
                    await expect(tx)
                      .to.emit(sortitionPool, "OperatorsRemoved")
                      .withArgs([member16.address, member17.address])
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
                    // eslint-disable-next-line @typescript-eslint/no-unused-expressions
                    expect(await randomBeacon.isRelayRequestInProgress()).to.be
                      .false
                  })
                }
              )

              context(
                "when first eligible member submits after the soft timeout",
                () => {
                  let tx: ContractTransaction
                  let receipt: ContractReceipt

                  beforeEach(async () => {
                    // Current request start block is 10. The `softTimeoutBlock`
                    // is `currentRequestStartBlock + groupSize *
                    // relayEntrySubmissionEligibilityDelay = 10 + 64 * 10 = 650`
                    // The hard timeout is `5760` blocks. Let's assume we want
                    // to submit the relay entry after 75% of the soft timeout
                    // period elapses. If so we need to do it at block
                    // `650 + (0.75 * 5760) = 4970`. However, we need to mine
                    // to one block before because the relay entry submission
                    // transaction will move the blockchain ahead by one block
                    // due to the Hardhat auto-mine feature.
                    await mineBlocksTo(4969)

                    // When determining the eligibility queue, the
                    // `(groupSignature % 64) + 1` equation points member `16`
                    // as the first eligible one. This is why we use that
                    // index as `submitRelayEntry` parameter. The `submitter`
                    // signer represents that member too.
                    tx = await randomBeacon
                      .connect(member16)
                      .submitRelayEntry(16, blsData.groupSignature)

                    receipt = await tx.wait()
                  })

                  it("should not remove any members from the sortition pool", async () => {
                    await expect(tx).to.not.emit(
                      sortitionPool,
                      "OperatorsRemoved"
                    )
                  })

                  it("should slash 75% of slashing amount for all members ", async () => {
                    // `relayEntrySubmissionFailureSlashingAmount = 1000e18`.
                    // 75% of the soft timeout period elapsed so we expect
                    // `750e18` to be slashed.
                    await expect(tx)
                      .to.emit(staking, "Slashed")
                      .withArgs(to1e18(750), group.members)
                  })

                  it("should emit RelayEntrySubmitted event", async () => {
                    await expect(tx)
                      .to.emit(randomBeacon, "RelayEntrySubmitted")
                      .withArgs(1, blsData.groupSignature)
                  })

                  it("should terminate the relay request", async () => {
                    // eslint-disable-next-line @typescript-eslint/no-unused-expressions
                    expect(await randomBeacon.isRelayRequestInProgress()).to.be
                      .false
                  })
                }
              )
            })

            context("when entry is not valid", () => {
              it("should revert", async () => {
                // In that case `(nextGroupSignature % 64) + 1` gives 3 so that
                // member needs to submit the wrong relay entry.
                await expect(
                  randomBeacon
                    .connect(member3)
                    .submitRelayEntry(3, blsData.nextGroupSignature)
                ).to.be.revertedWith("Invalid entry")
              })
            })
          })

          context("when submitter is not eligible", () => {
            it("should revert", async () => {
              await expect(
                randomBeacon
                  .connect(member17)
                  .submitRelayEntry(17, blsData.groupSignature)
              ).to.be.revertedWith("Submitter is not eligible")
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
              .submitRelayEntry(16, blsData.nextGroupSignature)
          ).to.be.revertedWith("Relay request timed out")
        })
      })
    })

    context("when relay request is not in progress", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(member16)
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
