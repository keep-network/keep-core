import {bls} from './helpers/data';
import increaseTime, { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import exceptThrow from './helpers/expectThrow';
import encodeCall from './helpers/encodeCall';
const KeepRandomBeaconFrontendProxy = artifacts.require('./KeepRandomBeaconFrontendProxy.sol');
const KeepRandomBeaconFrontendImplV1 = artifacts.require('./KeepRandomBeaconFrontendImplV1.sol');
const KeepRandomBeaconBackend = artifacts.require('./KeepRandomBeaconBackendStub.sol');

contract('TestKeepRandomBeaconViaProxy', function(accounts) {

  let frontendImplV1, frontendProxy, frontend, backend,
    account_one = accounts[0],
    account_two = accounts[1],
    account_three = accounts[2];

  beforeEach(async () => {
    frontendImplV1 = await KeepRandomBeaconFrontendImplV1.new();
    frontendProxy = await KeepRandomBeaconFrontendProxy.new(frontendImplV1.address);
    frontend = await KeepRandomBeaconFrontendImplV1.at(frontendProxy.address);
    backend = await KeepRandomBeaconBackend.new()
    await frontend.initialize(100, duration.days(30), bls.previousEntry, bls.groupPubKey, backend.address);
  });

  it("should be able to check if the implementation contract was initialized", async function() {
    let result = await frontend.initialized();
    assert.equal(result, true, "Implementation contract should be initialized.");
  });

  it("should fail to request relay entry with not enough ether", async function() {
    await exceptThrow(frontend.requestRelayEntry(0, {from: account_two, value: 99}));
  });

  it("should be able to request relay entry via implementation contract with enough ether", async function() {
    await frontend.requestRelayEntry(0, {from: account_two, value: 100})

    assert.equal((await frontend.getPastEvents())[0].event, 'RelayEntryRequested', "RelayEntryRequested event should occur on the implementation contract.");

    let contractBalance = await web3.eth.getBalance(frontend.address);
    assert.equal(contractBalance, 100, "Keep Random Beacon contract should receive ether.");

    let contractBalanceViaProxy = await web3.eth.getBalance(frontendProxy.address);
    assert.equal(contractBalanceViaProxy, 100, "Keep Random Beacon contract new balance should be visible via frontendProxy.");

  });

  it("should be able to request relay entry via frontendProxy contract with enough ether", async function() {
    await exceptThrow(frontendProxy.sendTransaction({from: account_two, value: 1000}));

    await web3.eth.sendTransaction({
      from: account_two, value: 100, gas: 200000, to: frontendProxy.address,
      data: encodeCall('requestRelayEntry', ['uint256'], [0])
    });

    assert.equal((await frontend.getPastEvents())[0].event, 'RelayEntryRequested', "RelayEntryRequested event should occur on the frontendProxy contract.");

    let contractBalance = await web3.eth.getBalance(frontend.address);
    assert.equal(contractBalance, 100, "Keep Random Beacon contract should receive ether.");

    let contractBalanceViaProxy = await web3.eth.getBalance(frontendProxy.address);
    assert.equal(contractBalanceViaProxy, 100, "Keep Random Beacon contract new balance should be visible via frontendProxy.");
  });

  it("owner should be able to withdraw ether from random beacon contract", async function() {

    let amount = web3.utils.toWei('1', 'ether');
    await web3.eth.sendTransaction({
      from: account_two, value: amount, gas: 200000, to: frontendProxy.address,
      data: encodeCall('requestRelayEntry', ['uint256'], [0])
    });

    // should fail to withdraw if not owner
    await exceptThrow(frontend.initiateWithdrawal({from: account_two}));
    await exceptThrow(frontend.finishWithdrawal(account_two, {from: account_two}));

    await frontend.initiateWithdrawal({from: account_one});
    await exceptThrow(frontend.finishWithdrawal(account_three, {from: account_one}));

    let contractStartBalance = await web3.eth.getBalance(frontend.address);
    // jump in time, full withdrawal delay
    await increaseTimeTo(await latestTime()+duration.days(30));

    let receiverStartBalance = web3.utils.fromWei(await web3.eth.getBalance(account_three), 'ether');
    await frontend.finishWithdrawal(account_three, {from: account_one});
    let receiverEndBalance = web3.utils.fromWei(await web3.eth.getBalance(account_three), 'ether');
    assert(receiverEndBalance > receiverStartBalance, "Receiver updated balance should include received ether.");

    let contractEndBalance = await web3.eth.getBalance(frontend.address);
    assert.equal(contractEndBalance, contractStartBalance - amount, "Keep Random Beacon contract should send all ether.");
    let contractEndBalanceViaProxy = await web3.eth.getBalance(frontendProxy.address);
    assert.equal(contractEndBalanceViaProxy, contractStartBalance - amount, "Keep Random Beacon contract updated balance should be visible via frontendProxy.");

  });

  it("should fail to update minimum payment by non owner", async function() {
    await exceptThrow(frontend.setMinimumPayment(123, {from: account_two}));
  });

  it("should be able to update minimum payment by the owner", async function() {
    await frontend.setMinimumPayment(123);
    let newMinPayment = await frontend.minimumPayment();
    assert.equal(newMinPayment, 123, "Should be able to get updated minimum payment.");
  });
});
