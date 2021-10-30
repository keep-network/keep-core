import { ethers } from "hardhat"
import { expect } from "chai"
import type { ContractTransaction } from "ethers"
import blsData from "./data/bls"
import { noMisbehaved, getDkgGroupSigners } from "./utils/dkg"
import { constants } from "./fixtures"
import type { GroupsStub } from "../typechain"
import type { DkgGroupSigners } from "./utils/dkg"

const { keccak256 } = ethers.utils

describe("Groups", () => {
  const groupPublicKey: string = ethers.utils.hexValue(blsData.groupPubKey)

  let signers: DkgGroupSigners
  let groups: GroupsStub
  let members: string[]

  before(async () => {
    signers = await getDkgGroupSigners(constants.groupSize)
    members = Array.from(signers.values())
  })

  beforeEach("load test fixture", async () => {
    const GroupsStub = await ethers.getContractFactory("GroupsStub")
    groups = await GroupsStub.deploy()
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
          expect(storedGroup.activationTimestamp).to.be.equal(0)
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

      let existingGroupMembers: string[]
      let newGroupMembers: string[]

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
            expect(storedGroup.activationTimestamp).to.be.equal(0)
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
            expect(storedGroup.activationTimestamp).to.be.equal(0)
            expect(storedGroup.members).to.be.deep.equal(newGroupMembers)
          })

          it("should not update existing group", async () => {
            const groupsRegistry = await groups.getGroupsRegistry()

            expect(groupsRegistry[0]).to.deep.equal(
              keccak256(existingGroupPublicKey)
            )

            const storedGroup = await groups.getGroup(existingGroupPublicKey)

            expect(storedGroup.groupPubKey).to.be.equal(existingGroupPublicKey)
            expect(storedGroup.activationTimestamp).to.be.equal(0)
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
            expect(storedGroup.activationTimestamp).to.be.equal(0)
            expect(storedGroup.members).to.be.deep.equal(newGroupMembers)
          })

          it("should not update existing group", async () => {
            const groupsRegistry = await groups.getGroupsRegistry()

            expect(groupsRegistry[0]).to.deep.equal(
              keccak256(existingGroupPublicKey)
            )

            const storedGroup = await groups.getGroup(existingGroupPublicKey)

            expect(storedGroup.groupPubKey).to.be.equal(existingGroupPublicKey)
            expect(storedGroup.activationTimestamp).to.be.equal(
              existingGroup.activationTimestamp
            )
            expect(storedGroup.members).to.be.deep.equal(existingGroupMembers)
          })
        })
      })
    })

    context("when existing group was popped", async () => {
      const existingGroupPublicKey = "0x1234567890"

      let existingGroupMembers: string[]
      let newGroupMembers: string[]

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
          expect(storedGroup.activationTimestamp).to.be.equal(0)
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
          expect(storedGroup.activationTimestamp).to.be.equal(0)
          expect(storedGroup.members).to.be.deep.equal(newGroupMembers)
        })

        it("should not update existing group", async () => {
          const storedGroup = await groups.getGroup(existingGroupPublicKey)

          expect(storedGroup.groupPubKey).to.be.equal(existingGroupPublicKey)
          expect(storedGroup.activationTimestamp).to.be.equal(0)
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

        it("should set activation timestamp for the group", async () => {
          // FIXME: Unclear why `tx.timestamp` is undefined
          const expectedActivationTimestamp = (
            await ethers.provider.getBlock(tx.blockHash)
          ).timestamp

          expect(
            (await groups.getGroup(groupPublicKey)).activationTimestamp
          ).to.be.equal(expectedActivationTimestamp)
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
      let members1: string[]
      let members2: string[]

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

          it("should not set activation timestamp for the other group", async () => {
            expect(
              (await groups.getGroup(groupPublicKey1)).activationTimestamp
            ).to.be.equal(0)
          })

          it("should set activation timestamp for the activated group", async () => {
            // FIXME: Unclear why `tx.timestamp` is undefined
            const expectedActivationTimestamp = (
              await ethers.provider.getBlock(tx.blockHash)
            ).timestamp

            expect(
              (await groups.getGroup(groupPublicKey2)).activationTimestamp
            ).to.be.equal(expectedActivationTimestamp)
          })

          it("should increase number of active groups", async () => {
            expect(await groups.numberOfActiveGroups()).to.be.equal(1)
          })
        })

        context("when the other group is active", async () => {
          let activationTimestamp1: number
          let tx: ContractTransaction

          // TODO: Update as the latest group got actiavted
          beforeEach(async () => {
            await groups.addCandidateGroup(
              groupPublicKey1,
              members1,
              noMisbehaved
            )
            const tx1 = await groups.activateCandidateGroup()
            activationTimestamp1 = (
              await ethers.provider.getBlock(tx1.blockHash)
            ).timestamp

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

          it("should not update activation timestamp for the other group", async () => {
            expect(
              (await groups.getGroup(groupPublicKey1)).activationTimestamp
            ).to.be.equal(activationTimestamp1)
          })

          it("should set activation timestamp for the activated group", async () => {
            // FIXME: Unclear why `tx.timestamp` is undefined
            const expectedActivationTimestamp = (
              await ethers.provider.getBlock(tx.blockHash)
            ).timestamp

            expect(
              (await groups.getGroup(groupPublicKey2)).activationTimestamp
            ).to.be.equal(expectedActivationTimestamp)
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

          it("should set activation timestamp for the group", async () => {
            // FIXME: Unclear why `tx.timestamp` is undefined
            const expectedActivationTimestamp = (
              await ethers.provider.getBlock(tx.blockHash)
            ).timestamp

            expect(
              (await groups.getGroup(groupPublicKey2)).activationTimestamp
            ).to.be.equal(expectedActivationTimestamp)
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
          expect(storedGroup.activationTimestamp).to.be.equal(0)
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
      let members1: string[]
      let members2: string[]

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
            expect(storedGroup1.activationTimestamp).to.be.equal(0)
            expect(storedGroup1.members).to.be.deep.equal(members1)

            const storedGroup2 = await groups.getGroup(groupPublicKey2)

            expect(storedGroup2.groupPubKey).to.be.equal(groupPublicKey2)
            expect(storedGroup2.activationTimestamp).to.be.equal(0)
            expect(storedGroup2.members).to.be.deep.equal(members2)
          })
        })

        context("when the other group is active", async () => {
          let activationTimestamp1: number
          let tx: ContractTransaction

          // TODO: Update as the latest group got actiavted
          beforeEach(async () => {
            await groups.addCandidateGroup(
              groupPublicKey1,
              members1,
              noMisbehaved
            )
            const tx1 = await groups.activateCandidateGroup()
            activationTimestamp1 = (
              await ethers.provider.getBlock(tx1.blockHash)
            ).timestamp

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
            expect(storedGroup1.activationTimestamp).to.be.equal(
              activationTimestamp1
            )
            expect(storedGroup1.members).to.be.deep.equal(members1)

            const storedGroup2 = await groups.getGroup(groupPublicKey2)

            expect(storedGroup2.groupPubKey).to.be.equal(groupPublicKey2)
            expect(storedGroup2.activationTimestamp).to.be.equal(0)
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
            expect(storedGroup1.activationTimestamp).to.be.equal(0)
            expect(storedGroup1.members).to.be.deep.equal(members2)

            const storedGroup2 = await groups.getGroup(groupPublicKey2)

            expect(storedGroup2.groupPubKey).to.be.equal(groupPublicKey2)
            expect(storedGroup2.activationTimestamp).to.be.equal(0)
            expect(storedGroup2.members).to.be.deep.equal(members2)
          })
        })
      })
    })
  })
})

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
