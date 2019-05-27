import {bls} from './helpers/data';
import { duration } from './helpers/increaseTime';
import exceptThrow from './helpers/expectThrow';
import {initContracts} from './helpers/initContracts';
const FrontendProxy = artifacts.require('./KeepRandomBeaconFrontendProxy.sol');
const FrontendImplV2 = artifacts.require('./examples/KeepRandomBeaconFrontendUpgradeExample.sol');


contract('TestKeepRandomBeaconUpgrade', function(accounts) {

  let backend, frontendProxy, frontend, frontendImplV2, frontendV2,
    account_two = accounts[1];

  before(async () => {
    let contracts = await initContracts(
      accounts,
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./StakingProxy.sol'),
      artifacts.require('./TokenStaking.sol'),
      FrontendProxy,
      artifacts.require('./KeepRandomBeaconFrontendImplV1.sol'),
      artifacts.require('./KeepRandomBeaconBackendStub.sol')
    );

    backend = contracts.backend;
    frontend = contracts.frontend;
    frontendProxy = await FrontendProxy.at(frontend.address);

    frontendImplV2 = await FrontendImplV2.new();
    frontendV2 = await FrontendImplV2.at(frontendProxy.address);

    // Using stub method to add first group to help testing.
    await backend.registerNewGroup(bls.groupPubKey);

    // Modify state so we can test later that eternal storage works as expected after upgrade
    await frontend.requestRelayEntry(bls.seed, {value: 10});
    await backend.relayEntry(2, bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

  });

  it("should be able to check if the implementation contract was initialized", async function() {
    assert.isTrue(await frontend.initialized(), "Implementation contract should be initialized.");
  });

  it("should fail to upgrade implementation if called by not contract owner", async function() {
    await exceptThrow(frontendProxy.upgradeTo(frontendImplV2.address, {from: account_two}));
  });

  it("should be able to upgrade implementation and initialize it with new data", async function() {
    await frontendProxy.upgradeTo(frontendImplV2.address);
    await frontendV2.initialize(100, duration.days(0), backend.address, 0);

    assert.isTrue(await frontendV2.initialized(), "Implementation contract should be initialized.");

    let newVar = await frontendV2.getNewVar();
    assert.equal(newVar, 1234, "Should be able to get new data from upgraded contract.");

    assert.isTrue(bls.groupSignature.eq(await frontendV2.previousEntry()), "Should keep previous storage after upgrade.");
  });

});
