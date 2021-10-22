import { ethers } from "hardhat"
import { expect } from "chai"
import type { ContractTransaction } from "ethers"
import { BigNumber } from "ethers"
import blsData from "./data/bls"
import { noMisbehaved, getDkgGroupSigners } from "./utils/dkg"
import { constants } from "./fixtures"
import type { TestGroups } from "../typechain"
import type { DkgGroupSigners } from "./utils/dkg"

describe("Groups", () => {
  const groupPublicKey: string = ethers.utils.hexValue(blsData.groupPubKey)

  let signers: DkgGroupSigners
  let groups: TestGroups
  let members: string[]

  before(async () => {
    signers = await getDkgGroupSigners(constants.groupSize)
    members = Array.from(signers.values())
  })

  beforeEach("load test fixture", async () => {
    const TestGroups = await ethers.getContractFactory("TestGroups")
    groups = await TestGroups.deploy()
  })

  describe("addPendingGroup", async () => {
    context("when no groups are registered", async () => {
      let tx: ContractTransaction

      context("with no misbehaved members", async () => {
        beforeEach(async () => {
          tx = await groups.addPendingGroup(
            groupPublicKey,
            members,
            noMisbehaved
          )
        })

        it("should emit PendingGroupRegistered event", async () => {
          expect(tx)
            .to.emit(groups, "PendingGroupRegistered")
            .withArgs(0, groupPublicKey)
        })

        it("should register a pending group", async () => {
          const storedGroup = await groups.getGroup(groupPublicKey)

          expect(storedGroup.groupPubKey).to.be.equal(groupPublicKey)
          expect(storedGroup.activationTimestamp).to.be.equal(0)
          expect(storedGroup.members).to.be.deep.equal(members)

          const groupsData = await groups.getGroups()

          expect(groupsData).to.be.lengthOf(1)
          expect(groupsData[0]).to.deep.equal(storedGroup)
        })

        it("should store a flagged group index", async () => {
          const expectedFlaggedIndex = calculateFlaggedIndex(0)

          expect(await groups.getFlaggedGroupIndex(groupPublicKey)).to.equal(
            expectedFlaggedIndex
          )
        })
      })

      context("with misbehaved members", async () => {
        context("with first member misbehaved", async () => {
          const misbehavedIndices: number[] = [1]

          beforeEach(async () => {
            const misbehaved = ethers.utils.hexlify(misbehavedIndices)

            tx = await groups.addPendingGroup(
              groupPublicKey,
              members,
              misbehaved
            )
          })

          it("should filter out misbehaved members", async () => {
            const expectedMembers = [...members]
            expectedMembers[0] = expectedMembers.pop()

            expect(
              (await groups.getGroup(groupPublicKey)).members
            ).to.be.deep.equal(expectedMembers)
          })
        })

        context("with last member misbehaved", async () => {
          const misbehavedIndices: number[] = [constants.groupSize]

          beforeEach(async () => {
            const misbehaved = ethers.utils.hexlify(misbehavedIndices)

            tx = await groups.addPendingGroup(
              groupPublicKey,
              members,
              misbehaved
            )
          })

          it("should filter out misbehaved members", async () => {
            const expectedMembers = [...members]
            expectedMembers.pop()

            expect(
              (await groups.getGroup(groupPublicKey)).members
            ).to.be.deep.equal(expectedMembers)
          })
        })

        context("with middle member misbehaved", async () => {
          const misbehavedIndices: number[] = [24]

          beforeEach(async () => {
            const misbehaved = ethers.utils.hexlify(misbehavedIndices)

            tx = await groups.addPendingGroup(
              groupPublicKey,
              members,
              misbehaved
            )
          })

          it("should filter out misbehaved members", async () => {
            const expectedMembers = [...members]
            expectedMembers[24 - 1] = expectedMembers.pop()

            expect(
              (await groups.getGroup(groupPublicKey)).members
            ).to.be.deep.equal(expectedMembers)
          })
        })

        context("with multiple members misbehaved", async () => {
          const misbehavedIndices: number[] = [1, 16, 35, constants.groupSize]

          beforeEach(async () => {
            const misbehaved = ethers.utils.hexlify(misbehavedIndices)

            tx = await groups.addPendingGroup(
              groupPublicKey,
              members,
              misbehaved
            )
          })

          it("should filter out misbehaved members", async () => {
            const expectedMembers = filterMisbehaved(members, misbehavedIndices)

            expect(
              (await groups.getGroup(groupPublicKey)).members
            ).to.be.deep.equal(expectedMembers)
          })
        })

        context("with misbehaved member index 0", async () => {
          const misbehavedIndices: number[] = [0]

          it("should panic", async () => {
            const misbehaved = ethers.utils.hexlify(misbehavedIndices)

            await expect(
              groups.addPendingGroup(groupPublicKey, members, misbehaved)
            ).to.be.revertedWith(
              "reverted with panic code 0x11 (Arithmetic operation underflowed or overflowed outside of an unchecked block)"
            )
          })
        })

        context(
          "with misbehaved member index greater than group size",
          async () => {
            const misbehavedIndices: number[] = [constants.groupSize + 1]

            it("should panic", async () => {
              const misbehaved = ethers.utils.hexlify(misbehavedIndices)

              await expect(
                groups.addPendingGroup(groupPublicKey, members, misbehaved)
              ).to.be.revertedWith(
                "reverted with panic code 0x32 (Array accessed at an out-of-bounds or negative index)"
              )
            })
          }
        )
      })
    })

    context("when existing group is already registered", async () => {
      const existingGroupPublicKey = "0x1234567890"

      let exsitingGroupMembers: string[]
      let newGroupMembers: string[]

      beforeEach(async () => {
        exsitingGroupMembers = members.slice(30)
        newGroupMembers = members.slice(-30)

        await groups.addPendingGroup(
          existingGroupPublicKey,
          exsitingGroupMembers,
          noMisbehaved
        )
      })

      context("when existing group is pending", async () => {
        let existingGroup

        beforeEach(async () => {
          existingGroup = await groups.getGroup(existingGroupPublicKey)
        })

        context("with the same group public key", async () => {
          const newGroupPublicKey = existingGroupPublicKey

          let tx: ContractTransaction

          beforeEach(async () => {
            tx = await groups.addPendingGroup(
              newGroupPublicKey,
              newGroupMembers,
              noMisbehaved
            )
          })

          it("should emit PendingGroupRegistered event", async () => {
            expect(tx)
              .to.emit(groups, "PendingGroupRegistered")
              .withArgs(1, newGroupPublicKey)
          })

          it("should register a pending group", async () => {
            const storedGroup = await groups.getGroup(newGroupPublicKey)

            expect(storedGroup.groupPubKey).to.be.equal(newGroupPublicKey)
            expect(storedGroup.activationTimestamp).to.be.equal(0)
            expect(storedGroup.members).to.be.deep.equal(newGroupMembers)

            const groupsData = await groups.getGroups()

            expect(groupsData).to.be.lengthOf(2)
            expect(groupsData[1]).to.deep.equal(storedGroup)
          })

          it("should not update existing group", async () => {
            const groupsData = await groups.getGroups()

            expect(groupsData[0]).to.deep.equal(existingGroup)
          })

          it("should update stored flagged index for existing group", async () => {
            expect(
              await groups.getFlaggedGroupIndex(existingGroupPublicKey)
            ).to.equal(calculateFlaggedIndex(1))
          })

          it("should store a flagged group index", async () => {
            expect(
              await groups.getFlaggedGroupIndex(newGroupPublicKey)
            ).to.equal(calculateFlaggedIndex(1))
          })
        })

        context("with unique group public key", async () => {
          const newGroupPublicKey = groupPublicKey

          let tx: ContractTransaction

          beforeEach(async () => {
            tx = await groups.addPendingGroup(
              newGroupPublicKey,
              newGroupMembers,
              noMisbehaved
            )
          })

          it("should emit PendingGroupRegistered event", async () => {
            expect(tx)
              .to.emit(groups, "PendingGroupRegistered")
              .withArgs(1, newGroupPublicKey)
          })

          it("should register a pending group", async () => {
            const storedGroup = await groups.getGroup(newGroupPublicKey)

            expect(storedGroup.groupPubKey).to.be.equal(newGroupPublicKey)
            expect(storedGroup.activationTimestamp).to.be.equal(0)
            expect(storedGroup.members).to.be.deep.equal(newGroupMembers)

            const groupsData = await groups.getGroups()

            expect(groupsData).to.be.lengthOf(2)
            expect(groupsData[1]).to.deep.equal(storedGroup)
          })

          it("should not update existing group", async () => {
            const groupsData = await groups.getGroups()

            expect(groupsData[0]).to.deep.equal(existingGroup)

            expect(
              await groups.getFlaggedGroupIndex(existingGroupPublicKey)
            ).to.equal(calculateFlaggedIndex(0))
          })

          it("should store a flagged group index", async () => {
            expect(
              await groups.getFlaggedGroupIndex(newGroupPublicKey)
            ).to.equal(calculateFlaggedIndex(1))
          })
        })
      })

      context("when existing group is active", async () => {
        let existingGroup

        beforeEach(async () => {
          await groups.activateGroup(existingGroupPublicKey)

          existingGroup = await groups.getGroup(existingGroupPublicKey)
        })

        context("with the same group public key", async () => {
          const newGroupPublicKey = existingGroupPublicKey

          it("should revert with 'group was already activated' error", async () => {
            expect(
              groups.addPendingGroup(
                newGroupPublicKey,
                newGroupMembers,
                noMisbehaved
              )
            ).to.be.revertedWith("group was already activated")
          })
        })

        context("with unique group public key", async () => {
          const newGroupPublicKey = groupPublicKey

          let tx: ContractTransaction

          beforeEach(async () => {
            tx = await groups.addPendingGroup(
              newGroupPublicKey,
              newGroupMembers,
              noMisbehaved
            )
          })

          it("should emit PendingGroupRegistered event", async () => {
            expect(tx)
              .to.emit(groups, "PendingGroupRegistered")
              .withArgs(1, newGroupPublicKey)
          })

          it("should register a pending group", async () => {
            const storedGroup = await groups.getGroup(newGroupPublicKey)

            expect(storedGroup.groupPubKey).to.be.equal(newGroupPublicKey)
            expect(storedGroup.activationTimestamp).to.be.equal(0)
            expect(storedGroup.members).to.be.deep.equal(newGroupMembers)

            const groupsData = await groups.getGroups()

            expect(groupsData).to.be.lengthOf(2)
            expect(groupsData[1]).to.deep.equal(storedGroup)
          })

          it("should not update existing group", async () => {
            const groupsData = await groups.getGroups()

            expect(groupsData[0]).to.deep.equal(existingGroup)

            expect(
              await groups.getFlaggedGroupIndex(existingGroupPublicKey)
            ).to.equal(calculateFlaggedIndex(0))
          })

          it("should store a flagged group index", async () => {
            expect(
              await groups.getFlaggedGroupIndex(newGroupPublicKey)
            ).to.equal(calculateFlaggedIndex(1))
          })
        })
      })
    })
  })
})

function calculateFlaggedIndex(index: number): BigNumber {
  // eslint-disable-next-line no-bitwise
  return BigNumber.from(index).xor(BigNumber.from(1).shl(255))
}

function filterMisbehaved(
  members: string[],
  misbehavedIndices: number[]
): string[] {
  const expectedMembers = [...members]
  misbehavedIndices.reverse().forEach((value) => {
    expectedMembers[value - 1] = expectedMembers[expectedMembers.length - 1]
    expectedMembers.pop()
  })

  return expectedMembers
}
