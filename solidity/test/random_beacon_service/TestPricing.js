const stakeDelegate = require('../helpers/stakeDelegate')
const blsData = require("../helpers/data.js")
const {initContracts} = require('../helpers/initContracts')
const assert = require('chai').assert
const {contract, accounts, web3} = require("@openzeppelin/test-environment")
const {time} = require("@openzeppelin/test-helpers")
const CallbackContract = contract.fromArtifact('CallbackContract')

describe('TestKeepRandomBeaconService/Pricing', function() {

  let token, stakingContract, operatorContract, serviceContract, callbackContract, entryFee, groupSize, group,
    owner = accounts[0],
    requestor = accounts[1],
    operator1 = accounts[2],
    operator2 = accounts[3],
    operator3 = accounts[4],
    beneficiary1 = accounts[5],
    beneficiary2 = accounts[6],
    beneficiary3 = accounts[7];

  beforeEach(async () => {
    let contracts = await initContracts(
      contract.fromArtifact('TokenStaking'),
      contract.fromArtifact('KeepRandomBeaconService'),
      contract.fromArtifact('KeepRandomBeaconServiceImplV1'),
      contract.fromArtifact('KeepRandomBeaconOperatorServicePricingStub')
    );

    token = contracts.token;
    stakingContract = contracts.stakingContract;
    operatorContract = contracts.operatorContract;
    serviceContract = contracts.serviceContract;
    callbackContract = await CallbackContract.new();

    // Using stub method to add first group to help testing.
    await operatorContract.registerNewGroup(blsData.groupPubKey);

    groupSize = web3.utils.toBN(3);
    group = await operatorContract.getGroupPublicKey(0);
    await operatorContract.setGroupMembers(group, [operator1, operator2, operator3])
    let minimumStake = await stakingContract.minimumStake()

    await stakeDelegate(stakingContract, token, owner, operator1, beneficiary1, operator1, minimumStake);
    await stakeDelegate(stakingContract, token, owner, operator2, beneficiary2, operator2, minimumStake);
    await stakeDelegate(stakingContract, token, owner, operator3, beneficiary3, operator3, minimumStake);

    entryFee = await serviceContract.entryFeeBreakdown()
  });

  it("should successfully refund callback surplus for a lower submission gas price", async () => {
    let gasPriceCeiling = web3.utils.toBN(web3.utils.toWei('20', 'gwei'))
    await operatorContract.setGasPriceCeiling(gasPriceCeiling)

    let callbackGas = web3.utils.toBN(await callbackContract.__beaconCallback.estimateGas(blsData.groupSignature))
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    
    await serviceContract.methods['requestRelayEntry(address,uint256)'](
      callbackContract.address,
      callbackGas,
      {value: entryFeeEstimate, from: requestor}
    );

    let submissionGasPrice = web3.utils.toBN(web3.utils.toWei('5', 'gwei'))
    let gasPriceDiff = gasPriceCeiling.sub(submissionGasPrice)

    let requestorBalance = await web3.eth.getBalance(requestor);
    await operatorContract.relayEntry(blsData.groupSignature, {gasPrice: submissionGasPrice})
    let updatedRequestorBalance = await web3.eth.getBalance(requestor)

    let refund = web3.utils.toBN(updatedRequestorBalance).sub(web3.utils.toBN(requestorBalance))

    let baseCallbackGas = await serviceContract.baseCallbackGas()
    let expectedSurplus = (callbackGas.add(baseCallbackGas)).mul(gasPriceDiff)
    
    assert.isTrue(expectedSurplus.eq(refund), "Callback gas surplus should be refunded to the requestor.");
  })

  it("should send group reward to each operator.", async function() {
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
    let tx = await serviceContract.methods['requestRelayEntry(address,uint256)'](
      callbackContract.address,
      0,
      {value: entryFeeEstimate, from: requestor}
    );

    let currentRequestStartBlock = web3.utils.toBN(tx.receipt.blockNumber);
    let relayEntryTimeout = await operatorContract.relayEntryTimeout();
    let deadlineBlock = currentRequestStartBlock.add(relayEntryTimeout);
    let entryReceivedBlock = currentRequestStartBlock.addn(1);
    let remainingBlocks = deadlineBlock.sub(entryReceivedBlock);
    let submissionWindow = deadlineBlock.sub(entryReceivedBlock);
    let decimalPoints = web3.utils.toBN(1e16);
    let delayFactor = (remainingBlocks.mul(decimalPoints).div(submissionWindow)).pow(web3.utils.toBN(2));
    let memberBaseReward = entryFee.groupProfitFee.div(groupSize)
    let expectedGroupMemberReward = memberBaseReward.mul(delayFactor).div(decimalPoints.pow(web3.utils.toBN(2)));

    await operatorContract.relayEntry(blsData.groupSignature);

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
    // currentRequestStartBlock: 38
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
    let tx = await serviceContract.methods['requestRelayEntry(address,uint256)'](
      callbackContract.address,
      0,
      {value: entryFeeEstimate, from: requestor}
    );

    let currentRequestStartBlock = web3.utils.toBN(tx.receipt.blockNumber);
    let relayEntryTimeout = await operatorContract.relayEntryTimeout();
    let deadlineBlock = currentRequestStartBlock.add(relayEntryTimeout).addn(1);
    let submissionStartBlock = currentRequestStartBlock.addn(1);
    let decimalPoints = web3.utils.toBN(1e16);

    await time.advanceBlockTo(web3.utils.toBN(await web3.eth.getBlockNumber()).addn(1));

    let entryReceivedBlock = web3.utils.toBN(await web3.eth.getBlockNumber()).add(web3.utils.toBN(1)); // web3.eth.getBlockNumber is 1 block behind solidity 'block.number'.
    let remainingBlocks = deadlineBlock.sub(entryReceivedBlock);
    let submissionWindow = deadlineBlock.sub(submissionStartBlock);
    let delayFactor = (remainingBlocks.mul(decimalPoints).div(submissionWindow)).pow(web3.utils.toBN(2));

    let memberBaseReward = entryFee.groupProfitFee.div(groupSize)
    let expectedGroupMemberReward = memberBaseReward.mul(delayFactor).div(decimalPoints.pow(web3.utils.toBN(2)));
    let expectedDelayPenalty = memberBaseReward.sub(memberBaseReward.mul(delayFactor).div(decimalPoints.pow(web3.utils.toBN(2))));
    let expectedSubmitterExtraReward = expectedDelayPenalty.mul(groupSize).muln(5).div(web3.utils.toBN(100));
    let requestSubsidy = entryFee.groupProfitFee.sub(expectedGroupMemberReward.mul(groupSize)).sub(expectedSubmitterExtraReward);

    let serviceContractBalance = web3.utils.toBN(await web3.eth.getBalance(serviceContract.address));

    await operatorContract.relayEntry(blsData.groupSignature);

    let groupMemberRewards = await operatorContract.getGroupMemberRewards(group);
    assert.isTrue(groupMemberRewards.eq(expectedGroupMemberReward), "Unexpected group member reward.");
    assert.isTrue(serviceContractBalance.add(requestSubsidy).eq(web3.utils.toBN(await web3.eth.getBalance(serviceContract.address))), "Service contract should receive request subsidy.");
  });
});
