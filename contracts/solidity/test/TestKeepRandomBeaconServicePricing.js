import mineBlocks from './helpers/mineBlocks';
import {bls} from './helpers/data';
import stakeDelegate from './helpers/stakeDelegate';
import {initContracts} from './helpers/initContracts';
const CallbackContract = artifacts.require('./examples/CallbackContract.sol');

contract('TestKeepRandomBeaconServicePricing', function(accounts) {

  let token, stakingContract, operatorContract, serviceContract, callbackContract, entryFee, groupSize, group,
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
    group = await operatorContract.getGroupPublicKey(0);
    await operatorContract.addGroupMember(group, operator1);
    await operatorContract.addGroupMember(group, operator2);
    await operatorContract.addGroupMember(group, operator3);

    await stakeDelegate(stakingContract, token, owner, operator1, magpie1, operator1, 0);
    await stakeDelegate(stakingContract, token, owner, operator2, magpie2, operator2, 0);
    await stakeDelegate(stakingContract, token, owner, operator3, magpie3, operator3, 0);

    entryFee = await serviceContract.entryFeeBreakdown()
  });

  it("should successfully refund callback gas surplus to the requestor if gas price was high", async function() {

    let defaultPriceFeedEstimate = await serviceContract.priceFeedEstimate();

    // Set higher gas price
    await serviceContract.setPriceFeedEstimate(defaultPriceFeedEstimate.mul(web3.utils.toBN(10)));
    await operatorContract.setPriceFeedEstimate(defaultPriceFeedEstimate.mul(web3.utils.toBN(10)));
    let callbackGas = await callbackContract.callback.estimateGas(bls.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    let excessCallbackFee = await serviceContract.callbackFee(callbackGas)

    await serviceContract.methods['requestRelayEntry(address,string,uint256)'](
      callbackContract.address,
      "callback(uint256)",
      callbackGas,
      {value: entryFeeEstimate, from: requestor}
    );

    let requestorBalance = await web3.eth.getBalance(requestor);

    await operatorContract.relayEntry(bls.groupSignature);

    // Put back the default gas price
    await serviceContract.setPriceFeedEstimate(defaultPriceFeedEstimate);
    await operatorContract.setPriceFeedEstimate(defaultPriceFeedEstimate);
    let expectedCallbackFee = await serviceContract.callbackFee((callbackGas/1.5).toFixed()) // Remove 1.5 fluctuation safety margin
    let updatedRequestorBalance = await web3.eth.getBalance(requestor)

    // Ethereum transaction min cost varies i.e. 20864-21000 Gas resulting slightly different
    // eth amounts: Surplus 0.00219018 vs Refund 0.00218752 so rounding up those for the tests
    let surplus = web3.utils.fromWei(web3.utils.toBN(excessCallbackFee).sub(web3.utils.toBN(expectedCallbackFee)), 'ether')
    let refund = web3.utils.fromWei(web3.utils.toBN(updatedRequestorBalance).sub(web3.utils.toBN(requestorBalance)), 'ether')
    assert.isTrue(Math.round(surplus*100)/100 === Math.round(refund*100)/100, "Callback gas surplus should be refunded to the requestor.");
  });

  it("should successfully refund callback gas surplus to the requestor if gas estimation was high", async function() {
    let callbackGas = await callbackContract.callback.estimateGas(bls.groupSignature);
    let expectedCallbackFee = await serviceContract.callbackFee((callbackGas/1.5).toFixed()); // Remove 1.5 fluctuation safety margin

    let excessCallbackGas = web3.utils.toBN(callbackGas).mul(web3.utils.toBN(2)); // Set higher callback gas estimate.
    let excessCallbackFee = await serviceContract.callbackFee(excessCallbackGas);

    let entryFeeEstimate = await serviceContract.entryFeeEstimate(excessCallbackGas)
    await serviceContract.methods['requestRelayEntry(address,string,uint256)'](
      callbackContract.address,
      "callback(uint256)",
      excessCallbackGas,
      {value: entryFeeEstimate, from: requestor}
    );

    let requestorBalance = await web3.eth.getBalance(requestor);
    await operatorContract.relayEntry(bls.groupSignature);
    let updatedRequestorBalance = await web3.eth.getBalance(requestor)

    // Ethereum transaction min cost varies i.e. 20864-21000 Gas resulting slightly different
    // eth amounts: Surplus 0.00219018 vs Refund 0.00218752 so rounding up those for the tests
    let surplus = web3.utils.fromWei(web3.utils.toBN(excessCallbackFee).sub(web3.utils.toBN(expectedCallbackFee)), 'ether')
    let refund = web3.utils.fromWei(web3.utils.toBN(updatedRequestorBalance).sub(web3.utils.toBN(requestorBalance)), 'ether')
    assert.isTrue(Math.round(surplus*100)/100 === Math.round(refund*100)/100, "Callback gas surplus should be refunded to the requestor.");
  });

  it("should send group reward to each operator.", async function() {
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
    let tx = await serviceContract.methods['requestRelayEntry(address,string,uint256)'](
      callbackContract.address,
      "callback(uint256)",
      0,
      {value: entryFeeEstimate, from: requestor}
    );

    let currentEntryStartBlock = web3.utils.toBN(tx.receipt.blockNumber);
    let relayEntryTimeout = await operatorContract.relayEntryTimeout();
    let relayEntryGenerationTime = await operatorContract.relayEntryGenerationTime();
    let deadlineBlock = currentEntryStartBlock.add(relayEntryTimeout);
    let entryReceivedBlock = currentEntryStartBlock.add(relayEntryGenerationTime).add(web3.utils.toBN(1));
    let remainingBlocks = deadlineBlock.sub(entryReceivedBlock);
    let submissionWindow = deadlineBlock.sub(entryReceivedBlock);
    let decimalPoints = web3.utils.toBN(1e16);
    let delayFactor = (remainingBlocks.mul(decimalPoints).div(submissionWindow)).pow(web3.utils.toBN(2));
    let memberBaseReward = entryFee.groupProfitFee.div(groupSize)
    let expectedGroupMemberReward = memberBaseReward.mul(delayFactor).div(decimalPoints.pow(web3.utils.toBN(2)));

    await operatorContract.relayEntry(bls.groupSignature);

    assert.isTrue(delayFactor.eq(web3.utils.toBN(1e16).pow(web3.utils.toBN(2))), "Delay factor expected to be 1 * 1e16 ^ 2.");

    let groupMemberRewards = await operatorContract.getGroupMemberRewards(group);
    assert.isTrue(web3.utils.toBN(groupMemberRewards).eq(web3.utils.toBN(expectedGroupMemberReward)), "Unexpected group member reward.");
  });

  it("should send part of the group reward to request subsidy pool based on the submission block.", async function() {
    // Example rewards breakdown:
    // entryVerificationGasEstimate: 1240000
    // groupCreationGasEstimate: 2260000
    // dkgContributionMargin: 10%
    // groupMemberBaseReward: 1050000000000000
    // groupSize: 5
    // entry fee estimate: 49230000000000000 wei
    // signing fee: 37200000000000000 wei
    // DKG fee: 6780000000000000 wei
    // relayEntryTimeout: 10 blocks
    // currentEntryStartBlock: 38
    // relay entry submission block: 44
    // decimals: 1e16
    // groupProfitFee: 42450000000000000 - 37200000000000000 = 5250000000000000 wei
    // memberBaseReward: 5250000000000000 / 5 = 1050000000000000 wei
    // entryTimeout: 38 + 10 = 48
    // delayFactor: ((48 - 44) * 1e16 / (10 - 1)) ^ 2 = 19753086419753082469135802469136
    // groupMemberDelayPenalty: 1050000000000000 * 80246913580246917530864197530864 / (1e16 ^ 2) = 842592592592592
    // groupMemberReward: 1050000000000000 * 19753086419753082469135802469136) / (1e16 ^ 2) = 207407407407407 wei
    // submitterExtraReward: 842592592592592 * 5 * 5 / 100 = 210648148148148 wei
    // submitterReward: 37200000000000000 + 210648148148148 = 37410648148148148 wei
    // subsidy = 5250000000000000 - 207407407407407 * 5 - 210648148148148 = 4002314814814817 wei
  
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
    let tx = await serviceContract.methods['requestRelayEntry(address,string,uint256)'](
      callbackContract.address,
      "callback(uint256)",
      0,
      {value: entryFeeEstimate, from: requestor}
    );

    let currentEntryStartBlock = web3.utils.toBN(tx.receipt.blockNumber);
    let relayEntryTimeout = await operatorContract.relayEntryTimeout();
    let relayEntryGenerationTime = await operatorContract.relayEntryGenerationTime();
    let deadlineBlock = currentEntryStartBlock.add(relayEntryTimeout).addn(1);
    let submissionStartBlock = currentEntryStartBlock.add(relayEntryGenerationTime).add(web3.utils.toBN(1));
    let decimalPoints = web3.utils.toBN(1e16);

    mineBlocks(relayEntryGenerationTime.toNumber() + 1);

    let entryReceivedBlock = web3.utils.toBN(await web3.eth.getBlockNumber()).add(web3.utils.toBN(1)); // web3.eth.getBlockNumber is 1 block behind solidity 'block.number'.
    let remainingBlocks = deadlineBlock.sub(entryReceivedBlock);
    let submissionWindow = deadlineBlock.sub(submissionStartBlock);
    let delayFactor = (remainingBlocks.mul(decimalPoints).div(submissionWindow)).pow(web3.utils.toBN(2));

    let memberBaseReward = entryFee.groupProfitFee.div(groupSize)
    let expectedGroupMemberReward = memberBaseReward.mul(delayFactor).div(decimalPoints.pow(web3.utils.toBN(2)));
    let expectedDelayPenalty = memberBaseReward.sub(memberBaseReward.mul(delayFactor).div(decimalPoints.pow(web3.utils.toBN(2))));
    let expectedSubmitterExtraReward = expectedDelayPenalty.mul(groupSize).mul(web3.utils.toBN(5)).div(web3.utils.toBN(100));
    let requestSubsidy = entryFee.groupProfitFee.sub(expectedGroupMemberReward.mul(groupSize)).sub(expectedSubmitterExtraReward);

    let serviceContractBalance = web3.utils.toBN(await web3.eth.getBalance(serviceContract.address));

    await operatorContract.relayEntry(bls.groupSignature);

    let groupMemberRewards = await operatorContract.getGroupMemberRewards(group);
    assert.isTrue(groupMemberRewards.eq(expectedGroupMemberReward), "Unexpected group member reward.");
    assert.isTrue(serviceContractBalance.add(requestSubsidy).eq(web3.utils.toBN(await web3.eth.getBalance(serviceContract.address))), "Service contract should receive request subsidy.");
  });
});
