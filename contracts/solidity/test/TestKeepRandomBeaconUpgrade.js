import {bls} from './helpers/data';
import increaseTime, { duration } from './helpers/increaseTime';
import exceptThrow from './helpers/expectThrow';
const KeepRandomBeaconFrontendProxy = artifacts.require('./KeepRandomBeaconFrontendProxy.sol');
const KeepRandomBeaconFrontendImplV1 = artifacts.require('./KeepRandomBeaconFrontendImplV1.sol');
const Upgrade = artifacts.require('./examples/KeepRandomBeaconFrontendUpgradeExample.sol');
const KeepRandomBeaconBackend = artifacts.require('./KeepRandomBeaconBackendStub.sol');


contract('TestKeepRandomBeaconUpgrade', function(accounts) {

  let frontendImplV1, frontendImplV2, frontendProxy, frontendV1, frontendV2, backend,
    account_one = accounts[0],
    account_two = accounts[1];

  beforeEach(async () => {
    frontendImplV1 = await KeepRandomBeaconFrontendImplV1.new();
    frontendImplV2 = await Upgrade.new();
    frontendProxy = await KeepRandomBeaconFrontendProxy.new(frontendImplV1.address);
    frontendV1 = await KeepRandomBeaconFrontendImplV1.at(frontendProxy.address);
    backend = await KeepRandomBeaconBackend.new()
    await frontendV1.initialize(100, duration.days(0), bls.previousEntry, bls.groupPubKey, backend.address);

    // Add a few calls that modify state so we can test later that eternal storage works as expected after upgrade
    await frontendV1.requestRelayEntry(0, {from: account_two, value: 100});
    await frontendV1.requestRelayEntry(0, {from: account_two, value: 100});
    await frontendV1.requestRelayEntry(0, {from: account_two, value: 100});

  });

  it("should be able to check if the implementation contract was initialized", async function() {
    let result = await frontendV1.initialized();
    assert.equal(result, true, "Implementation contract should be initialized.");
  });

  it("should fail to upgrade implementation if called by not contract owner", async function() {
    await exceptThrow(frontendProxy.upgradeTo(frontendImplV2.address, {from: account_two}));
  });

  it("should be able to upgrade implementation and initialize it with new data", async function() {
    await frontendProxy.upgradeTo(frontendImplV2.address);
    
    frontendV2 = await Upgrade.at(frontendProxy.address);
    await frontendV2.initialize(100, duration.days(0), bls.previousEntry, bls.groupPubKey, backend.address);

    let result = await frontendV2.initialized();
    assert.equal(result, true, "Implementation contract should be initialized.");

    let newVar = await frontendV2.getNewVar();
    assert.equal(newVar, 1234, "Should be able to get new data from upgraded contract.");

    await frontendV2.requestRelayEntry(0, {from: account_two, value: 100})

    assert.equal((await frontendV2.getPastEvents())[0].args['requestID'], 6, "requestID should not be reset and should continue to increment where it was left in previous implementation.");

  });

});
