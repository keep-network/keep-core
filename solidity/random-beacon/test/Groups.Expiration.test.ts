/* eslint-disable no-await-in-loop */

import { ethers, waffle, helpers } from "hardhat"
import { expect } from "chai"
import { noMisbehaved } from "./utils/dkg"
import type { GroupsStub } from "../typechain"

const fixture = async () => {
  const GroupsStub = await ethers.getContractFactory("GroupsStub")
  const groups = await GroupsStub.deploy()

  return groups
}

const { mineBlocksTo } = helpers.time

describe("Groups", () => {
  const members: number[] = []
  const groupLifetime = 20
  const relayEntryTimeout = 10

  let groups: GroupsStub

  beforeEach("load test fixture", async () => {
    groups = await waffle.loadFixture(fixture)
  })

  describe("expireOldGroups", async () => {
    beforeEach(async () => {
      await groups.setGroupLifetime(groupLifetime)
      await groups.setRelayEntryTimeout(relayEntryTimeout)
    })

    context("when active groups were created", async () => {
      it("should be able to count the number of active groups", async () => {
        const expectedGroupCount = 23
        await addGroups(expectedGroupCount)
        const numberOfGroups = await groups.numberOfActiveGroups()
        expect(numberOfGroups).to.be.equal(expectedGroupCount)
      })
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

    context("when all groups expired", async () => {
      it("should revert group selection", async () => {
        await addGroups(5)

        const currentBlock = await ethers.provider.getBlock("latest")
        await mineBlocksTo(currentBlock.number + groupLifetime)

        await expect(groups.selectGroup(0)).to.be.revertedWith(
          "No active groups"
        )
      })

      it("should allow to add and select new group", async () => {
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

    context("when a group is active", async () => {
      // - we start with [AAAAAA]
      // - we check whether the first group is stale and assert it is not since
      //   an active group cannot be stale
      it("should not mark group as stale", async () => {
        await addGroups(6)

        const isStale = await groups.isStaleGroup(ethers.utils.hexlify(1))

        expect(isStale).to.be.equal(false)
      })

      context("when there are other expired groups", async () => {
        // - we start with [AAAAAAAAAAAAAAA]
        // - we expire the first 10 groups so that we have [EEEEEEEEEEAAAAA]
        // - we check whether any of active groups is stale and assert it's not
        it("should not mark group as stale", async () => {
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
      })

      context("when there are other stale groups", async () => {
        // - we start with [AAAAAAAAAAAAAAA]
        // - we expire the first 10 groups so that we have [EEEEEEEEEEAAAAA]
        // - we mine as many blocks as needed to mark expired groups as stale
        // - we check whether any of active groups is stale and assert it's not
        it("should not mark group as stale", async () => {
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
      })
    })

    context("when a group is expired", async () => {
      context("when a group can still sign a relay entry", async () => {
        // - we start with [AAAAAA]
        // - we mine as many blocks as needed to qualify the first group as expired
        //   and we run group selection to mark it as such; we have [EAAAAA]
        // - we check whether this group is a stale group and assert it is not since
        //   relay request timeout did not pass since the group expiration block
        it("should not mark group as stale", async () => {
          await addGroups(6)

          await expireGroup(0)
          await groups.selectGroup(0)

          const isStale = await groups.isStaleGroupById(0)

          expect(isStale).to.be.equal(false)
        })
      })

      context("when a group can no longer sign a relay entry", async () => {
        // - we start with [AAAAAA]
        // - we mine as many blocks as needed to qualify the first group as expired
        //   and we run group selection to mark it as such; we have [EAAAAA]
        // - we mine as many blocks as defined by relay request timeout
        // - we check whether this group is a stale group and assert it is stale since
        //   relay request timeout did pass since the group expiration block
        it("should mark group as stale", async () => {
          await addGroups(6)

          await expireGroup(0)
          await groups.selectGroup(0)

          const currentBlock = await ethers.provider.getBlock("latest")
          await mineBlocksTo(currentBlock.number + relayEntryTimeout)

          const isStale = await groups.isStaleGroupById(0)

          expect(isStale).to.be.equal(true)
        })
      })
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
