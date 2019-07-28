import { duration } from './increaseTime';
import { bls } from './data';

async function initContracts(accounts, KeepToken, StakingProxy, TokenStaking, KeepRandomBeaconService,
  KeepRandomBeaconServiceImplV1, KeepRandomBeaconOperator) {

  const gasPrice = web3.utils.toBN(20).mul(web3.utils.toBN(10**9)); // (20 Gwei) TODO: Use historical average of recently served requests?
  const signingCostEstimate = web3.utils.toBN(1240000).mul(gasPrice); // (Wei) TODO: Update once alt_bn128 gas costs reduction is implemented.
  const createGroupCostEstimate = web3.utils.toBN(2260000).mul(gasPrice); // Wei
  const minimumCallbackAllowance = web3.utils.toBN(200000).mul(gasPrice); // Wei
  const profitMargin = web3.utils.toBN(100); // %
  const entryFeeEstimate = signingCostEstimate.add(createGroupCostEstimate); // Wei
  const profitMarginEstimate = entryFeeEstimate.mul(profitMargin).div(web3.utils.toBN(100)); // Wei

  let token, stakingProxy, stakingContract,
    serviceContractImplV1, serviceContractProxy, serviceContract,
    operatorContract;

  let minimumStake = web3.utils.toBN(200000),
    groupThreshold = 15,
    groupSize = 20,
    timeoutInitial = 20,
    timeoutSubmission = 100,
    timeoutChallenge = 60,
    timeDKG = 20,
    resultPublicationBlockStep = 3,
    groupActiveTime = 300,
    activeGroupsThreshold = 5,
    minimumPayment = entryFeeEstimate.add(profitMarginEstimate).add(minimumCallbackAllowance), // Wei
    withdrawalDelay = 1,
    relayRequestTimeout = 10;

  // Initialize Keep token contract
  token = await KeepToken.new();

  // Initialize staking contract under proxy
  stakingProxy = await StakingProxy.new();
  stakingContract = await TokenStaking.new(token.address, stakingProxy.address, duration.days(30));
  await stakingProxy.authorizeContract(stakingContract.address, {from: accounts[0]})

  // Initialize Keep Random Beacon service contract
  serviceContractImplV1 = await KeepRandomBeaconServiceImplV1.new();
  serviceContractProxy = await KeepRandomBeaconService.new(serviceContractImplV1.address);
  serviceContract = await KeepRandomBeaconServiceImplV1.at(serviceContractProxy.address)

  // Initialize Keep Random Beacon operator contract
  operatorContract = await KeepRandomBeaconOperator.new();
  await operatorContract.initialize(
    stakingProxy.address, serviceContract.address, minimumStake, groupThreshold,
    groupSize, timeoutInitial, timeoutSubmission, timeoutChallenge, timeDKG, resultPublicationBlockStep,
    activeGroupsThreshold, groupActiveTime, relayRequestTimeout,
    bls.groupSignature, bls.groupPubKey
  );

  await serviceContract.initialize(minimumPayment, withdrawalDelay, operatorContract.address);

  // TODO: replace with a secure authorization protocol (addressed in RFC 4).
  await operatorContract.authorizeStakingContract(stakingContract.address);
  await operatorContract.relayEntry(bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

  return {
    config: {
      minimumStake: minimumStake,
      groupThreshold: groupThreshold,
      groupSize: groupSize,
      timeoutInitial: timeoutInitial,
      timeoutSubmission: timeoutSubmission,
      timeoutChallenge: timeoutChallenge,
      timeDKG: timeDKG,
      resultPublicationBlockStep: resultPublicationBlockStep,
      groupActiveTime: groupActiveTime,
      activeGroupsThreshold: activeGroupsThreshold,
      minimumPayment: minimumPayment,
      withdrawalDelay: withdrawalDelay,
      relayRequestTimeout: relayRequestTimeout
    },
    token: token,
    stakingProxy: stakingProxy,
    stakingContract: stakingContract,
    serviceContract: serviceContract,
    operatorContract: operatorContract
  };
};

module.exports.initContracts = initContracts;
