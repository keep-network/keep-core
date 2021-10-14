import { ethers, waffle, helpers } from "hardhat"
import { expect } from "chai"
import { blsData } from "./helpers/data"
import { to1e18 } from "./helpers/functions"
import { constants, testDeployment } from "./helpers/fixtures"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { RandomBeacon, TestToken, MaintenancePool } from "../typechain"
import {exec} from "child_process";

const { time } = helpers
const { mineBlocks } = time

interface GroupMember {
  index: number
  signer: SignerWithAddress
}

describe.only("RandomBeacon - Relay", function () {
  const relayRequestFee = to1e18(100)

  let requester: SignerWithAddress
  let member1: GroupMember
  let member2: GroupMember
  let member3: GroupMember
  let member4: GroupMember
  let member5: GroupMember
  let member6: GroupMember
  let member7: GroupMember
  let member8: GroupMember

  let randomBeacon: RandomBeacon
  let testToken: TestToken
  let maintenancePool: MaintenancePool

  before(async () => {
    let signer1: SignerWithAddress
    let signer2: SignerWithAddress
    let signer3: SignerWithAddress
    let signer4: SignerWithAddress
    let signer5: SignerWithAddress
    let signer6: SignerWithAddress
    let signer7: SignerWithAddress
    let signer8: SignerWithAddress

    [
      requester,
      signer1,
      signer2,
      signer3,
      signer4,
      signer5,
      signer6,
      signer7,
      signer8,
    ] = await ethers.getSigners()

    member1 = {index: 1, signer: signer1}
    member2 = {index: 2, signer: signer2}
    member3 = {index: 3, signer: signer3}
    member4 = {index: 4, signer: signer4}
    member5 = {index: 5, signer: signer5}
    member6 = {index: 6, signer: signer6}
    member7 = {index: 7, signer: signer7}
    member8 = {index: 8, signer: signer8}

    // Use smaller group size to make testing easier.
    constants.groupSize = 8
    constants.signatureThreshold = 5
  })

  beforeEach("load test fixture", async () => {
    const contracts = await waffle.loadFixture(testDeployment)

    randomBeacon = contracts.randomBeacon as RandomBeacon
    testToken = contracts.testToken as TestToken
    maintenancePool = contracts.maintenancePool as MaintenancePool
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
            previousMaintenancePoolBalance = await testToken.balanceOf(maintenancePool.address)
            await approveTestToken()
            tx = await randomBeacon.connect(requester).requestRelayEntry(blsData.previousEntry)
          })

          it("should deposit relay request fee to the maintenance pool", async () => {
            const actualMaintenancePoolBalance = await testToken.balanceOf(maintenancePool.address)
            expect(actualMaintenancePoolBalance.sub(previousMaintenancePoolBalance)).to.be.equal(relayRequestFee)
          })

          it("should emit RelayEntryRequested event", async () => {
            await expect(tx).to
              .emit(randomBeacon, "RelayEntryRequested")
              .withArgs(1, blsData.groupPubKey, blsData.previousEntry)
          })
        })

        context("when the requester doesn't pay the relay request fee", () => {
          it("should revert", async () => {
            await expect(
              randomBeacon.connect(requester).requestRelayEntry(blsData.previousEntry)
            ).to.be.revertedWith("Transfer amount exceeds allowance")
          })
        })
      })

      context("when there is an other relay entry in progress", () => {
        beforeEach(async () => {
          await approveTestToken()
          await randomBeacon.connect(requester).requestRelayEntry(blsData.previousEntry)
        })

        it("should revert", async () => {
          await expect(
            randomBeacon.connect(requester).requestRelayEntry(blsData.previousEntry)
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
        await randomBeacon.connect(requester).requestRelayEntry(blsData.previousEntry)
      })

      context("when relay entry is not timed out", () => {
        context("when submitter index is valid", () => {
          context("when entry is valid", () => {
            it("should correctly manage the eligibility queue", async () => {
              // At the beginning only member 8 is eligible because
              // (blsData.groupSignature % groupSize) + 1 = 8.
              await assertMembersEligible([member8])
              await assertMembersNotEligible([member1, member2, member3, member4, member5, member6, member7])

              await mineBlocks(10)

              await assertMembersEligible([member8, member1])
              await assertMembersNotEligible([member2, member3, member4, member5, member6, member7])

              await mineBlocks(10)

              await assertMembersEligible([member8, member1, member2])
              await assertMembersNotEligible([member3, member4, member5, member6, member7])

              await mineBlocks(10)

              await assertMembersEligible([member8, member1, member2, member3])
              await assertMembersNotEligible([member4, member5, member6, member7])

              await mineBlocks(10)

              await assertMembersEligible([member8, member1, member2, member3, member4])
              await assertMembersNotEligible([member5, member6, member7])

              await mineBlocks(10)

              await assertMembersEligible([member8, member1, member2, member3, member4, member5])
              await assertMembersNotEligible([member6, member7])

              await mineBlocks(10)

              await assertMembersEligible([member8, member1, member2, member3, member4, member5, member6])
              await assertMembersNotEligible([member7])

              await mineBlocks(10)

              await assertMembersEligible([member8, member1, member2, member3, member4, member5, member6, member7])
            })

            it("should emit RelayEntrySubmitted event", async () => {
              await expect(
                randomBeacon.connect(member8.signer)
                  .submitRelayEntry(member8.index, blsData.groupSignature)
              ).to
                .emit(randomBeacon, "RelayEntrySubmitted")
                .withArgs(1, blsData.groupSignature)
            })
          })

          context("when entry is not valid", () => {
            it("should revert", async () => {
              await expect(
                randomBeacon.connect(member8.signer)
                  .submitRelayEntry(member8.index, blsData.nextGroupSignature)
              ).to.be.revertedWith("Invalid entry")
            })
          })
        })

        context("when submitter index is beyond valid range", () => {
          it("should revert", async () => {
            await expect(
              randomBeacon.connect(member8.signer)
                .submitRelayEntry(0, blsData.nextGroupSignature)
            ).to.be.revertedWith("Invalid submitter index")

            await expect(
              randomBeacon.connect(member8.signer)
                .submitRelayEntry(9, blsData.nextGroupSignature)
            ).to.be.revertedWith("Invalid submitter index")
          })
        })

        context("when submitter index does not correspond to sender address", () => {
          it("should revert", async () => {
            await expect(
              randomBeacon.connect(member8.signer)
                .submitRelayEntry(7, blsData.nextGroupSignature)
            ).to.be.revertedWith("Unexpected submitter index")
          })
        })
      })

      context("when relay entry is timed out", () => {
        it("should revert", async () => {
          // groupSize * relayEntrySubmissionEligibilityDelay + relayEntryHardTimeout
          await mineBlocks(8 * 10 + 5760)

          await expect(
            randomBeacon.connect(member8.signer)
              .submitRelayEntry(member8.index, blsData.nextGroupSignature)
          ).to.be.revertedWith("Relay request timed out")
        })
      })
    })

    context("when relay request is not in progress", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon.connect(member8.signer)
            .submitRelayEntry(member8.index, blsData.nextGroupSignature)
        ).to.be.revertedWith("No relay request in progress")
      })
    })
  })

  async function approveTestToken() {
    await testToken.mint(requester.address, relayRequestFee)
    await testToken.connect(requester).approve(randomBeacon.address, relayRequestFee)
  }

  async function assertMembersEligible(members: GroupMember[]) {
    for(let i = 0; i < members.length; i++) {
      const member = members[i]
      await expect(
        randomBeacon.connect(member.signer).callStatic
          .submitRelayEntry(member.index, blsData.groupSignature)
      ).not.to.be.reverted
    }
  }

  async function assertMembersNotEligible(members: GroupMember[]) {
    for(let i = 0; i < members.length; i++) {
      const member = members[i]
      await expect(
        randomBeacon.connect(member.signer).callStatic
          .submitRelayEntry(member.index, blsData.groupSignature)
      ).to.be.revertedWith("Submitter is not eligible")
    }
  }
})