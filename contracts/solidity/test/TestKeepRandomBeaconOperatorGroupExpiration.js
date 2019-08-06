import mineBlocks from './helpers/mineBlocks';
import {initContracts} from './helpers/initContracts';

contract('TestKeepRandomBeaconOperatorGroupExpiration', function() {

  let operatorContract;
  
  const testGroupsNumber = 10;

  const groupActiveTime = 300;
  const activeGroupsThreshold = 5;
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

  async function addGroups(numberOfGroups) {
    for (var i = 1; i <= numberOfGroups; i++)
      await operatorContract.registerNewGroup([i]);
  }

  async function expireGroup(groupIndex) {
    let groupRegistrationBlock = await operatorContract.getGroupRegistrationBlockHeight(groupIndex);
    let currentBlock = await web3.eth.getBlockNumber();
    // If current block is larger than group registration block by group active time then
    // it is not necessary to mine any blocks cause the group is already expired
    if (currentBlock - groupRegistrationBlock <= groupActiveTime)
      await mineBlocks(groupActiveTime - (currentBlock - groupRegistrationBlock) + 1);
  }

  it("it should be able to count the number of active groups", async function() {

    await addGroups(testGroupsNumber);

    let numberOfGroups = await operatorContract.numberOfGroups();
    assert.equal(Number(numberOfGroups), testGroupsNumber, "Number of groups is not equal to number of test groups");
  });

  // - we start with [AAAAAAAAAA]
  // - threshold is equal 5
  // - we mine as many blocks as needed to make groups [0, 4] expired
  // - we select group at position 4 which is expired
  // - we should end up with [EEEEEAAAAA]
  it("should mark all groups as expired except active threshold when\
 the selected group is right before of threshold section and it is expired", async function() {

    await addGroups(testGroupsNumber);
    await expireGroup(4);

    await operatorContract.selectGroup(4) // 4 % 10 = 4

    let expiredOffset = await operatorContract.expiredGroupOffset();
    let numberOfGroups = await operatorContract.numberOfGroups();

    assert.equal(expiredOffset, activeGroupsThreshold, "Unexpected expired offset");
    assert.equal(Number(numberOfGroups), activeGroupsThreshold, "Number of groups is not equal to active groups threshold");
  });

  // - we start with [AAAAAAAAAA]
  // - threshold is equal 5
  // - we mine as many blocks as needed to make groups [0, 5] expired
  // - we select group at position 5 which is expired
  // - we should end up with [EEEEEAAAAA]
  it("should mark all groups as expired except active threshold when\
 the selected group is right at the beginning of threshold section and it is expired", async function() {

    await addGroups(testGroupsNumber);
    await expireGroup(5);

    await operatorContract.selectGroup(5) // 5 % 10 = 5

    let expiredOffset = await operatorContract.expiredGroupOffset();
    let numberOfGroups = await operatorContract.numberOfGroups();

    assert.equal(expiredOffset, activeGroupsThreshold, "Unexpected expired offset");
    assert.equal(Number(numberOfGroups), activeGroupsThreshold, "Number of groups is not equal to active groups threshold");
  });

  // - we start with [AAAAAAAAAA]
  // - threshold is equal 5
  // - we mine as many blocks as needed to make groups [0, 6] expired
  // - we select group at position 6 which is expired
  // - we should end up with [EEEEEAAAAA]
  it("should mark all groups as expired except active threshold when\
 the selected group is right after the beginning of threshold section and it is expired", async function() {

    await addGroups(testGroupsNumber);
    await expireGroup(6);

    await operatorContract.selectGroup(6) // 6 % 10 = 6

    let expiredOffset = await operatorContract.expiredGroupOffset();
    let numberOfGroups = await operatorContract.numberOfGroups();

    assert.equal(expiredOffset, activeGroupsThreshold, "Unexpected expired offset");
    assert.equal(Number(numberOfGroups), activeGroupsThreshold, "Number of groups is not equal to active groups threshold");
  });

  // - we start with [AAAAAAAAAA]
  // - threshold is equal 5
  // - we mine as many blocks as needed to mark all groups as expired
  // - we select group at position 0 which is expired
  // - we should end up with [EEEEEAAAAA]
  it("should mark all groups as expired except active threshold when\
 selected the very first group", async function() {

    await addGroups(testGroupsNumber);
    await expireGroup(9);

    await operatorContract.selectGroup(0);

    let expiredOffset = await operatorContract.expiredGroupOffset();
    let numberOfGroups = await operatorContract.numberOfGroups();

    assert.equal(expiredOffset, activeGroupsThreshold, "Unexpected expired offset");
    assert.equal(Number(numberOfGroups), activeGroupsThreshold, "Number of groups is not equal to active groups threshold");
  });

  // - we start with [AAAAAAAAAA]
  // - threshold is equal 5
  // - we mine as many blocks as needed to mark all groups as expired
  // - we select group at position 9 (testGroupsNumber - 1) which is expired
  // - we should end up with [EEEEEAAAAA]
  it("should mark all groups as expired except active threshold when\
 selected the very last group", async function() {

    await addGroups(testGroupsNumber);
    await expireGroup(9);

    await operatorContract.selectGroup(testGroupsNumber - 1); // 9

    let expiredOffset = await operatorContract.expiredGroupOffset();
    let numberOfGroups = await operatorContract.numberOfGroups();

    assert.equal(expiredOffset, activeGroupsThreshold, "Unexpected expired offset");
    assert.equal(Number(numberOfGroups), activeGroupsThreshold, "Number of groups is not equal to active groups threshold");
  });

  // - we start with [AAAAAAAAAA]
  // - threshold is equal 5
  // - we mine as many blocks as needed to mark all groups as expired
  // - we select group at position 10 (testGroupsNumber) which is expired
  // - we should end up with [EEEEEAAAAA]
  it("should mark all groups as expired except active threshold when\
 selected the one after the last group (modulo operation check)", async function() {

    await addGroups(testGroupsNumber);
    await expireGroup(9);

    await operatorContract.selectGroup(testGroupsNumber); // 10

    let expiredOffset = await operatorContract.expiredGroupOffset();
    let numberOfGroups = await operatorContract.numberOfGroups();

    assert.equal(expiredOffset, activeGroupsThreshold, "Unexpected expired offset");
    assert.equal(Number(numberOfGroups), activeGroupsThreshold, "Number of groups is not equal to active groups threshold");
  });

  // - we start with [AAAAAAAAAA]
  // - threshold is equal 5
  // - we mine as many blocks as needed to mark all groups as expired
  // - we add more groups so we have [AAAAAAAAAAAAAAAAAAAA]
  // - we select group at position 1 which is expired
  // - we should end up with [EEEEEEEEEEAAAAAAAAAA]
  it("it should be able to mark only a subset of groups as expired", async function() {

    await addGroups(testGroupsNumber);
    await expireGroup(9);

    for (var i = 1; i <= testGroupsNumber; i++)
      await operatorContract.registerNewGroup([i]);

    await operatorContract.selectGroup(1);

    let after = await operatorContract.numberOfGroups();

    assert.equal(Number(after), testGroupsNumber, "Number of groups should not fall below the test groups number");
  });

  // - we start with [A]
  // - threshold is equal 5
  // - we mine as many blocks as needed to mark the group as expired
  // - we select group at position 0 which is expired
  // - we should end up with [A]
  it("should not mark group as expired when\
 there is just one group and it is expired", async function() {

    await addGroups(1);
    await expireGroup(0) // indexed from 0

    await operatorContract.selectGroup(0);

    let expiredOffset = await operatorContract.expiredGroupOffset();
    let numberOfGroups = await operatorContract.numberOfGroups();

    assert.equal(expiredOffset, 0, "Unexpected expired offset");
    assert.equal(Number(numberOfGroups), 1, "Unexpected number of groups");
  });

  // - we start with [AAAA]
  // - threshold is equal 5
  // - we mine as many blocks as needed to mark all the groups as expired
  // - we select group at position 0 which is expired
  // - we should end up with [AAAA]
  it("should not mark groups as expired when there is less groups than threshold\
 and they are all expired ", async function() {
    let groupsCount = activeGroupsThreshold - 1

    await addGroups(groupsCount);
    await expireGroup(groupsCount - 1) // indexed from 0

    await operatorContract.selectGroup(0);

    let expiredOffset = await operatorContract.expiredGroupOffset();
    let numberOfGroups = await operatorContract.numberOfGroups();

    assert.equal(expiredOffset, 0, "Unexpected expired offset");
    assert.equal(Number(numberOfGroups), groupsCount, "Unexpected number of groups");
  });

  // - we start with [AAAAA]
  // - threshold is equal 5
  // - we mine as many blocks as needed to mark all the groups as expired
  // - we select group at position 0 which is expired
  // - we should end up with [AAAAA]
  it("should not mark groups as expired when there is threshold number of groups\
 and they are all expired ", async function() {
    let groupsCount = activeGroupsThreshold
    await addGroups(groupsCount);
    await expireGroup(groupsCount - 1) // indexed from 0

    await operatorContract.selectGroup(0);

    let expiredOffset = await operatorContract.expiredGroupOffset();
    let numberOfGroups = await operatorContract.numberOfGroups();

    assert.equal(expiredOffset, 0, "Unexpected expired offset");
    assert.equal(Number(numberOfGroups), groupsCount, "Unexpected number of groups");
  });

  // - we start with [AAAAAA]
  // - we check whether the first group is stale and assert it is not since
  //   an active group cannot be stale
  it("should not mark group as stale if it is active", async function() {
    let groupsCount = activeGroupsThreshold + 1
    await addGroups(groupsCount);

    let pubKey = await operatorContract.getGroupPublicKey(0);

    let isStale = await operatorContract.isStaleGroup(pubKey);

    assert.equal(isStale, false, "Group should not be marked as stale");
  });

  // - we start with [AAAAAAAAAAAAAAA]
  // - we expire the first 10 groups so that we have [EEEEEEEEEEAAAAA]
  // - we check whether any of active groups is stale and assert it's not
  it("should not mark group as stale if it is active and \
 there are other expired groups", async function() {
    let groupsCount = 15
    await addGroups(groupsCount);
    await expireGroup(9); // expire first 10 groups (we index from 0)

    await operatorContract.selectGroup(0);

    for (var i = 10; i < groupsCount; i++) {
      let pubKey = await operatorContract.getGroupPublicKey(i);
      let isStale = await operatorContract.isStaleGroup(pubKey);

      assert.equal(isStale, false, "Group should not be marked as stale")
    }
  });

  // - we start with [AAAAAAAAAAAAAAA]
  // - we expire the first 10 groups so that we have [EEEEEEEEEEAAAAA]
  // - we mine as many blocks as needed to mark expired groups as stale
  // - we check whether any of active groups is stale and assert it's not
  it("should not mark group as stale if it is active and \
 there are other stale groups", async function() {
    let groupsCount = 15
    await addGroups(groupsCount);
    await expireGroup(9); // expire first 10 groups (we index from 0)

    await operatorContract.selectGroup(0);

    await mineBlocks(relayEntryTimeout);

    for (var i = 10; i < groupsCount; i++) {
      let pubKey = await operatorContract.getGroupPublicKey(i);
      let isStale = await operatorContract.isStaleGroup(pubKey);

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
  it("should not mark group as stale if its expiration time passed but \
 it is not marked as such", async function() {
    let groupsCount = activeGroupsThreshold + 1
    await addGroups(groupsCount);

    let pubKey = await operatorContract.getGroupPublicKey(0);

    // mine blocks but do not select group so it's not marked as expired
    await mineBlocks(groupActiveTime + relayEntryTimeout);

    let isStale  = await operatorContract.isStaleGroup(pubKey);

    assert.equal(isStale, false, "Group should not be marked as stale");
  });

  // - we start with [AAAAAA]
  // - we mine as many blocks as needed to qualify the first group as expired 
  //   and we run group selection to mark it as such; we have [EAAAAA]
  // - we check whether this group is a stale group and assert it is not since
  //   relay request timeout did not pass since the group expiration block
  it("should not mark group as stale if it is expired but \
 can be still signing relay entry", async function() {
    let groupsCount = activeGroupsThreshold + 1
    await addGroups(groupsCount);

    let pubKey = await operatorContract.getGroupPublicKey(0);

    await expireGroup(0);
    await operatorContract.selectGroup(0);

    let isStale  = await operatorContract.isStaleGroup(pubKey);

    assert.equal(isStale, false, "Group should not be marked as stale");
  });

  // - we start with [AAAAAA]
  // - we mine as many blocks as needed to qualify the first group as expired
  //   and we run group selection to mark it as such; we have [EAAAAA]
  // - we mine as many blocks as defined by relay request timeout
  // - we check whether this group is a stale group and assert it is stale since
  //   relay request timeout did pass since the group expiration block
  it("should mark group as stale if it is expired and \
 can be no longer signing relay entry", async function() {
     let groupsCount = activeGroupsThreshold + 1
     await addGroups(groupsCount);
 
     let pubKey = await operatorContract.getGroupPublicKey(0);
 
     await expireGroup(0);
     await operatorContract.selectGroup(0);
 
     await mineBlocks(relayEntryTimeout);

     let isStale  = await operatorContract.isStaleGroup(pubKey);

     assert.equal(isStale, true, "Group should be marked as stale");
   });

   // - we start with [AAAAAA]
   // - we check whether group with a non-existing public key is stale and
   //   we assert it is, since we assume all non-existing groups are stale
   it("should say group is stale if it could not be found", async function() {
    let groupsCount = activeGroupsThreshold + 1
    await addGroups(groupsCount);

    let pubKey = "0x1337"; // group with such pub key does not exist

    let isStale  = await operatorContract.isStaleGroup(pubKey);

    assert.equal(isStale, true, "Group should be marked as stale");
  });
});
