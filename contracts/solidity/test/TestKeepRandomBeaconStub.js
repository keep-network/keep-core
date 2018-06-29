import { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
const Proxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconStub = artifacts.require('./KeepRandomBeaconStub.sol');

contract('TestKeepRandomBeaconStub', function(accounts) {

  let implV1, proxy, implViaProxy,
    account_one = accounts[0];

  beforeEach(async () => {
    implV1 = await KeepRandomBeaconStub.new();
    proxy = await Proxy.new('v1', implV1.address);
    implViaProxy = await KeepRandomBeaconStub.at(proxy.address);
    await implViaProxy.initialize();
  });

  it("should be able to request relay entry and get response", async function() {
    const relayEntryRequestedEvent = implViaProxy.RelayEntryRequested();
    const relayEntryGeneratedEvent = implViaProxy.RelayEntryGenerated();
    let previousRandomNumber;

    await implViaProxy.requestRelayEntry(10, 123456789, {from: account_one, value: 100});

    relayEntryRequestedEvent.get(function(error, result){
      assert.equal(result[0].event, 'RelayEntryRequested', "RelayEntryRequested event should occur on the implementation contract.");
    });

    relayEntryGeneratedEvent.get(function(error, result){
      previousRandomNumber = result[0].args['requestResponse'].toNumber();
      assert.equal(result[0].event, 'RelayEntryGenerated', "RelayEntryGenerated event should occur on the implementation contract.");
    });

    await increaseTimeTo(latestTime()+duration.seconds(1));
    await implViaProxy.requestRelayEntry(10, 123456789, {from: account_one, value: 100});

    relayEntryGeneratedEvent.get(function(error, result){
      assert.equal(result[0].args['previousEntry'].toNumber(), previousRandomNumber, "Previous entry should be present in the event.");
      assert.notEqual(result[0].args['requestResponse'].toNumber(), previousRandomNumber, "New number should be different from previous.");
    });

  });

});
