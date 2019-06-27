import {bls} from './helpers/data';
import { duration } from './helpers/increaseTime';
import expectThrow from './helpers/expectThrow';
import mineBlocks from './helpers/mineBlocks';
const Proxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconImplV1 = artifacts.require('./KeepRandomBeaconImplV1.sol');
const Upgrade = artifacts.require('./examples/KeepRandomBeaconUpgradeExample.sol');
const KeepGroup = artifacts.require('./KeepGroupStub.sol');


contract('TestKeepRandomBeaconUpgrade', function(accounts) {
  const relayRequestTimeout = 10;
  const blocksForward = 20;

  let implV1, implV2, proxy, implViaProxy, impl2ViaProxy, keepGroup,
    account_two = accounts[1];

  beforeEach(async () => {
    implV1 = await KeepRandomBeaconImplV1.new();
    implV2 = await Upgrade.new();
    proxy = await Proxy.new(implV1.address);
    implViaProxy = await KeepRandomBeaconImplV1.at(proxy.address);
    keepGroup = await KeepGroup.new();
    await implViaProxy.initialize(100, duration.days(0), bls.previousEntry, bls.groupPubKey, keepGroup.address, 
      relayRequestTimeout);

    // Add a few calls that modify state so we can test later that eternal storage works as expected after upgrade
    await implViaProxy.requestRelayEntry(0, {from: account_two, value: 100});
    await implViaProxy.relayEntry(1, bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

    await mineBlocks(blocksForward)
    await implViaProxy.requestRelayEntry(0, {from: account_two, value: 100});
    await implViaProxy.relayEntry(2, bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

    await mineBlocks(blocksForward)
    await implViaProxy.requestRelayEntry(0, {from: account_two, value: 100});
    await implViaProxy.relayEntry(3, bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

  });

  it("should be able to check if the implementation contract was initialized", async function() {
    let result = await implViaProxy.initialized();
    assert.equal(result, true, "Implementation contract should be initialized.");
  });

  it("should fail to upgrade implementation if called by not contract owner", async function() {
    await expectThrow(proxy.upgradeTo(implV2.address, {from: account_two}));
  });

  it("should be able to upgrade implementation and initialize it with new data", async function() {
    await proxy.upgradeTo(implV2.address);
    
    impl2ViaProxy = await Upgrade.at(proxy.address);
    await impl2ViaProxy.initialize(100, duration.days(0), bls.previousEntry, bls.groupPubKey, keepGroup.address,
      relayRequestTimeout);

    let result = await impl2ViaProxy.initialized();
    assert.equal(result, true, "Implementation contract should be initialized.");

    let newVar = await impl2ViaProxy.getNewVar();
    assert.equal(newVar, 1234, "Should be able to get new data from upgraded contract.");

    await mineBlocks(blocksForward)
    await impl2ViaProxy.requestRelayEntry(0, {from: account_two, value: 100})

    assert.equal((await impl2ViaProxy.getPastEvents())[0].args['requestID'], 6, "requestID should not be reset and should continue to increment where it was left in previous implementation.");

  });

});
