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

  let stakingProxy, minimumStake, groupThreshold, groupSize,
    timeoutInitial, timeoutSubmission, timeoutChallenge,
    groupExpirationTimeout, numberOfActiveGroups, testGroupsNumber,
    keepRandomBeaconImplV1, keepRandomBeaconProxy,
    keepGroupImplV1, keepGroupProxy, keepGroupImplViaProxy

  beforeEach(async () => {
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
    groupExpirationTimeout = 300;
    numberOfActiveGroups = 5;

    keepGroupImplV1 = await KeepGroupImplV1.new();
    keepGroupProxy = await KeepGroupProxy.new(keepGroupImplV1.address);
    keepGroupImplViaProxy = await KeepGroupImplV1.at(keepGroupProxy.address);
    await keepGroupImplViaProxy.initialize(
      stakingProxy.address, keepRandomBeaconProxy.address, minimumStake, groupThreshold, groupSize, timeoutInitial, timeoutSubmission, timeoutChallenge, groupExpirationTimeout, numberOfActiveGroups
    );

    testGroupsNumber = 100;

    for (var i = 1; i <= testGroupsNumber; i++)
      await keepGroupImplViaProxy.groupAdd([i]);
  });

  it("it should be able to count the number of active groups", async function() {
    let numberOfGroups = await keepGroupImplViaProxy.numberOfGroups();
    assert.equal(Number(numberOfGroups), testGroupsNumber, "Number of groups not equals to number of test groups");
  });

  it("should be able to check if at least one group is marked as expired", async function() {
    let numberOfGroups = await keepGroupImplViaProxy.numberOfGroups();
    
    for (var i = 1; i <= testGroupsNumber; i++) {
      mineBlocks(groupExpirationTimeout);
      await keepGroupImplViaProxy.selectGroup("1");
      numberOfGroups = await keepGroupImplViaProxy.numberOfGroups();

      if (Number(numberOfGroups) < testGroupsNumber)
        break;
    }

    assert.notEqual(Number(numberOfGroups), testGroupsNumber, "Some groups should be marked as expired");
  });

  it("should be able to check that groups are marked as expired except the minimal active groups number", async function() {
    let after = await keepGroupImplViaProxy.numberOfGroups();

    for (var i = 1; i <= testGroupsNumber; i++) {
      mineBlocks(groupExpirationTimeout);
      await keepGroupImplViaProxy.selectGroup("1");
      after = await keepGroupImplViaProxy.numberOfGroups();

      if (Number(after) == numberOfActiveGroups)
        break;
    }
    
    assert.equal(Number(after), numberOfActiveGroups, "Number of groups should be at least 2 below the test group numbers");
  });
});