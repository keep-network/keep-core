import {bls} from './helpers/data';
import latestTime from './helpers/latestTime';
import { duration, increaseTimeTo } from './helpers/increaseTime';
import expectThrow from './helpers/expectThrow';
import expectThrowWithMessage from './helpers/expectThrowWithMessage';
import encodeCall from './helpers/encodeCall';
import mineBlocks from './helpers/mineBlocks';
const Proxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconImplV1 = artifacts.require('./KeepRandomBeaconImplV1.sol');
const KeepGroup = artifacts.require('./KeepGroupStub.sol');

contract('TestKeepRandomBeaconViaProxy', function(accounts) {

  const relayRequestTimeout = 10;
  const blocksForward = 20;

  let implV1, proxy, implViaProxy, keepGroup,
    account_one = accounts[0],
    account_two = accounts[1],
    account_three = accounts[2];

  beforeEach(async () => {
    implV1 = await KeepRandomBeaconImplV1.new();
    proxy = await Proxy.new(implV1.address);
    implViaProxy = await KeepRandomBeaconImplV1.at(proxy.address);
    keepGroup = await KeepGroup.new();
  });

  describe('Logic for threshold random number generation', function() {
    beforeEach(async () => {
      await implViaProxy.initialize(100, duration.days(30), bls.previousEntry, bls.groupPubKey, keepGroup.address, 
        relayRequestTimeout);
    });
  
    it("should be able to check if the implementation contract was initialized", async function() {
      let result = await implViaProxy.initialized();
      assert.equal(result, true, "Implementation contract should be initialized.");
    });
  
    it("should fail to request relay entry with not enough ether", async function() {
      await expectThrow(implViaProxy.requestRelayEntry(0, {from: account_two, value: 99}));
    });
  
    it("should be able to request relay entry via implementation contract with enough ether", async function() {
      await mineBlocks(blocksForward);
      await implViaProxy.requestRelayEntry(0, {from: account_two, value: 100})
  
      assert.equal((await implViaProxy.getPastEvents())[0].event, 'RelayEntryRequested', "RelayEntryRequested event should occur on the implementation contract.");
  
      let contractBalance = await web3.eth.getBalance(implViaProxy.address);
      assert.equal(contractBalance, 100, "Keep Random Beacon contract should receive ether.");
  
      let contractBalanceViaProxy = await web3.eth.getBalance(proxy.address);
      assert.equal(contractBalanceViaProxy, 100, "Keep Random Beacon contract new balance should be visible via proxy.");
    });
  
    it("should be able to request relay entry via proxy contract with enough ether", async function() {
      await expectThrow(proxy.sendTransaction({from: account_two, value: 1000}));
  
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
      await expectThrow(implViaProxy.initiateWithdrawal({from: account_two}));
      await expectThrow(implViaProxy.finishWithdrawal(account_two, {from: account_two}));
  
      await implViaProxy.initiateWithdrawal({from: account_one});
      await expectThrow(implViaProxy.finishWithdrawal(account_three, {from: account_one}));
  
      let contractStartBalance = await web3.eth.getBalance(implViaProxy.address);
      // jump in time, full withdrawal delay
      await increaseTimeTo(await latestTime()+duration.days(30));
  
      let receiverStartBalance = web3.utils.fromWei(await web3.eth.getBalance(account_three), 'ether');
      await implViaProxy.finishWithdrawal(account_three, {from: account_one});
      let receiverEndBalance = web3.utils.fromWei(await web3.eth.getBalance(account_three), 'ether');
      assert(Number(receiverEndBalance) > Number(receiverStartBalance), "Receiver updated balance should include received ether.");
  
      let contractEndBalance = await web3.eth.getBalance(implViaProxy.address);
      assert.equal(contractEndBalance, contractStartBalance - amount, "Keep Random Beacon contract should send all ether.");
      let contractEndBalanceViaProxy = await web3.eth.getBalance(proxy.address);
      assert.equal(contractEndBalanceViaProxy, contractStartBalance - amount, "Keep Random Beacon contract updated balance should be visible via proxy.");
  
    });
  
    it("should fail to update minimum payment by non owner", async function() {
      await expectThrow(implViaProxy.setMinimumPayment(123, {from: account_two}));
    });
  
    it("should be able to update minimum payment by the owner", async function() {
      await implViaProxy.setMinimumPayment(123);
      let newMinPayment = await implViaProxy.minimumPayment();
      assert.equal(newMinPayment, 123, "Should be able to get updated minimum payment.");
    });
  })

  describe("Relay request timeout when expecting an error", function() {
    it("should throw an error when signing is not in progress and a block number is lower than the relay entry timeout", async function() {
      let currentBlockNumber = await web3.eth.getBlockNumber()
      let timeout = currentBlockNumber + blocksForward
      await implViaProxy.initialize(100, duration.days(30), bls.previousEntry, bls.groupPubKey, keepGroup.address, timeout)
  
      await expectThrowWithMessage(implViaProxy.requestRelayEntry(0, {from: account_two, value: 100}), 'Relay entry request is in progress.')
    });
  
    it("should throw an error when signing is in progess and a block number is lower than the relay entry timeout", async function() {
      await implViaProxy.initialize(100, duration.days(30), bls.previousEntry, bls.groupPubKey, keepGroup.address, relayRequestTimeout)
      await mineBlocks(blocksForward);
  
      await implViaProxy.requestRelayEntry(0, {from: account_two, value: 100})
  
      await expectThrowWithMessage(implViaProxy.requestRelayEntry(0, {from: account_two, value: 100}), 'Relay entry request is in progress.')
    });
  
    it("should throw an error when signing is in progess and a block number is higher than the relay entry timeout", async function() {
      await implViaProxy.initialize(100, duration.days(30), bls.previousEntry, bls.groupPubKey, keepGroup.address, relayRequestTimeout)
      await implViaProxy.requestRelayEntry(0, {from: account_two, value: 100})
      
      await mineBlocks(blocksForward);

      await expectThrowWithMessage(implViaProxy.requestRelayEntry(0, {from: account_two, value: 100}), 'Relay entry request is in progress.')
    });

  })

});
