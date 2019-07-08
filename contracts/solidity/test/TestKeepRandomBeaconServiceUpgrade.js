import {bls} from './helpers/data';
import { duration } from './helpers/increaseTime';
import exceptThrow from './helpers/expectThrow';
import {initContracts} from './helpers/initContracts';
const ServiceContractProxy = artifacts.require('./KeepRandomBeaconService.sol');
const ServiceContractImplV2 = artifacts.require('./examples/KeepRandomBeaconServiceUpgradeExample.sol');


contract('TestKeepRandomBeaconServiceUpgrade', function(accounts) {

  let operatorContract, serviceContractProxy, serviceContract, serviceContractImplV2, serviceContractV2,
    account_two = accounts[1];

  before(async () => {
    let contracts = await initContracts(
      accounts,
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./StakingProxy.sol'),
      artifacts.require('./TokenStaking.sol'),
      ServiceContractProxy,
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./KeepRandomBeaconOperatorStub.sol')
    );

    operatorContract = contracts.operatorContract;
    serviceContract = contracts.serviceContract;
    serviceContractProxy = await ServiceContractProxy.at(serviceContract.address);

    serviceContractImplV2 = await ServiceContractImplV2.new();
    serviceContractV2 = await ServiceContractImplV2.at(serviceContractProxy.address);

    // Using stub method to add first group to help testing.
    await operatorContract.registerNewGroup(bls.groupPubKey);

    // Modify state so we can test later that eternal storage works as expected after upgrade
    await serviceContract.requestRelayEntry(bls.seed, {value: 10});
    await operatorContract.relayEntry(1, bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

  });

  it("should be able to check if the implementation contract was initialized", async function() {
    assert.isTrue(await serviceContract.initialized(), "Implementation contract should be initialized.");
  });

  it("should fail to upgrade implementation if called by not contract owner", async function() {
    await exceptThrow(serviceContractProxy.upgradeTo(serviceContractImplV2.address, {from: account_two}));
  });

  it("should be able to upgrade implementation and initialize it with new data", async function() {
    await serviceContractProxy.upgradeTo(serviceContractImplV2.address);
    await serviceContractV2.initialize(100, duration.days(0), operatorContract.address);

    assert.isTrue(await serviceContractV2.initialized(), "Implementation contract should be initialized.");

    let newVar = await serviceContractV2.getNewVar();
    assert.equal(newVar, 1234, "Should be able to get new data from upgraded contract.");

    assert.isTrue(bls.groupSignature.eq(await serviceContractV2.previousEntry()), "Should keep previous storage after upgrade.");
  });

});
