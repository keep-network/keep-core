import { duration } from './increaseTime';
import { bls } from './data';

async function initContracts(accounts, KeepToken, StakingProxy, TokenStaking, KeepRandomBeaconFrontend,
  KeepRandomBeaconFrontendImplV1, KeepRandomBeaconOperator) {

  let token, stakingProxy, stakingContract,
    frontendImplV1, frontendProxy, frontend,
    operatorContract;

  let minimumStake = 200000,
    groupThreshold = 15,
    groupSize = 20,
    timeoutInitial = 20,
    timeoutSubmission = 100,
    timeoutChallenge = 60,
    timeDKG = 20,
    resultPublicationBlockStep = 3,
    groupActiveTime = 300,
    activeGroupsThreshold = 5,
    minPayment = 1,
    withdrawalDelay = 1,
    relayRequestTimeout = 10;

  // Initialize Keep token contract
  token = await KeepToken.new();

  // Initialize staking contract under proxy
  stakingProxy = await StakingProxy.new();
  stakingContract = await TokenStaking.new(token.address, stakingProxy.address, duration.days(30));
  await stakingProxy.authorizeContract(stakingContract.address, {from: accounts[0]})

  // Initialize Keep Random Beacon frontend contract
  frontendImplV1 = await KeepRandomBeaconFrontendImplV1.new();
  frontendProxy = await KeepRandomBeaconFrontend.new(frontendImplV1.address);
  frontend = await KeepRandomBeaconFrontendImplV1.at(frontendProxy.address)

  // Initialize Keep Random Beacon operator contract
  operatorContract = await KeepRandomBeaconOperator.new();
  await operatorContract.initialize(
    stakingProxy.address, frontend.address, minimumStake, groupThreshold,
    groupSize, timeoutInitial, timeoutSubmission, timeoutChallenge, timeDKG, resultPublicationBlockStep,
    activeGroupsThreshold, groupActiveTime,
    bls.groupSignature, bls.groupPubKey
  );

  await frontend.initialize(minPayment, withdrawalDelay, operatorContract.address, relayRequestTimeout);

  // TODO: replace with a secure authorization protocol (addressed in RFC 4).
  await operatorContract.authorizeStakingContract(stakingContract.address);
  await operatorContract.relayEntry(1, bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

  return {
    token: token,
    frontend: frontend,
    operatorContract: operatorContract
  };
};

module.exports.initContracts = initContracts;
