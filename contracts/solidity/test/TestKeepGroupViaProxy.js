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
    keepGroupImplV1, keepGroupProxy, keepGroupImplViaProxy, groupOnePubKey, groupTwoPubKey,
    account_one = accounts[0],
    account_two = accounts[1];

  beforeEach(async () => {
    token = await KeepToken.new();
    
    // Initialize staking contract under proxy
    stakingProxy = await StakingProxy.new();
    stakingContract = await TokenStaking.new(token.address, stakingProxy.address, duration.days(30));
    await stakingProxy.authorizeContract(stakingContract.address, {from: account_one})

    // Initialize Keep Random Beacon
    let minimumStake = 200;
    keepRandomBeaconImplV1 = await KeepRandomBeaconImplV1.new();
    keepRandomBeaconProxy = await KeepRandomBeaconProxy.new('v1', keepRandomBeaconImplV1.address);
    keepRandomBeaconImplViaProxy = await KeepRandomBeaconImplV1.at(keepRandomBeaconProxy.address);
    await keepRandomBeaconImplViaProxy.initialize(stakingProxy.address, 100, minimumStake, duration.days(30));

    // Initialize Keep Group contract
    keepGroupImplV1 = await KeepGroupImplV1.new();
    keepGroupProxy = await KeepGroupProxy.new('v1', keepGroupImplV1.address);
    keepGroupImplViaProxy = await KeepGroupImplV1.at(keepGroupProxy.address);
    await keepGroupImplViaProxy.initialize(6, 10, keepRandomBeaconProxy.address);

    // Create test groups.
    groupOnePubKey = "0x1000000000000000000000000000000000000000000000000000000000000000";
    await keepGroupImplViaProxy.createGroup(groupOnePubKey);
    groupTwoPubKey = "0x2000000000000000000000000000000000000000000000000000000000000000";
    await keepGroupImplViaProxy.createGroup(groupTwoPubKey);

    // Stake tokens as account one so it has minimum stake to be able to get into a group.
    await token.approveAndCall(stakingContract.address, minimumStake, "", {from: account_one});

    // Add member to the first group. 
    await keepGroupImplViaProxy.addMemberToGroup(groupOnePubKey, account_one);
  
  });

  it("should be able to check if the implementation contract was initialized", async function() {
    assert.equal(await keepGroupImplViaProxy.initialized(), true, "Implementation contract should be initialized.");
  });

  it("should be able to return a total number of created group", async function() {
    assert.equal(await keepGroupImplViaProxy.numberOfGroups(), 2, "Should get correct total group count.");
  });

  it("should be able to get a group index number with provided group public key", async function() {
    assert.equal(await keepGroupImplViaProxy.getGroupIndex(groupOnePubKey), 0, "Should get correct group index number for group one.");
    assert.equal(await keepGroupImplViaProxy.getGroupIndex(groupTwoPubKey), 1, "Should get correct group index number for group two.");
  });

  it("should be able to get group public key by group index number", async function() {
    assert.equal(await keepGroupImplViaProxy.getGroupPubKey(0), groupOnePubKey, "Should get group public key.");
  });

  it("should be able to add a member to specified group", async function() {
    assert.equal(await keepGroupImplViaProxy.isMember(groupOnePubKey, account_one), true, "Member should be added to the group.");
  });

  // it("should not be able to add a member to specified group if member has no minimum stake", async function() {
  //   await exceptThrow(keepGroupImplViaProxy.addMemberToGroup(groupOnePubKey, account_two));
  // });

});
