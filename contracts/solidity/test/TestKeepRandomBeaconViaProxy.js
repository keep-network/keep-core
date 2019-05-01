import {bls} from './helpers/data';
import increaseTime, { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import exceptThrow from './helpers/expectThrow';
import encodeCall from './helpers/encodeCall';
const Proxy = artifacts.require('./KeepRandomBeaconFrontendProxy.sol');
const KeepRandomBeaconFrontendImplV1 = artifacts.require('./KeepRandomBeaconFrontendImplV1.sol');
const KeepRandomBeaconBackend = artifacts.require('./KeepRandomBeaconBackendStub.sol');

contract('TestKeepRandomBeaconViaProxy', function(accounts) {

  let implV1, proxy, implViaProxy, keepRandomBeaconBackend,
    account_one = accounts[0],
    account_two = accounts[1],
    account_three = accounts[2];

  beforeEach(async () => {
    implV1 = await KeepRandomBeaconFrontendImplV1.new();
    proxy = await Proxy.new(implV1.address);
    implViaProxy = await KeepRandomBeaconFrontendImplV1.at(proxy.address);
    keepRandomBeaconBackend = await KeepRandomBeaconBackend.new()
    await implViaProxy.initialize(100, duration.days(30), bls.previousEntry, bls.groupPubKey, keepRandomBeaconBackend.address);
  });

  it("should be able to check if the implementation contract was initialized", async function() {
    let result = await implViaProxy.initialized();
    assert.equal(result, true, "Implementation contract should be initialized.");
  });

  it("should fail to request relay entry with not enough ether", async function() {
    await exceptThrow(implViaProxy.requestRelayEntry(0, {from: account_two, value: 99}));
  });

  it("should be able to request relay entry via implementation contract with enough ether", async function() {
    await implViaProxy.requestRelayEntry(0, {from: account_two, value: 100})

    assert.equal((await implViaProxy.getPastEvents())[0].event, 'RelayEntryRequested', "RelayEntryRequested event should occur on the implementation contract.");

    let contractBalance = await web3.eth.getBalance(implViaProxy.address);
    assert.equal(contractBalance, 100, "Keep Random Beacon contract should receive ether.");

    let contractBalanceViaProxy = await web3.eth.getBalance(proxy.address);
    assert.equal(contractBalanceViaProxy, 100, "Keep Random Beacon contract new balance should be visible via proxy.");

  });

  it("should be able to request relay entry via proxy contract with enough ether", async function() {
    await exceptThrow(proxy.sendTransaction({from: account_two, value: 1000}));

    await web3.eth.sendTransaction({
      from: account_two, value: 100, gas: 200000, to: proxy.address,
      data: encodeCall('requestRelayEntry', ['uint256'], [0])
    });

    assert.equal((await implViaProxy.getPastEvents())[0].event, 'RelayEntryRequested', "RelayEntryRequested event should occur on the proxy contract.");

    let contractBalance = await web3.eth.getBalance(implViaProxy.address);
    assert.equal(contractBalance, 100, "Keep Random Beacon contract should receive ether.");

    let contractBalanceViaProxy = await web3.eth.getBalance(proxy.address);
    assert.equal(contractBalanceViaProxy, 100, "Keep Random Beacon contract new balance should be visible via proxy.");
  });

  it("owner should be able to withdraw ether from random beacon contract", async function() {

    let amount = web3.utils.toWei('1', 'ether');
    await web3.eth.sendTransaction({
      from: account_two, value: amount, gas: 200000, to: proxy.address,
      data: encodeCall('requestRelayEntry', ['uint256'], [0])
    });

    // should fail to withdraw if not owner
    await exceptThrow(implViaProxy.initiateWithdrawal({from: account_two}));
    await exceptThrow(implViaProxy.finishWithdrawal(account_two, {from: account_two}));

    await implViaProxy.initiateWithdrawal({from: account_one});
    await exceptThrow(implViaProxy.finishWithdrawal(account_three, {from: account_one}));

    let contractStartBalance = await web3.eth.getBalance(implViaProxy.address);
    // jump in time, full withdrawal delay
    await increaseTimeTo(await latestTime()+duration.days(30));

    let receiverStartBalance = web3.utils.fromWei(await web3.eth.getBalance(account_three), 'ether');
    await implViaProxy.finishWithdrawal(account_three, {from: account_one});
    let receiverEndBalance = web3.utils.fromWei(await web3.eth.getBalance(account_three), 'ether');
    assert(receiverEndBalance > receiverStartBalance, "Receiver updated balance should include received ether.");

    let contractEndBalance = await web3.eth.getBalance(implViaProxy.address);
    assert.equal(contractEndBalance, contractStartBalance - amount, "Keep Random Beacon contract should send all ether.");
    let contractEndBalanceViaProxy = await web3.eth.getBalance(proxy.address);
    assert.equal(contractEndBalanceViaProxy, contractStartBalance - amount, "Keep Random Beacon contract updated balance should be visible via proxy.");

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
