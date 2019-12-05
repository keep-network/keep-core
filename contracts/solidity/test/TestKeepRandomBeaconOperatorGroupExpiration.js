import mineBlocks from './helpers/mineBlocks';
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot";
const GroupsExpirationStub = artifacts.require('./stubs/GroupsExpirationStub.sol')
import expectThrowWithMessage from './helpers/expectThrowWithMessage';
const Groups = artifacts.require("./libraries/operator/Groups.sol");

contract('GroupsExpirationStub', function(accounts) {

  let groups;

  const groupActiveTime = 20;
  const activeGroupsThreshold = 5;
  const relayEntryTimeout = 10;

  before(async () => {
    const groupsLibrary = await Groups.new();
    await GroupsExpirationStub.link("Groups", groupsLibrary.address);
    groups = await GroupsExpirationStub.new();
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

  async function addGroups(numberOfGroups) {
    for (var i = 1; i <= numberOfGroups; i++)
      await groups.addGroup([i]);
  }

  async function expireGroup(groupIndex) {
    let groupRegistrationBlock = await groups.getGroupRegistrationBlockHeight(groupIndex);
    let currentBlock = await web3.eth.getBlockNumber();

    // If current block is larger than group registration block by group active time then
    // it is not necessary to mine any blocks cause the group is already expired
    if (currentBlock - groupRegistrationBlock <= groupActiveTime) {
      await mineBlocks(groupActiveTime - (currentBlock - groupRegistrationBlock));
    }
  }

  async function runExpirationTest(groupSize, expiredCount, beaconValue) {
    await addGroups(groupSize);
    if (expiredCount > 0) {
      // expire group accepts group index, we need to subtract one from the 
      // count since we index from 0.
      await expireGroup(expiredCount - 1); 
    }
    return groups.selectGroup.call(beaconValue);
  }

  it("should be able to count the number of active groups", async function() {
    let expectedGroupCount = 23;
    await addGroups(expectedGroupCount);
    let numberOfGroups = await groups.numberOfGroups();
    assert.equal(Number(numberOfGroups), expectedGroupCount, "Unexpected number of groups");
  });

  describe("should expire old groups and select active one", async () => {
    it("A beacon_value = 0", async function() {
      let selectedIndex = await runExpirationTest(1, 0, 0);
      assert.equal(0, selectedIndex);
    });
    it("A beacon_value = 1", async function() {
      let selectedIndex = await runExpirationTest(1, 0, 1);
      assert.equal(0, selectedIndex);
    });
    it("AAA beacon_value = 0", async function() {
      let selectedIndex = await runExpirationTest(3, 0, 0);
      assert.equal(0, selectedIndex);
    });
    it("AAA beacon_value = 1", async function() {
      let selectedIndex = await runExpirationTest(3, 0, 1);
      assert.equal(1, selectedIndex);
    });
    it("AAA beacon_value = 2", async function() {
      let selectedIndex = await runExpirationTest(3, 0, 2);
      assert.equal(2, selectedIndex);
    });
    it("AAA beacon_value = 3", async function() {
      let selectedIndex = await runExpirationTest(3, 0, 3);
      assert.equal(0, selectedIndex);
    });
    it("EAA beacon_value = 0", async function() {
      let selectedIndex = await runExpirationTest(3, 1, 0);
      assert.equal(0, selectedIndex); // min active threshold does not let to move offset
    });
    it("EAAAA beacon_value = 0", async function() {
      let selectedIndex = await runExpirationTest(5, 1, 0);
      assert.equal(0, selectedIndex); // min active threshold does not let to move offset
    });
    it("EAAAAA beacon_value = 0", async function() {
      let selectedIndex = await runExpirationTest(6, 1, 0);
      assert.equal(1, selectedIndex); // min active threshold does allow to move offset
    });
    it("EEEEAAAAAA beacon_value = 0", async function() { 
      let selectedIndex = await runExpirationTest(10, 4, 0);
      assert.equal(4, selectedIndex);
    });
    it("EEEEAAAAAA beacon_value = 1", async function() {
      let selectedIndex = await runExpirationTest(10, 4, 1);
      assert.equal(5, selectedIndex);
    });
    it("EEEEAAAAAA beacon_value = 2", async function() {
      let selectedIndex = await runExpirationTest(10, 4, 2);
      assert.equal(6, selectedIndex);
    });
    it("EEEEAAAAAA beacon_value = 3", async function() {
      let selectedIndex = await runExpirationTest(10, 4, 3);
      assert.equal(7, selectedIndex);
    });
    it("EEEEAAAAAA beacon_value = 4", async function() {
      let selectedIndex = await runExpirationTest(10, 4, 4);
      assert.equal(8, selectedIndex);
    });
    it("EEEEAAAAAA beacon_value = 5", async function() {
      let selectedIndex = await runExpirationTest(10, 4, 5);
      assert.equal(9, selectedIndex);
    });
    it("EEEEAAAAAA beacon_value = 6", async function() {
      let selectedIndex = await runExpirationTest(10, 4, 6);
      assert.equal(4, selectedIndex);
    });
    it("EEEEAAAAAA beacon_value = 7", async function() {
      let selectedIndex = await runExpirationTest(10, 4, 7);
      assert.equal(5, selectedIndex);
    });
    it("EEEEAAAAAA beacon_value = 8", async function() {
      let selectedIndex = await runExpirationTest(10, 4, 8);
      assert.equal(6, selectedIndex);
    });
    it("EEEEAAAAAA beacon_value = 9", async function() {
      let selectedIndex = await runExpirationTest(10, 4, 9);
      assert.equal(7, selectedIndex);
    });
    it("EEEEAAAAAA beacon_value = 10", async function() {
      let selectedIndex = await runExpirationTest(10, 4, 10);
      assert.equal(8, selectedIndex);
    });
    it("EEEEAAAAAA beacon_value = 11", async function() {
      let selectedIndex = await runExpirationTest(10, 4, 11);
      assert.equal(9, selectedIndex);
    });
    it("EEEEEEEEEA beacon_value = 0", async function() {
      let selectedIndex = await runExpirationTest(10, 9, 0);
      assert.equal(5, selectedIndex); // min active threshold does not let to move offset further than to 5
    });
    it("EEEEEEEEEA beacon_value = 1", async function() {
      let selectedIndex = await runExpirationTest(10, 9, 1);
      assert.equal(6, selectedIndex); // min active threshold does not let to move offset further than to 5
    });
    it("EEEEEEEEEA beacon_value = 10", async function() {
      let selectedIndex = await runExpirationTest(10, 9, 10);
      assert.equal(5, selectedIndex); // min active threshold does not let to move offset further than to 5
    });
    it("EEEEEEEEEA beacon_value = 11", async function() {
      let selectedIndex = await runExpirationTest(10, 9, 11);
      assert.equal(6, selectedIndex); // min active threshold does not let to move offset further than to 5
    });
  });
  
  // - we start with [AAAAAA]
  // - we check whether the first group is stale and assert it is not since
  //   an active group cannot be stale
  it("should not mark group as stale if it is active", async function() {
    let groupsCount = activeGroupsThreshold + 1
    await addGroups(groupsCount);

    let pubKey = await groups.getGroupPublicKey(0);

    let isStale = await groups.isStaleGroup(pubKey);

    assert.equal(isStale, false, "Group should not be marked as stale");
  });

  // - we start with [AAAAAAAAAAAAAAA]
  // - we expire the first 10 groups so that we have [EEEEEEEEEEAAAAA]
  // - we check whether any of active groups is stale and assert it's not
  it("should not mark group as stale if it is active and\
 there are other expired groups", async function() {
    let groupsCount = 15
    await addGroups(groupsCount);
    await expireGroup(9); // expire first 10 groups (we index from 0)

    await groups.selectGroup(0);

    for (var i = 10; i < groupsCount; i++) {
      let pubKey = await groups.getGroupPublicKey(i);
      let isStale = await groups.isStaleGroup(pubKey);

      assert.equal(isStale, false, "Group should not be marked as stale")
    }
  });

  // - we start with [AAAAAAAAAAAAAAA]
  // - we expire the first 10 groups so that we have [EEEEEEEEEEAAAAA]
  // - we mine as many blocks as needed to mark expired groups as stale
  // - we check whether any of active groups is stale and assert it's not
  it("should not mark group as stale if it is active and\
 there are other stale groups", async function() {
    let groupsCount = 15
    await addGroups(groupsCount);
    await expireGroup(9); // expire first 10 groups (we index from 0)

    await groups.selectGroup(0);

    await mineBlocks(relayEntryTimeout);

    for (var i = 10; i < groupsCount; i++) {
      let pubKey = await groups.getGroupPublicKey(i);
      let isStale = await groups.isStaleGroup(pubKey);

      assert.equal(isStale, false, "Group should not be marked as stale")
    }
  });

  // - we start with [AAAAA]
  // - we mine as many blocks as needed to have all the groups qualify as stale
  // - we check whether the group at position 0 is stale
  // - group should not be marked as stale since it is not marked as expired
  //   (no group selection was triggered); group can be stale only if it has
  //   been marked as expired - `selectGroup` may decide not to mark group as
  //   expired even though it reached its expiration time (minimum threshold)
  it("should not mark group as stale if its expiration time passed but\
 it is not marked as such", async function() {
    let groupsCount = activeGroupsThreshold + 1
    await addGroups(groupsCount);

    let pubKey = await groups.getGroupPublicKey(0);

    // mine blocks but do not select group so it's not marked as expired
    await mineBlocks(groupActiveTime + relayEntryTimeout);

    let isStale  = await groups.isStaleGroup(pubKey);

    assert.equal(isStale, false, "Group should not be marked as stale");
  });

  // - we start with [AAAAAA]
  // - we mine as many blocks as needed to qualify the first group as expired 
  //   and we run group selection to mark it as such; we have [EAAAAA]
  // - we check whether this group is a stale group and assert it is not since
  //   relay request timeout did not pass since the group expiration block
  it("should not mark group as stale if it is expired but\
 can be still signing relay entry", async function() {
    let groupsCount = activeGroupsThreshold + 1
    await addGroups(groupsCount);

    let pubKey = await groups.getGroupPublicKey(0);

    await expireGroup(0);
    await groups.selectGroup(0);

    let isStale  = await groups.isStaleGroup(pubKey);

    assert.equal(isStale, false, "Group should not be marked as stale");
  });

  // - we start with [AAAAAA]
  // - we mine as many blocks as needed to qualify the first group as expired
  //   and we run group selection to mark it as such; we have [EAAAAA]
  // - we mine as many blocks as defined by relay request timeout
  // - we check whether this group is a stale group and assert it is stale since
  //   relay request timeout did pass since the group expiration block
  it("should mark group as stale if it is expired and\
 can be no longer signing relay entry", async function() {
     let groupsCount = activeGroupsThreshold + 1
     await addGroups(groupsCount);
 
     let pubKey = await groups.getGroupPublicKey(0);
 
     await expireGroup(0);
     await groups.selectGroup(0);
 
     await mineBlocks(relayEntryTimeout);

     let isStale  = await groups.isStaleGroup(pubKey);

     assert.equal(isStale, true, "Group should be marked as stale");
   });

   // - we start with [AAAAAA]
   // - we check whether group with a non-existing public key is stale and
   //   we assert the check should fail
   it("should fail stale check if group could not be found", async function() {
    let groupsCount = activeGroupsThreshold + 1
    await addGroups(groupsCount);

    let pubKey = "0x1337"; // group with such pub key does not exist
    await expectThrowWithMessage(
      groups.isStaleGroup(pubKey),
      "Group does not exist"
    );
  });
});