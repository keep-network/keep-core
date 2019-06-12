import {bls} from './helpers/data';
const Proxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeacon = artifacts.require('./KeepRandomBeaconImplV1.sol');
const KeepGroupStub = artifacts.require('./KeepGroupStub.sol');
const CallbackContract = artifacts.require('./examples/CallbackContract.sol');

contract('TestKeepRandomBeaconCallback', function() {
  const relayRequestTimeout = 10;
  let impl, proxy, keepRandomBeacon, callbackContract, keepGroupStub;

  before(async () => {
    impl = await KeepRandomBeacon.new();
    proxy = await Proxy.new(impl.address);
    keepRandomBeacon = await KeepRandomBeacon.at(proxy.address);

    keepGroupStub = await KeepGroupStub.new();
    await keepRandomBeacon.initialize(1, 1, bls.previousEntry, bls.groupPubKey, keepGroupStub.address,
      relayRequestTimeout);
    callbackContract = await CallbackContract.new();
  });

  it("should produce entry if callback contract was not provided", async function() {
    let tx = await keepRandomBeacon.requestRelayEntry(bls.seed, "0x0000000000000000000000000000000000000000", "", {value: 10});
    let requestId = tx.logs[0].args.requestID;
    await keepRandomBeacon.relayEntry(requestId, bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

    let result = await callbackContract.lastEntry();
    assert.isFalse(result.eq(bls.groupSignature), "Value should not change on the callback contract.");
  });

  it("should successfully call method on a callback contract", async function() {
    let tx = await keepRandomBeacon.requestRelayEntry(bls.seed, callbackContract.address, "callback(uint256)", {value: 10});
    let requestId = tx.logs[0].args.requestID;
    await keepRandomBeacon.relayEntry(requestId, bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

    let result = await callbackContract.lastEntry();
    assert.isTrue(result.eq(bls.groupSignature), "Value updated by the callback should be the same as relay entry.");
  });

});
