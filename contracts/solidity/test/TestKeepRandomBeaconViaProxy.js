import { duration, increaseTimeTo } from './helpers/increaseTime';
import {bls} from './helpers/data';
import latestTime from './helpers/latestTime';
import exceptThrow from './helpers/expectThrow';
import encodeCall from './helpers/encodeCall';
import {initContracts} from './helpers/initContracts';
const FrontendProxy = artifacts.require('./KeepRandomBeaconFrontendProxy.sol')

contract('TestKeepRandomBeaconViaProxy', function(accounts) {

  let frontend, frontendProxy, backend,
    account_one = accounts[0],
    account_two = accounts[1],
    account_three = accounts[2];

  before(async () => {
    let contracts = await initContracts(
      accounts,
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./StakingProxy.sol'),
      artifacts.require('./TokenStaking.sol'),
      FrontendProxy,
      artifacts.require('./KeepRandomBeaconFrontendImplV1.sol'),
      artifacts.require('./KeepRandomBeaconBackendStub.sol')
    );
  
    backend = contracts.backend;
    frontend = contracts.frontend;
    frontendProxy = await FrontendProxy.at(frontend.address);

    // Using stub method to add first group to help testing.
    await backend.registerNewGroup(bls.groupPubKey);
  });

  
  it("should be able to check if the frontend contract was initialized", async function() {
    assert.isTrue(await frontend.initialized(), "Frontend contract should be initialized.");
  });

  it("should fail to request relay entry with not enough ether", async function() {
    await exceptThrow(frontend.requestRelayEntry(0, {from: account_two, value: 0}));
  });

  it("should be able to request relay with enough ether", async function() {
    await frontend.requestRelayEntry(0, {from: account_two, value: 100})

    assert.equal((await backend.getPastEvents())[0].event, 'RelayEntryRequested', "RelayEntryRequested event should occur on backend contract.");

    let contractBalance = await web3.eth.getBalance(frontend.address);
    assert.equal(contractBalance, 100, "Keep Random Beacon frontend contract should receive ether.");

    let contractBalanceViaProxy = await web3.eth.getBalance(frontendProxy.address);
    assert.equal(contractBalanceViaProxy, 100, "Keep Random Beacon frontend contract new balance should be visible via frontendProxy.");

  });

  it("should be able to request relay entry via frontendProxy contract with enough ether", async function() {
    await exceptThrow(frontendProxy.sendTransaction({from: account_two, value: 1000}));

    await web3.eth.sendTransaction({
      from: account_two, value: 100, gas: 200000, to: frontendProxy.address,
      data: encodeCall('requestRelayEntry', ['uint256'], [0])
    });

    assert.equal((await backend.getPastEvents())[0].event, 'RelayEntryRequested', "RelayEntryRequested event should occur on the backend contract.");

    let contractBalance = await web3.eth.getBalance(frontend.address);
    assert.equal(contractBalance, 200, "Keep Random Beacon frontend contract should receive ether.");

    let contractBalanceFrontend = await web3.eth.getBalance(frontendProxy.address);
    assert.equal(contractBalanceFrontend, 200, "Keep Random Beacon contract new balance should be visible via frontendProxy.");
  });

  it("owner should be able to withdraw ether from random beacon frontend contract", async function() {

    // should fail to withdraw if not owner
    await exceptThrow(frontend.initiateWithdrawal({from: account_two}));
    await exceptThrow(frontend.finishWithdrawal(account_two, {from: account_two}));

    await frontend.initiateWithdrawal({from: account_one});
    await exceptThrow(frontend.finishWithdrawal(account_three, {from: account_one}));

    // jump in time, full withdrawal delay
    await increaseTimeTo(await latestTime()+duration.days(30));

    let receiverStartBalance = web3.utils.fromWei(await web3.eth.getBalance(account_three), 'ether');
    await frontend.finishWithdrawal(account_three, {from: account_one});
    let receiverEndBalance = web3.utils.fromWei(await web3.eth.getBalance(account_three), 'ether');
    assert(Number(receiverEndBalance) > Number(receiverStartBalance), "Receiver updated balance should include received ether.");

    let contractEndBalance = await web3.eth.getBalance(frontend.address);
    assert.equal(contractEndBalance, 0, "Keep Random Beacon contract should send all ether.");
    let contractEndBalanceViaProxy = await web3.eth.getBalance(frontendProxy.address);
    assert.equal(contractEndBalanceViaProxy, 0, "Keep Random Beacon contract updated balance should be visible via frontendProxy.");

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
