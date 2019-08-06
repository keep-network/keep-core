import { duration } from './increaseTime';
import { bls } from './data';

async function initContracts(accounts, KeepToken, StakingProxy, TokenStaking, KeepRandomBeaconService,
  KeepRandomBeaconServiceImplV1, KeepRandomBeaconOperator) {

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
    minimumGasPrice = web3.utils.toBN(20).mul(web3.utils.toBN(10**9)), // (20 Gwei) TODO: Use historical average of recently served requests?
    minimumCallbackAllowance = web3.utils.toBN(200000), // Minimum gas required for relay request callback.
    profitMargin = 1, // Signing group reward per each member in % of the entry fee.
    createGroupFee = 10, // Fraction in % of the estimated cost of group creation that is included in relay request payment.
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

  await serviceContract.initialize(minimumGasPrice, minimumCallbackAllowance, profitMargin, createGroupFee, withdrawalDelay, operatorContract.address);

  // TODO: replace with a secure authorization protocol (addressed in RFC 4).
  await operatorContract.authorizeStakingContract(stakingContract.address);

  // Add initial funds to the fee pool to trigger group creation on relay entry without waiting for fee accumulation
  let createGroupGasEstimateCost = await operatorContract.createGroupGasEstimate();
  await serviceContract.fundCreateGroupFeePool({value: createGroupGasEstimateCost.mul(minimumGasPrice)});

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
      minimumGasPrice: minimumGasPrice,
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
