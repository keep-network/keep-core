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
    keepGroupImplV1, keepGroupImplProxy, keepGroupViaProxy,
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
  });

  it("should be able to check if the implementation contract was initialized", async function() {
    let result = await keepGroupViaProxy.initialized();
    assert.equal(result, true, "Implementation contract should be initialized.");
  });
});
