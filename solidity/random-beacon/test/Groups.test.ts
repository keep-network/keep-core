import { ethers } from "hardhat"
import { expect } from "chai"
import type { ContractTransaction } from "ethers"
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

      beforeEach(async () => {
        tx = await groups.addPendingGroup(groupPublicKey, members, noMisbehaved)
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
          it("should revert with 'group was already activated' error", async () => {
            expect(
              groups.addPendingGroup(
                existingGroupPublicKey,
                newGroupMembers,
                noMisbehaved
              )
            ).to.be.revertedWith("group was already activated")
          })
        })

        context("with unique group public key", async () => {
          let tx: ContractTransaction

          beforeEach(async () => {
            tx = await groups.addPendingGroup(
              groupPublicKey,
              newGroupMembers,
              noMisbehaved
            )
          })

          it("should emit PendingGroupRegistered event", async () => {
            expect(tx)
              .to.emit(groups, "PendingGroupRegistered")
              .withArgs(1, groupPublicKey)
          })

          it("should register a pending group", async () => {
            const storedGroup = await groups.getGroup(groupPublicKey)

            expect(storedGroup.groupPubKey).to.be.equal(groupPublicKey)
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
        })
      })
    })

    // TODO: Add tests for setGroupMembers to remove misbehaved
  })
})
