import increaseTime, { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import exceptThrow from './helpers/expectThrow';
import encodeCall from './helpers/encodeCall';
const KeepToken = artifacts.require('./KeepToken.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const Proxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconImplV1 = artifacts.require('./KeepRandomBeaconImplV1.sol');

contract('TestKeepRandomBeaconViaProxy', function(accounts) {

  let token, stakingProxy, stakingContract, implV1, proxy, implViaProxy,
    account_one = accounts[0],
    account_two = accounts[1];

  beforeEach(async () => {
    token = await KeepToken.new();
    stakingProxy = await StakingProxy.new();
    stakingContract = await TokenStaking.new(token.address, stakingProxy.address, duration.days(30));
    implV1 = await KeepRandomBeaconImplV1.new();
    proxy = await Proxy.new('v1', implV1.address);
    implViaProxy = await KeepRandomBeaconImplV1.at(proxy.address);
    await implViaProxy.initialize(stakingProxy.address, 100, 200);
  });

  it("should be able to check if the implementation contract was initialized", async function() {
    let result = await implViaProxy.initialized();
    assert.equal(result, true, "Implementation contract should be initialized.");
  });

  it("should fail to request relay entry with not enough ether", async function() {
    await exceptThrow(implViaProxy.requestRelay(0, 0, {from: account_two, value: 99}));
  });

  it("should be able to request relay entry via implementation contract with enough ether", async function() {
    const relayEntryRequestedEvent = implViaProxy.RelayEntryRequested();
    await implViaProxy.requestRelay(0, 0, {from: account_two, value: 100})

    relayEntryRequestedEvent.get(function(error, result){
      assert.equal(result[0].event, 'RelayEntryRequested', "RelayEntryRequested event should occur on the implementation contract.");
    });
  });

  it("should be able to request relay entry via proxy contract with enough ether", async function() {
    const relayEntryRequestedEvent = proxy.RelayEntryRequested();

    await web3.eth.sendTransaction({
      from: account_two, value: 100, gas: 200000, to: proxy.address,
      data: encodeCall('requestRelay', ['uint256', 'uint256'], [0,0])
    });

    relayEntryRequestedEvent.get(function(error, result){
      assert.equal(result[0].event, 'RelayEntryRequested', "RelayEntryRequested event should occur on the proxy contract.");
    });
  });

  it("should fail to update minimum stake and minimum payments by non owner", async function() {
    await exceptThrow(implViaProxy.setMinimumPayment(123, {from: account_two}));
    await exceptThrow(implViaProxy.setMinimumStake(123, {from: account_two}));
  });

  it("should be able to update minimum stake and minimum payments by the owner", async function() {
    await implViaProxy.setMinimumPayment(123);
    let newMinPayment = await implViaProxy.minimumPayment();
    assert.equal(newMinPayment, 123, "Should be able to get updated minimum payment.");

    await implViaProxy.setMinimumStake(123);
    let newMinStake = await implViaProxy.minimumStake();
    assert.equal(newMinStake, 123, "Should be able to get updated minimum stake.");
  });
});
