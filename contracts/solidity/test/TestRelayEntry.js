import exceptThrow from './helpers/expectThrow';
import {bls} from './helpers/data';
const KeepRandomBeaconProxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconImplV1 = artifacts.require('./KeepRandomBeaconImplV1.sol');
const KeepGroupStub = artifacts.require('./KeepGroupStub.sol');


contract('TestRelayEntry', function() {
  const relayRequestTimeout = 10;

  let keepRandomBeaconImplV1, keepRandomBeaconProxy, keepRandomBeaconImplViaProxy, keepGroupStub;

  beforeEach(async () => {

    // Initialize Keep Random Beacon contract
    keepRandomBeaconImplV1 = await KeepRandomBeaconImplV1.new();
    keepRandomBeaconProxy = await KeepRandomBeaconProxy.new(keepRandomBeaconImplV1.address);
    keepRandomBeaconImplViaProxy = await KeepRandomBeaconImplV1.at(keepRandomBeaconProxy.address);

    keepGroupStub = await KeepGroupStub.new();
    await keepRandomBeaconImplViaProxy.initialize(1,1, bls.previousEntry, bls.groupPubKey, keepGroupStub.address,
      relayRequestTimeout);
    await keepRandomBeaconImplViaProxy.requestRelayEntry(bls.seed, {value: 10});
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
