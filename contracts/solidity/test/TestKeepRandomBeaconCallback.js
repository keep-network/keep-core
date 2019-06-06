import {bls} from './helpers/data';
const Proxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeacon = artifacts.require('./KeepRandomBeaconStub.sol');
const CallbackContract = artifacts.require('./examples/CallbackContract.sol');

contract('TestKeepRandomBeaconCallback', function() {

  let impl, proxy, keepRandomBeacon, callbackContract;

  before(async () => {
    impl = await KeepRandomBeacon.new();
    proxy = await Proxy.new(impl.address);
    keepRandomBeacon = await KeepRandomBeacon.at(proxy.address);
    await keepRandomBeacon.initialize();
    callbackContract = await CallbackContract.new();
  });

  it("should produce entry if callback contract was not provided", async function() {

    await keepRandomBeacon.requestRelayEntry(0, "0x0000000000000000000000000000000000000000", "");
    await keepRandomBeacon.relayEntry(0, bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

    let result = await callbackContract.lastEntry();
    assert.isFalse(result.eq(bls.groupSignature), "Value should not change on the callback contract.");
  });

  it("should successfully call method on a callback contract", async function() {

    await keepRandomBeacon.requestRelayEntry(0, callbackContract.address, "callback(uint256)");
    await keepRandomBeacon.relayEntry(1, bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

    let result = await callbackContract.lastEntry();
    assert.isTrue(result.eq(bls.groupSignature), "Value updated by the callback should be the same as relay entry.");
  });

});
