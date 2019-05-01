import { duration } from './helpers/increaseTime';
import mineBlocks from './helpers/mineBlocks';
const StakingProxy = artifacts.require('./StakingProxy.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const KeepToken = artifacts.require('./KeepToken.sol');
const KeepRandomBeaconProxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconImplV1 = artifacts.require('./KeepRandomBeaconImplV1.sol');
const KeepGroupProxy = artifacts.require('./KeepGroup.sol');
const KeepGroupImplV1 = artifacts.require('./KeepGroupImplV1Test.sol');

const minimumStake = 200000;
const groupThreshold = 15;
const groupSize = 20;
const timeoutInitial = 20;
const timeoutSubmission = 50;
const timeoutChallenge = 60;
const timeDKG = 20;
const resultPublicationBlockStep = 3;
const groupExpirationTime = 300;
const activeGroupsThreshold = 5;
const testGroupsNumber = 10;
const expirationStepTime = groupExpirationTime / 10;

contract('Test Group Expirations', function(accounts) {

  let token, stakingProxy, stakingContract,
    keepRandomBeaconImplV1, keepRandomBeaconProxy,
    keepGroupImplV1, keepGroupProxy, keepGroupImplViaProxy,
    owner = accounts[0]

  var selected_counter = -1;
  var gs_counter = -1;
  var expire_counter = -1;
  var print_header = ["selectGroupV0 -> ", "selectGroupV1 -> ", "selectGroupV2 -> ", "selectGroupV3 -> "];
  var print_buffer = ["", "", "", ""];

  beforeEach(async () => {
    token = await KeepToken.new();
    // Initialize staking contract under proxy
    stakingProxy = await StakingProxy.new();
    stakingContract = await TokenStaking.new(token.address, stakingProxy.address, duration.days(30));
    await stakingProxy.authorizeContract(stakingContract.address, {from: owner})
    
    keepRandomBeaconImplV1 = await KeepRandomBeaconImplV1.new();
    keepRandomBeaconProxy = await KeepRandomBeaconProxy.new(keepRandomBeaconImplV1.address);

    // Initialize Keep Group contract

    keepGroupImplV1 = await KeepGroupImplV1.new();
    keepGroupProxy = await KeepGroupProxy.new(keepGroupImplV1.address);
    keepGroupImplViaProxy = await KeepGroupImplV1.at(keepGroupProxy.address);

    await keepGroupImplViaProxy.initialize(
      stakingProxy.address, keepRandomBeaconProxy.address, minimumStake,
      groupThreshold, groupSize, timeoutInitial, timeoutSubmission,
      timeoutChallenge, timeDKG, resultPublicationBlockStep, activeGroupsThreshold,
      groupExpirationTime
    );

    for (var i = 1; i <= testGroupsNumber; i++) {
      await keepGroupImplViaProxy.registerNewGroup([i]); // 2 blocks
      mineBlocks(8);
    }
    gs_counter++
    var different = false;

    if (gs_counter % 4 == 0) {
      for (var i = 0; i < 4; i++) {
        for (var j = 0; j < 4; j++) {
          if (print_buffer[i] != print_buffer[j])
            different = true;
        }
      }
      if (different) {
        for (var i = 0; i < 4; i++) {
          console.log(print_header[i], print_buffer[i]);
        }
      }
      print_buffer = ["", "", "", ""];
      selected_counter++;  
      if (selected_counter % 10 == 0) {
        expire_counter++;
      }
    }

  });

  async function printExpirations(select, expire, selected) {
    let gstx;

    await mineBlocks(expirationStepTime*expire);
    switch(select) {
      case 0: gstx = await keepGroupImplViaProxy.selectGroupV0(selected); break;
      case 1: gstx = await keepGroupImplViaProxy.selectGroupV1(selected); break;
      case 2: gstx = await keepGroupImplViaProxy.selectGroupV2(selected); break;
      case 3: gstx = await keepGroupImplViaProxy.selectGroupV3(selected); break;
    }

    var groups = "";
    let expiredOffset = await keepGroupImplViaProxy.getExpiredOffset();
    
    print_buffer[select] = print_buffer[select].concat("Expiration: " + Number(expirationStepTime*expire).toString() + " Selected: " + Number(selected).toString() + " expiredOffset: " + Number(expiredOffset).toString() + " Groups: " );

    for (var i = 0; i < testGroupsNumber; i++) {
      if (i >= expiredOffset)
        groups = groups.concat("A");
      else
        groups = groups.concat("E");
    }
    print_buffer[select] = print_buffer[select].concat(groups);
  }

  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
  it("", async function() {await printExpirations(gs_counter % 4, selected_counter % 10, expire_counter % 10);});
});
