import { duration } from './increaseTime';
import { bls } from './data';

async function initContracts(accounts, KeepToken, StakingProxy, TokenStaking, KeepRandomBeaconService,
  KeepRandomBeaconServiceImplV1, KeepRandomBeaconOperator) {

  let token, stakingProxy, stakingContract,
    serviceContractImplV1, serviceContractProxy, serviceContract,
    operatorContract;

  let minPayment = 1
  let withdrawalDelay = 1

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
    stakingProxy.address, serviceContract.address,
    bls.previousEntry, bls.seed, bls.groupPubKey
  );

  await serviceContract.initialize(minPayment, withdrawalDelay, operatorContract.address);

  // TODO: replace with a secure authorization protocol (addressed in RFC 4).
  await operatorContract.authorizeStakingContract(stakingContract.address);
  await operatorContract.relayEntry(bls.groupSignature);

  return {
    token: token,
    stakingProxy: stakingProxy,
    stakingContract: stakingContract,
    serviceContract: serviceContract,
    operatorContract: operatorContract
  };
};

module.exports.initContracts = initContracts;
