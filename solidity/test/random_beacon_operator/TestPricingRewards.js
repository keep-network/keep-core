const blsData = require("../helpers/data.js")
const {initContracts} = require('../helpers/initContracts')
const {createSnapshot, restoreSnapshot} = require("../helpers/snapshot.js")
const {contract, accounts, web3} = require("@openzeppelin/test-environment")
const {time} = require("@openzeppelin/test-helpers")

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

describe('KeepRandomBeaconOperator/PricingRewards', function() {
  let serviceContract;
  let operatorContract;

  before(async () => {
    let contracts = await initContracts(
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
    expect(delayFactor).to.eq.BN('10000000000000000') 
  });

  it("should correctly evaluate delay factor at the first submission block", async () => {
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: accounts[0]});

    await time.advanceBlockTo(web3.utils.toBN(await web3.eth.getBlockNumber()).addn(1));

    let delayFactor = await operatorContract.delayFactor.call();
    expect(delayFactor).to.eq.BN('10000000000000000') 
  });

  it("should correctly evaluate delay factor at the second submission block", async () => {
    let startBlock = await operatorContract.currentRequestStartBlock()
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: accounts[0]});

    await time.advanceBlockTo(web3.utils.toBN(await web3.eth.getBlockNumber()).addn(2));

    let delayFactor = await operatorContract.delayFactor.call();
    // currentRequestStartBlock = 0
    // T_received = 2
    // T_deadline = 0 + 384 + 1 = 385
    // T_begin = 0 + 1 = 1
    // [(T_deadline - T_received) / (T_deadline - T_begin)]^2 = [(385 - 2) / (385 - 1)]^2
    expect(delayFactor).to.eq.BN('9947984483506943')                                  
  });

  it("should correctly evaluate delay factor in the last block before timeout", async () => {
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: accounts[0]});

    let relayEntryTimeout = await operatorContract.relayEntryTimeout();
    await time.advanceBlockTo(relayEntryTimeout.addn(await web3.eth.getBlockNumber()));

    let delayFactor = await operatorContract.delayFactor.call();        
    // currentRequestStartBlock = 0
    // T_received = 384
    // T_deadline = 0 + 384 + 1 = 385
    // T_begin = 0 + 1 = 1
    // [(T_deadline - T_received) / (T_deadline - T_begin)]^2 = [(385 - 384) / (385 - 1)]^2
    expect(delayFactor).to.eq.BN('67816840277')       
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

    expect(rewards.groupMemberReward).to.eq.BN(expectedGroupMemberReward)
    expect(rewards.submitterReward).to.eq.BN(expectedSubmitterReward)
    expect(rewards.subsidy).to.eq.BN(expectedSubsidy)
  });

  it("should correctly evaluate rewards for entry submitted " +
     "at the first submission block", async() => {
    await operatorContract.setGroupMemberBaseReward(966);
    await operatorContract.setEntryVerificationGasEstimate(10050);
    await operatorContract.setGasPriceCeiling(150000, {from: accounts[0]});

    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: accounts[0]});  

    await time.advanceBlockTo(web3.utils.toBN(await web3.eth.getBlockNumber()).addn(1));

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
    
    expect(rewards.groupMemberReward).to.eq.BN(expectedGroupMemberReward)
    expect(rewards.submitterReward).to.eq.BN(expectedSubmitterReward)
    expect(rewards.subsidy).to.eq.BN(expectedSubsidy)
  });  

  it("should correctly evaluate rewards for the entry submitted " + 
     "at the second submission block", async () => {
    await operatorContract.setGroupMemberBaseReward(1987000);
    await operatorContract.setEntryVerificationGasEstimate(50050);
    await operatorContract.setGasPriceCeiling(1400000, {from: accounts[0]});

    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: accounts[0]});  

    await time.advanceBlockTo(web3.utils.toBN(await web3.eth.getBlockNumber()).addn(2));

    // There is one block of delay so the delay factor is 0.9947984483506943.
    // Group member reward should be scaled by the delay factor: 
    // 1987000 * 0.9947984483506943 = ~1966355
    let expectedGroupMemberReward = web3.utils.toBN("1976664");

    // The entire entry verification fee is paid to the submitter 
    // regardless of their gas expenditure. The submitter is free to spend 
    // less or more, keeping the surplus or paying the difference:
    // 50050 * 1400000 = 70070000000
    // 
    // To incentivize a race for the submitter position, the submitter 
    // receives delay penalty * group size * 0.05 as an extra reward:
    // 1987000 * (1 - 0.9947984483506943) * 64 * 5% = ~33073
    //
    // 70070000000 + 33073 = 70070033073          
    let expectedSubmitterReward = web3.utils.toBN("70070033073");  

    // If the amount paid out to the signing group in group rewards and the 
    // submitter’s extra reward is less than the profit margin, the 
    // difference is added to the beacon’s request subsidy pool to 
    // incentivize customers to request entries.
    //
    // profit margin: 1987000 * 64 = 127168000
    // paid member rewards: 1976664 * 64 = 126506496
    // submitter extra reward: 33073
    // 
    // 127168000 - 126506496 - 33073 = 628431
    let expectedSubsidy = web3.utils.toBN("628431");              

    let rewards = await operatorContract.getNewEntryRewardsBreakdown.call(); 
    
    expect(rewards.groupMemberReward).to.eq.BN(expectedGroupMemberReward)
    expect(rewards.submitterReward).to.eq.BN(expectedSubmitterReward)
    expect(rewards.subsidy).to.eq.BN(expectedSubsidy)
  });  

  it("should correctly evaluate rewards for the entry submitted " + 
     "in the last block before timeout", async () => {
    await operatorContract.setGroupMemberBaseReward(1382000000);
    await operatorContract.setEntryVerificationGasEstimate(50020);
    await operatorContract.setGasPriceCeiling(2000000, {from: accounts[0]});

    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: accounts[0]}); 

    let relayEntryTimeout = await operatorContract.relayEntryTimeout();
    await time.advanceBlockTo(relayEntryTimeout.addn(await web3.eth.getBlockNumber()));

    // There is one block left before the timeout so the delay factor is 
    // 0.000067816840277.
    // Group member reward should be scaled by the delay factor: 
    // 1382000000 * 0.0000067816840277 = ~9372
    let expectedGroupMemberReward = web3.utils.toBN("9372");  

    // The entire entry verification fee is paid to the submitter 
    // regardless of their gas expenditure. The submitter is free to spend 
    // less or more, keeping the surplus or paying the difference:
    // 50020 * 2000000 = 100040000000
    // 
    // To incentivize a race for the submitter position, the submitter 
    // receives delay penalty * group size * 0.05 as an extra reward:
    // 1382000000 * (1 - 0.0000067816840277) * 64 * 5% = ~4422370008
    //
    // 100040000000 + 4422370008 = 104462370008          
    let expectedSubmitterReward = web3.utils.toBN("104462370008"); 

    // If the amount paid out to the signing group in group rewards and the 
    // submitter’s extra reward is less than the profit margin, the 
    // difference is added to the beacon’s request subsidy pool to 
    // incentivize customers to request entries.
    //
    // profit margin: 1382000000 * 64 = 88448000000
    // paid member rewards: 9372 * 64 = 599808
    // submitter extra reward: 4422280034
    // 
    // 88448000000 - 599808 - 4422370008 = 84025030184
    let expectedSubsidy = web3.utils.toBN("84025030184");    

    let rewards = await operatorContract.getNewEntryRewardsBreakdown.call();
    
    expect(rewards.groupMemberReward).to.eq.BN(expectedGroupMemberReward)
    expect(rewards.submitterReward).to.eq.BN(expectedSubmitterReward)
    expect(rewards.subsidy).to.eq.BN(expectedSubsidy)
  });
});
