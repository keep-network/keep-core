import mineBlocks from './helpers/mineBlocks';
import {bls} from './helpers/data';
import {initContracts} from './helpers/initContracts';
const CallbackContract = artifacts.require('./examples/CallbackContract.sol');

contract('TestKeepRandomBeaconServicePricing', function(accounts) {

  let operatorContract, serviceContract, callbackContract,entryFee, groupSize,
    requestor = accounts[1],
    operator1 = accounts[2],
    operator2 = accounts[3],
    operator3 = accounts[4];

  beforeEach(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./KeepRandomBeaconOperatorStub.sol')
    );

    operatorContract = contracts.operatorContract;
    serviceContract = contracts.serviceContract;
    callbackContract = await CallbackContract.new();

    // Using stub method to add first group to help testing.
    await operatorContract.registerNewGroup(bls.groupPubKey);

    groupSize = web3.utils.toBN(3);
    await operatorContract.setGroupSize(groupSize);
    let group = await operatorContract.getGroupPublicKey(0);
    await operatorContract.addGroupMember(group, operator1);
    await operatorContract.addGroupMember(group, operator2);
    await operatorContract.addGroupMember(group, operator3);

    entryFee = await serviceContract.entryFeeBreakdown()
  });

  it("should successfully refund callback gas surplus to the requestor", async function() {

    // Set higher gas price
    await serviceContract.setMinimumGasPrice(web3.utils.toWei(web3.utils.toBN(200), 'gwei'));

    let minimumPayment = await serviceContract.minimumPayment()
    await serviceContract.methods['requestRelayEntry(uint256,address,string)'](
      bls.seed,
      callbackContract.address,
      "callback(uint256)",
      {value: minimumPayment, from: requestor}
    );

    let minimumCallbackPayment = await serviceContract.minimumCallbackPayment()
    let requestorBalance = await web3.eth.getBalance(requestor);

    await operatorContract.relayEntry(bls.nextGroupSignature);

    // Put back the default gas price
    await serviceContract.setMinimumGasPrice(web3.utils.toWei(web3.utils.toBN(20), 'gwei'));

    let updatedMinimumCallbackPayment = await serviceContract.minimumCallbackPayment()
    let updatedRequestorBalance = await web3.eth.getBalance(requestor)

    let surplus = web3.utils.toBN(minimumCallbackPayment).sub(web3.utils.toBN(updatedMinimumCallbackPayment))
    let refund = web3.utils.toBN(updatedRequestorBalance).sub(web3.utils.toBN(requestorBalance))

    assert.isTrue(refund.eq(surplus), "Callback gas surplus should be refunded to the requestor.");

  });

  it("should send group reward to each operator.", async function() {

    let operator1balance = web3.utils.toBN(await web3.eth.getBalance(operator1));
    let operator2balance = web3.utils.toBN(await web3.eth.getBalance(operator2));
    let operator3balance = web3.utils.toBN(await web3.eth.getBalance(operator3));

    let minimumPayment = await serviceContract.minimumPayment()
    await serviceContract.methods['requestRelayEntry(uint256,address,string)'](
      bls.seed,
      callbackContract.address,
      "callback(uint256)",
      {value: minimumPayment, from: requestor}
    );

    let currentEntryStartBlock = await operatorContract.currentEntryStartBlock();
    let relayEntryTimeout = await operatorContract.relayEntryTimeout();
    let deadlineBlock = currentEntryStartBlock.add(relayEntryTimeout);
    let currentBlock = web3.utils.toBN(await web3.eth.getBlockNumber()).add(web3.utils.toBN(1)); // web3.eth.getBlockNumber is 1 block behind solidity 'block.number'.

    let decimalPoints = web3.utils.toBN(100);
    let delayFactor = (deadlineBlock.sub(currentBlock)).mul(decimalPoints).div(relayEntryTimeout.sub(web3.utils.toBN(1))).pow(web3.utils.toBN(2));
    let baseReward = entryFee.profitMargin.div(groupSize)
    let expectedGroupReward = baseReward.mul(delayFactor).div(decimalPoints.pow(web3.utils.toBN(2)));

    await operatorContract.relayEntry(bls.nextGroupSignature);

    assert.isTrue(operator1balance.add(expectedGroupReward).eq(web3.utils.toBN(await web3.eth.getBalance(operator1))), "Operator should receive group reward.");
    assert.isTrue(operator2balance.add(expectedGroupReward).eq(web3.utils.toBN(await web3.eth.getBalance(operator2))), "Operator should receive group reward.");
    assert.isTrue(operator3balance.add(expectedGroupReward).eq(web3.utils.toBN(await web3.eth.getBalance(operator3))), "Operator should receive group reward.");
  });

  it("should send part of the group reward to request subsidy pool based on the submission block .", async function() {

    let operator1balance = web3.utils.toBN(await web3.eth.getBalance(operator1));
    let operator2balance = web3.utils.toBN(await web3.eth.getBalance(operator2));
    let operator3balance = web3.utils.toBN(await web3.eth.getBalance(operator3));

    let minimumPayment = await serviceContract.minimumPayment()
    await serviceContract.methods['requestRelayEntry(uint256,address,string)'](
      bls.seed,
      callbackContract.address,
      "callback(uint256)",
      {value: minimumPayment, from: requestor}
    );

    let currentEntryStartBlock = await operatorContract.currentEntryStartBlock();
    let relayEntryTimeout = await operatorContract.relayEntryTimeout();
    let deadlineBlock = currentEntryStartBlock.add(relayEntryTimeout);
    let decimalPoints = web3.utils.toBN(100);

    mineBlocks(relayEntryTimeout.toNumber()/2);

    let currentBlock = web3.utils.toBN(await web3.eth.getBlockNumber()).add(web3.utils.toBN(1)); // web3.eth.getBlockNumber is 1 block behind solidity 'block.number'.
    let delayFactor = (deadlineBlock.sub(currentBlock)).mul(decimalPoints).div(relayEntryTimeout.sub(web3.utils.toBN(1))).pow(web3.utils.toBN(2));

    let baseReward = entryFee.profitMargin.div(groupSize)
    let expectedGroupReward = baseReward.mul(delayFactor).div(decimalPoints.pow(web3.utils.toBN(2)));
    let requestSubsidy = entryFee.profitMargin.sub(expectedGroupReward.mul(groupSize));

    let serviceContractBalance = web3.utils.toBN(await web3.eth.getBalance(serviceContract.address));

    await operatorContract.relayEntry(bls.nextGroupSignature);

    assert.isTrue(operator1balance.add(expectedGroupReward).eq(web3.utils.toBN(await web3.eth.getBalance(operator1))), "Operator should receive reduced group reward.");
    assert.isTrue(operator2balance.add(expectedGroupReward).eq(web3.utils.toBN(await web3.eth.getBalance(operator2))), "Operator should receive reduced group reward.");
    assert.isTrue(operator3balance.add(expectedGroupReward).eq(web3.utils.toBN(await web3.eth.getBalance(operator3))), "Operator should receive reduced group reward.");
    assert.isTrue(serviceContractBalance.add(requestSubsidy).eq(web3.utils.toBN(await web3.eth.getBalance(serviceContract.address))), "Service contract should receive request subsidy.");
  });
});
