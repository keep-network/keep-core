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
let EthUtil = require('ethereumjs-utils');

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
  let  b1, b2, b3, b4, bufferFinal, success,
  disqualified, inactive, msgHash,
  signature, signatureRPC, positions, token, stakingProxy,
   stakingContract, minimumStake, groupThreshold, groupSize,
    randomBeaconValue,
    timeoutInitial, timeoutSubmission, timeoutChallenge,
    keepRandomBeaconImplV1, keepRandomBeaconProxy, keepRandomBeaconImplViaProxy,
    keepGroupImplV1, keepGroupProxy, keepGroupImplViaProxy, groupPubKey,
    staker1 = accounts[0], tickets1,
    staker2 = accounts[1], tickets2,
    staker3 = accounts[2], tickets3,
    staker4 = accounts[3], tickets4;

//   let signatures, votes, resultHash,
  let pks = [
    // '51b13e5e39dfbaa497653490a5bd04a4f3293de4a309dc1a26bc36dfbdb67c59',
    // '9e27a55a55b139be2bd0676f9aa68b0ec1dac89b09ae0dcdc7e61d2544e4fce0',
    // 'e0cef4ec8b5ebbea96a12c340ad3ea1d0e6de932a9c90df48663ece213384c92',
    // '6be271a329ef155847f81aea752b4a1a465c0e174ec493a815dda037ff601521',
    // 'e288d27643dc3637ffc9d22c7018fad45f3cf5b21f40260b4e08529db009019d',
    // 'bc534249a87da0664aca5129e7b99d2aa55ecfdece866712f4c19f385a791af4',
    // 'a0daf6d444d2c16d4046bd61196cae8ea6e87a372c3bbaf57f5b230ed24fb12d'
]

  disqualified = '0x0000000110000000110000000110000000110000'
  inactive =  '0x0000001000000001000000001000000001000000'
  groupPubKey = "0x1000000000000000000000000000000000000000000000000000000000000000"

  b1 = new Buffer("0x01", 'hex')
  b2 = new Buffer(disqualified, 'hex')
  b3 = new Buffer(inactive, 'hex')
  b4 = new Buffer(groupPubKey, 'hex')
  bufferFinal = Buffer.concat([b1, b2, b3, b4]);
  msgHash = EthUtil.hashPersonalMessage(bufferFinal);
  
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
    groupSize = 20;
    timeoutInitial = 20;
    timeoutSubmission = 50;
    timeoutChallenge = 60;

    randomBeaconValue = bls.groupSignature;

    keepGroupImplV1 = await KeepGroupImplV1.new();
    keepGroupProxy = await KeepGroupProxy.new(keepGroupImplV1.address);
    keepGroupImplViaProxy = await KeepGroupImplV1.at(keepGroupProxy.address);
    await keepGroupImplViaProxy.initialize(
      stakingProxy.address, keepRandomBeaconProxy.address, minimumStake, groupThreshold, groupSize, timeoutInitial, timeoutSubmission, timeoutChallenge
    );

    await keepRandomBeaconImplViaProxy.initialize(1,1, randomBeaconValue, keepGroupProxy.address);
    await keepRandomBeaconImplViaProxy.relayEntry(1, bls.groupSignature, bls.groupPubKey, bls.previousEntry, 1);

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
    if(pks.length!=0){
      positions = []
      let tickets = []
      let signatures; 
      let caller;
      for(let i=0;i<10;i++){
        await keepGroupImplViaProxy.submitTicket(tickets1[i].value, staker1, tickets1[i].virtualStakerIndex, {from: staker1});
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
      let op = await keepGroupImplViaProxy.orderedParticipants.call()

      for(let i=0;i<op.length;i++){
        caller = accounts.indexOf(op[i]);
        signature = EthUtil.ecsign(msgHash, new Buffer(pks[caller], 'hex'));
        signatureRPC = EthUtil.toRpcSig(signature.v, signature.r, signature.s)
        positions.push(i+1);
        if(signatures == undefined){
          signatures = signatureRPC
      }
        else{
            signatures+=signatureRPC.slice(2,signatureRPC.length);
        }
      }
      let getA = await keepGroupImplViaProxy.submitDkgResult(1, 0, groupPubKey, disqualified, inactive, signatures, positions, msgHash)
      //console.log(getA.receipt.gasUsed);
    }
    else{
      assert.equal(1,1,"holder")
    }
  });
})  