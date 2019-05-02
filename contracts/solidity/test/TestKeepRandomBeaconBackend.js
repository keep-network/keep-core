import {bls} from './helpers/data';
import { duration } from './helpers/increaseTime';
import exceptThrow from './helpers/expectThrow';
const KeepToken = artifacts.require('./KeepToken.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const KeepRandomBeaconFrontendProxy = artifacts.require('./KeepRandomBeaconFrontendProxy.sol');
const KeepRandomBeaconFrontendImplV1 = artifacts.require('./KeepRandomBeaconFrontendImplV1.sol');
const KeepRandomBeaconBackend = artifacts.require('./KeepRandomBeaconBackend.sol');

contract('TestKeepRandomBeaconBackend', function(accounts) {

  let token, stakingProxy, stakingContract, minimumStake, groupThreshold, groupSize,
    timeoutInitial, timeoutSubmission, timeoutChallenge, timeDKG, resultPublicationBlockStep,
    backend,
    frontendImplV1, frontendProxy,
    account_one = accounts[0],
    account_two = accounts[1];

  beforeEach(async () => {
    token = await KeepToken.new();

    // Initialize staking contract under proxy
    stakingProxy = await StakingProxy.new();
    stakingContract = await TokenStaking.new(token.address, stakingProxy.address, duration.days(30));
    await stakingProxy.authorizeContract(stakingContract.address, {from: account_one})

    // Initialize Keep Random Beacon contract
    frontendImplV1 = await KeepRandomBeaconFrontendImplV1.new(1,1);
    frontendProxy = await KeepRandomBeaconFrontendProxy.new(frontendImplV1.address);

    // Initialize Keep Random Beacon backend contract
    minimumStake = 200;
    groupThreshold = 150;
    groupSize = 200;
    timeoutInitial = 20;
    timeoutSubmission = 100;
    timeoutChallenge = 60;
    timeDKG = 20;
    resultPublicationBlockStep = 3;

    backend = await KeepRandomBeaconBackend.new();
    await backend.initialize(
      stakingProxy.address, frontendProxy.address, minimumStake, groupThreshold,
      groupSize, timeoutInitial, timeoutSubmission, timeoutChallenge, timeDKG, resultPublicationBlockStep,
      bls.groupSignature, bls.groupPubKey
    );
  });

  it("should fail to update minimum stake by non owner", async function() {
    await exceptThrow(backend.setMinimumStake(123, {from: account_two}));
  });

  it("should be able to update minimum stake by the owner", async function() {
    await backend.setMinimumStake(123);
    let newMinStake = await backend.minimumStake();
    assert.equal(newMinStake, 123, "Should be able to get updated minimum stake.");
  });

  it("should be able to check if the implementation contract was initialized", async function() {
    assert.isTrue(await backend.initialized(), "Implementation contract should be initialized.");
  });
});
