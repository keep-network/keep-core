import { duration } from './helpers/increaseTime';
import exceptThrow from './helpers/expectThrow';
const KeepToken = artifacts.require('./KeepToken.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const KeepRandomBeaconProxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconImplV1 = artifacts.require('./KeepRandomBeaconImplV1.sol');
const KeepRandomBeaconBackend = artifacts.require('./KeepRandomBeaconBackend.sol');

contract('TestKeepRandomBeaconBackend', function(accounts) {

  let token, stakingProxy, stakingContract, minimumStake, groupThreshold, groupSize,
    timeoutInitial, timeoutSubmission, timeoutChallenge, timeDKG, resultPublicationBlockStep,
    keepRandomBeaconBackend,
    keepRandomBeaconImplV1, keepRandomBeaconProxy,
    account_one = accounts[0],
    account_two = accounts[1];

  beforeEach(async () => {
    token = await KeepToken.new();

    // Initialize staking contract under proxy
    stakingProxy = await StakingProxy.new();
    stakingContract = await TokenStaking.new(token.address, stakingProxy.address, duration.days(30));
    await stakingProxy.authorizeContract(stakingContract.address, {from: account_one})

    // Initialize Keep Random Beacon contract
    keepRandomBeaconImplV1 = await KeepRandomBeaconImplV1.new(1,1);
    keepRandomBeaconProxy = await KeepRandomBeaconProxy.new(keepRandomBeaconImplV1.address);

    // Initialize Keep Random Beacon backend contract
    minimumStake = 200;
    groupThreshold = 150;
    groupSize = 200;
    timeoutInitial = 20;
    timeoutSubmission = 100;
    timeoutChallenge = 60;
    timeDKG = 20;
    resultPublicationBlockStep = 3;

    keepRandomBeaconBackend = await KeepRandomBeaconBackend.new();
    await keepRandomBeaconBackend.initialize(
      stakingProxy.address, keepRandomBeaconProxy.address, minimumStake, groupThreshold,
      groupSize, timeoutInitial, timeoutSubmission, timeoutChallenge, timeDKG, resultPublicationBlockStep
    );
  });

  it("should fail to update minimum stake by non owner", async function() {
    await exceptThrow(keepRandomBeaconBackend.setMinimumStake(123, {from: account_two}));
  });

  it("should be able to update minimum stake by the owner", async function() {
    await keepRandomBeaconBackend.setMinimumStake(123);
    let newMinStake = await keepRandomBeaconBackend.minimumStake();
    assert.equal(newMinStake, 123, "Should be able to get updated minimum stake.");
  });

  it("should be able to check if the implementation contract was initialized", async function() {
    assert.equal(await keepRandomBeaconBackend.initialized(), true, "Implementation contract should be initialized.");
  });
});
