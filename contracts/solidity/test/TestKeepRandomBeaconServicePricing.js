import mineBlocks from './helpers/mineBlocks';
import {bls} from './helpers/data';
import stakeDelegate from './helpers/stakeDelegate';
import {initContracts} from './helpers/initContracts';
const CallbackContract = artifacts.require('./examples/CallbackContract.sol');

contract('TestKeepRandomBeaconServicePricing', function(accounts) {

  let token, stakingContract, operatorContract, serviceContract, callbackContract, entryFee, groupSize,
    owner = accounts[0],
    requestor = accounts[1],
    operator1 = accounts[2],
    operator2 = accounts[3],
    operator3 = accounts[4],
    magpie1 = accounts[5],
    magpie2 = accounts[6],
    magpie3 = accounts[7];

  beforeEach(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./stubs/KeepRandomBeaconOperatorStub.sol')
    );

    token = contracts.token;
    stakingContract = contracts.stakingContract;
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

    await stakeDelegate(stakingContract, token, owner, operator1, magpie1, 0);
    await stakeDelegate(stakingContract, token, owner, operator2, magpie2, 0);
    await stakeDelegate(stakingContract, token, owner, operator3, magpie3, 0);

    entryFee = await serviceContract.entryFeeBreakdown()
  });

  it("should successfully refund callback gas surplus to the requestor if gas price was high", async function() {

    // Set higher gas price
    await serviceContract.setMinimumGasPrice(web3.utils.toWei(web3.utils.toBN(200), 'gwei'));

    let callbackGas = await callbackContract.callback.estimateGas(bls.nextGroupSignature);
    let minimumPayment = await serviceContract.minimumPayment(callbackGas)
    await serviceContract.methods['requestRelayEntry(uint256,address,string,uint256)'](
      bls.seed,
      callbackContract.address,
      "callback(uint256)",
      0,
      {value: minimumPayment, from: requestor}
    );

    let minimumCallbackPayment = await serviceContract.minimumCallbackPayment(0)
    let requestorBalance = await web3.eth.getBalance(requestor);

    await operatorContract.relayEntry(bls.nextGroupSignature);

    // Put back the default gas price
    await serviceContract.setMinimumGasPrice(web3.utils.toWei(web3.utils.toBN(20), 'gwei'));

    let updatedMinimumCallbackPayment = await serviceContract.minimumCallbackPayment(0)
    let updatedRequestorBalance = await web3.eth.getBalance(requestor)

    let surplus = web3.utils.toBN(minimumCallbackPayment).sub(web3.utils.toBN(updatedMinimumCallbackPayment))
    let refund = web3.utils.toBN(updatedRequestorBalance).sub(web3.utils.toBN(requestorBalance))

    assert.isTrue(refund.eq(surplus), "Callback gas surplus should be refunded to the requestor.");

  });

  it("should successfully refund callback gas surplus to the requestor if gas estimation was high", async function() {

    let callbackGas = await callbackContract.callback.estimateGas(bls.nextGroupSignature);
    let expectedCallbackPayment = await serviceContract.minimumCallbackPayment(callbackGas);

    let excessCallbackGas = web3.utils.toBN(callbackGas).mul(web3.utils.toBN(2)); // Set higher callback gas estimate.
    let excessCallbackPayment = await serviceContract.minimumCallbackPayment(excessCallbackGas);

    let minimumPayment = await serviceContract.minimumPayment(excessCallbackGas)
    await serviceContract.methods['requestRelayEntry(uint256,address,string,uint256)'](
      bls.seed,
      callbackContract.address,
      "callback(uint256)",
      excessCallbackGas,
      {value: minimumPayment, from: requestor}
    );

    let requestorBalance = await web3.eth.getBalance(requestor);
    await operatorContract.relayEntry(bls.nextGroupSignature);
    let updatedRequestorBalance = await web3.eth.getBalance(requestor)

    // Ethereum transaction min cost varies i.e. 20864-21000 Gas resulting slightly different
    // eth amounts: Surplus 0.00219018 vs Refund 0.00218752 so rounding up those for the tests
    let surplus = web3.utils.fromWei(web3.utils.toBN(excessCallbackPayment).sub(web3.utils.toBN(expectedCallbackPayment)), 'ether')
    let refund = web3.utils.fromWei(web3.utils.toBN(updatedRequestorBalance).sub(web3.utils.toBN(requestorBalance)), 'ether')
    assert.isTrue(Math.round(surplus*10000)/10000 === Math.round(refund*10000)/10000, "Callback gas surplus should be refunded to the requestor.");
  });

  it("should send group reward to each operator.", async function() {

    let magpie1balance = web3.utils.toBN(await web3.eth.getBalance(magpie1));
    let magpie2balance = web3.utils.toBN(await web3.eth.getBalance(magpie2));
    let magpie3balance = web3.utils.toBN(await web3.eth.getBalance(magpie3));

    let minimumPayment = await serviceContract.minimumPayment(0)
    await serviceContract.methods['requestRelayEntry(uint256,address,string,uint256)'](
      bls.seed,
      callbackContract.address,
      "callback(uint256)",
      0,
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

    assert.isTrue(magpie1balance.add(expectedGroupReward).eq(web3.utils.toBN(await web3.eth.getBalance(magpie1))), "Beneficiary should receive group reward.");
    assert.isTrue(magpie2balance.add(expectedGroupReward).eq(web3.utils.toBN(await web3.eth.getBalance(magpie2))), "Beneficiary should receive group reward.");
    assert.isTrue(magpie3balance.add(expectedGroupReward).eq(web3.utils.toBN(await web3.eth.getBalance(magpie3))), "Beneficiary should receive group reward.");
  });

  it("should send part of the group reward to request subsidy pool based on the submission block .", async function() {

    let magpie1balance = web3.utils.toBN(await web3.eth.getBalance(magpie1));
    let magpie2balance = web3.utils.toBN(await web3.eth.getBalance(magpie2));
    let magpie3balance = web3.utils.toBN(await web3.eth.getBalance(magpie3));

    let minimumPayment = await serviceContract.minimumPayment(0)
    await serviceContract.methods['requestRelayEntry(uint256,address,string,uint256)'](
      bls.seed,
      callbackContract.address,
      "callback(uint256)",
      0,
      {value: minimumPayment, from: requestor}
    );

    let currentEntryStartBlock = await operatorContract.currentEntryStartBlock();
    let relayEntryTimeout = await operatorContract.relayEntryTimeout();
    let deadlineBlock = currentEntryStartBlock.add(relayEntryTimeout);
    let decimalPoints = web3.utils.toBN(100);

    mineBlocks(relayEntryTimeout.toNumber()/2);

    let currentBlock = web3.utils.toBN(await web3.eth.getBlockNumber()).add(web3.utils.toBN(1)); // web3.eth.getBlockNumber is 1 block behind solidity 'block.number'.
    let delayFactor = (deadlineBlock.sub(currentBlock)).mul(decimalPoints).div(relayEntryTimeout.sub(web3.utils.toBN(1))).pow(web3.utils.toBN(2));
    let delayFactorInverse = decimalPoints.pow(web3.utils.toBN(2)).sub(delayFactor);

    let baseReward = entryFee.profitMargin.div(groupSize)
    let expectedGroupReward = baseReward.mul(delayFactor).div(decimalPoints.pow(web3.utils.toBN(2)));
    let expectedDelayPenalty = baseReward.mul(delayFactorInverse).div(decimalPoints.pow(web3.utils.toBN(2)));
    let expectedExtraReward = expectedDelayPenalty.mul(groupSize).mul(web3.utils.toBN(5)).div(web3.utils.toBN(100));
    let requestSubsidy = entryFee.profitMargin.sub(expectedGroupReward.mul(groupSize)).sub(expectedExtraReward);

    let serviceContractBalance = web3.utils.toBN(await web3.eth.getBalance(serviceContract.address));

    await operatorContract.relayEntry(bls.nextGroupSignature);

    assert.isTrue(magpie1balance.add(expectedGroupReward).eq(web3.utils.toBN(await web3.eth.getBalance(magpie1))), "Beneficiary should receive reduced group reward.");
    assert.isTrue(magpie2balance.add(expectedGroupReward).eq(web3.utils.toBN(await web3.eth.getBalance(magpie2))), "Beneficiary should receive reduced group reward.");
    assert.isTrue(magpie3balance.add(expectedGroupReward).eq(web3.utils.toBN(await web3.eth.getBalance(magpie3))), "Beneficiary should receive reduced group reward.");
    assert.isTrue(serviceContractBalance.add(requestSubsidy).eq(web3.utils.toBN(await web3.eth.getBalance(serviceContract.address))), "Service contract should receive request subsidy.");
  });
});
