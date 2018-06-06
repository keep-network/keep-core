import increaseTime, { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import exceptThrow from './helpers/expectThrow';
import encodeCall from './helpers/encodeCall';
const KeepToken = artifacts.require('./KeepToken.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const KeepRandomBeaconProxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconImplV1 = artifacts.require('./KeepRandomBeaconImplV1.sol');
const KeepGroupProxy = artifacts.require('./KeepGroup.sol');
const KeepGroupImplV1 = artifacts.require('./KeepGroupImplV1.sol');

contract('TestKeepGroupViaProxy', function(accounts) {

  let token, stakingProxy, stakingContract, 
    keepRandomBeaconImplV1, keepRandomBeaconProxy, keepRandomBeaconImplViaProxy,
    keepGroupImplV1, keepGroupImplProxy, keepGroupViaProxy, groupOnePubKey, groupTwoPubKey,
    account_one = accounts[0],
    account_two = accounts[1];

  beforeEach(async () => {
    token = await KeepToken.new();
    stakingProxy = await StakingProxy.new();
    stakingContract = await TokenStaking.new(token.address, stakingProxy.address, duration.days(30));

    // Initialize Keep Random Beacon
    keepRandomBeaconImplV1 = await KeepRandomBeaconImplV1.new();
    keepRandomBeaconProxy = await KeepRandomBeaconProxy.new('v1', keepRandomBeaconImplV1.address);
    keepRandomBeaconImplViaProxy = await KeepRandomBeaconImplV1.at(keepRandomBeaconProxy.address);
    await keepRandomBeaconImplViaProxy.initialize(stakingProxy.address, 100, 200, duration.days(30));

    // Initialize Keep Group contract
    keepGroupImplV1 = await KeepGroupImplV1.new();
    keepGroupImplProxy = await KeepGroupProxy.new('v1', keepGroupImplV1.address);
    keepGroupViaProxy = await KeepGroupImplV1.at(keepGroupImplProxy.address);
    await keepGroupViaProxy.initialize(6, 10, keepGroupViaProxy.address);

    // Create test groups.
    groupOnePubKey = "0x1000000000000000000000000000000000000000000000000000000000000000";
    await keepGroupViaProxy.createGroup(groupOnePubKey);
    groupTwoPubKey = "0x2000000000000000000000000000000000000000000000000000000000000000";
    await keepGroupViaProxy.createGroup(groupTwoPubKey);

    // Add member to the first group.
    await keepGroupViaProxy.addMemberToGroup(groupOnePubKey, account_one);
  });

  it("should be able to check if the implementation contract was initialized", async function() {
    assert.equal(await keepGroupViaProxy.initialized(), true, "Implementation contract should be initialized.");
  });

  it("should be able to create a group with provided group public key", async function() {
    assert.equal(await keepGroupViaProxy.groupExistsView(groupOnePubKey), true, "Created group should exist.");
  });

  it("should be able to return a total number of created group", async function() {
    assert.equal(await keepGroupViaProxy.getNumberOfGroups(), 2, "Should get correct total group count.");
  });

  it("should be able to get a group index number with provided group public key", async function() {
    assert.equal(await keepGroupViaProxy.getGroupNumber(groupOnePubKey), 1, "Should get correct group index number for group one.");
    assert.equal(await keepGroupViaProxy.getGroupNumber(groupTwoPubKey), 2, "Should get correct group index number for group two.");
  });

  it("should be able to get group public key by group index number", async function() {
    assert.equal(await keepGroupViaProxy.getGroupPubKey(1), groupOnePubKey, "Should get group public key.");
  });

  it("should be able to add a member to specified group", async function() {
    assert.equal(await keepGroupViaProxy.isMember(groupOnePubKey, account_one), true, "Member should be added to the group.");
  });

  it("should be able to get number of members in a group by providing a group index number", async function() {
    assert.equal(await keepGroupViaProxy.getGroupNMembers(1), 1, "Group one should have 1 member.");
    assert.equal(await keepGroupViaProxy.getGroupNMembers(2), 0, "Group two should have 0 members.");
  });

});
