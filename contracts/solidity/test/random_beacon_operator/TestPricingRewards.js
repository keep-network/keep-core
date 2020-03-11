import {initContracts} from '../helpers/initContracts';
import {createSnapshot, restoreSnapshot} from '../helpers/snapshot';
import {bls} from '../helpers/data';
import mineBlocks from '../helpers/mineBlocks';

contract('KeepRandomBeaconOperator', function(accounts) {
  let serviceContract;
  let operatorContract;

  before(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./stubs/KeepRandomBeaconOperatorPricingStub.sol')
    );
    
    serviceContract = contracts.serviceContract;
    operatorContract = contracts.operatorContract;

    await operatorContract.registerNewGroup(bls.groupPubKey);
  });

  beforeEach(async () => {
    await createSnapshot()
  });
    
  afterEach(async () => {
    await restoreSnapshot()
  });

  it("should correctly evaluate delay factor right after the request", async () => {
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate});

    let delayFactor = await operatorContract.delayFactor();        

    let expectedDelayFactor = web3.utils.toBN(10000000000000000);
    assert.isTrue(expectedDelayFactor.eq(delayFactor));
  });

  it("should correctly evaluate delay factor at the first submission block", async () => {
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate});

    mineBlocks(await operatorContract.relayEntryGenerationTime());

    let delayFactor = await operatorContract.delayFactor();

    let expectedDelayFactor = web3.utils.toBN(10000000000000000);
    assert.isTrue(expectedDelayFactor.eq(delayFactor));
  });

  it("should correctly evaluate delay factor at the second submission block", async () => {
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate});

    mineBlocks((await operatorContract.relayEntryGenerationTime()).addn(1));

    let delayFactor = await operatorContract.delayFactor();

    let expectedDelayFactor = web3.utils.toBN('9896104600694443');
    assert.isTrue(expectedDelayFactor.eq(delayFactor));
  });

  it("should correctly evaluate delay factor in the last block before timeout", async () => {
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate});

    mineBlocks((await operatorContract.relayEntryTimeout()).subn(1));

    let delayFactor = await operatorContract.delayFactor();        

    let expectedDelayFactor = web3.utils.toBN('271267361111');
    assert.isTrue(expectedDelayFactor.eq(delayFactor));        
  });

  it("should correctly evaluate rewards for entry submitted " + 
     "right after the request", async () => {
    await operatorContract.setGroupMemberBaseReward(1410);
    await operatorContract.setEntryVerificationGasEstimate(10020);
    await serviceContract.setPriceFeedEstimate(140000);
    await operatorContract.setPriceFeedEstimate(140000);

    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate});

    // No delay so entire group member base reward is paid and nothing
    // goes to the subsidy pool.
    let expectedGroupMemberReward = web3.utils.toBN("1410");
    let expectedSubsidy = web3.utils.toBN("0");
    // The entire entry verification fee is paid to the submitter 
    // regardless of their gas expenditure. The submitter is free to spend 
    // less or more, keeping the surplus or paying the difference:
    // 10020 * 140000 * 150% (fluctuation margin)
    let expectedSubmitterReward = web3.utils.toBN("2104200000");
    
    let rewards = await operatorContract.getNewEntryRewardsBreakdown(); 

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
    await serviceContract.setPriceFeedEstimate(150000);
    await operatorContract.setPriceFeedEstimate(150000);  

    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate});  

    mineBlocks(await operatorContract.relayEntryGenerationTime()); 

    // No delay so entire group member base reward is paid and nothing
    // goes to the subsidy pool.
    let expectedGroupMemberReward = web3.utils.toBN("966");
    let expectedSubsidy = web3.utils.toBN("0");
    // The entire entry verification fee is paid to the submitter 
    // regardless of their gas expenditure. The submitter is free to spend 
    // less or more, keeping the surplus or paying the difference:
    // 10050 * 150000 * 150% (fluctuation margin)
    let expectedSubmitterReward = web3.utils.toBN("2261250000");
    
    let rewards = await operatorContract.getNewEntryRewardsBreakdown(); 
    
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
    await serviceContract.setPriceFeedEstimate(1400000);
    await operatorContract.setPriceFeedEstimate(1400000);  

    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate});  

    mineBlocks((await operatorContract.relayEntryGenerationTime()).addn(1));  

    // There is one block of delay so the delay factor is 0.9896104600694443.
    // Group member reward should be scaled by the delay factor: 
    // 1987000 * 0.9896104600694443 = ~1966355
    let expectedGroupMemberReward = web3.utils.toBN("1966355");

    // The entire entry verification fee is paid to the submitter 
    // regardless of their gas expenditure. The submitter is free to spend 
    // less or more, keeping the surplus or paying the difference:
    // 50050 * 1400000 * 150% (fluctuation margin) = 105105000000
    // 
    // To incentivize a race for the submitter position, the submitter 
    // receives delay penalty * group size * 0.05 as an extra reward:
    // 1987000 * (1 - 0.9896104600694443) * 64 * 5% = ~66060
    //
    // 105105000000 + 66060 = 105105066060          
    let expectedSubmitterReward = web3.utils.toBN("105105066060");  

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

    let rewards = await operatorContract.getNewEntryRewardsBreakdown(); 
    
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
    await serviceContract.setPriceFeedEstimate(2000000);
    await operatorContract.setPriceFeedEstimate(2000000);  

    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate}); 

    mineBlocks((await operatorContract.relayEntryTimeout()).subn(1));

    // There is one block left before the timeout so the delay factor is 
    // 0.0000271267361111.
    // Group member reward should be scaled by the delay factor: 
    // 1382000000 * 0.0000271267361111 = ~37489
    let expectedGroupMemberReward = web3.utils.toBN("37489");  

    // The entire entry verification fee is paid to the submitter 
    // regardless of their gas expenditure. The submitter is free to spend 
    // less or more, keeping the surplus or paying the difference:
    // 50020 * 2000000 * 150% (fluctuation margin) = 150060000000
    // 
    // To incentivize a race for the submitter position, the submitter 
    // receives delay penalty * group size * 0.05 as an extra reward:
    // 1382000000 * (1 - 0.0000271267361111) * 64 * 5% = ~4422280034
    //
    // 150060000000 + 4422280034 = 154482280034          
    let expectedSubmitterReward = web3.utils.toBN("154482280034"); 

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

    let rewards = await operatorContract.getNewEntryRewardsBreakdown();
    
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
