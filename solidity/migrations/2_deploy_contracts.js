const KeepToken = artifacts.require("./KeepToken.sol");
const ModUtils = artifacts.require("./utils/ModUtils.sol");
const AltBn128 = artifacts.require("./cryptography/AltBn128.sol");
const BLS = artifacts.require("./cryptography/BLS.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const PermissiveStakingPolicy = artifacts.require('./PermissiveStakingPolicy.sol');
const GuaranteedMinimumStakingPolicy = artifacts.require('./GuaranteedMinimumStakingPolicy.sol');
const TokenGrant = artifacts.require("./TokenGrant.sol");
const KeepRandomBeaconService = artifacts.require("./KeepRandomBeaconService.sol");
const KeepRandomBeaconServiceImplV1 = artifacts.require("./KeepRandomBeaconServiceImplV1.sol");
const KeepRandomBeaconOperator = artifacts.require("./KeepRandomBeaconOperator.sol");
const GroupSelection = artifacts.require("./libraries/operator/GroupSelection.sol");
const Groups = artifacts.require("./libraries/operator/Groups.sol");
const DKGResultVerification = artifacts.require("./libraries/operator/DKGResultVerification.sol");
const Reimbursements = artifacts.require("./libraries/operator/Reimbursements.sol");
const Registry = artifacts.require("./Registry.sol");

let governance, registryKeeper, panicButton, operatorContractUpgrader;

let initializationPeriod = 518400; // ~6 days
const undelegationPeriod = 7776000; // ~3 months
const withdrawalDelay = 86400; // 1 day
const dkgContributionMargin = 1; // 1%

module.exports = async function(deployer, network, accounts) {
  // Set the stake initialization period to 1 block for local development and testnet.
  if (network === 'local' || network === 'ropsten' || network === 'keep_dev') {
    initializationPeriod = 1;
  }

  if (network === 'mainnet') {
    if (accounts.length < 4) {
      throw Error("Not enough accounts for mainnet deployment")
    }
    governance = accounts[0]
    registryKeeper = accounts[1]
    panicButton = accounts[2]
    operatorContractUpgrader = accounts[3]
  } else {
    // Set all roles to the default account for simplicity
    governance = accounts[0]
    registryKeeper = accounts[0]
    panicButton = accounts[0]
    operatorContractUpgrader = accounts[0]
  }

  // Deployer 'from' address is used in the contract constructor
  await deployer.deploy(KeepToken, {from: governance}); // All token supply goes to the governance
  await deployer.deploy(Registry, {from: governance}); // All roles set to governance by default

  // Non ownable contracts (deployer 'from' is not relevant)
  await deployer.deploy(ModUtils);
  await deployer.link(ModUtils, AltBn128);
  await deployer.deploy(AltBn128);
  await deployer.link(AltBn128, BLS);
  await deployer.deploy(BLS);
  await deployer.deploy(TokenStaking, KeepToken.address, Registry.address, initializationPeriod, undelegationPeriod);
  await deployer.deploy(PermissiveStakingPolicy);
  await deployer.deploy(GuaranteedMinimumStakingPolicy, TokenStaking.address);
  await deployer.deploy(TokenGrant, KeepToken.address);
  await deployer.deploy(GroupSelection);
  await deployer.link(GroupSelection, KeepRandomBeaconOperator);
  await deployer.link(BLS, Groups);
  await deployer.deploy(Groups);
  await deployer.link(Groups, KeepRandomBeaconOperator);
  await deployer.deploy(DKGResultVerification);
  await deployer.link(DKGResultVerification, KeepRandomBeaconOperator);
  await deployer.deploy(Reimbursements);
  await deployer.link(Reimbursements, KeepRandomBeaconOperator);
  await deployer.link(BLS, KeepRandomBeaconOperator);

  // Non ownable contract (deployer 'from' is not relevant)
  const keepRandomBeaconServiceImplV1 = await deployer.deploy(KeepRandomBeaconServiceImplV1);

  const initialize = keepRandomBeaconServiceImplV1.contract.methods
      .initialize(
          dkgContributionMargin,
          withdrawalDelay,
          Registry.address
      ).encodeABI();

  // Deployer 'from' address is used in the contract constructor to set
  // 'operatorContractUpgrader' role. It can upgrade service contract
  // implementation, add operator contracts to it and withdraw eth from
  // the contract.
  await deployer.deploy(
      KeepRandomBeaconService,
      KeepRandomBeaconServiceImplV1.address,
      initialize,
      {from: operatorContractUpgrader}
  );

  // Deployer 'from' address is used in the contract constructor to set
  // 'operatorContractUpgrader' role. It сan add/remove service contracts
  // in the contract
  await deployer.deploy(KeepRandomBeaconOperator, KeepRandomBeaconService.address, TokenStaking.address,
    {from: operatorContractUpgrader}
  );
};
