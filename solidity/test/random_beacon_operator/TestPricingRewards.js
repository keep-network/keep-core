const blsData = require("../helpers/data.js")
const initContracts = require('../helpers/initContracts')
const assert = require('chai').assert
const mineBlocks = require("../helpers/mineBlocks")
const {createSnapshot, restoreSnapshot} = require("../helpers/snapshot.js")
const {contract, accounts, web3} = require("@openzeppelin/test-environment")

describe('KeepRandomBeaconOperator/PricingRewards', function() {
  let serviceContract;
  let operatorContract;

  before(async () => {
    let contracts = await initContracts(
      contract.fromArtifact('KeepToken'),
      contract.fromArtifact('TokenStaking'),
      contract.fromArtifact('KeepRandomBeaconService'),
      contract.fromArtifact('KeepRandomBeaconServiceImplV1'),
      contract.fromArtifact('KeepRandomBeaconOperatorPricingStub')
    );
    
    serviceContract = contracts.serviceContract;
    operatorContract = contracts.operatorContract;

    await operatorContract.registerNewGroup(blsData.groupPubKey);
  });

  beforeEach(async () => {
    await createSnapshot()
  });
    
  afterEach(async () => {
    await restoreSnapshot()
  });

  it("should correctly evaluate delay factor right after the request", async () => {
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: accounts[0]});

    let delayFactor = await operatorContract.delayFactor.call();        

    let expectedDelayFactor = web3.utils.toBN(10000000000000000);
    assert.isTrue(expectedDelayFactor.eq(delayFactor));
  });

  it("should correctly evaluate delay factor at the first submission block", async () => {
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: accounts[0]});

    mineBlocks(1);

    let delayFactor = await operatorContract.delayFactor.call();

    let expectedDelayFactor = web3.utils.toBN(10000000000000000);
    assert.isTrue(expectedDelayFactor.eq(delayFactor));
  });

  it("should correctly evaluate delay factor at the second submission block", async () => {
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: accounts[0]});

    mineBlocks(2);

    let delayFactor = await operatorContract.delayFactor.call();

    let expectedDelayFactor = web3.utils.toBN('9896104600694443');
    assert.isTrue(expectedDelayFactor.eq(delayFactor));
  });

  it("should correctly evaluate delay factor in the last block before timeout", async () => {
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: accounts[0]});

    mineBlocks(await operatorContract.relayEntryTimeout());

    let delayFactor = await operatorContract.delayFactor.call();        

    let expectedDelayFactor = web3.utils.toBN('271267361111');
    assert.isTrue(expectedDelayFactor.eq(delayFactor));        
  });

  it("should correctly evaluate rewards for entry submitted " + 
     "right after the request", async () => {
    await operatorContract.setGroupMemberBaseReward(1410);
    await operatorContract.setEntryVerificationGasEstimate(10020);
    await operatorContract.setGasPriceCeiling(140000, {from: accounts[0]});

    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: accounts[0]});

    // No delay so entire group member base reward is paid and nothing
    // goes to the subsidy pool.
    let expectedGroupMemberReward = web3.utils.toBN("1410");
    let expectedSubsidy = web3.utils.toBN("0");
    // The entire entry verification fee is paid to the submitter 
    // regardless of their gas expenditure. The submitter is free to spend 
    // less or more, keeping the surplus or paying the difference:
    // 10020 * 140000
    let expectedSubmitterReward = web3.utils.toBN("1402800000");
    
    let rewards = await operatorContract.getNewEntryRewardsBreakdown.call(); 

    assert.isTrue(
      expectedGroupMemberReward.eq(rewards.groupMemberReward),
      "unexpected group member reward"
    );
    assert.isTrue(
      expectedSubmitterReward.eq(rewards.submitterReward),
      "unexpected submitter reward"
    );
    assert.isTrue(
      expectedSubsidy.eq(rewards.subsidy),
      "unexpected subsidy"
    );
  });

  it("should correctly evaluate rewards for entry submitted " +
     "at the first submission block", async() => {
    await operatorContract.setGroupMemberBaseReward(966);
    await operatorContract.setEntryVerificationGasEstimate(10050);
    await operatorContract.setGasPriceCeiling(150000, {from: accounts[0]});

    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: accounts[0]});  

    mineBlocks(1);

    // No delay so entire group member base reward is paid and nothing
    // goes to the subsidy pool.
    let expectedGroupMemberReward = web3.utils.toBN("966");
    let expectedSubsidy = web3.utils.toBN("0");
    // The entire entry verification fee is paid to the submitter 
    // regardless of their gas expenditure. The submitter is free to spend 
    // less or more, keeping the surplus or paying the difference:
    // 10050 * 150000
    let expectedSubmitterReward = web3.utils.toBN("1507500000");
    
    let rewards = await operatorContract.getNewEntryRewardsBreakdown.call(); 
    
    assert.isTrue(
      expectedGroupMemberReward.eq(rewards.groupMemberReward),
      "unexpected group member reward"
    );
    assert.isTrue(
      expectedSubmitterReward.eq(rewards.submitterReward),
      "unexpected submitter reward"
    );
    assert.isTrue(
      expectedSubsidy.eq(rewards.subsidy),
      "unexpected subsidy"
    );
  });  

  it("should correctly evaluate rewards for the entry submitted " + 
     "at the second submission block", async () => {
    await operatorContract.setGroupMemberBaseReward(1987000);
    await operatorContract.setEntryVerificationGasEstimate(50050);
    await operatorContract.setGasPriceCeiling(1400000, {from: accounts[0]});

    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: accounts[0]});  

    mineBlocks(2); 

    // There is one block of delay so the delay factor is 0.9896104600694443.
    // Group member reward should be scaled by the delay factor: 
    // 1987000 * 0.9896104600694443 = ~1966355
    let expectedGroupMemberReward = web3.utils.toBN("1966355");

    // The entire entry verification fee is paid to the submitter 
    // regardless of their gas expenditure. The submitter is free to spend 
    // less or more, keeping the surplus or paying the difference:
    // 50050 * 1400000 = 70070000000
    // 
    // To incentivize a race for the submitter position, the submitter 
    // receives delay penalty * group size * 0.05 as an extra reward:
    // 1987000 * (1 - 0.9896104600694443) * 64 * 5% = ~66060
    //
    // 70070000000 + 66060 = 70070066060          
    let expectedSubmitterReward = web3.utils.toBN("70070066060");  

    // If the amount paid out to the signing group in group rewards and the 
    // submitter’s extra reward is less than the profit margin, the 
    // difference is added to the beacon’s request subsidy pool to 
    // incentivize customers to request entries.
    //
    // profit margin: 1987000 * 64 = 127168000
    // paid member rewards: 1966355 * 64 = 125846720
    // submitter extra reward: 66060
    // 
    // 127168000 - 125846720 - 66060 = 1255220
    let expectedSubsidy = web3.utils.toBN("1255220");              

    let rewards = await operatorContract.getNewEntryRewardsBreakdown.call(); 
    
    assert.isTrue(
      expectedGroupMemberReward.eq(rewards.groupMemberReward),
      "unexpected group member reward"
    );
    assert.isTrue(
      expectedSubmitterReward.eq(rewards.submitterReward),
      "unexpected submitter reward"
    );
    assert.isTrue(
      expectedSubsidy.eq(rewards.subsidy),
      "unexpected subsidy"
    );
  });  

  it("should correctly evaluate rewards for the entry submitted " + 
     "in the last block before timeout", async () => {
    await operatorContract.setGroupMemberBaseReward(1382000000);
    await operatorContract.setEntryVerificationGasEstimate(50020);
    await operatorContract.setGasPriceCeiling(2000000, {from: accounts[0]});

    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: accounts[0]}); 

    mineBlocks(await operatorContract.relayEntryTimeout());

    // There is one block left before the timeout so the delay factor is 
    // 0.0000271267361111.
    // Group member reward should be scaled by the delay factor: 
    // 1382000000 * 0.0000271267361111 = ~37489
    let expectedGroupMemberReward = web3.utils.toBN("37489");  

    // The entire entry verification fee is paid to the submitter 
    // regardless of their gas expenditure. The submitter is free to spend 
    // less or more, keeping the surplus or paying the difference:
    // 50020 * 2000000 = 100040000000
    // 
    // To incentivize a race for the submitter position, the submitter 
    // receives delay penalty * group size * 0.05 as an extra reward:
    // 1382000000 * (1 - 0.0000271267361111) * 64 * 5% = ~4422280034
    //
    // 100040000000 + 4422280034 = 104462280034          
    let expectedSubmitterReward = web3.utils.toBN("104462280034"); 

    // If the amount paid out to the signing group in group rewards and the 
    // submitter’s extra reward is less than the profit margin, the 
    // difference is added to the beacon’s request subsidy pool to 
    // incentivize customers to request entries.
    //
    // profit margin: 1382000000 * 64 = 88448000000
    // paid member rewards: 37489 * 64 = 2399296
    // submitter extra reward: 4422280034
    // 
    // 88448000000 - 2399296 - 4422280034 = 84023320670
    let expectedSubsidy = web3.utils.toBN("84023320670");    

    let rewards = await operatorContract.getNewEntryRewardsBreakdown.call();
    
    assert.isTrue(
      expectedGroupMemberReward.eq(rewards.groupMemberReward),
      "unexpected group member reward"
    );
    assert.isTrue(
      expectedSubmitterReward.eq(rewards.submitterReward),
      "unexpected submitter reward"
    );
    assert.isTrue(
      expectedSubsidy.eq(rewards.subsidy),
      "unexpected subsidy"
    );
  });
});
