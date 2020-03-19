import {bls} from '../helpers/data';
import {duration, increaseTimeTo} from '../helpers/increaseTime';
import expectThrow from '../helpers/expectThrow';
import {initContracts} from '../helpers/initContracts';
import latestTime from "../helpers/latestTime";
import {createSnapshot, restoreSnapshot} from "../helpers/snapshot";
const ServiceContractProxy = artifacts.require('./KeepRandomBeaconService.sol');
const ServiceContractImplV2 = artifacts.require('./examples/KeepRandomBeaconServiceUpgradeExample.sol');
const {expectEvent, time} = require("@openzeppelin/test-helpers");

contract('KeepRandomBeaconService/Upgrade', function(accounts) {

  let operatorContract, serviceContractProxy, serviceContract, serviceContractImplV2, serviceContractV2,
    account_one = accounts[0],
    account_two = accounts[1];

  before(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      ServiceContractProxy,
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./stubs/KeepRandomBeaconOperatorStub.sol')
    );

    operatorContract = contracts.operatorContract;
    serviceContract = contracts.serviceContract;
    serviceContractProxy = await ServiceContractProxy.at(serviceContract.address);

    serviceContractImplV2 = await ServiceContractImplV2.new();
    serviceContractV2 = await ServiceContractImplV2.at(serviceContractProxy.address);

    // Using stub method to add first group to help testing.
    await operatorContract.registerNewGroup(bls.groupPubKey);
    operatorContract.setGroupSize(3);
    let group = await operatorContract.getGroupPublicKey(0);
    await operatorContract.setGroupMembers(group, [accounts[0], accounts[1], accounts[2]]);

    // Modify state so we can test later that eternal storage works as expected after upgrade
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate});
    await operatorContract.relayEntry(bls.groupSignature);
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

  it("should set first account as admin", async function() {
    assert.equal(
        await serviceContractProxy.admin(),
        account_one,
        "Account one should be set as admin"
    );
  });

  it("upgrade time delay should be set", async function() {
    assert.equal(
        (await serviceContractProxy.upgradeTimeDelay()).toNumber(),
        86400, // 1 day
        "Upgrade time delay should be one day"
    );
  });

  it("should be able to check if the implementation contract was initialized", async function() {
    assert.isTrue(
        await serviceContract.initialized(),
        "Implementation contract should be initialized."
    );
  });

  it("should fail to upgrade implementation if called by not contract owner", async function() {
    const initialize = serviceContractV2.contract.methods
        .initialize(
            100,
            duration.days(0),
            '0x0000000000000000000000000000000000000001'
        ).encodeABI();

    await expectThrow(serviceContractProxy.upgradeToAndCall(
        serviceContractImplV2.address,
        initialize,
        {from: account_two}
    ));
  });

  it("should be able to upgrade implementation and initialize it with new data", async function() {
    let previousEntryBefore = await serviceContractV2.previousEntry();
    const firstImplAddress = await serviceContractProxy.implementation();

    assert.notEqual(
        firstImplAddress,
        serviceContractImplV2.address,
        "Implementation should be other than V2 address at the beginning"
    );

    assert.equal(
        await serviceContractProxy.upgradeInitiatedTimestamp(),
        0,
        "Upgrade initiated timestamp should be 0 at the beginning"
    );

    const initialize = serviceContractV2.contract.methods
        .initialize(
            100,
            duration.days(0),
            '0x0000000000000000000000000000000000000001'
        ).encodeABI();

    let receipt = await serviceContractProxy.upgradeToAndCall(
        serviceContractImplV2.address,
        initialize
    );

    const upgradeStartedTime = await time.latest();

    expectEvent(receipt, "UpgradeStarted", {
      implementation: serviceContractImplV2.address,
      timestamp: upgradeStartedTime
    });

    assert.equal(
        await serviceContractProxy.implementation(),
        firstImplAddress,
        "Implementation should remain the same before upgrade is completed"
    );

    assert.equal(
        await serviceContractProxy.newImplementation(),
        serviceContractImplV2.address,
        "New implementation should be set to V2 address"
    );

    assert.equal(
        (await serviceContractProxy.upgradeInitiatedTimestamp()).toNumber(),
        upgradeStartedTime,
        "Upgrade initiated timestamp should be set correctly"
    );

    // Must wait upgrade time delay before complete upgrade.
    await expectThrow(serviceContractProxy.completeUpgrade());

    await increaseTimeTo(await latestTime()+duration.days(1));

    // Getting data from new contract shouldn't
    // be possible before upgrade is completed.
    await expectThrow(serviceContractV2.getNewVar());

    receipt = await serviceContractProxy.completeUpgrade();

    expectEvent(receipt, "UpgradeCompleted", {
      implementation: serviceContractImplV2.address,
    });

    assert.equal(
        await serviceContractProxy.implementation(),
        serviceContractImplV2.address,
        "Implementation should be changed to V2 address"
    );

    assert.equal(
        await serviceContractProxy.newImplementation(),
        serviceContractImplV2.address,
        "New implementation should remain set to V2 address"
    );

    assert.equal(
        await serviceContractProxy.upgradeInitiatedTimestamp(),
        0,
        "Upgrade initiated timestamp should be 0 at the end"
    );

    assert.isTrue(
        await serviceContractV2.initialized(),
        "Implementation contract should be initialized."
    );

    let newVar = await serviceContractV2.getNewVar();
    assert.equal(
        newVar,
        1234,
        "Should be able to get new data from upgraded contract."
    );

    let previousEntryAfter = await serviceContractV2.previousEntry()
    assert.equal(
        previousEntryBefore,
        previousEntryAfter,
        "Should keep previous storage after upgrade.")
    ;
  });
});
