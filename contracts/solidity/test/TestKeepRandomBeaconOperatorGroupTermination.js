import mineBlocks from './helpers/mineBlocks';
import {initContracts} from './helpers/initContracts';
import expireGroup from './helpers/expireGroup';

contract('TestKeepRandomBeaconOperatorGroupTermination', function() {

    let operatorContract;

    const groupActiveTime = 100;
    const activeGroupsThreshold = 1;
    const relayEntryTimeout = 10;

    beforeEach(async () => {
      let contracts = await initContracts(
        artifacts.require('./KeepToken.sol'),
        artifacts.require('./TokenStaking.sol'),
        artifacts.require('./KeepRandomBeaconService.sol'),
        artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
        artifacts.require('./KeepRandomBeaconOperatorStub.sol')
      );
      
      operatorContract = contracts.operatorContract;
      
      operatorContract.setGroupActiveTime(groupActiveTime);
      operatorContract.setActiveGroupsThreshold(activeGroupsThreshold);
      operatorContract.setRelayEntryTimeout(relayEntryTimeout);
    });

    async function runTerminationTest(groupsCount, expiredCount, terminatedGroups, beaconValue ) {
      for (var i = 1; i <= groupsCount; i++) {
        await operatorContract.registerNewGroup([i]);
        await mineBlocks(10);
      }

      if (expiredCount > 0) {
        // expire group accepts group index, we need to subtract one from the 
        // count since we index from 0.
        await expireGroup(operatorContract, expiredCount - 1); 
      }

      for (const groupIndex of terminatedGroups) {
        await operatorContract.terminateGroup(groupIndex);
      }

      return operatorContract.selectGroup.call(beaconValue);      
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
});