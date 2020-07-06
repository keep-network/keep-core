const {createSnapshot, restoreSnapshot} = require("../helpers/snapshot.js")
const blsData = require("../helpers/data.js")
const {initContracts} = require('../helpers/initContracts')
const assert = require('chai').assert
const {contract, web3, accounts} = require("@openzeppelin/test-environment")
const {expectRevert, time} = require("@openzeppelin/test-helpers")

const ServiceContractProxy = contract.fromArtifact('KeepRandomBeaconService')

describe('TestKeepRandomBeaconService/ViaProxy', function() {

  let serviceContract, serviceContractProxy, operatorContract,
    account_one = accounts[0],
    account_two = accounts[1],
    account_three = accounts[2],
    entryFeeEstimate, entryFeeBreakdown;

  before(async () => {
    let contracts = await initContracts(
      contract.fromArtifact('TokenStaking'),
      ServiceContractProxy,
      contract.fromArtifact('KeepRandomBeaconServiceImplV1'),
      contract.fromArtifact('KeepRandomBeaconOperatorStub')
    );

    operatorContract = contracts.operatorContract;
    serviceContract = contracts.serviceContract;
    serviceContractProxy = await ServiceContractProxy.at(serviceContract.address);

    // Using stub method to add first group to help testing.
    await operatorContract.registerNewGroup(blsData.groupPubKey);
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
    await expectRevert(
      serviceContract.methods['requestRelayEntry()']({from: account_two, value: 0}),
      "Payment is less than required minimum."
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
});
