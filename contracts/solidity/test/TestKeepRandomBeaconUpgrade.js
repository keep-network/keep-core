import {bls} from './helpers/data';
import increaseTime, { duration } from './helpers/increaseTime';
import exceptThrow from './helpers/expectThrow';
const KeepRandomBeaconFrontendProxy = artifacts.require('./KeepRandomBeaconFrontendProxy.sol');
const KeepRandomBeaconFrontendImplV1 = artifacts.require('./KeepRandomBeaconFrontendImplV1.sol');
const KeepRandomBeaconFrontendImplV2 = artifacts.require('./examples/KeepRandomBeaconFrontendUpgradeExample.sol');
const KeepRandomBeaconBackend = artifacts.require('./KeepRandomBeaconBackendStub.sol');


contract('TestKeepRandomBeaconUpgrade', function(accounts) {
  const relayRequestTimeout = 10;

  let frontendImplV1, frontendImplV2, frontendProxy, frontendV1, frontendV2, backend,
    account_one = accounts[0],
    account_two = accounts[1];

  beforeEach(async () => {
    frontendImplV1 = await KeepRandomBeaconFrontendImplV1.new();
    frontendImplV2 = await KeepRandomBeaconFrontendImplV2.new();
    frontendProxy = await KeepRandomBeaconFrontendProxy.new(frontendImplV1.address);
    frontendV1 = await KeepRandomBeaconFrontendImplV1.at(frontendProxy.address);
    frontendV2 = await KeepRandomBeaconFrontendImplV2.at(frontendProxy.address);
    backend = await KeepRandomBeaconBackend.new()
    await backend.authorizeFrontendContract(frontendProxy.address);
    await frontendV1.initialize(100, duration.days(0), backend.address, relayRequestTimeout);

    // Modify state so we can test later that eternal storage works as expected after upgrade
    await backend.relayEntry(1, bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

  });

  it("should be able to check if the implementation contract was initialized", async function() {
    assert.isTrue(await frontendV1.initialized(), "Implementation contract should be initialized.");
  });

  it("should fail to upgrade implementation if called by not contract owner", async function() {
    await exceptThrow(frontendProxy.upgradeTo(frontendImplV2.address, {from: account_two}));
  });

  it("should be able to upgrade implementation and initialize it with new data", async function() {
    await frontendProxy.upgradeTo(frontendImplV2.address);
    await frontendV2.initialize(100, duration.days(0), backend.address);

    assert.isTrue(await frontendV2.initialized(), "Implementation contract should be initialized.");

    let newVar = await frontendV2.getNewVar();
    assert.equal(newVar, 1234, "Should be able to get new data from upgraded contract.");

    assert.isTrue(bls.groupSignature.eq(await frontendV2.previousEntry()), "Should keep previous storage after upgrade.");
  });

});
