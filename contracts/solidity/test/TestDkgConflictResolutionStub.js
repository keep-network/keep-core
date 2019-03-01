import { duration } from './helpers/increaseTime';
import exceptThrow from './helpers/expectThrow';
import BigNumber from 'bignumber.js';
import abi from 'ethereumjs-abi';
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

function mineBlocks(blocks) {
  for (let i = 0; i <= blocks; i++) {
    web3.currentProvider.send({
      jsonrpc: "2.0",
      method: "evm_mine",
      id: 12345
    }, function(err, _) {
      if (err) console.log("Error mining a block.")
    });
  }
}

contract('TestDkgConflictResolution', function(accounts) {

    let ordered, blockNum, submissionStart, DkgSubmissionT,tickets, success,
    disqualified, inactive, token, stakingProxy, stakingContract, 
    minimumStake, groupThreshold, groupSize, randomBeaconValue,
    timeoutInitial, timeoutSubmission, timeoutChallenge,
    keepRandomBeaconImplV1, keepRandomBeaconProxy, keepRandomBeaconImplViaProxy,
    keepGroupImplV1, keepGroupProxy, keepGroupImplViaProxy, groupPubKey,
    staker1 = accounts[0], tickets1,
    staker2 = accounts[1], tickets2,
    staker3 = accounts[2], tickets3

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
    await keepRandomBeaconImplViaProxy.initialize(1,1);

    // Initialize Keep Group contract
    minimumStake = 200000;
    groupThreshold = 15;
    groupSize = 20;
    timeoutInitial = 20;
    timeoutSubmission = 40;
    timeoutChallenge = 60;

    randomBeaconValue = 123456789;

    keepGroupImplV1 = await KeepGroupImplV1.new();
    keepGroupProxy = await KeepGroupProxy.new(keepGroupImplV1.address);
    keepGroupImplViaProxy = await KeepGroupImplV1.at(keepGroupProxy.address);
    await keepGroupImplViaProxy.initialize(
      stakingProxy.address, keepRandomBeaconProxy.address, minimumStake, groupThreshold, groupSize, timeoutInitial, timeoutSubmission, timeoutChallenge
    );

    tickets = []
    success = true
    disqualified = '0x0000000110000000110000000110000000110000'
    inactive =  '0x0000001000000001000000001000000001000000'
    groupPubKey = "0x1000000000000000000000000000000000000000000000000000000000000000";

    // Stake tokens as account one so it has minimum stake to be able to get into a group.
    await token.approveAndCall(stakingContract.address, minimumStake*1000, "0x00", {from: staker1});
    tickets1 = generateTickets(randomBeaconValue, staker1, 1000);

    // Send tokens to staker2 and stake
    await token.transfer(staker2, minimumStake*2000, {from: staker1});
    await token.approveAndCall(stakingContract.address, minimumStake*2000, "0x00", {from: staker2});
    tickets2 = generateTickets(randomBeaconValue, staker2, 2000);

    // Send tokens to staker3 and stake
    await token.transfer(staker3, minimumStake*3000, {from: staker1});
    await token.approveAndCall(stakingContract.address, minimumStake*3000, "0x00", {from: staker3});
    tickets3 = generateTickets(randomBeaconValue, staker3, 3000);

    await keepRandomBeaconImplViaProxy.setGroupContract(keepGroupProxy.address);
    await keepRandomBeaconImplViaProxy.relayEntry(1, randomBeaconValue, "0x01", 1, 1);

    //submit tickets for a groupSize of 20.
    

    //TODO:
    //Remove from here and implement only when multiple group member logic is needed
    
    for(let i=0;i<10;i++){
        await keepGroupImplViaProxy.submitTicket(tickets1[i].value, staker1, tickets1[i].virtualStakerIndex);
        tickets.push(tickets1[i].value);
    }
    for(let i=0;i<7;i++){
        await keepGroupImplViaProxy.submitTicket(tickets2[i].value, staker2, tickets2[i].virtualStakerIndex, {from: staker2});
        tickets.push(tickets2[i].value);
    }
    for(let i=0;i<3;i++){
        await keepGroupImplViaProxy.submitTicket(tickets3[i].value, staker3, tickets3[i].virtualStakerIndex, {from: staker3});
        tickets.push(tickets3[i].value);
    }
    tickets = tickets.sort(function(a, b){return a-b});

    //skip to the first block where ticket challenges is over.
    ordered = await keepGroupImplViaProxy.orderedTickets();
    blockNum =  await web3.eth.getBlockNumber();
    submissionStart = await keepGroupImplViaProxy.ticketSubmissionStartBlock();
    DkgSubmissionT = ((submissionStart.toNumber() + timeoutChallenge) - blockNum)+1
    mineBlocks(DkgSubmissionT);


  });
   it("should fail to submit result if invalid index is provided", async function() {
    let finalParticipants = await keepGroupImplViaProxy.selectedParticipants();
    let GoodSender = finalParticipants[0];
    let badSender;
    for(let i=1;i<finalParticipants.length;i++){
        if(finalParticipants[i] != GoodSender) {
            badSender = finalParticipants[i];
            break;
        }
    } 
    assert.equal(tickets.length, ordered.length, "should have 20 tickets");
    await exceptThrow(keepGroupImplViaProxy.submitDkgResult(0,0,success,groupPubKey,disqualified,inactive, {from: badSender}));

});

  it("should accept first DKG submission with correct index and fail subsequent attempts", async function() {

    let finalParticipants = await keepGroupImplViaProxy.selectedParticipants();
    let sender = finalParticipants[0];
    
    await keepGroupImplViaProxy.submitDkgResult(1,1,success,groupPubKey,disqualified,inactive, {from: sender});
    await exceptThrow(keepGroupImplViaProxy.submitDkgResult(1,1,success,groupPubKey,disqualified,inactive, {from: sender}));
    assert.equal(tickets.length, ordered.length, "should have 20 tickets");
    assert.equal(tickets.length, groupSize, "ticket number equal to group size");

});

it("Should accept and reject votes properly and return correct DKG result", async function() {
  let finalParticipants = await keepGroupImplViaProxy.selectedParticipants();
  let sender;
 
  for(let i=1;i<=tickets.length;i++){
    sender = finalParticipants[i - 1];
    await keepGroupImplViaProxy.submitDkgResult(i+2,i,success,groupPubKey,disqualified,inactive, {from: sender});
    
  }
  let votes = await keepGroupImplViaProxy.getDkgResultSubmissions();
  assert.equal(await keepGroupImplViaProxy.getFinalResult.call(22), groupPubKey , "get correct final result");
  assert.equal(votes[1][0], tickets.length , "get corrent number of votes for result");
});

});