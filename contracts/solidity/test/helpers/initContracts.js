const BLS = artifacts.require('./cryptography/BLS.sol');
const GroupSelection = artifacts.require('./libraries/operator/GroupSelection.sol');
const Groups = artifacts.require('./libraries/operator/Groups.sol');
const DKGResultVerification = artifacts.require("./libraries/operator/DKGResultVerification.sol");
const Reimbursements = artifacts.require("./libraries/operator/Reimbursements.sol");
const Registry = artifacts.require("./Registry.sol");

async function initContracts(KeepToken, TokenStaking, KeepRandomBeaconService,
  KeepRandomBeaconServiceImplV1, KeepRandomBeaconOperator) {

  let token, registry, stakingContract,
    serviceContractImplV1, serviceContractProxy, serviceContract,
    operatorContract;

  let dkgContributionMargin = 1, // 1%
    withdrawalDelay = 1,
    stakeInitializationPeriod = 1,
    stakeUndelegationPeriod = 30;

  // Initialize Keep token contract
  token = await KeepToken.new();

  // Initialize registry contract
  registry = await Registry.new();

  // Initialize staking contract
  stakingContract = await TokenStaking.new(token.address, registry.address, stakeInitializationPeriod, stakeUndelegationPeriod);

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
  const reimbursements = await Reimbursements.new();
  await KeepRandomBeaconOperator.link("GroupSelection", groupSelection.address);
  await KeepRandomBeaconOperator.link("Groups", groups.address);
  await KeepRandomBeaconOperator.link("DKGResultVerification", dkgResultVerification.address);
  await KeepRandomBeaconOperator.link("Reimbursements", reimbursements.address);

  operatorContract = await KeepRandomBeaconOperator.new(serviceContractProxy.address, stakingContract.address);

  await registry.approveOperatorContract(operatorContract.address);
  await serviceContract.initialize(dkgContributionMargin, withdrawalDelay, registry.address);

  // Set service contract owner as operator contract upgrader by default
  const operatorContractUpgrader = await serviceContract.owner()
  await registry.setOperatorContractUpgrader(serviceContract.address, operatorContractUpgrader);

  await serviceContract.addOperatorContract(operatorContract.address);

  let dkgGasEstimate = await operatorContract.dkgGasEstimate();

  // Genesis should include payment to cover DKG cost to create first group
  let gasPriceCeiling = await operatorContract.gasPriceCeiling();
  await operatorContract.genesis({value: dkgGasEstimate.mul(gasPriceCeiling)});

  return {
    registry: registry,
    token: token,
    stakingContract: stakingContract,
    serviceContract: serviceContract,
    operatorContract: operatorContract
  };
};

module.exports.initContracts = initContracts;
