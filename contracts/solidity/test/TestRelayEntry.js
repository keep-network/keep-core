import BigNumber from 'bignumber.js';
const KeepRandomBeaconProxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconImplV1 = artifacts.require('./KeepRandomBeaconImplV1.sol');
const KeepGroupStub = artifacts.require('./KeepGroupStub.sol');


contract('TestRelayEntry', function() {

  let keepRandomBeaconImplV1, keepRandomBeaconProxy, keepRandomBeaconImplViaProxy, relayEntryGeneratedEvent,
    keepGroupStub;

  beforeEach(async () => {

    // Initialize Keep Random Beacon contract
    keepRandomBeaconImplV1 = await KeepRandomBeaconImplV1.new();
    keepRandomBeaconProxy = await KeepRandomBeaconProxy.new(keepRandomBeaconImplV1.address);
    keepRandomBeaconImplViaProxy = await KeepRandomBeaconImplV1.at(keepRandomBeaconProxy.address);
    await keepRandomBeaconImplViaProxy.initialize(1,1);

    keepGroupStub = await KeepGroupStub.new();
    await keepRandomBeaconImplViaProxy.setGroupContract(keepGroupStub.address);

    relayEntryGeneratedEvent = keepRandomBeaconImplViaProxy.RelayEntryGenerated();
  });

  it("should be able to submit valid relay entry", async function() {
    let requestID = 1;

    // Data generated using client Go code with master secret key 123
    let groupPubKey = "0x1f1954b33144db2b5c90da089e8bde287ec7089d5d6433f3b6becaefdb678b1b2a9de38d14bef2cf9afc3c698a4211fa7ada7b4f036a2dfef0dc122b423259d0";
    let previousEntry = new BigNumber('0x884b130ed81751b63d0f5882483d4a24a7640bdf371f23b78dbeb520c84e3a85');
    let groupSignature = new BigNumber('0x1fb34abfa2a9844a58776650e399bca3e08ab134e42595e03e3efc5a0472bcd8');

    await keepRandomBeaconImplViaProxy.relayEntry(requestID, groupSignature, groupPubKey, previousEntry);

    relayEntryGeneratedEvent.get(function(_, result){
      assert.equal(result[0].args.requestResponse.equals(groupSignature), true, "Should emit event with successfully submitted groupSignature.");
    });

  });

});
