import { duration } from './helpers/increaseTime';
import exceptThrow from './helpers/expectThrow';
import {bls} from './helpers/data';
import { AssertionError } from 'assert';
const KeepToken = artifacts.require('./KeepToken.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const KeepRandomBeaconProxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconImplV1 = artifacts.require('./KeepRandomBeaconImplV1.sol');
const KeepGroupProxy = artifacts.require('./KeepGroup.sol');
const KeepGroupImplV1 = artifacts.require('./KeepGroupImplV1.sol');

function generateTickets(randomBeaconValue, stakerValue, stakerWeight) {
  let tickets = [];
  for (let i = 1; i <= stakerWeight; i++) {
    let ticketValue = web3.utils.toBN(
      web3.utils.soliditySha3({t: 'uint', v: randomBeaconValue}, {t: 'uint', v: stakerValue}, {t: 'uint', v: i})
    );
    let ticket = {
      value: ticketValue,
      virtualStakerIndex: i
    }
    tickets.push(ticket);
  }
  return tickets
}

contract('TestPublishDkgResult', function(accounts) {
  let  disqualified, inactive, resultHash,
  signature, positions, token, stakingProxy,
  stakingContract, minimumStake, groupThreshold, groupSize,
  randomBeaconValue, requestId,
  timeoutInitial, timeoutSubmission, timeoutChallenge,
  keepRandomBeaconImplV1, keepRandomBeaconProxy, keepRandomBeaconImplViaProxy,
  keepGroupImplV1, keepGroupProxy, keepGroupImplViaProxy, groupPubKey,
  staker1 = accounts[0], tickets1,
  staker2 = accounts[1], tickets2,
  staker3 = accounts[2], tickets3;
  requestId = 0;
  disqualified = '0x00000001100000001100000001100000001100000000000110000000110000000110000000110000000000011000000011000000011000000011000000000001100000001100000001100000001100000000'
  inactive =  '0x00000010000000010000000010000000010000000000001000000001000000001000000001000000000000100000000100000000100000000100000000000010000000010000000010000000010000000000'
  groupPubKey = "0x1000000000000000000000000000000000000000000000000000000000000000"

  
  resultHash = web3.utils.soliditySha3(disqualified, inactive, groupPubKey);
  
  beforeEach(async () => {
    token = await KeepToken.new();
    
    // Initialize staking contract under proxy
    stakingProxy = await StakingProxy.new();
    stakingContract = await TokenStaking.new(token.address, stakingProxy.address, duration.days(30));
    await stakingProxy.authorizeContract(stakingContract.address, {from: staker1})

    // Initialize Keep Random Beacon contract
    keepRandomBeaconImplV1 = await KeepRandomBeaconImplV1.new();
    keepRandomBeaconProxy = await KeepRandomBeaconProxy.new(keepRandomBeaconImplV1.address);
    keepRandomBeaconImplViaProxy = await KeepRandomBeaconImplV1.at(keepRandomBeaconProxy.address);

    // Initialize Keep Group contract
    minimumStake = 200000;
    groupThreshold = 15;
    groupSize = 75;
    timeoutInitial = 20;
    timeoutSubmission = 100;
    timeoutChallenge = 60;

    randomBeaconValue = bls.groupSignature;

    keepGroupImplV1 = await KeepGroupImplV1.new();
    keepGroupProxy = await KeepGroupProxy.new(keepGroupImplV1.address);
    keepGroupImplViaProxy = await KeepGroupImplV1.at(keepGroupProxy.address);
    await keepGroupImplViaProxy.initialize(
      stakingProxy.address, keepRandomBeaconProxy.address, minimumStake, groupThreshold, groupSize, timeoutInitial, timeoutSubmission, timeoutChallenge
    );

    await keepRandomBeaconImplViaProxy.initialize(1,1, randomBeaconValue, bls.groupPubKey, keepGroupProxy.address);
    await keepRandomBeaconImplViaProxy.relayEntry(1, bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

    // Stake tokens as account one so it has minimum stake to be able to get into a group.
    await token.approveAndCall(stakingContract.address, minimumStake*2000, "0x00", {from: staker1});
    tickets1 = generateTickets(randomBeaconValue, staker1, 2000);

    // Send tokens to staker2 and stake
    await token.transfer(staker2, minimumStake*2000, {from: staker1});
    await token.approveAndCall(stakingContract.address, minimumStake*2000, "0x00", {from: staker2});
    tickets2 = generateTickets(randomBeaconValue, staker2, 2000);

    // Send tokens to staker3 and stake
    await token.transfer(staker3, minimumStake*3000, {from: staker1});
    await token.approveAndCall(stakingContract.address, minimumStake*3000, "0x00", {from: staker3});
    tickets3 = generateTickets(randomBeaconValue, staker3, 3000);    
  })
  it("should generate signatures and submit a correct result", async function() {

    positions = []
    let signatures; 
    let callerIndex;

    for(let i=0;i<5;i++){
      await keepGroupImplViaProxy.submitTicket(tickets1[i].value, staker1, tickets1[i].virtualStakerIndex, {from: staker1});
    }
    for(let i=0;i<2;i++){
      await keepGroupImplViaProxy.submitTicket(tickets2[i].value, staker2, tickets2[i].virtualStakerIndex, {from: staker2});
    }
    for(let i=0;i<1;i++){
      await keepGroupImplViaProxy.submitTicket(tickets3[i].value, staker3, tickets3[i].virtualStakerIndex, {from: staker3});
    }
    let orderedParticipants = await keepGroupImplViaProxy.orderedParticipants.call()

    for(let i=0;i<orderedParticipants.length;i++){    
      callerIndex = accounts.indexOf(orderedParticipants[i]);    
      signature =  await web3.eth.sign(resultHash, accounts[callerIndex]);
      positions.push(i+1);
      if (signatures == undefined) signatures = signature
      else signatures+=signature.slice(2,signature.length);
    }
    
    await keepGroupImplViaProxy.submitDkgResult(requestId, 1, groupPubKey, disqualified, inactive, signatures, positions)
    let submitted = await keepGroupImplViaProxy.isDkgResultSubmitted.call(requestId);
    assert.equal(submitted, true, "DkgResult should should be submitted");
  });
})  