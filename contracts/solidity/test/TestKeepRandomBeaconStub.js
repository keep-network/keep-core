import { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
const Proxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconStub = artifacts.require('./KeepRandomBeaconStub.sol');

contract('TestKeepRandomBeaconStub', function(accounts) {

  let implV1, proxy, implViaProxy, blockReward, seed,
    account_one = accounts[0];

  beforeEach(async () => {
    implV1 = await KeepRandomBeaconStub.new();
    proxy = await Proxy.new(implV1.address);
    implViaProxy = await KeepRandomBeaconStub.at(proxy.address);
    await implViaProxy.initialize();
    blockReward = 10;
    seed = 123456789;
  });

  it("should be able to request relay entry and get response", async function() {
    await implViaProxy.requestRelayEntry(blockReward, seed, {from: account_one, value: 100});

    assert.equal((await implViaProxy.getPastEvents())[0].event, 'RelayEntryRequested', "RelayEntryRequested event should occur on the implementation contract.");
    assert.equal((await implViaProxy.getPastEvents())[1].event, 'RelayEntryGenerated', "RelayEntryGenerated event should occur on the implementation contract.");

    let previousRandomNumber = (await implViaProxy.getPastEvents())[1].args['requestResponse'].toString();  
    await increaseTimeTo(await latestTime()+duration.seconds(1));
    await implViaProxy.requestRelayEntry(blockReward, seed, {from: account_one, value: 100});

    assert.equal((await implViaProxy.getPastEvents())[1].args['previousEntry'].toString(), previousRandomNumber, "Previous entry should be present in the event.");
    assert.notEqual((await implViaProxy.getPastEvents())[1].args['requestResponse'].toString(), previousRandomNumber, "New number should be different from previous.");

  });

});
