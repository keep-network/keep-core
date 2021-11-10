/* eslint-disable no-await-in-loop */

import { ethers, waffle, helpers } from "hardhat"
import { expect } from "chai"
import type { ContractTransaction } from "ethers"
import blsData from "./data/bls"
import { noMisbehaved } from "./utils/dkg"
import { constants } from "./fixtures"
import type { GroupsStub } from "../typechain"

const { keccak256 } = ethers.utils

const fixture = async () => {
  const GroupsStub = await ethers.getContractFactory("GroupsStub")
  const groups = await GroupsStub.deploy()

  return groups
}

const { mineBlocksTo } = helpers.time

describe("Groups", () => {
  const groupPublicKey: string = ethers.utils.hexValue(blsData.groupPubKey)
  const members: number[] = []
  const groupLifetime = 20
  const relayEntryTimeout = 10

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
          expect(storedGroup.members).to.be.deep.equal(members)
        })
      })

      context("with misbehaved members", async () => {
        context("with first member misbehaved", async () => {
          const misbehavedIndices: number[] = [1]

          beforeEach(async () => {
            tx = await groups.addCandidateGroup(
              groupPublicKey,
              members,
              misbehavedIndices
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
            tx = await groups.addCandidateGroup(
              groupPublicKey,
              members,
              misbehavedIndices
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
            tx = await groups.addCandidateGroup(
              groupPublicKey,
              members,
              misbehavedIndices
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
            tx = await groups.addCandidateGroup(
              groupPublicKey,
              members,
              misbehavedIndices
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
            expect(storedGroup.members).to.be.deep.equal(newGroupMembers)
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
            expect(storedGroup.members).to.be.deep.equal(newGroupMembers)
          })

          it("should not update existing group", async () => {
            const groupsRegistry = await groups.getGroupsRegistry()

            expect(groupsRegistry[0]).to.deep.equal(
              keccak256(existingGroupPublicKey)
            )

            const storedGroup = await groups.getGroup(existingGroupPublicKey)

            expect(storedGroup.groupPubKey).to.be.equal(existingGroupPublicKey)
            expect(storedGroup.activationBlockNumber).to.be.equal(0)
            expect(storedGroup.members).to.be.deep.equal(existingGroupMembers)
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

          it("should revert with 'group with this public key was already activated' error", async () => {
            await expect(
              groups.addCandidateGroup(
                newGroupPublicKey,
                newGroupMembers,
                noMisbehaved
              )
            ).to.be.revertedWith(
              "group with this public key was already activated"
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
            expect(storedGroup.members).to.be.deep.equal(newGroupMembers)
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
            expect(storedGroup.members).to.be.deep.equal(existingGroupMembers)
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
          expect(storedGroup.members).to.be.deep.equal(newGroupMembers)
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
          expect(storedGroup.members).to.be.deep.equal(newGroupMembers)
        })

        it("should not update existing group", async () => {
          const storedGroup = await groups.getGroup(existingGroupPublicKey)

          expect(storedGroup.groupPubKey).to.be.equal(existingGroupPublicKey)
          expect(storedGroup.activationBlockNumber).to.be.equal(0)
          expect(storedGroup.members).to.be.deep.equal(existingGroupMembers)
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

        it("should revert with 'the latest registered group was already activated' error", async () => {
          await expect(groups.activateCandidateGroup()).to.be.revertedWith(
            "the latest registered group was already activated"
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

        context("when both groups are candidate", async () => {
          let tx: ContractTransaction

          beforeEach(async () => {
            await groups.addCandidateGroup(
              groupPublicKey1,
              members1,
              noMisbehaved
            )
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
              .withArgs(1, groupPublicKey2)
          })

          it("should not set activation block number for the other group", async () => {
            expect(
              (await groups.getGroup(groupPublicKey1)).activationBlockNumber
            ).to.be.equal(0)
          })

          it("should set activation block number for the activated group", async () => {
            expect(
              (await groups.getGroup(groupPublicKey2)).activationBlockNumber
            ).to.be.equal(tx.blockNumber)
          })

          it("should increase number of active groups", async () => {
            expect(await groups.numberOfActiveGroups()).to.be.equal(1)
          })
        })

        context("when the other group is active", async () => {
          let activationBlockNumber1: number
          let tx: ContractTransaction

          // TODO: Update as the latest group got actiavted
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

            tx = await groups.activateCandidateGroup()
          })

          it("should emit GroupActivated event", async () => {
            await expect(tx)
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
            ).to.be.equal(tx.blockNumber)
          })

          it("should increase number of active groups", async () => {
            expect(await groups.numberOfActiveGroups()).to.be.equal(2)
          })
        })
      })

      context("with the same group public key", async () => {
        const groupPublicKey1 = groupPublicKey
        const groupPublicKey2 = groupPublicKey

        context("when both groups are candidate", async () => {
          let tx: ContractTransaction

          beforeEach(async () => {
            await groups.addCandidateGroup(
              groupPublicKey1,
              members1,
              noMisbehaved
            )
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
              .withArgs(1, groupPublicKey2)
          })

          it("should set activation block number for the group", async () => {
            expect(
              (await groups.getGroup(groupPublicKey2)).activationBlockNumber
            ).to.be.equal(tx.blockNumber)
          })

          it("should increase number of active groups", async () => {
            expect(await groups.numberOfActiveGroups()).to.be.equal(1)
          })
        })
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
          expect(storedGroup.members).to.be.deep.equal(members)
        })
      })

      context("when the group is active", async () => {
        beforeEach(async () => {
          await groups.activateCandidateGroup()
        })

        it("should revert with 'the latest registered group was already activated' error", async () => {
          await expect(groups.activateCandidateGroup()).to.be.revertedWith(
            "the latest registered group was already activated"
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

        context("when both groups are candidate", async () => {
          let tx: ContractTransaction

          beforeEach(async () => {
            await groups.addCandidateGroup(
              groupPublicKey1,
              members1,
              noMisbehaved
            )
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
            expect(storedGroup1.activationBlockNumber).to.be.equal(0)
            expect(storedGroup1.members).to.be.deep.equal(members1)

            const storedGroup2 = await groups.getGroup(groupPublicKey2)

            expect(storedGroup2.groupPubKey).to.be.equal(groupPublicKey2)
            expect(storedGroup2.activationBlockNumber).to.be.equal(0)
            expect(storedGroup2.members).to.be.deep.equal(members2)
          })
        })

        context("when the other group is active", async () => {
          let activationBlockNumber1: number
          let tx: ContractTransaction

          // TODO: Update as the latest group got actiavted
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
            expect(storedGroup1.members).to.be.deep.equal(members1)

            const storedGroup2 = await groups.getGroup(groupPublicKey2)

            expect(storedGroup2.groupPubKey).to.be.equal(groupPublicKey2)
            expect(storedGroup2.activationBlockNumber).to.be.equal(0)
            expect(storedGroup2.members).to.be.deep.equal(members2)
          })
        })
      })

      context("with the same group public key", async () => {
        const groupPublicKey1 = groupPublicKey
        const groupPublicKey2 = groupPublicKey

        context("when both groups are candidate", async () => {
          let tx: ContractTransaction

          beforeEach(async () => {
            await groups.addCandidateGroup(
              groupPublicKey1,
              members1,
              noMisbehaved
            )
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
            expect(storedGroup1.activationBlockNumber).to.be.equal(0)
            expect(storedGroup1.members).to.be.deep.equal(members2)

            const storedGroup2 = await groups.getGroup(groupPublicKey2)

            expect(storedGroup2.groupPubKey).to.be.equal(groupPublicKey2)
            expect(storedGroup2.activationBlockNumber).to.be.equal(0)
            expect(storedGroup2.members).to.be.deep.equal(members2)
          })
        })
      })
    })
  })

  describe.only("expireOldGroups", async () => {
    beforeEach(async () => {
      await groups.setGroupLifetime(groupLifetime)
      await groups.setRelayEntryTimeout(relayEntryTimeout)
    })

    context("when expiring old groups and selecting active ones", async () => {
      it("A beacon_value = 0", async () => {
        const selectedIndex = await runExpirationTest(1, 0, 0)
        expect(selectedIndex).to.be.equal(0)
      })
      it("A beacon_value = 1", async () => {
        const selectedIndex = await runExpirationTest(1, 0, 1)
        expect(selectedIndex).to.be.equal(0)
      })
      it("AAA beacon_value = 0", async () => {
        const selectedIndex = await runExpirationTest(3, 0, 0)
        expect(selectedIndex).to.be.equal(0)
      })
      it("AAA beacon_value = 1", async () => {
        const selectedIndex = await runExpirationTest(3, 0, 1)
        expect(selectedIndex).to.be.equal(1)
      })
      it("AAA beacon_value = 2", async () => {
        const selectedIndex = await runExpirationTest(3, 0, 2)
        expect(selectedIndex).to.be.equal(2)
      })
      it("AAA beacon_value = 3", async () => {
        const selectedIndex = await runExpirationTest(3, 0, 3)
        expect(selectedIndex).to.be.equal(0)
      })
      it("EAA beacon_value = 0", async () => {
        const selectedIndex = await runExpirationTest(3, 1, 0)
        expect(selectedIndex).to.be.equal(1)
      })
      it("EEEEAAAAAA beacon_value = 0", async () => {
        const selectedIndex = await runExpirationTest(10, 4, 0)
        expect(selectedIndex).to.be.equal(4)
      })
      it("EEEEAAAAAA beacon_value = 1", async () => {
        const selectedIndex = await runExpirationTest(10, 4, 1)
        expect(selectedIndex).to.be.equal(5)
      })
      it("EEEEAAAAAA beacon_value = 2", async () => {
        const selectedIndex = await runExpirationTest(10, 4, 2)
        expect(selectedIndex).to.be.equal(6)
      })
      it("EEEEAAAAAA beacon_value = 3", async () => {
        const selectedIndex = await runExpirationTest(10, 4, 3)
        expect(selectedIndex).to.be.equal(7)
      })
      it("EEEEAAAAAA beacon_value = 4", async () => {
        const selectedIndex = await runExpirationTest(10, 4, 4)
        expect(selectedIndex).to.be.equal(8)
      })
      it("EEEEAAAAAA beacon_value = 5", async () => {
        const selectedIndex = await runExpirationTest(10, 4, 5)
        expect(selectedIndex).to.be.equal(9)
      })
      it("EEEEAAAAAA beacon_value = 6", async () => {
        const selectedIndex = await runExpirationTest(10, 4, 6)
        expect(selectedIndex).to.be.equal(4)
      })
      it("EEEEAAAAAA beacon_value = 7", async () => {
        const selectedIndex = await runExpirationTest(10, 4, 7)
        expect(selectedIndex).to.be.equal(5)
      })
      it("EEEEAAAAAA beacon_value = 8", async () => {
        const selectedIndex = await runExpirationTest(10, 4, 8)
        expect(selectedIndex).to.be.equal(6)
      })
      it("EEEEAAAAAA beacon_value = 9", async () => {
        const selectedIndex = await runExpirationTest(10, 4, 9)
        expect(selectedIndex).to.be.equal(7)
      })
      it("EEEEAAAAAA beacon_value = 10", async () => {
        const selectedIndex = await runExpirationTest(10, 4, 10)
        expect(selectedIndex).to.be.equal(8)
      })
      it("EEEEAAAAAA beacon_value = 11", async () => {
        const selectedIndex = await runExpirationTest(10, 4, 11)
        expect(selectedIndex).to.be.equal(9)
      })
      it("EEEEEEEEEA beacon_value = 0", async () => {
        const selectedIndex = await runExpirationTest(10, 9, 0)
        expect(selectedIndex).to.be.equal(9)
      })
      it("EEEEEEEEEA beacon_value = 1", async () => {
        const selectedIndex = await runExpirationTest(10, 9, 1)
        expect(selectedIndex).to.be.equal(9)
      })
      it("EEEEEEEEEA beacon_value = 10", async () => {
        const selectedIndex = await runExpirationTest(10, 9, 10)
        expect(selectedIndex).to.be.equal(9)
      })
      it("EEEEEEEEEA beacon_value = 11", async () => {
        const selectedIndex = await runExpirationTest(10, 9, 11)
        expect(selectedIndex).to.be.equal(9)
      })
    })

    it("should be able to count the number of active groups", async () => {
      const expectedGroupCount = 23
      await addGroups(expectedGroupCount)
      const numberOfGroups = await groups.numberOfActiveGroups()
      expect(numberOfGroups).to.be.equal(expectedGroupCount)
    })

    it("should revert group selection when all groups expired", async () => {
      await addGroups(5)

      const currentBlock = await ethers.provider.getBlock("latest")
      await mineBlocksTo(currentBlock.number + groupLifetime)

      await expect(groups.selectGroup(0)).to.be.revertedWith("No active groups")
    })

    // - we start with [AAAAAA]
    // - we check whether the first group is stale and assert it is not since
    //   an active group cannot be stale
    it("should not mark group as stale if it is active", async () => {
      await addGroups(6)

      const isStale = await groups.isStaleGroup(ethers.utils.hexlify(1))

      expect(isStale).to.be.equal(false)
    })

    // - we start with [AAAAAAAAAAAAAAA]
    // - we expire the first 10 groups so that we have [EEEEEEEEEEAAAAA]
    // - we check whether any of active groups is stale and assert it's not
    it("should not mark group as stale if it is active and there are other expired groups", async () => {
      const groupsCount = 15
      await addGroups(groupsCount)
      await expireGroup(8) // move height to expire first 9 groups (we index from 0)

      // this will move height by one and expire 9 + 1 groups
      await groups.selectGroup(0)

      for (let i = 10; i < groupsCount; i++) {
        const isStale = await groups.isStaleGroupById(i)

        expect(isStale).to.be.equal(false)
      }
    })

    // - we start with [AAAAAAAAAAAAAAA]
    // - we expire the first 10 groups so that we have [EEEEEEEEEEAAAAA]
    // - we mine as many blocks as needed to mark expired groups as stale
    // - we check whether any of active groups is stale and assert it's not
    it("should not mark group as stale if it is active and there are other stale groups", async () => {
      const groupsCount = 15
      await addGroups(groupsCount)
      await expireGroup(8) // move height to expire first 9 groups (we index from 0)

      // this will move height by one and expire 9 + 1 groups
      await groups.selectGroup(0)
      const currentBlock = await ethers.provider.getBlock("latest")

      await mineBlocksTo(relayEntryTimeout + currentBlock.number)

      for (let i = 10; i < groupsCount; i++) {
        const isStale = await groups.isStaleGroupById(i)

        expect(isStale).to.be.equal(false)
      }
    })

    // - we start with [AAAAAA]
    // - we mine as many blocks as needed to qualify the first group as expired
    //   and we run group selection to mark it as such; we have [EAAAAA]
    // - we check whether this group is a stale group and assert it is not since
    //   relay request timeout did not pass since the group expiration block
    it("should not mark group as stale if it is expired but can be still signing relay entry", async () => {
      await addGroups(6)

      await expireGroup(0)
      await groups.selectGroup(0)

      const isStale = await groups.isStaleGroupById(0)

      expect(isStale).to.be.equal(false)
    })

    // - we start with [AAAAAA]
    // - we mine as many blocks as needed to qualify the first group as expired
    //   and we run group selection to mark it as such; we have [EAAAAA]
    // - we mine as many blocks as defined by relay request timeout
    // - we check whether this group is a stale group and assert it is stale since
    //   relay request timeout did pass since the group expiration block
    it("should mark group as stale if it is expired and can be no longer signing relay entry", async () => {
      await addGroups(6)

      await expireGroup(0)
      await groups.selectGroup(0)

      const currentBlock = await ethers.provider.getBlock("latest")
      await mineBlocksTo(currentBlock.number + relayEntryTimeout)

      const isStale = await groups.isStaleGroupById(0)

      expect(isStale).to.be.equal(true)
    })

    it("should allow to add and select new group even if all other groups expired", async () => {
      await addGroups(5)
      const currentBlock = await ethers.provider.getBlock("latest")
      await mineBlocksTo(currentBlock.number + groupLifetime)

      await groups.addCandidateGroup(
        ethers.utils.hexlify(6),
        members,
        noMisbehaved
      )
      await groups.activateCandidateGroup()

      const selected = await groups.callStatic.selectGroup(0)
      await groups.selectGroup(0)
      const numberOfGroups = await groups.numberOfActiveGroups()

      expect(numberOfGroups).to.be.equal(1)
      expect(selected).to.be.equal(5)
    })
  })

  async function addGroups(numberOfGroups) {
    for (let i = 1; i <= numberOfGroups; i++) {
      await groups.addCandidateGroup(
        ethers.utils.hexlify(i),
        members,
        noMisbehaved
      )
      await groups.activateCandidateGroup()
    }
  }

  async function expireGroup(groupIndex) {
    const group = await groups.getGroupById(groupIndex)
    const activationBlock = group.activationBlockNumber
    const currentBlock = await ethers.provider.getBlock("latest")

    if (currentBlock.number - activationBlock.toNumber() <= groupLifetime) {
      const minedBlocksToExpireGroup =
        currentBlock.number +
        (groupLifetime - (currentBlock.number - activationBlock.toNumber())) +
        1
      await mineBlocksTo(minedBlocksToExpireGroup)
    }
  }

  async function runExpirationTest(numberOfGroups, expiredCount, beaconValue) {
    await addGroups(numberOfGroups)
    if (expiredCount > 0) {
      // expire group accepts group index, we need to subtract one from the
      // count since we index from 0.
      await expireGroup(expiredCount - 1)
    }
    const selectedGroup = await groups.callStatic.selectGroup(beaconValue)
    return selectedGroup
  }
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
