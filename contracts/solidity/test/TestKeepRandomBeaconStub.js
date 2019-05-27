import { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
const KeepRandomBeaconFrontend = artifacts.require('./KeepRandomBeaconFrontend.sol');
const KeepRandomBeaconFrontendStub = artifacts.require('./KeepRandomBeaconFrontendStub.sol');

contract('TestKeepRandomBeaconFrontendStub', function(accounts) {

  let frontendImplV1, frontendProxy, frontend, seed,
    account_one = accounts[0];

  before(async () => {
    frontendImplV1 = await KeepRandomBeaconFrontendStub.new();
    frontendProxy = await KeepRandomBeaconFrontend.new(frontendImplV1.address);
    frontend = await KeepRandomBeaconFrontendStub.at(frontendProxy.address);
    await frontend.initialize();
    seed = 123456789;
  });

  it("should be able to request relay entry and get response", async function() {
    await frontend.requestRelayEntry(seed, {from: account_one, value: 100});

    assert.equal((await frontend.getPastEvents())[0].event, 'RelayEntryRequested', "RelayEntryRequested event should occur on the implementation contract.");
    assert.equal((await frontend.getPastEvents())[1].event, 'RelayEntryGenerated', "RelayEntryGenerated event should occur on the implementation contract.");

    let previousRandomNumber = (await frontend.getPastEvents())[1].args['requestResponse'].toString();  
    await increaseTimeTo(await latestTime()+duration.seconds(1));
    await frontend.requestRelayEntry(seed, {from: account_one, value: 100});

    assert.equal((await frontend.getPastEvents())[1].args['previousEntry'].toString(), previousRandomNumber, "Previous entry should be present in the event.");
    assert.notEqual((await frontend.getPastEvents())[1].args['requestResponse'].toString(), previousRandomNumber, "New number should be different from previous.");

  });

});
