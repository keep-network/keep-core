import { duration } from './helpers/increaseTime';
import exceptThrow from './helpers/expectThrow';
import {bls} from './helpers/data';
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

contract('TestKeepGroupExpiration', function(accounts) {

  let token, stakingProxy, minimumStake, groupThreshold, groupSize,
    timeoutInitial, timeoutSubmission, timeoutChallenge,
    keepRandomBeaconImplV1, keepRandomBeaconProxy,
    keepGroupImplV1, keepGroupProxy, keepGroupImplViaProxy

  beforeEach(async () => {
    token = await KeepToken.new();
    
    // Initialize staking contract under proxy
    stakingProxy = await StakingProxy.new();
    
    keepRandomBeaconImplV1 = await KeepRandomBeaconImplV1.new();
    keepRandomBeaconProxy = await KeepRandomBeaconProxy.new(keepRandomBeaconImplV1.address);
    
    // Initialize Keep Group contract
    minimumStake = 200000;
    groupThreshold = 15;
    groupSize = 20;
    timeoutInitial = 20;
    timeoutSubmission = 50;
    timeoutChallenge = 60;

    keepGroupImplV1 = await KeepGroupImplV1.new();
    keepGroupProxy = await KeepGroupProxy.new(keepGroupImplV1.address);
    keepGroupImplViaProxy = await KeepGroupImplV1.at(keepGroupProxy.address);
    await keepGroupImplViaProxy.initialize(
      stakingProxy.address, keepRandomBeaconProxy.address, minimumStake, groupThreshold, groupSize, timeoutInitial, timeoutSubmission, timeoutChallenge
    );

    for (var i = 1; i <= 7; i++)
      await keepGroupImplV1.submitGroupPublicKey([i], i);
  });

  it("should be able to check if groups were added", async function() {
    let numberOfGroups = await keepGroupImplV1.numberOfGroups();
    assert.equal(numberOfGroups.toString(), "7", "Number of groups not equal 7");
  });

  it("should be able to check if groups expire", async function() {
    let before = await keepGroupImplV1.numberOfGroups();
    assert.equal(before.toString(), "7", "Number of groups should be equal 7"); 
    await keepGroupImplV1.selectGroup("1");
    let after = await keepGroupImplV1.numberOfGroups();
    assert.notEqual(after.toString(), "7", "Number of groups after `selectGroup()` should not be equal 7"); 
    assert.notEqual(after.toString(), before.toString(), "Number of groups should not be equal");
  });

  it("should be able to check if last group is not expiring", async function() {
    await keepGroupImplV1.selectGroup("1");
    await keepGroupImplV1.selectGroup("1");
    await keepGroupImplV1.selectGroup("1");
    await keepGroupImplV1.selectGroup("1");
  
    let after = await keepGroupImplV1.numberOfGroups();
    assert.equal(after.toString(), "1", "Number of groups should be equal 1");
  });

  it("should be able to survive this stress test", async function() {
    for (var i = 1; i <= 100; i++)
      await keepGroupImplV1.submitGroupPublicKey([i], i);

    for (var i = 1; i <= 101; i++)
      await keepGroupImplV1.selectGroup(i);

    let after = await keepGroupImplV1.numberOfGroups();
    assert.equal(after.toString(), "1", "Number of groups should be equal 1");
  });
});
