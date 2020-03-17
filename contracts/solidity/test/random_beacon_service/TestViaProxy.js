import { duration, increaseTimeTo } from '../helpers/increaseTime';
import {bls} from '../helpers/data';
import latestTime from '../helpers/latestTime';
import expectThrow from '../helpers/expectThrow';
import {initContracts} from '../helpers/initContracts';
import {createSnapshot, restoreSnapshot} from "../helpers/snapshot";
const ServiceContractProxy = artifacts.require('./KeepRandomBeaconService.sol')

contract('TestKeepRandomBeaconService/ViaProxy', function(accounts) {

  let serviceContract, serviceContractProxy, operatorContract,
    account_one = accounts[0],
    account_two = accounts[1],
    account_three = accounts[2],
    entryFeeEstimate, entryFeeBreakdown;

  before(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      ServiceContractProxy,
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./stubs/KeepRandomBeaconOperatorStub.sol')
    );

    operatorContract = contracts.operatorContract;
    serviceContract = contracts.serviceContract;
    serviceContractProxy = await ServiceContractProxy.at(serviceContract.address);

    // Using stub method to add first group to help testing.
    await operatorContract.registerNewGroup(bls.groupPubKey);
    let group = await operatorContract.getGroupPublicKey(0);
    await operatorContract.setGroupMembers(group, [accounts[0]]);

    entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    entryFeeBreakdown = await serviceContract.entryFeeBreakdown();
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

  it("should be able to check if the service contract was initialized", async function() {
    assert.isTrue(
      await serviceContract.initialized(),
      "Service contract should be initialized."
    );
  });

  it("should fail to request relay entry with not enough ether", async function() {
    await expectThrow(
      serviceContract.methods['requestRelayEntry()']({from: account_two, value: 0})
    );
  });

  it("should be able to request relay with enough ether", async function() {
    let initialRequesterBalance = await web3.eth.getBalance(account_two);
    await serviceContract.fundRequestSubsidyFeePool({from: account_one, value: 100});
    let requestorSubsidy = web3.utils.toBN(1); // 1% is returned to the requestor.

    let initialServiceContractBalance = web3.utils.toBN(
      await web3.eth.getBalance(serviceContract.address)
    );
    let dkgSubmitterReimbursementFee = await operatorContract.dkgSubmitterReimbursementFee()

    let tx = await serviceContract.methods['requestRelayEntry()'](
      {from: account_two, value: entryFeeEstimate}
    )
    let transactionCost = web3.utils
      .toBN(tx.receipt.gasUsed)
      .mul(web3.utils.toWei(web3.utils.toBN(20), 'gwei')); // 20 default gasPrice

    assert.equal(
      (await operatorContract.getPastEvents())[0].event, 
      'RelayEntryRequested', 
      "RelayEntryRequested event should occur on operator contract."
    );

    assert.isTrue(
      web3.utils.toBN(initialRequesterBalance)
        .sub(entryFeeEstimate)
        .sub(transactionCost)
        .add(requestorSubsidy)
        .eq(web3.utils.toBN(await web3.eth.getBalance(account_two))), 
      "Requestor should receive 1% subsidy."
    );

    let serviceContractBalance = await web3.eth.getBalance(serviceContract.address);
    assert.isTrue(
      web3.utils.toBN(serviceContractBalance)
      .eq(initialServiceContractBalance
        .add(entryFeeBreakdown.dkgContributionFee)
        .sub(requestorSubsidy)
      ), 
      "Keep Random Beacon service contract should receive DKG fee fraction."
    );

    let serviceContractBalanceViaProxy = await web3.eth.getBalance(serviceContractProxy.address);
    assert.isTrue(
      web3.utils.toBN(serviceContractBalanceViaProxy)
      .eq(initialServiceContractBalance
        .add(entryFeeBreakdown.dkgContributionFee)
        .sub(requestorSubsidy)
      ), 
      "Keep Random Beacon service contract new balance should be visible via serviceContractProxy."
    );

    let operatorContractBalance = await web3.eth.getBalance(operatorContract.address);
    assert.isTrue(
      web3.utils.toBN(operatorContractBalance)
      .eq(entryFeeBreakdown.entryVerificationFee
        .add(entryFeeBreakdown.groupProfitFee)
        .add(dkgSubmitterReimbursementFee)
      ), 
      "Keep Random Beacon operator contract should receive entry fee, " +
      "group profit fee and dkg submitter reimbursement."
    );
  });

  it("should be able to request relay entry via serviceContractProxy contract with enough ether", async function() {
    let initialRequesterBalance = await web3.eth.getBalance(account_two);
    await serviceContract.fundRequestSubsidyFeePool({from: account_one, value: 100});
    let requestorSubsidy = web3.utils.toBN(1); // 1% is returned to the requestor.

    let initialServiceContractBalance = web3.utils.toBN(
      await web3.eth.getBalance(serviceContract.address)
    );
    let dkgSubmitterReimbursementFee = await operatorContract.dkgSubmitterReimbursementFee()

    let gasPrice = web3.utils.toWei(web3.utils.toBN(20), 'gwei');
    let transactionCost; 

    await web3.eth.sendTransaction({
      // if you see a plain 'revert' error, it's probably because of not enough gas
      from: account_two, 
      value: entryFeeEstimate, 
      gas: 500000, 
      gasPrice: gasPrice,
      to: serviceContractProxy.address,
      data: web3.eth.abi.encodeFunctionSignature('requestRelayEntry()')
    }).then(function(receipt){
      transactionCost = web3.utils.toBN(receipt.gasUsed).mul(gasPrice);
    });

    assert.equal(
      (await operatorContract.getPastEvents())[0].event, 
      'RelayEntryRequested', 
      "RelayEntryRequested event should occur on the operator contract."
    );
    
    assert.isTrue(
      web3.utils.toBN(initialRequesterBalance)
        .sub(entryFeeEstimate)
        .sub(transactionCost)
        .add(requestorSubsidy)
        .eq(web3.utils.toBN(await web3.eth.getBalance(account_two))), 
      "Requestor should receive 1% subsidy."
    );

    let contractBalance = await web3.eth.getBalance(serviceContract.address);
    assert.isTrue(
      web3.utils.toBN(contractBalance)
      .eq(initialServiceContractBalance
        .add(entryFeeBreakdown.dkgContributionFee)
        .sub(requestorSubsidy)
      ), 
      "Keep Random Beacon service contract should receive DKG fee fraction."
    );

    let contractBalanceServiceContract = await web3.eth.getBalance(serviceContractProxy.address);
    assert.isTrue(
      web3.utils.toBN(contractBalanceServiceContract)
      .eq(initialServiceContractBalance
        .add(entryFeeBreakdown.dkgContributionFee)
        .sub(requestorSubsidy)
      ), 
      "Keep Random Beacon service contract new balance should be visible via serviceContractProxy."
    );

    let operatorContractBalance = await web3.eth.getBalance(operatorContract.address);
    assert.isTrue(
      web3.utils.toBN(operatorContractBalance)
      .eq(entryFeeBreakdown.entryVerificationFee
        .add(entryFeeBreakdown.groupProfitFee)
        .add(dkgSubmitterReimbursementFee)
      ), 
      "Keep Random Beacon operator contract should receive entry fee, " + 
      "group profit fee and dkg submitter reimbursement."
    );
  });

  it("owner should be able to withdraw ether from random beacon service contract", async function() {
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0)

    // Send higher fee than entryFeeEstimate
    await serviceContract.methods['requestRelayEntry()'](
      {from: account_one, value: entryFeeEstimate.mul(web3.utils.toBN(2))}
    )

    // should fail to withdraw if not owner
    await expectThrow(serviceContract.initiateWithdrawal({from: account_two}));
    await expectThrow(serviceContract.finishWithdrawal(account_two, {from: account_two}));

    await serviceContract.initiateWithdrawal({from: account_one});
    await expectThrow(serviceContract.finishWithdrawal(account_three, {from: account_one}));

    // jump in time, full undelegation period
    await increaseTimeTo(await latestTime()+duration.days(30));

    let receiverStartBalance = await web3.eth.getBalance(account_three);
    await serviceContract.finishWithdrawal(account_three, {from: account_one});
    let receiverEndBalance = await web3.eth.getBalance(account_three);
    assert.isTrue(
      web3.utils.toBN(receiverEndBalance).gt(web3.utils.toBN(receiverStartBalance)),
      "Receiver updated balance should include received ether."
    );

    let contractEndBalance = await web3.eth.getBalance(serviceContract.address);
    assert.equal(
      contractEndBalance, 
      0, 
      "Keep Random Beacon contract should send all ether."
    );

    let contractEndBalanceViaProxy = await web3.eth.getBalance(serviceContractProxy.address);
    assert.equal(
      contractEndBalanceViaProxy,
      0, 
      "Keep Random Beacon contract updated balance should be visible via serviceContractProxy."
    );
  });
});
