import increaseTime, { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import exceptThrow from './helpers/expectThrow';
const KeepToken = artifacts.require('./KeepToken.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const Proxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconImplV1 = artifacts.require('./KeepRandomBeaconImplV1.sol');
const Upgrade = artifacts.require('./examples/KeepRandomBeaconUpgradeExample.sol');


contract('TestKeepRandomBeaconViaProxy', function(accounts) {

  let token, stakingProxy, stakingContract, implV1, implV2, proxy, implViaProxy, impl2ViaProxy,
    account_one = accounts[0],
    account_two = accounts[1];

  beforeEach(async () => {
    token = await KeepToken.new();
    stakingProxy = await StakingProxy.new();
    stakingContract = await TokenStaking.new(token.address, stakingProxy.address, duration.days(30));
    implV1 = await KeepRandomBeaconImplV1.new();
    proxy = await Proxy.new('v1', implV1.address);
    implViaProxy = await KeepRandomBeaconImplV1.at(proxy.address);
    await implViaProxy.initialize(stakingProxy.address, 100, 200);
  });

  it("should be able to check if the implementation contract was initialized", async function() {
    let result = await implViaProxy.initialized();
    assert.equal(result, true, "Implementation contract should be initialized.");
  });

  it("should be able to upgrade implementation and initialize it with new data", async function() {

    implV2 = await Upgrade.new();
    await proxy.upgradeTo('v2', implV2.address);
    
    impl2ViaProxy = await Upgrade.at(proxy.address);
    await impl2ViaProxy.initialize(stakingProxy.address, 100, 200);

    let result = await impl2ViaProxy.initialized();
    assert.equal(result, true, "Implementation contract should be initialized.");

    let newVar = await impl2ViaProxy.getNewVar();
    assert.equal(newVar, 1234, "Should be able to get new data from upgraded contract.");
  });

});
