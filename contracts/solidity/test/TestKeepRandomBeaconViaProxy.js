import increaseTime, { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import exceptThrow from './helpers/expectThrow';
import encodeCall from './helpers/encodeCall';
const Proxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconImplV1 = artifacts.require('./KeepRandomBeaconImplV1.sol');

contract('TestKeepRandomBeaconViaProxy', function(accounts) {

  let implV1, proxy, implViaProxy,
    account_one = accounts[0],
    account_two = accounts[1];

  beforeEach(async () => {
    implV1 = await KeepRandomBeaconImplV1.new();
    proxy = await Proxy.new(implV1.address);
    implViaProxy = await KeepRandomBeaconImplV1.at(proxy.address);
    await implViaProxy.initialize(100, duration.days(30));
  });

  it("should be able to check if the implementation contract was initialized", async function() {
    let result = await implViaProxy.initialized();
    assert.equal(result, true, "Implementation contract should be initialized.");
  });

  it("should fail to request relay entry with not enough ether", async function() {
    await exceptThrow(implViaProxy.requestRelayEntry(0, 0, {from: account_two, value: 99}));
  });

  it("should be able to request relay entry via implementation contract with enough ether", async function() {
    const relayEntryRequestedEvent = implViaProxy.RelayEntryRequested();
    await implViaProxy.requestRelayEntry(0, 0, {from: account_two, value: 100})

    relayEntryRequestedEvent.get(function(error, result){
      assert.equal(result[0].event, 'RelayEntryRequested', "RelayEntryRequested event should occur on the implementation contract.");
    });

    let contractBalance = await web3.eth.getBalance(implViaProxy.address).toNumber();
    assert.equal(contractBalance, 100, "Keep Random Beacon contract should receive ether.");

    let contractBalanceViaProxy = await web3.eth.getBalance(proxy.address).toNumber();
    assert.equal(contractBalanceViaProxy, 100, "Keep Random Beacon contract new balance should be visible via proxy.");

  });

  it("should be able to request relay entry via proxy contract with enough ether", async function() {
    const relayEntryRequestedEvent = implViaProxy.RelayEntryRequested();

    await exceptThrow(proxy.sendTransaction({from: account_two, value: 1000}));

    await web3.eth.sendTransaction({
      from: account_two, value: 100, gas: 200000, to: proxy.address,
      data: encodeCall('requestRelayEntry', ['uint256', 'uint256'], [0,0])
    });

    relayEntryRequestedEvent.get(function(error, result){
      assert.equal(result[0].event, 'RelayEntryRequested', "RelayEntryRequested event should occur on the proxy contract.");
    });

    let contractBalance = await web3.eth.getBalance(implViaProxy.address).toNumber();
    assert.equal(contractBalance, 100, "Keep Random Beacon contract should receive ether.");

    let contractBalanceViaProxy = await web3.eth.getBalance(proxy.address).toNumber();
    assert.equal(contractBalanceViaProxy, 100, "Keep Random Beacon contract new balance should be visible via proxy.");
  });

  it("owner should be able to withdraw ether from random beacon contract", async function() {

    await web3.eth.sendTransaction({
      from: account_two, value: web3.toWei(1, 'ether'), gas: 200000, to: proxy.address,
      data: encodeCall('requestRelayEntry', ['uint256', 'uint256'], [0,0])
    });

    let ownerStartBalance = web3.fromWei(await web3.eth.getBalance(account_one).toNumber(), 'ether');

    // should fail to withdraw if not owner
    await exceptThrow(implViaProxy.initiateWithdrawal({from: account_two}));
    await exceptThrow(implViaProxy.finishWithdrawal({from: account_two}));

    await implViaProxy.initiateWithdrawal({from: account_one});
    await exceptThrow(implViaProxy.finishWithdrawal({from: account_one}));

    // jump in time, full withdrawal delay
    await increaseTimeTo(latestTime()+duration.days(30));
    await implViaProxy.finishWithdrawal({from: account_one});

    let contractBalance = await web3.eth.getBalance(implViaProxy.address).toNumber();
    assert.equal(contractBalance, 0, "Keep Random Beacon contract should send all ether.");
    let contractBalanceViaProxy = await web3.eth.getBalance(proxy.address).toNumber();
    assert.equal(contractBalanceViaProxy, 0, "Keep Random Beacon contract updated balance should be visible via proxy.");

    let ownerEndBalance = web3.fromWei(await web3.eth.getBalance(account_one).toNumber(), 'ether');
    assert(ownerEndBalance > ownerStartBalance, "Owner updated balance should include received ether.");
  });

  it("should fail to update minimum payment by non owner", async function() {
    await exceptThrow(implViaProxy.setMinimumPayment(123, {from: account_two}));
  });

  it("should be able to update minimum payment by the owner", async function() {
    await implViaProxy.setMinimumPayment(123);
    let newMinPayment = await implViaProxy.minimumPayment();
    assert.equal(newMinPayment, 123, "Should be able to get updated minimum payment.");
  });
});
