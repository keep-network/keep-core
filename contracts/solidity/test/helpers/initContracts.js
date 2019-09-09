import { duration } from './increaseTime';
const BLS = artifacts.require('./cryptography/BLS.sol');

async function initContracts(KeepToken, TokenStaking, KeepRandomBeaconService,
  KeepRandomBeaconServiceImplV1, KeepRandomBeaconOperator, KeepRandomBeaconOperatorGroups) {

  let token, stakingContract,
    serviceContractImplV1, serviceContractProxy, serviceContract,
    operatorContract, groupContract;

  const minPayment = 1
  const withdrawalDelay = 1

  // Initialize Keep token contract
  token = await KeepToken.new();

  // Initialize staking contract
  stakingContract = await TokenStaking.new(token.address, duration.days(30));

  // Initialize Keep Random Beacon service contract
  serviceContractImplV1 = await KeepRandomBeaconServiceImplV1.new();
  serviceContractProxy = await KeepRandomBeaconService.new(serviceContractImplV1.address);
  serviceContract = await KeepRandomBeaconServiceImplV1.at(serviceContractProxy.address)

  // Initialize Keep Random Beacon operator contract
  const bls = await BLS.new();
  await KeepRandomBeaconOperator.link("BLS", bls.address);
  groupContract = await KeepRandomBeaconOperatorGroups.new();
  operatorContract = await KeepRandomBeaconOperator.new(serviceContractProxy.address, stakingContract.address, groupContract.address);
  await groupContract.setOperatorContract(operatorContract.address);

  await serviceContract.initialize(minPayment, withdrawalDelay, operatorContract.address);

  await operatorContract.genesis();

  return {
    token: token,
    stakingContract: stakingContract,
    serviceContract: serviceContract,
    operatorContract: operatorContract,
    groupContract: groupContract
  };
};

module.exports.initContracts = initContracts;
