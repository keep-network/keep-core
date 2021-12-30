/* eslint-disable no-await-in-loop */

import { ethers, waffle } from "hardhat"
import { expect } from "chai"
import type { ContractTransaction } from "ethers"
import blsData from "./data/bls"
import { constants } from "./fixtures"
import type { GroupsStub } from "../typechain"
import { noMisbehaved } from "./utils/dkg"
import { hashUint32Array } from "./utils/groups"

const { keccak256 } = ethers.utils

const fixture = async () => {
  const GroupsStub = await ethers.getContractFactory("GroupsStub")
  const groups = await GroupsStub.deploy()

  return groups
}

describe("Groups", () => {
  const groupPublicKey: string = ethers.utils.hexValue(blsData.groupPubKey)
  const members: number[] = []

  let groups: GroupsStub

  before(async () => {
    for (let i = 0; i < constants.groupSize; i++) members.push(10000 + i)
  })

  beforeEach("load test fixture", async () => {
    groups = await waffle.loadFixture(fixture)
  })

  describe("addCandidateGroup", async () => {
    context("when no groups are registered", async () => {
      let tx: ContractTransaction

      context("with no misbehaved members", async () => {
        beforeEach(async () => {
          tx = await groups.addCandidateGroup(
            groupPublicKey,
            members,
            noMisbehaved
          )
        })

        it("should emit CandidateGroupRegistered event", async () => {
          await expect(tx)
            .to.emit(groups, "CandidateGroupRegistered")
            .withArgs(groupPublicKey)
        })

        it("should register group", async () => {
          const groupsRegistry = await groups.getGroupsRegistry()

          expect(groupsRegistry).to.be.lengthOf(1)
          expect(groupsRegistry[0]).to.deep.equal(keccak256(groupPublicKey))
        })

        it("should store group data", async () => {
          const storedGroup = await groups.getGroup(groupPublicKey)

          expect(storedGroup.groupPubKey).to.be.equal(groupPublicKey)
          expect(storedGroup.activationBlockNumber).to.be.equal(0)
          expect(storedGroup.membersHash).to.be.equal(hashUint32Array(members))
        })
      })

      context("with misbehaved members", async () => {
        context("with first member misbehaved", async () => {
          it("should filter out misbehaved members", async () => {
            const misbehavedIndices: number[] = [1]

            await groups.addCandidateGroup(
              groupPublicKey,
              members,
              misbehavedIndices
            )

            const expectedMembers = [...members]
            expectedMembers.splice(0, 1)

            expect(
              (await groups.getGroup(groupPublicKey)).membersHash
            ).to.be.equal(hashUint32Array(expectedMembers))
          })
        })

        context("with last member misbehaved", async () => {
          it("should filter out misbehaved members", async () => {
            const misbehavedIndices: number[] = [constants.groupSize]
            await groups.addCandidateGroup(
              groupPublicKey,
              members,
              misbehavedIndices
            )
            const expectedMembers = [...members]
            expectedMembers.pop()

            expect(
              (await groups.getGroup(groupPublicKey)).membersHash
            ).to.be.equal(hashUint32Array(expectedMembers))
          })
        })

        context("with middle member misbehaved", async () => {
          it("should filter out misbehaved members", async () => {
            const misbehavedIndices: number[] = [24]
            await groups.addCandidateGroup(
              groupPublicKey,
              members,
              misbehavedIndices
            )

            const expectedMembers = [...members]
            expectedMembers.splice(23, 1)

            expect(
              (await groups.getGroup(groupPublicKey)).membersHash
            ).to.be.equal(hashUint32Array(expectedMembers))
          })
        })

        context("with multiple members misbehaved", async () => {
          it("should filter out misbehaved members", async () => {
            const misbehavedIndices: number[] = [1, 16, 35, constants.groupSize]
            await groups.addCandidateGroup(
              groupPublicKey,
              members,
              misbehavedIndices
            )
            const expectedMembers = [...members]
            expectedMembers.splice(0, 1) // index -1
            expectedMembers.splice(14, 1) // index -2 (cause expectedMembers already shrinked)
            expectedMembers.splice(32, 1) // index -3
            expectedMembers.splice(constants.groupSize - 4, 1) // index -4

            expect(
              (await groups.getGroup(groupPublicKey)).membersHash
            ).to.be.equal(hashUint32Array(expectedMembers))
          })
        })

        context("with misbehaved member index 0", async () => {
          const misbehavedIndices: number[] = [0]

          it("should panic", async () => {
            await expect(
              groups.addCandidateGroup(
                groupPublicKey,
                members,
                misbehavedIndices
              )
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
              await expect(
                groups.addCandidateGroup(
                  groupPublicKey,
                  members,
                  misbehavedIndices
                )
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

      let existingGroupMembers: number[]
      let newGroupMembers: number[]

      beforeEach(async () => {
        existingGroupMembers = members.slice(30)
        newGroupMembers = members.slice(-30)

        await groups.addCandidateGroup(
          existingGroupPublicKey,
          existingGroupMembers,
          noMisbehaved
        )
      })

      context("when existing group is candidate", async () => {
        beforeEach(async () => {
          await groups.getGroup(existingGroupPublicKey)
        })

        context("with the same group public key", async () => {
          const newGroupPublicKey = existingGroupPublicKey

          let tx: ContractTransaction

          beforeEach(async () => {
            tx = await groups.addCandidateGroup(
              newGroupPublicKey,
              newGroupMembers,
              noMisbehaved
            )
          })

          it("should emit CandidateGroupRegistered event", async () => {
            await expect(tx)
              .to.emit(groups, "CandidateGroupRegistered")
              .withArgs(newGroupPublicKey)
          })

          it("should register group", async () => {
            const groupsRegistry = await groups.getGroupsRegistry()

            expect(groupsRegistry).to.be.lengthOf(2)
            expect(groupsRegistry[1]).to.deep.equal(
              keccak256(newGroupPublicKey)
            )
          })

          it("should store group data", async () => {
            const storedGroup = await groups.getGroup(newGroupPublicKey)

            expect(storedGroup.groupPubKey).to.be.equal(newGroupPublicKey)
            expect(storedGroup.activationBlockNumber).to.be.equal(0)
            expect(storedGroup.membersHash).to.be.equal(
              hashUint32Array(newGroupMembers)
            )
          })

          it("should not update existing group entry", async () => {
            const groupsRegistry = await groups.getGroupsRegistry()

            expect(groupsRegistry[0]).to.deep.equal(
              keccak256(existingGroupPublicKey)
            )
          })
        })

        context("with unique group public key", async () => {
          const newGroupPublicKey = groupPublicKey

          let tx: ContractTransaction

          beforeEach(async () => {
            tx = await groups.addCandidateGroup(
              newGroupPublicKey,
              newGroupMembers,
              noMisbehaved
            )
          })

          it("should emit CandidateGroupRegistered event", async () => {
            await expect(tx)
              .to.emit(groups, "CandidateGroupRegistered")
              .withArgs(newGroupPublicKey)
          })

          it("should register group", async () => {
            const groupsRegistry = await groups.getGroupsRegistry()

            expect(groupsRegistry).to.be.lengthOf(2)
            expect(groupsRegistry[1]).to.deep.equal(
              keccak256(newGroupPublicKey)
            )
          })

          it("should store group data", async () => {
            const storedGroup = await groups.getGroup(newGroupPublicKey)

            expect(storedGroup.groupPubKey).to.be.equal(newGroupPublicKey)
            expect(storedGroup.activationBlockNumber).to.be.equal(0)
            expect(storedGroup.membersHash).to.be.equal(
              hashUint32Array(newGroupMembers)
            )
          })

          it("should not update existing group", async () => {
            const groupsRegistry = await groups.getGroupsRegistry()

            expect(groupsRegistry[0]).to.deep.equal(
              keccak256(existingGroupPublicKey)
            )

            const storedGroup = await groups.getGroup(existingGroupPublicKey)

            expect(storedGroup.groupPubKey).to.be.equal(existingGroupPublicKey)
            expect(storedGroup.activationBlockNumber).to.be.equal(0)
            expect(storedGroup.membersHash).to.be.equal(
              hashUint32Array(existingGroupMembers)
            )
          })
        })
      })

      context("when existing group is active", async () => {
        let existingGroup

        beforeEach(async () => {
          await groups.activateCandidateGroup()

          existingGroup = await groups.getGroup(existingGroupPublicKey)
        })

        context("with the same group public key", async () => {
          const newGroupPublicKey = existingGroupPublicKey

          it("should revert with 'Group with this public key was already activated' error", async () => {
            await expect(
              groups.addCandidateGroup(
                newGroupPublicKey,
                newGroupMembers,
                noMisbehaved
              )
            ).to.be.revertedWith(
              "Group with this public key was already activated"
            )
          })
        })

        context("with unique group public key", async () => {
          const newGroupPublicKey = groupPublicKey

          let tx: ContractTransaction

          beforeEach(async () => {
            tx = await groups.addCandidateGroup(
              newGroupPublicKey,
              newGroupMembers,
              noMisbehaved
            )
          })

          it("should emit CandidateGroupRegistered event", async () => {
            await expect(tx)
              .to.emit(groups, "CandidateGroupRegistered")
              .withArgs(newGroupPublicKey)
          })

          it("should register group", async () => {
            const groupsRegistry = await groups.getGroupsRegistry()

            expect(groupsRegistry).to.be.lengthOf(2)
            expect(groupsRegistry[1]).to.deep.equal(
              keccak256(newGroupPublicKey)
            )
          })

          it("should store group data", async () => {
            const storedGroup = await groups.getGroup(newGroupPublicKey)

            expect(storedGroup.groupPubKey).to.be.equal(newGroupPublicKey)
            expect(storedGroup.activationBlockNumber).to.be.equal(0)
            expect(storedGroup.membersHash).to.be.equal(
              hashUint32Array(newGroupMembers)
            )
          })

          it("should not update existing group", async () => {
            const groupsRegistry = await groups.getGroupsRegistry()

            expect(groupsRegistry[0]).to.deep.equal(
              keccak256(existingGroupPublicKey)
            )

            const storedGroup = await groups.getGroup(existingGroupPublicKey)

            expect(storedGroup.groupPubKey).to.be.equal(existingGroupPublicKey)
            expect(storedGroup.activationBlockNumber).to.be.equal(
              existingGroup.activationBlockNumber
            )
            expect(storedGroup.membersHash).to.be.equal(
              hashUint32Array(existingGroupMembers)
            )
          })
        })
      })
    })

    context("when existing group was popped", async () => {
      const existingGroupPublicKey = "0x1234567890"

      let existingGroupMembers: number[]
      let newGroupMembers: number[]

      beforeEach(async () => {
        existingGroupMembers = members.slice(30)
        newGroupMembers = members.slice(-30)

        await groups.addCandidateGroup(
          existingGroupPublicKey,
          existingGroupMembers,
          noMisbehaved
        )

        await groups.popCandidateGroup()
      })

      context("with the same group public key", async () => {
        const newGroupPublicKey = existingGroupPublicKey

        let tx: ContractTransaction

        beforeEach(async () => {
          tx = await groups.addCandidateGroup(
            newGroupPublicKey,
            newGroupMembers,
            noMisbehaved
          )
        })

        it("should emit CandidateGroupRegistered event", async () => {
          await expect(tx)
            .to.emit(groups, "CandidateGroupRegistered")
            .withArgs(newGroupPublicKey)
        })

        it("should register group", async () => {
          const groupsRegistry = await groups.getGroupsRegistry()

          expect(groupsRegistry).to.be.lengthOf(1)
          expect(groupsRegistry[0]).to.deep.equal(keccak256(newGroupPublicKey))
        })

        it("should store group data", async () => {
          const storedGroup = await groups.getGroup(newGroupPublicKey)

          expect(storedGroup.groupPubKey).to.be.equal(newGroupPublicKey)
          expect(storedGroup.activationBlockNumber).to.be.equal(0)
          expect(storedGroup.membersHash).to.be.equal(
            hashUint32Array(newGroupMembers)
          )
        })
      })

      context("with unique group public key", async () => {
        const newGroupPublicKey = groupPublicKey

        let tx: ContractTransaction

        beforeEach(async () => {
          tx = await groups.addCandidateGroup(
            newGroupPublicKey,
            newGroupMembers,
            noMisbehaved
          )
        })

        it("should emit CandidateGroupRegistered event", async () => {
          await expect(tx)
            .to.emit(groups, "CandidateGroupRegistered")
            .withArgs(newGroupPublicKey)
        })

        it("should register group", async () => {
          const groupsRegistry = await groups.getGroupsRegistry()

          expect(groupsRegistry).to.be.lengthOf(1)
          expect(groupsRegistry[0]).to.deep.equal(keccak256(newGroupPublicKey))
        })

        it("should store group data", async () => {
          const storedGroup = await groups.getGroup(newGroupPublicKey)

          expect(storedGroup.groupPubKey).to.be.equal(newGroupPublicKey)
          expect(storedGroup.activationBlockNumber).to.be.equal(0)
          expect(storedGroup.membersHash).to.be.equal(
            hashUint32Array(newGroupMembers)
          )
        })

        it("should not update existing group", async () => {
          const storedGroup = await groups.getGroup(existingGroupPublicKey)

          expect(storedGroup.groupPubKey).to.be.equal(existingGroupPublicKey)
          expect(storedGroup.activationBlockNumber).to.be.equal(0)
          expect(storedGroup.membersHash).to.be.equal(
            hashUint32Array(existingGroupMembers)
          )
        })
      })
    })
  })

  describe("activateCandidateGroup", async () => {
    context("when no groups are registered", async () => {
      it("should revert with 'group does not exist' error", async () => {
        await expect(groups.activateCandidateGroup()).to.be.revertedWith(
          "reverted with panic code 0x11 (Arithmetic operation underflowed or overflowed outside of an unchecked block)"
        )
      })
    })

    context("when one group is registered", async () => {
      beforeEach(async () => {
        await groups.addCandidateGroup(groupPublicKey, members, noMisbehaved)
      })

      context("when the group is candidate", async () => {
        let tx: ContractTransaction

        beforeEach(async () => {
          tx = await groups.activateCandidateGroup()
        })

        it("should emit GroupActivated event", async () => {
          await expect(tx)
            .to.emit(groups, "GroupActivated")
            .withArgs(0, groupPublicKey)
        })

        it("should set activation block number for the group", async () => {
          expect(
            (await groups.getGroup(groupPublicKey)).activationBlockNumber
          ).to.be.equal(tx.blockNumber)
        })

        it("should increase number of active groups", async () => {
          expect(await groups.numberOfActiveGroups()).to.be.equal(1)
        })
      })

      context("when the group is active", async () => {
        beforeEach(async () => {
          await groups.activateCandidateGroup()
        })

        it("should revert with 'The latest registered group was already activated' error", async () => {
          await expect(groups.activateCandidateGroup()).to.be.revertedWith(
            "The latest registered group was already activated"
          )
        })
      })
    })

    context("when two groups are registered", async () => {
      let members1: number[]
      let members2: number[]

      beforeEach(async () => {
        members1 = members.slice(30)
        members2 = members.slice(-30)
      })

      context("with unique group public keys", async () => {
        const groupPublicKey1 = "0x0001"
        const groupPublicKey2 = "0x0002"

        context("when candidate group added after activation", async () => {
          let tx: ContractTransaction

          beforeEach(async () => {
            await groups.addCandidateGroup(
              groupPublicKey1,
              members1,
              noMisbehaved
            )

            tx = await groups.activateCandidateGroup()

            await groups.addCandidateGroup(
              groupPublicKey2,
              members2,
              noMisbehaved
            )
          })

          it("should emit GroupActivated event", async () => {
            await expect(tx)
              .to.emit(groups, "GroupActivated")
              .withArgs(0, groupPublicKey1)
          })

          it("should set activation block number for the activated group", async () => {
            expect(
              (await groups.getGroup(groupPublicKey1)).activationBlockNumber
            ).to.be.equal(tx.blockNumber)
          })

          it("should not set activation block number for the other group", async () => {
            expect(
              (await groups.getGroup(groupPublicKey2)).activationBlockNumber
            ).to.be.equal(0)
          })

          it("should increase number of active groups", async () => {
            expect(await groups.numberOfActiveGroups()).to.be.equal(1)
          })
        })

        context("when the other group is active", async () => {
          let activationBlockNumber1: number
          let tx1: ContractTransaction
          let tx2: ContractTransaction

          // TODO: Update as the latest group got activated
          beforeEach(async () => {
            await groups.addCandidateGroup(
              groupPublicKey1,
              members1,
              noMisbehaved
            )

            tx1 = await groups.activateCandidateGroup()
            activationBlockNumber1 = (
              await ethers.provider.getBlock(tx1.blockHash)
            ).number

            await groups.addCandidateGroup(
              groupPublicKey2,
              members2,
              noMisbehaved
            )

            tx2 = await groups.activateCandidateGroup()
          })

          it("should emit GroupActivated event for the first group", async () => {
            await expect(tx1)
              .to.emit(groups, "GroupActivated")
              .withArgs(0, groupPublicKey1)
          })

          it("should emit GroupActivated event for the second group", async () => {
            await expect(tx2)
              .to.emit(groups, "GroupActivated")
              .withArgs(1, groupPublicKey2)
          })

          it("should not update activation block number for the other group", async () => {
            expect(
              (await groups.getGroup(groupPublicKey1)).activationBlockNumber
            ).to.be.equal(activationBlockNumber1)
          })

          it("should set activation block number for the activated group", async () => {
            expect(
              (await groups.getGroup(groupPublicKey2)).activationBlockNumber
            ).to.be.equal(tx2.blockNumber)
          })

          it("should increase number of active groups", async () => {
            expect(await groups.numberOfActiveGroups()).to.be.equal(2)
          })
        })
      })

      context("with the same group public key", async () => {
        const groupPublicKey1 = groupPublicKey
        const groupPublicKey2 = groupPublicKey

        context(
          "when the first group was challenged and replaced",
          async () => {
            let tx: ContractTransaction

            beforeEach(async () => {
              await groups.addCandidateGroup(
                groupPublicKey1,
                members1,
                noMisbehaved
              )

              await groups.popCandidateGroup()

              await groups.addCandidateGroup(
                groupPublicKey2,
                members2,
                noMisbehaved
              )

              tx = await groups.activateCandidateGroup()
            })

            it("should emit GroupActivated event", async () => {
              await expect(tx)
                .to.emit(groups, "GroupActivated")
                .withArgs(0, groupPublicKey2)
            })

            it("should set activation block number for the group", async () => {
              expect(
                (await groups.getGroup(groupPublicKey2)).activationBlockNumber
              ).to.be.equal(tx.blockNumber)
            })

            it("should increase number of active groups", async () => {
              expect(await groups.numberOfActiveGroups()).to.be.equal(1)
            })
          }
        )
      })
    })
  })

  describe("popCandidateGroup", async () => {
    context("when no groups are registered", async () => {
      it("should revert with 'group does not exist' error", async () => {
        await expect(groups.popCandidateGroup()).to.be.revertedWith(
          "reverted with panic code 0x11 (Arithmetic operation underflowed or overflowed outside of an unchecked block)"
        )
      })
    })

    context("when one group is registered", async () => {
      beforeEach(async () => {
        await groups.addCandidateGroup(groupPublicKey, members, noMisbehaved)
      })

      context("when the group is candidate", async () => {
        let tx: ContractTransaction

        beforeEach(async () => {
          tx = await groups.popCandidateGroup()
        })

        it("should emit CandidateGroupRemoved event", async () => {
          await expect(tx)
            .to.emit(groups, "CandidateGroupRemoved")
            .withArgs(groupPublicKey)
        })

        it("should remove registered group", async () => {
          const groupsRegistry = await groups.getGroupsRegistry()
          expect(groupsRegistry).to.be.lengthOf(0)
        })

        it("should not update stored group data", async () => {
          const storedGroup = await groups.getGroup(groupPublicKey)
          expect(storedGroup.groupPubKey).to.be.equal(groupPublicKey)
          expect(storedGroup.activationBlockNumber).to.be.equal(0)
          expect(storedGroup.membersHash).to.be.equal(hashUint32Array(members))
        })
      })

      context("when the group is active", async () => {
        beforeEach(async () => {
          await groups.activateCandidateGroup()
        })

        it("should revert with 'The latest registered group was already activated' error", async () => {
          await expect(groups.activateCandidateGroup()).to.be.revertedWith(
            "The latest registered group was already activated"
          )
        })
      })
    })

    context("when two groups are registered", async () => {
      let members1: number[]
      let members2: number[]

      beforeEach(async () => {
        members1 = members.slice(30)
        members2 = members.slice(-30)
      })

      context("with unique group public keys", async () => {
        const groupPublicKey1 = "0x0001"
        const groupPublicKey2 = "0x0002"

        context("when the other group is active", async () => {
          let activationBlockNumber1: number
          let tx: ContractTransaction

          beforeEach(async () => {
            await groups.addCandidateGroup(
              groupPublicKey1,
              members1,
              noMisbehaved
            )
            const tx1 = await groups.activateCandidateGroup()
            activationBlockNumber1 = (
              await ethers.provider.getBlock(tx1.blockHash)
            ).number

            await groups.addCandidateGroup(
              groupPublicKey2,
              members2,
              noMisbehaved
            )

            tx = await groups.popCandidateGroup()
          })

          it("should emit CandidateGroupRemoved event", async () => {
            await expect(tx)
              .to.emit(groups, "CandidateGroupRemoved")
              .withArgs(groupPublicKey2)
          })

          it("should remove registered group", async () => {
            const groupsRegistry = await groups.getGroupsRegistry()

            expect(groupsRegistry).to.be.lengthOf(1)
            expect(groupsRegistry[0]).to.deep.equal(keccak256(groupPublicKey1))
          })

          it("should not update stored group data", async () => {
            const storedGroup1 = await groups.getGroup(groupPublicKey1)

            expect(storedGroup1.groupPubKey).to.be.equal(groupPublicKey1)
            expect(storedGroup1.activationBlockNumber).to.be.equal(
              activationBlockNumber1
            )
            expect(storedGroup1.membersHash).to.be.equal(
              hashUint32Array(members1)
            )

            const storedGroup2 = await groups.getGroup(groupPublicKey2)

            expect(storedGroup2.groupPubKey).to.be.equal(groupPublicKey2)
            expect(storedGroup2.activationBlockNumber).to.be.equal(0)
            expect(storedGroup2.membersHash).to.be.equal(
              hashUint32Array(members2)
            )
          })
        })
      })
    })
  })
})

function filterMisbehaved(
  members: number[],
  misbehavedIndices: number[]
): number[] {
  const expectedMembers = [...members]
  misbehavedIndices.reverse().forEach((value) => {
    expectedMembers[value - 1] = expectedMembers[expectedMembers.length - 1]
    expectedMembers.pop()
  })

  return expectedMembers
}
