const {contract, accounts} = require("@openzeppelin/test-environment")

const BLS = contract.fromArtifact('BLS');
const GroupSelection = contract.fromArtifact('GroupSelection');
const Groups = contract.fromArtifact('Groups');
const DKGResultVerification = contract.fromArtifact("DKGResultVerification");
const Reimbursements = contract.fromArtifact("Reimbursements");
const Registry = contract.fromArtifact("Registry");

async function initContracts(KeepToken, TokenStaking, KeepRandomBeaconService,
  KeepRandomBeaconServiceImplV1, KeepRandomBeaconOperator) {

  let token, registry, stakingContract,
    serviceContractImplV1, serviceContractProxy, serviceContract,
    operatorContract;

  let dkgContributionMargin = 1, // 1%
    withdrawalDelay = 1,
    stakeInitializationPeriod = 30, // In seconds
    stakeUndelegationPeriod = 300; // In seconds

  // Initialize Keep token contract
  token = await KeepToken.new({from: accounts[0]});

  // Initialize registry contract
  registry = await Registry.new({from: accounts[0]});

  // Initialize staking contract
  stakingContract = await TokenStaking.new(token.address, registry.address, stakeInitializationPeriod, stakeUndelegationPeriod, {from: accounts[0]});

  // Initialize Keep Random Beacon service contract
  serviceContractImplV1 = await KeepRandomBeaconServiceImplV1.new({from: accounts[0]});

  const initialize = serviceContractImplV1.contract.methods
      .initialize(
          dkgContributionMargin,
          withdrawalDelay,
          registry.address,
      ).encodeABI();

  serviceContractProxy = await KeepRandomBeaconService.new(serviceContractImplV1.address, initialize, {from: accounts[0]});

  serviceContract = await KeepRandomBeaconServiceImplV1.at(serviceContractProxy.address);
  // Initialize Keep Random Beacon operator contract
  const bls = await BLS.new({from: accounts[0]});
  await KeepRandomBeaconOperator.detectNetwork()
  await KeepRandomBeaconOperator.link("BLS", bls.address);
  const groupSelection = await GroupSelection.new({from: accounts[0]});
  await Groups.detectNetwork()
  await Groups.link("BLS", bls.address);
  const groups = await Groups.new({from: accounts[0]});

  const dkgResultVerification = await DKGResultVerification.new({from: accounts[0]});

  const reimbursements = await Reimbursements.new({from: accounts[0]});

  await KeepRandomBeaconOperator.detectNetwork()
  await KeepRandomBeaconOperator.link("GroupSelection", groupSelection.address);
  await KeepRandomBeaconOperator.detectNetwork()
  await KeepRandomBeaconOperator.link("Groups", groups.address);
  await KeepRandomBeaconOperator.detectNetwork()
  await KeepRandomBeaconOperator.link("DKGResultVerification", dkgResultVerification.address);
  await KeepRandomBeaconOperator.detectNetwork()
  await KeepRandomBeaconOperator.link("Reimbursements", reimbursements.address);
  operatorContract = await KeepRandomBeaconOperator.new(serviceContractProxy.address, stakingContract.address, {from: accounts[0]});

  await registry.approveOperatorContract(operatorContract.address, {from: accounts[0]});

  // Set service contract owner as operator contract upgrader by default
  const operatorContractUpgrader = await serviceContractProxy.admin({from: accounts[0]})
  await registry.setOperatorContractUpgrader(serviceContract.address, operatorContractUpgrader, {from: accounts[0]});

  await serviceContract.addOperatorContract(operatorContract.address, {from: accounts[0]});

  let dkgGasEstimate = await operatorContract.dkgGasEstimate({from: accounts[0]});

  // Genesis should include payment to cover DKG cost to create first group
  let gasPriceCeiling = await operatorContract.gasPriceCeiling({from: accounts[0]});
  await operatorContract.genesis({value: dkgGasEstimate.mul(gasPriceCeiling), from: accounts[0]});

  return {
    registry: registry,
    token: token,
    stakingContract: stakingContract,
    serviceContract: serviceContract,
    operatorContract: operatorContract
  };
};

module.exports = initContracts;
