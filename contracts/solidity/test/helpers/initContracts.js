import { duration } from './increaseTime';
const BLS = artifacts.require('./cryptography/BLS.sol');
const GroupSelection = artifacts.require('./libraries/operator/GroupSelection.sol');
const Groups = artifacts.require('./libraries/operator/Groups.sol');
const DKGResultVerification = artifacts.require("./libraries/operator/DKGResultVerification.sol");
const RegistryKeeper = artifacts.require("./RegistryKeeper.sol");

async function initContracts(KeepToken, TokenStaking, KeepRandomBeaconService,
  KeepRandomBeaconServiceImplV1, KeepRandomBeaconOperator) {

  let token, registryKeeper, stakingContract,
    serviceContractImplV1, serviceContractProxy, serviceContract,
    operatorContract;

  let priceFeedEstimate = web3.utils.toBN(20).mul(web3.utils.toBN(10**9)), // (20 Gwei = 20 * 10^9 wei)
    fluctuationMargin = 50, // 50%
    dkgContributionMargin = 1, // 1%
    withdrawalDelay = 1;

  // Initialize Keep token contract
  token = await KeepToken.new();

  // Initialize registry keeper contract
  registryKeeper = await RegistryKeeper.new();

  // Initialize staking contract
  stakingContract = await TokenStaking.new(token.address, registryKeeper.address, duration.days(30));

  // Initialize Keep Random Beacon service contract
  serviceContractImplV1 = await KeepRandomBeaconServiceImplV1.new();
  serviceContractProxy = await KeepRandomBeaconService.new(serviceContractImplV1.address);
  serviceContract = await KeepRandomBeaconServiceImplV1.at(serviceContractProxy.address)

  // Initialize Keep Random Beacon operator contract
  const bls = await BLS.new();
  await KeepRandomBeaconOperator.link("BLS", bls.address);
  const groupSelection = await GroupSelection.new();
  const groups = await Groups.new();
  const dkgResultVerification = await DKGResultVerification.new();
  await KeepRandomBeaconOperator.link("GroupSelection", groupSelection.address);
  await KeepRandomBeaconOperator.link("Groups", groups.address);
  await KeepRandomBeaconOperator.link("DKGResultVerification", dkgResultVerification.address);
  operatorContract = await KeepRandomBeaconOperator.new(serviceContractProxy.address, stakingContract.address);

  await registryKeeper.approveOperatorContract(operatorContract.address);
  await serviceContract.initialize(priceFeedEstimate, fluctuationMargin, dkgContributionMargin, withdrawalDelay, registryKeeper.address);
  await serviceContract.addOperatorContract(operatorContract.address);

  let dkgGasEstimate = await operatorContract.dkgGasEstimate();
  let gasPriceWithFluctuationMargin = priceFeedEstimate.add(priceFeedEstimate.mul(web3.utils.toBN(fluctuationMargin)).div(web3.utils.toBN(100)));

  // Genesis should include payment to cover DKG cost to create first group
  await operatorContract.genesis({value: dkgGasEstimate.mul(gasPriceWithFluctuationMargin)});

  return {
    registryKeeper: registryKeeper,
    token: token,
    stakingContract: stakingContract,
    serviceContract: serviceContract,
    operatorContract: operatorContract
  };
};

module.exports.initContracts = initContracts;
