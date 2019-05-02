import { duration } from './helpers/increaseTime';
import mineBlocks from './helpers/mineBlocks';
const StakingProxy = artifacts.require('./StakingProxy.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const KeepToken = artifacts.require('./KeepToken.sol');
const KeepRandomBeaconProxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconImplV1 = artifacts.require('./KeepRandomBeaconImplV1.sol');
const KeepGroupProxy = artifacts.require('./KeepGroup.sol');
const KeepGroupImplV1 = artifacts.require('./KeepGroupImplV1Stub.sol');

const minimumStake = 200000;
const groupThreshold = 15;
const groupSize = 20;
const timeoutInitial = 20;
const timeoutSubmission = 50;
const timeoutChallenge = 60;
const timeDKG = 20;
const resultPublicationBlockStep = 3;
const groupActiveTime = 300;
const activeGroupsThreshold = 5;
const testGroupsNumber = 10;
const expirationStepTime = groupActiveTime / 10;
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
      groupActiveTime
    );
  });

  async function expireGroup(groupIndex) {
    let groupRegistrationBlock = await keepGroupImplViaProxy.getGroupRegistrationBlockHeight(groupIndex);
    let currentBlock = await web3.eth.getBlockNumber();
    // If current block is larger than group registration block by group active time then
    // it is not necessary to mine any blocks cause the group is already expired
    if (currentBlock - groupRegistrationBlock <= groupActiveTime)
      await mineBlocks(groupActiveTime - (currentBlock - groupRegistrationBlock) + 1);
  }

  it("it should be able to count the number of active groups", async function() {
    for (var i = 1; i <= testGroupsNumber; i++)
      await keepGroupImplViaProxy.registerNewGroup([i]);

    let numberOfGroups = await keepGroupImplViaProxy.numberOfGroups();
    assert.equal(Number(numberOfGroups), testGroupsNumber, "Number of groups is not equal to number of test groups");
  });

  // - we start with [AAAAAAAAAA]
  // - threshold is equal 5
  // - we mine as many blocks as needed to make groups [0, 4] expired
  // - we select group at position 4 which is expired
  // - we should end up with [EEEEEAAAAA]
  it("should mark all groups as expired except active threshold when\
 the selected group is right before of threshold section and it is expired", async function() {

    for (var i = 1; i <= testGroupsNumber; i++) {
      await keepGroupImplViaProxy.registerNewGroup([i]);
    }

    await expireGroup(4);

    await keepGroupImplViaProxy.selectGroup(4) // 4 % 10 = 4
    let expiredOffset = await keepGroupImplViaProxy.getExpiredOffset();

    assert.equal(expiredOffset, activeGroupsThreshold, "Unexpected expired offset");

    let numberOfGroups = await keepGroupImplViaProxy.numberOfGroups();
    assert.equal(Number(numberOfGroups), activeGroupsThreshold, "Number of groups is not equal to active groups threshold");
  });

  // - we start with [AAAAAAAAAA]
  // - threshold is equal 5
  // - we mine as many blocks as needed to make groups [0, 5] expired
  // - we select group at position 5 which is expired
  // - we should end up with [EEEEEAAAAA]
  it("should mark all groups as expired except active threshold when\
 the selected group is right at the beginning of threshold section and it is expired", async function() {

    for (var i = 1; i <= testGroupsNumber; i++) {
      await keepGroupImplViaProxy.registerNewGroup([i]);
    }

    await expireGroup(5);

    await keepGroupImplViaProxy.selectGroup(5) // 5 % 10 = 5
    let expiredOffset = await keepGroupImplViaProxy.getExpiredOffset();

    assert.equal(expiredOffset, activeGroupsThreshold, "Unexpected expired offset");

    let numberOfGroups = await keepGroupImplViaProxy.numberOfGroups();
    assert.equal(Number(numberOfGroups), activeGroupsThreshold, "Number of groups is not equal to active groups threshold");
  });

  // - we start with [AAAAAAAAAA]
  // - threshold is equal 5
  // - we mine as many blocks as needed to make groups [0, 6] expired
  // - we select group at position 6 which is expired
  // - we should end up with [EEEEEAAAAA]
  it("should mark all groups as expired except active threshold when\
 the selected group is right after the beginning of threshold section and it is expired", async function() {

    for (var i = 1; i <= testGroupsNumber; i++) {
      await keepGroupImplViaProxy.registerNewGroup([i]);
    }

    await expireGroup(6);

    await keepGroupImplViaProxy.selectGroup(6) // 6 % 10 = 6
    let expiredOffset = await keepGroupImplViaProxy.getExpiredOffset();

    assert.equal(expiredOffset, activeGroupsThreshold, "Unexpected expired offset");

    let numberOfGroups = await keepGroupImplViaProxy.numberOfGroups();
    assert.equal(Number(numberOfGroups), activeGroupsThreshold, "Number of groups is not equal to active groups threshold");
  });

  // - we start with [AAAAAAAAAA]
  // - threshold is equal 5
  // - we mine as many blocks as needed to mark all groups as expired
  // - we select group at position 0 which is expired
  // - we should end up with [EEEEEAAAAA]
  it("should mark all groups as expired except active threshold when\
 selected the very first group", async function() {
    for (var i = 1; i <= testGroupsNumber; i++)
      await keepGroupImplViaProxy.registerNewGroup([i]);

    await expireGroup(9);

    await keepGroupImplViaProxy.selectGroup(0);

    let expiredOffset = await keepGroupImplViaProxy.getExpiredOffset();
    let numberOfGroups = await keepGroupImplViaProxy.numberOfGroups();

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
    for (var i = 1; i <= testGroupsNumber; i++)
      await keepGroupImplViaProxy.registerNewGroup([i]);

    await expireGroup(9);

    await keepGroupImplViaProxy.selectGroup(testGroupsNumber - 1); // 9

    let expiredOffset = await keepGroupImplViaProxy.getExpiredOffset();
    let numberOfGroups = await keepGroupImplViaProxy.numberOfGroups();

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
    for (var i = 1; i <= testGroupsNumber; i++)
      await keepGroupImplViaProxy.registerNewGroup([i]);

    await expireGroup(9);

    await keepGroupImplViaProxy.selectGroup(testGroupsNumber); // 10

    let expiredOffset = await keepGroupImplViaProxy.getExpiredOffset();
    let numberOfGroups = await keepGroupImplViaProxy.numberOfGroups();

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
    for (var i = 1; i <= testGroupsNumber; i++)
      await keepGroupImplViaProxy.registerNewGroup([i]);

    let after = await keepGroupImplViaProxy.numberOfGroups();

    await expireGroup(9);

    for (var i = 1; i <= testGroupsNumber; i++)
      await keepGroupImplViaProxy.registerNewGroup([i]);

    await keepGroupImplViaProxy.selectGroup(1);

    after = await keepGroupImplViaProxy.numberOfGroups();

    assert.equal(Number(after), testGroupsNumber, "Number of groups should not fall below the test groups number");
  });
});
