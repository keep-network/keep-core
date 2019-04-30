import { duration } from './helpers/increaseTime';
import mineBlocks from './helpers/mineBlocks';
const StakingProxy = artifacts.require('./StakingProxy.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const KeepToken = artifacts.require('./KeepToken.sol');
const KeepRandomBeaconProxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconImplV1 = artifacts.require('./KeepRandomBeaconImplV1.sol');
const KeepGroupProxy = artifacts.require('./KeepGroup.sol');
const KeepGroupImplV1 = artifacts.require('./KeepGroupImplV1.sol');

const minimumStake = 200000;
const groupThreshold = 15;
const groupSize = 20;
const timeoutInitial = 20;
const timeoutSubmission = 50;
const timeoutChallenge = 60;
const timeDKG = 20;
const resultPublicationBlockStep = 3;
const groupExpirationTime = 300;
const activeGroupsThreshold = 5;
const testGroupsNumber = 10;
const expirationStepTime = groupExpirationTime / 10;
const expectedOffset = 5;

contract('TestKeepGroupExpiration', function(accounts) {

  let token, stakingProxy, stakingContract,
    keepRandomBeaconImplV1, keepRandomBeaconProxy,
    keepGroupImplV1, keepGroupProxy, keepGroupImplViaProxy,
    owner = accounts[0]

  beforeEach(async () => {
    token = await KeepToken.new();
    // Initialize staking contract under proxy
    stakingProxy = await StakingProxy.new();
    stakingContract = await TokenStaking.new(token.address, stakingProxy.address, duration.days(30));
    await stakingProxy.authorizeContract(stakingContract.address, {from: owner})
    
    keepRandomBeaconImplV1 = await KeepRandomBeaconImplV1.new();
    keepRandomBeaconProxy = await KeepRandomBeaconProxy.new(keepRandomBeaconImplV1.address);

    // Initialize Keep Group contract

    keepGroupImplV1 = await KeepGroupImplV1.new();
    keepGroupProxy = await KeepGroupProxy.new(keepGroupImplV1.address);
    keepGroupImplViaProxy = await KeepGroupImplV1.at(keepGroupProxy.address);

    await keepGroupImplViaProxy.initialize(
      stakingProxy.address, keepRandomBeaconProxy.address, minimumStake,
      groupThreshold, groupSize, timeoutInitial, timeoutSubmission,
      timeoutChallenge, timeDKG, resultPublicationBlockStep, activeGroupsThreshold,
      groupExpirationTime
    );
  });

  async function testExpiration(expire, selected) {
    await mineBlocks(expirationStepTime*expire);
    await keepGroupImplViaProxy.selectGroup(selected);
    let expiredOffset = await keepGroupImplViaProxy.getExpiredOffset();
    return Number(expiredOffset);
  }

  it("it should be able to count the number of active groups", async function() {
    for (var i = 1; i <= testGroupsNumber; i++)
      await keepGroupImplViaProxy.registerNewGroup([i]);

    let numberOfGroups = await keepGroupImplViaProxy.numberOfGroups();
    assert.equal(Number(numberOfGroups), testGroupsNumber, "Number of groups not equals to number of test groups");
  });

  /* Following 6 tests are manualy derived from the analyisis of the aggressive group marking.
   * There are finetuned for the following parameters:
   *
   *  groupExpirationTime = 300;
   *  activeGroupsThreshold = 5;
   *  testGroupsNumber = 10;
   *  expirationStepTime = groupExpirationTime / 10;
   *
   * After every change of the above parameters the following tests will need to be updated.
   */
  it("it should mark all groups as expired except activeGroupsThreshold #1", async function() {
    for (var i = 1; i <= testGroupsNumber; i++) {
      await keepGroupImplViaProxy.registerNewGroup([i]); // 2 blocks
      mineBlocks(8);
    }

    let expiredOffset = await testExpiration(4, 0);
    assert.equal(Number(expiredOffset), expectedOffset, "Expired offset should be equal expected expire offset");
  });

  it("it should mark all groups as expired except activeGroupsThreshold #2", async function() {
    for (var i = 1; i <= testGroupsNumber; i++) {
      await keepGroupImplViaProxy.registerNewGroup([i]); // 2 blocks
      mineBlocks(8);
    }

    let expiredOffset = await testExpiration(7, 1);
    assert.equal(Number(expiredOffset), expectedOffset, "Expired offset should be equal expected expire offset");
  });

  it("it should mark all groups as expired except activeGroupsThreshold #3", async function() {
    for (var i = 1; i <= testGroupsNumber; i++) {
      await keepGroupImplViaProxy.registerNewGroup([i]); // 2 blocks
      mineBlocks(8);
    }

    let expiredOffset = await testExpiration(5, 2);
    assert.equal(Number(expiredOffset), expectedOffset, "Expired offset should be equal expected expire offset");
  });

  it("it should mark all groups as expired except activeGroupsThreshold #4", async function() {
    for (var i = 1; i <= testGroupsNumber; i++) {
      await keepGroupImplViaProxy.registerNewGroup([i]); // 2 blocks
      mineBlocks(8);
    }

    let expiredOffset = await testExpiration(4, 4);
    assert.equal(Number(expiredOffset), expectedOffset, "Expired offset should be equal expected expire offset");
  });

  it("it should mark all groups as expired except activeGroupsThreshold #5", async function() {
    for (var i = 1; i <= testGroupsNumber; i++) {
      await keepGroupImplViaProxy.registerNewGroup([i]); // 2 blocks
      mineBlocks(8);
    }

    let expiredOffset = await testExpiration(5, 5);
    assert.equal(Number(expiredOffset), expectedOffset, "Expired offset should be equal expected expire offset");
    });

  it("it should mark all groups as expired except activeGroupsThreshold #6", async function() {
    for (var i = 1; i <= testGroupsNumber; i++) {
      await keepGroupImplViaProxy.registerNewGroup([i]); // 2 blocks
      mineBlocks(8);
    }

    let expiredOffset = await testExpiration(6, 6);
    assert.equal(Number(expiredOffset), expectedOffset, "Expired offset should be equal expected expire offset");
  });

  it("should be able to check if at least one group is marked as expired", async function() {
    for (var i = 1; i <= testGroupsNumber; i++)
      await keepGroupImplViaProxy.registerNewGroup([i]);

    let numberOfGroups = await keepGroupImplViaProxy.numberOfGroups();
    
    mineBlocks(groupExpirationTime);
    await keepGroupImplViaProxy.selectGroup(1);
    numberOfGroups = await keepGroupImplViaProxy.numberOfGroups();

    assert.equal(Number(numberOfGroups), activeGroupsThreshold, "Some groups should be marked as expired");
  });

  it("should be able to check that groups are marked as expired except the minimal active groups number", async function() {
    for (var i = 1; i <= testGroupsNumber; i++)
      await keepGroupImplViaProxy.registerNewGroup([i]);

    let after = await keepGroupImplViaProxy.numberOfGroups();

    for (var i = 1; i <= testGroupsNumber; i++) {
      mineBlocks(groupExpirationTime);
      await keepGroupImplViaProxy.selectGroup((testGroupsNumber - 1) % i);
      after = await keepGroupImplViaProxy.numberOfGroups();
    }
    
    assert.equal(Number(after), activeGroupsThreshold, "Number of groups should not fall below the threshold of active groups");
  });

  it("it should be able to mark only a subset of groups as expired", async function() {
    for (var i = 1; i <= testGroupsNumber; i++)
      await keepGroupImplViaProxy.registerNewGroup([i]);

    let after = await keepGroupImplViaProxy.numberOfGroups();

    mineBlocks(groupExpirationTime*2);

    for (var i = 1; i <= testGroupsNumber; i++)
      await keepGroupImplViaProxy.registerNewGroup([i]);

    await keepGroupImplViaProxy.selectGroup(1);

    after = await keepGroupImplViaProxy.numberOfGroups();

    assert.equal(Number(after), testGroupsNumber, "Number of groups should not fall below the test groups number");
  });
});
