import mineBlocks from './helpers/mineBlocks';
import expectThrowWithMessage from './helpers/expectThrowWithMessage';
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot";
const GroupsTerminationStub = artifacts.require('./stubs/GroupsTerminationStub.sol')
const Groups = artifacts.require("./libraries/operator/Groups.sol");

contract('GroupsTerminationStub', function(accounts) {
    let groups;

    const groupActiveTime = 5;
    const activeGroupsThreshold = 1;

    before(async () => {
      const groupsLibrary = await Groups.new();
      await GroupsTerminationStub.link("Groups", groupsLibrary.address);
      groups = await GroupsTerminationStub.new();
      groups.setActiveGroupsThreshold(activeGroupsThreshold);
    });

    beforeEach(async () => {
      await createSnapshot()
    });

    afterEach(async () => {
      await restoreSnapshot()
    });

    async function runTerminationTest(groupsCount, expiredCount, terminatedGroups, beaconValue ) {
      await groups.registerNewGroups(expiredCount);
      mineBlocks(groupActiveTime);
      await groups.registerNewGroups(groupsCount - expiredCount);

      for (const groupIndex of terminatedGroups) {
        await groups.terminateGroup(groupIndex);
      }

      return groups.selectGroup.call(beaconValue);
    }
    
    describe("should not select terminated groups", async () => {
      it("TA beacon_value = 0", async function() {
        let selectedIndex = await runTerminationTest(2, 0, [0], 0);
        assert.equal(1, selectedIndex)
      })
      it("TA beacon_value = 1", async function() {
        let selectedIndex = await runTerminationTest(2, 0, [0], 1);
        assert.equal(1, selectedIndex)
      })    
      it("TA beacon_value = 2", async function() { 
        let selectedIndex = await runTerminationTest(2, 0, [0], 2);
        assert.equal(1, selectedIndex);
      })
      it("AT beacon_value = 0", async function() {
        let selectedIndex = await runTerminationTest(2, 0, [1], 0);
        assert.equal(0, selectedIndex);
      }) 
      it("AT beacon_value = 1", async function() {
        let selectedIndex = await runTerminationTest(2, 0, [1], 1);
        assert.equal(0, selectedIndex);
      })
      it("AT beacon_value = 2", async function() {
        let selectedIndex = await runTerminationTest(2, 0, [1], 2);
        assert.equal(0, selectedIndex);
      })
      it("TAA beacon_value = 0", async function() {
        let selectedIndex = await runTerminationTest(3, 0, [0], 0);
        assert.equal(1, selectedIndex);
      })
      it("TAA beacon_value = 1", async function() {
        let selectedIndex = await runTerminationTest(3, 0, [0], 1);
        assert.equal(2, selectedIndex);
      })
      it("TAA beacon_value = 2", async function() {
        let selectedIndex = await runTerminationTest(3, 0, [0], 2);
        assert.equal(1, selectedIndex);
      })
      it("AAT beacon_value = 0", async function() {
        let selectedIndex = await runTerminationTest(3, 0, [2], 0);
        assert.equal(0, selectedIndex);
      })
      it("AAT beacon_value = 1", async function() {
        let selectedIndex = await runTerminationTest(3, 0, [2], 1);
        assert.equal(1, selectedIndex);
      })
      it("AAT beacon_value = 2", async function() {
        let selectedIndex = await runTerminationTest(3, 0, [2], 2);
        assert.equal(0, selectedIndex);
      })
      it("ATA beacon_value = 0", async function() {
        let selectedIndex = await runTerminationTest(3, 0, [1], 0);
        assert.equal(0, selectedIndex);
      })
      it("ATA beacon_value = 1", async function() {
        let selectedIndex = await runTerminationTest(3, 0, [1], 1);
        assert.equal(2, selectedIndex);
      })
      it("ATA beacon_value = 2", async function() {
        let selectedIndex = await runTerminationTest(3, 0, [1], 2);
        assert.equal(0, selectedIndex);
      })
      it("TTA beacon_value = 0", async function() {
        let selectedIndex = await runTerminationTest(3, 0, [0, 1], 0);
        assert.equal(2, selectedIndex);
      })
      it("TTA beacon_value = 1", async function() {
        let selectedIndex = await runTerminationTest(3, 0, [0, 1], 1);
        assert.equal(2, selectedIndex);
      })
      it("ATT beacon_value = 0", async function() {
        let selectedIndex = await runTerminationTest(3, 0, [1, 2], 0);
        assert.equal(0, selectedIndex);
      })
      it("ATT beacon_value = 1", async function() {
        let selectedIndex = await runTerminationTest(3, 0, [1, 2], 1);
        assert.equal(0, selectedIndex);
      })
    })
    
    describe("should not select expired or terminated groups", async () => {
      it("ETA beacon_value = 0", async function() {
        let selectedIndex = await runTerminationTest(3, 1, [1], 0);
        assert.equal(2, selectedIndex)
      })
      it("ETA beacon_value = 1", async function() {
        let selectedIndex = await runTerminationTest(3, 1, [1], 1);
        assert.equal(2, selectedIndex)
      })
      it("ETA beacon_value = 2", async function() {
        let selectedIndex = await runTerminationTest(3, 1, [1], 2);
        assert.equal(2, selectedIndex)
      })
      it("ETA beacon_value = 3", async function() {
        let selectedIndex = await runTerminationTest(3, 1, [1], 3);
        assert.equal(2, selectedIndex)
      })
      it("EAT beacon_value = 0", async function() {
        let selectedIndex = await runTerminationTest(3, 1, [2], 0);
        assert.equal(1, selectedIndex)
      })
      it("EAT beacon_value = 1", async function() {
        let selectedIndex = await runTerminationTest(3, 1, [2], 1);
        assert.equal(1, selectedIndex)
      })
      it("EAT beacon_value = 2", async function() {
        let selectedIndex = await runTerminationTest(3, 1, [2], 2);
        assert.equal(1, selectedIndex)
      })
      it("EAT beacon_value = 3", async function() {
        let selectedIndex = await runTerminationTest(3, 1, [2], 3);
        assert.equal(1, selectedIndex)
      })
      it("EEETTATAAT beacon_value = 0", async function() {
        let selectedIndex = await runTerminationTest(10, 3, [3, 4, 6, 9], 0);
        assert.equal(5, selectedIndex)
      })
      it("EEETTATAAT beacon_value = 1", async function() {
        let selectedIndex = await runTerminationTest(10, 3, [3, 4, 6, 9], 1);
        assert.equal(7, selectedIndex)
      })
      it("EEETTATAAT beacon_value = 2", async function() {
        let selectedIndex = await runTerminationTest(10, 3, [3, 4, 6, 9], 2);
        assert.equal(8, selectedIndex)
      })
      it("EEETTATAAT beacon_value = 3", async function() {
        let selectedIndex = await runTerminationTest(10, 3, [3, 4, 6, 9], 3);
        assert.equal(5, selectedIndex)
      })
      it("EEETTATAAT beacon_value = 4", async function() {
        let selectedIndex = await runTerminationTest(10, 3, [3, 4, 6, 9], 4);
        assert.equal(7, selectedIndex)
      })
      it("EEETTATAAT beacon_value = 5", async function() {
        let selectedIndex = await runTerminationTest(10, 3, [3, 4, 6, 9], 5);
        assert.equal(8, selectedIndex)
      })
    })

    describe("should include terminated groups when checking the minimum active groups threshold", async () => {    
      beforeEach(async () => {
        await groups.setActiveGroupsThreshold(5); 
      });      
      /*
        We do not expire any more groups because the minimum active threshold
        condition is not met (4 < 5) and we do not take the terminated group
        into account for group selection.
      */
      it("EEEET beacon_value = 0, active threshold = 5", async function() {
        let selectedIndex = await runTerminationTest(5, 4, [4], 0);
        assert.equal(0, selectedIndex)
      })
      it("EEEET beacon_value = 3, active threshold = 5", async function() {
        let selectedIndex = await runTerminationTest(5, 4, [4], 3);
        assert.equal(3, selectedIndex)
      })
      it("EEEET beacon_value = 4, active threshold = 5", async function() {
        let selectedIndex = await runTerminationTest(5, 4, [4], 4);
        assert.equal(0, selectedIndex)
      })
      /*
        We do not expire any more groups because the minimum active threshold
        condition would not be met (5 = 5) and we do not take the terminated
        group into account for group selection.
      */
      it("EEEEET beacon_value = 0, active threshold = 5", async function() {
        groups.setActiveGroupsThreshold(5); 
        let selectedIndex = await runTerminationTest(6, 5, [5], 0);
        assert.equal(0, selectedIndex)
      })
      it("EEEEET beacon_value = 4, active threshold = 5", async function() {
        groups.setActiveGroupsThreshold(5); 
        let selectedIndex = await runTerminationTest(6, 5, [5], 4);
        assert.equal(4, selectedIndex)
      })
      it("EEEEET beacon_value = 5, active threshold = 5", async function() {
        groups.setActiveGroupsThreshold(5); 
        let selectedIndex = await runTerminationTest(6, 5, [5], 5);
        assert.equal(0, selectedIndex)
      })
    })

    describe("should fail when there are no active groups", async () => {
      it("T", async function() {
        await expectThrowWithMessage(
          runTerminationTest(1, 0, [0], 0), 
          "At least one active group required"
        );
      })
      it("TT", async function() {
        await expectThrowWithMessage(
          runTerminationTest(2, 0, [0, 1], 0), 
          "At least one active group required"
        );
      })
    })
});