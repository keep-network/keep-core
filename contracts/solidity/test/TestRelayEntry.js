import exceptThrow from './helpers/expectThrow';
import {bls} from './helpers/data';
const KeepRandomBeaconFrontendProxy = artifacts.require('./KeepRandomBeaconFrontendProxy.sol');
const KeepRandomBeaconFrontendImplV1 = artifacts.require('./KeepRandomBeaconFrontendImplV1.sol');
const KeepRandomBeaconBackendStub = artifacts.require('./KeepRandomBeaconBackendStub.sol');


contract('TestRelayEntry', function() {
  const relayRequestTimeout = 10;

  let frontendImplV1, frontendProxy, frontend, backend;

  beforeEach(async () => {

    // Initialize Keep Random Beacon contract
    frontendImplV1 = await KeepRandomBeaconFrontendImplV1.new();
    frontendProxy = await KeepRandomBeaconFrontendProxy.new(frontendImplV1.address);
    frontend = await KeepRandomBeaconFrontendImplV1.at(frontendProxy.address);

    backend = await KeepRandomBeaconBackendStub.new();
    backend.authorizeFrontendContract(frontend.address);
    await frontend.initialize(1, 1, backend.address, relayRequestTimeout);
    await frontend.requestRelayEntry(bls.seed, {value: 10});

  });

  it("should not be able to submit invalid relay entry", async function() {
    let requestID = 1;

    // Invalid signature
    let groupSignature = web3.utils.toBN('0x0fb34abfa2a9844a58776650e399bca3e08ab134e42595e03e3efc5a0472bcd8');

    await exceptThrow(backend.relayEntry(requestID, groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed));
  });

  it("should be able to submit valid relay entry", async function() {
    let requestID = 1;

    await backend.relayEntry(requestID, bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

    assert.equal((await frontend.getPastEvents())[0].args['requestResponse'].toString(),
      bls.groupSignature.toString(), "Should emit event with successfully submitted groupSignature."
    );

  });

});
