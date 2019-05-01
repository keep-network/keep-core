import exceptThrow from './helpers/expectThrow';
import {bls} from './helpers/data';
const KeepRandomBeaconFrontendProxy = artifacts.require('./KeepRandomBeaconFrontendProxy.sol');
const KeepRandomBeaconFrontendImplV1 = artifacts.require('./KeepRandomBeaconFrontendImplV1.sol');
const KeepRandomBeaconBackendStub = artifacts.require('./KeepRandomBeaconBackendStub.sol');


contract('TestRelayEntry', function() {

  let keepRandomBeaconFrontendImplV1, keepRandomBeaconFrontendProxy, keepRandomBeaconFrontendImplViaProxy, relayEntryGeneratedEvent,
    keepRandomBeaconBackendStub;

  beforeEach(async () => {

    // Initialize Keep Random Beacon contract
    keepRandomBeaconFrontendImplV1 = await KeepRandomBeaconFrontendImplV1.new();
    keepRandomBeaconFrontendProxy = await KeepRandomBeaconFrontendProxy.new(keepRandomBeaconFrontendImplV1.address);
    keepRandomBeaconFrontendImplViaProxy = await KeepRandomBeaconFrontendImplV1.at(keepRandomBeaconFrontendProxy.address);

    keepRandomBeaconBackendStub = await KeepRandomBeaconBackendStub.new();
    await keepRandomBeaconFrontendImplViaProxy.initialize(1,1, bls.previousEntry, bls.groupPubKey, keepRandomBeaconBackendStub.address);
    await keepRandomBeaconFrontendImplViaProxy.requestRelayEntry(bls.seed, {value: 10});

    relayEntryGeneratedEvent = keepRandomBeaconFrontendImplViaProxy.RelayEntryGenerated();
  });

  it("should not be able to submit invalid relay entry", async function() {
    let requestID = 1;

    // Invalid signature
    let groupSignature = web3.utils.toBN('0x0fb34abfa2a9844a58776650e399bca3e08ab134e42595e03e3efc5a0472bcd8');

    await exceptThrow(keepRandomBeaconFrontendImplViaProxy.relayEntry(requestID, groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed));
  });

  it("should be able to submit valid relay entry", async function() {
    let requestID = 1;

    await keepRandomBeaconFrontendImplViaProxy.relayEntry(requestID, bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

    assert.equal((await keepRandomBeaconFrontendImplViaProxy.getPastEvents())[0].args['requestResponse'].toString(),
      bls.groupSignature.toString(), "Should emit event with successfully submitted groupSignature."
    );

  });

});
