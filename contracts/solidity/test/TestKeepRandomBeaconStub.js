import { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
const KeepRandomBeaconService = artifacts.require('./KeepRandomBeaconService.sol');
const KeepRandomBeaconServiceStub = artifacts.require('./KeepRandomBeaconServiceStub.sol');

contract('TestKeepRandomBeaconServiceStub', function(accounts) {

  let serviceContractImplV1, serviceContractProxy, serviceContract, seed,
    account_one = accounts[0];

  before(async () => {
    serviceContractImplV1 = await KeepRandomBeaconServiceStub.new();
    serviceContractProxy = await KeepRandomBeaconService.new(serviceContractImplV1.address);
    serviceContract = await KeepRandomBeaconServiceStub.at(serviceContractProxy.address);
    await serviceContract.initialize();
    seed = 123456789;
  });

  it("should be able to request relay entry and get response", async function() {
    await serviceContract.requestRelayEntry(seed, {from: account_one, value: 100});

    assert.equal((await serviceContract.getPastEvents())[0].event, 'RelayEntryRequested', "RelayEntryRequested event should occur on the implementation contract.");
    assert.equal((await serviceContract.getPastEvents())[1].event, 'RelayEntryGenerated', "RelayEntryGenerated event should occur on the implementation contract.");

    let previousRandomNumber = (await serviceContract.getPastEvents())[1].args['requestResponse'].toString();  
    await increaseTimeTo(await latestTime()+duration.seconds(1));
    await serviceContract.requestRelayEntry(seed, {from: account_one, value: 100});

    assert.equal((await serviceContract.getPastEvents())[1].args['previousEntry'].toString(), previousRandomNumber, "Previous entry should be present in the event.");
    assert.notEqual((await serviceContract.getPastEvents())[1].args['requestResponse'].toString(), previousRandomNumber, "New number should be different from previous.");

  });

});
