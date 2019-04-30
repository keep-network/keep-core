import exceptThrow from './helpers/expectThrow';
import {bls} from './helpers/data';
const KeepRandomBeaconProxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconImplV1 = artifacts.require('./KeepRandomBeaconImplV1.sol');
const KeepRandomBeaconBackendStub = artifacts.require('./KeepRandomBeaconBackendStub.sol');


contract('TestRelayEntry', function() {

  let keepRandomBeaconImplV1, keepRandomBeaconProxy, keepRandomBeaconImplViaProxy, relayEntryGeneratedEvent,
    keepRandomBeaconBackendStub;

  beforeEach(async () => {

    // Initialize Keep Random Beacon contract
    keepRandomBeaconImplV1 = await KeepRandomBeaconImplV1.new();
    keepRandomBeaconProxy = await KeepRandomBeaconProxy.new(keepRandomBeaconImplV1.address);
    keepRandomBeaconImplViaProxy = await KeepRandomBeaconImplV1.at(keepRandomBeaconProxy.address);

    keepRandomBeaconBackendStub = await KeepRandomBeaconBackendStub.new();
    await keepRandomBeaconImplViaProxy.initialize(1,1, bls.previousEntry, bls.groupPubKey, keepRandomBeaconBackendStub.address);
    await keepRandomBeaconImplViaProxy.requestRelayEntry(bls.seed, {value: 10});

    relayEntryGeneratedEvent = keepRandomBeaconImplViaProxy.RelayEntryGenerated();
  });

  it("should not be able to submit invalid relay entry", async function() {
    let requestID = 1;

    // Invalid signature
    let groupSignature = web3.utils.toBN('0x0fb34abfa2a9844a58776650e399bca3e08ab134e42595e03e3efc5a0472bcd8');

    await exceptThrow(keepRandomBeaconImplViaProxy.relayEntry(requestID, groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed));
  });

  it("should be able to submit valid relay entry", async function() {
    let requestID = 1;

    await keepRandomBeaconImplViaProxy.relayEntry(requestID, bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

    assert.equal((await keepRandomBeaconImplViaProxy.getPastEvents())[0].args['requestResponse'].toString(),
      bls.groupSignature.toString(), "Should emit event with successfully submitted groupSignature."
    );

  });

});
