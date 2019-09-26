import { duration, increaseTimeTo } from './helpers/increaseTime';
import {bls} from './helpers/data';
import latestTime from './helpers/latestTime';
import expectThrow from './helpers/expectThrow';
import encodeCall from './helpers/encodeCall';
import {initContracts} from './helpers/initContracts';
const ServiceContractProxy = artifacts.require('./KeepRandomBeaconService.sol')

contract('TestKeepRandomBeaconServiceViaProxy', function(accounts) {

  let serviceContract, serviceContractProxy, operatorContract, groupContract,
    account_one = accounts[0],
    account_two = accounts[1],
    account_three = accounts[2],
    entryFeeEstimate, minimumCallbackFee, entryFee;

  beforeEach(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      ServiceContractProxy,
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./stubs/KeepRandomBeaconOperatorStub.sol'),
      artifacts.require('./KeepRandomBeaconOperatorGroups.sol')
    );

    operatorContract = contracts.operatorContract;
    groupContract = contracts.groupContract;
    serviceContract = contracts.serviceContract;
    serviceContractProxy = await ServiceContractProxy.at(serviceContract.address);

    // Using stub method to add first group to help testing.
    await operatorContract.registerNewGroup(bls.groupPubKey);
    let group = await groupContract.getGroupPublicKey(0);
    await operatorContract.addGroupMember(group, accounts[0]);

    entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
    minimumCallbackFee = await serviceContract.minimumCallbackFee(0)
    entryFee = await serviceContract.entryFeeBreakdown()
  });

  it("should be able to check if the service contract was initialized", async function() {
    assert.isTrue(await serviceContract.initialized(), "Service contract should be initialized.");
  });

  it("should fail to request relay entry with not enough ether", async function() {
    await expectThrow(serviceContract.requestRelayEntry(0, {from: account_two, value: 0}));
  });

  it("should be able to request relay with enough ether", async function() {
    let account_two_balance = await web3.eth.getBalance(account_two);
    await serviceContract.fundRequestSubsidyFeePool({from: account_one, value: 100});
    let requestorSubsidy = web3.utils.toBN(1); // 1% is returned to the requestor.

    let contractPreviousBalance = web3.utils.toBN(await web3.eth.getBalance(serviceContract.address));
    let dkgSubmitterReimbursement = await operatorContract.dkgSubmitterReimbursement()

    let tx = await serviceContract.requestRelayEntry(0, {from: account_two, value: entryFeeEstimate})
    let transactionCost = web3.utils.toBN(tx.receipt.gasUsed).mul(web3.utils.toWei(web3.utils.toBN(20), 'gwei')); // 20 default gasPrice
    assert.isTrue(web3.utils.toBN(account_two_balance).sub(entryFeeEstimate).sub(transactionCost).add(requestorSubsidy).eq(web3.utils.toBN(await web3.eth.getBalance(account_two))), "Requestor should receive 1% subsidy.");

    assert.equal((await operatorContract.getPastEvents())[0].event, 'SignatureRequested', "SignatureRequested event should occur on operator contract.");

    let contractBalance = await web3.eth.getBalance(serviceContract.address);
    assert.isTrue(web3.utils.toBN(contractBalance).eq(contractPreviousBalance.add(entryFee.dkgContributionFee).sub(requestorSubsidy)), "Keep Random Beacon service contract should receive DKG fee fraction.");

    let contractBalanceViaProxy = await web3.eth.getBalance(serviceContractProxy.address);
    assert.isTrue(web3.utils.toBN(contractBalanceViaProxy).eq(contractPreviousBalance.add(entryFee.dkgContributionFee).sub(requestorSubsidy)), "Keep Random Beacon service contract new balance should be visible via serviceContractProxy.");

    let operatorContractBalance = await web3.eth.getBalance(operatorContract.address);
    assert.isTrue(web3.utils.toBN(operatorContractBalance).eq(entryFee.entryVerificationFee.add(minimumCallbackFee).add(entryFee.groupProfitFee).add(dkgSubmitterReimbursement)), "Keep Random Beacon operator contract should receive entry fee, callback payment, group profit fee and dkg submitter reimbursement.");
  
  });

  it("should be able to request relay entry via serviceContractProxy contract with enough ether", async function() {
    await expectThrow(serviceContractProxy.sendTransaction({from: account_two, value: entryFeeEstimate}));

    let contractPreviousBalance = web3.utils.toBN(await web3.eth.getBalance(serviceContract.address));
    let dkgSubmitterReimbursement = await operatorContract.dkgSubmitterReimbursement()

    await web3.eth.sendTransaction({
      // if you see a plain 'revert' error, it's probably because of not enough gas
      from: account_two, value: entryFeeEstimate, gas: 400000, to: serviceContractProxy.address,
      data: encodeCall('requestRelayEntry', ['uint256'], [0])
    });

    assert.equal((await operatorContract.getPastEvents())[0].event, 'SignatureRequested', "SignatureRequested event should occur on the operator contract.");

    let contractBalance = await web3.eth.getBalance(serviceContract.address);
    assert.isTrue(web3.utils.toBN(contractBalance).eq(contractPreviousBalance.add(entryFee.dkgContributionFee)), "Keep Random Beacon service contract should receive DKG fee fraction.");

    let contractBalanceServiceContract = await web3.eth.getBalance(serviceContractProxy.address);
    assert.isTrue(web3.utils.toBN(contractBalanceServiceContract).eq(contractPreviousBalance.add(entryFee.dkgContributionFee)), "Keep Random Beacon service contract new balance should be visible via serviceContractProxy.");

    let operatorContractBalance = await web3.eth.getBalance(operatorContract.address);
    assert.isTrue(web3.utils.toBN(operatorContractBalance).eq(entryFee.entryVerificationFee.add(minimumCallbackFee).add(entryFee.groupProfitFee).add(dkgSubmitterReimbursement)), "Keep Random Beacon operator contract should receive entry fee, callback payment, group profit fee and dkg submitter reimbursement.");
  });

  it("owner should be able to withdraw ether from random beacon service contract", async function() {
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
    await serviceContract.requestRelayEntry(0, {from: account_one, value: entryFeeEstimate})

    // should fail to withdraw if not owner
    await expectThrow(serviceContract.initiateWithdrawal({from: account_two}));
    await expectThrow(serviceContract.finishWithdrawal(account_two, {from: account_two}));

    await serviceContract.initiateWithdrawal({from: account_one});
    await expectThrow(serviceContract.finishWithdrawal(account_three, {from: account_one}));

    // jump in time, full withdrawal delay
    await increaseTimeTo(await latestTime()+duration.days(30));

    let receiverStartBalance = await web3.eth.getBalance(account_three);
    await serviceContract.finishWithdrawal(account_three, {from: account_one});
    let receiverEndBalance = await web3.eth.getBalance(account_three);
    assert.isTrue(receiverEndBalance > receiverStartBalance, "Receiver updated balance should include received ether.");

    let contractEndBalance = await web3.eth.getBalance(serviceContract.address);
    assert.equal(contractEndBalance, 0, "Keep Random Beacon contract should send all ether.");
    let contractEndBalanceViaProxy = await web3.eth.getBalance(serviceContractProxy.address);
    assert.equal(contractEndBalanceViaProxy, 0, "Keep Random Beacon contract updated balance should be visible via serviceContractProxy.");

  });

  it("should fail to update minimum gas price by non owner", async function() {
    await expectThrow(serviceContract.setMinimumGasPrice(123, {from: account_two}));
  });

  it("should be able to update minimum gas price by the owner", async function() {
    await serviceContract.setMinimumGasPrice(123);
    let newMinGasPrice = await serviceContract.minimumGasPrice();
    assert.equal(newMinGasPrice, 123, "Should be able to get updated minimum gas price.");
  });
});
