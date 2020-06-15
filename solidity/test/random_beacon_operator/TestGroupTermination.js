const {expectRevert, time} = require("@openzeppelin/test-helpers")
const assert = require('chai').assert
const {createSnapshot, restoreSnapshot} = require("../helpers/snapshot.js")
const {contract, accounts, web3} = require("@openzeppelin/test-environment")
const GroupsTerminationStub = contract.fromArtifact('GroupsTerminationStub')
const Groups = contract.fromArtifact('Groups');
const BLS = contract.fromArtifact('BLS');

describe('KeepRandomBeaconOperator/GroupTermination', function() {
    let groups;

    const groupActiveTime = web3.utils.toBN(5);

    before(async () => {
      const bls = await BLS.new({from: accounts[0]});
      await Groups.detectNetwork()
      await Groups.link("BLS", bls.address);
      const groupsLibrary = await Groups.new();
      await GroupsTerminationStub.detectNetwork()
      await GroupsTerminationStub.link("Groups", groupsLibrary.address);
      groups = await GroupsTerminationStub.new();
    });

    beforeEach(async () => {
      await createSnapshot()
    });

    afterEach(async () => {
      await restoreSnapshot()
    });

    async function runTerminationTest(groupsCount, expiredCount, terminatedGroups, beaconValue ) {
      await groups.registerNewGroups(expiredCount);
      await time.advanceBlockTo(groupActiveTime.addn(await web3.eth.getBlockNumber()))
      await groups.registerNewGroups(groupsCount - expiredCount);

      for (const groupIndex of terminatedGroups) {
        await groups.terminateGroup(groupIndex);
      }

      return groups.selectGroup.call(beaconValue);
    }

    describe("performs selection with terminated groups not in ascending order", async () => {
      it("AATAT beacon_value = 0", async () => {
        let selectedIndex = await runTerminationTest(5, 0, [4, 2], 0)
        assert.equal(0, selectedIndex)
      })
      it("AATAT beacon_value = 1", async () => {
        let selectedIndex = await runTerminationTest(5, 0, [4, 2], 1)
        assert.equal(1, selectedIndex)
      })
      it("AATAT beacon_value = 2", async () => {
        let selectedIndex = await runTerminationTest(5, 0, [4, 2], 2)
        assert.equal(3, selectedIndex)
      })
      it("TATATTA beacon_value = 0", async () => {
        let selectedIndex = await runTerminationTest(7, 0, [5, 2, 4, 0], 0)
        assert.equal(1, selectedIndex)
      })
      it("TATATTA beacon_value = 1", async () => {
        let selectedIndex = await runTerminationTest(7, 0, [5, 2, 4, 0], 1)
        assert.equal(3, selectedIndex)
      })
      it("TATATTA beacon_value = 2", async () => {
        let selectedIndex = await runTerminationTest(7, 0, [5, 2, 4, 0], 2)
        assert.equal(6, selectedIndex)
      })
      it("AATATTA beacon_value = 0", async () => {
        let selectedIndex = await runTerminationTest(7, 0, [5, 2, 4], 0)
        assert.equal(0, selectedIndex)
      })
      it("TATATAAT beacon_value = 0", async () => {
        let selectedIndex = await runTerminationTest(8, 0, [7, 0, 4, 2], 0)
        assert.equal(1, selectedIndex)
      })
      it("TATATAAT beacon_value = 1", async () => {
        let selectedIndex = await runTerminationTest(8, 0, [7, 0, 4, 2], 1)
        assert.equal(3, selectedIndex)
      })
      it("TATATAAT beacon_value = 2", async () => {
        let selectedIndex = await runTerminationTest(8, 0, [7, 0, 4, 2], 2)
        assert.equal(5, selectedIndex)
      })
      it("TATATAAT beacon_value = 3", async () => {
        let selectedIndex = await runTerminationTest(8, 0, [7, 0, 4, 2], 3)
        assert.equal(6, selectedIndex)
      })
    })
    
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

    describe("should fail when there are no active groups", async () => {
      it("T", async function() {
        await expectRevert(
          runTerminationTest(1, 0, [0], 0),
          "No active groups"
        );
      })
      it("TT", async function() {
        await expectRevert(
          runTerminationTest(2, 0, [0, 1], 0),
          "No active groups"
        );
      })
      it("ET", async function () {
        await expectRevert(
          runTerminationTest(2, 1, [1], 0),
          "No active groups"
        );
      })
    })
});
