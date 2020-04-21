const KeepToken = artifacts.require("./KeepToken.sol");
const ModUtils = artifacts.require("./utils/ModUtils.sol");
const AltBn128 = artifacts.require("./cryptography/AltBn128.sol");
const BLS = artifacts.require("./cryptography/BLS.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const PermissiveStakingPolicy = artifacts.require('./PermissiveStakingPolicy.sol');
const GuaranteedMinimumStakingPolicy = artifacts.require('./GuaranteedMinimumStakingPolicy.sol');
const TokenGrant = artifacts.require("./TokenGrant.sol");
const ManagedGrantFactory = artifacts.require("./ManagedGrantFactory.sol");
const KeepRandomBeaconService = artifacts.require("./KeepRandomBeaconService.sol");
const KeepRandomBeaconServiceImplV1 = artifacts.require("./KeepRandomBeaconServiceImplV1.sol");
const KeepRandomBeaconOperator = artifacts.require("./KeepRandomBeaconOperator.sol");
const KeepRandomBeaconOperatorStatistics = artifacts.require("./statistics/KeepRandomBeaconOperatorStatistics.sol");
const GroupSelection = artifacts.require("./libraries/operator/GroupSelection.sol");
const Groups = artifacts.require("./libraries/operator/Groups.sol");
const DKGResultVerification = artifacts.require("./libraries/operator/DKGResultVerification.sol");
const Reimbursements = artifacts.require("./libraries/operator/Reimbursements.sol");
const Registry = artifacts.require("./Registry.sol");

let initializationPeriod = 518400; // ~6 days
const undelegationPeriod = 7776000; // ~3 months
const withdrawalDelay = 86400; // 1 day
const dkgContributionMargin = 1; // 1%

module.exports = async function(deployer, network) {

  // Set the stake initialization period to 1 block for local development and testnet.
  if (network === 'local' || network === 'ropsten' || network === 'keep_dev') {
    initializationPeriod = 1;
  }

  await deployer.deploy(ModUtils);
  await deployer.link(ModUtils, AltBn128);
  await deployer.deploy(AltBn128);
  await deployer.link(AltBn128, BLS);
  await deployer.deploy(BLS);
  await deployer.deploy(KeepToken);
  await deployer.deploy(Registry);
  await deployer.deploy(TokenStaking, KeepToken.address, Registry.address, initializationPeriod, undelegationPeriod);
  await deployer.deploy(PermissiveStakingPolicy);
  await deployer.deploy(GuaranteedMinimumStakingPolicy, TokenStaking.address);
  await deployer.deploy(TokenGrant, KeepToken.address);
  await deployer.deploy(
    ManagedGrantFactory,
    KeepToken.address,
    TokenGrant.address
  );
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

  const keepRandomBeaconServiceImplV1 = await deployer.deploy(KeepRandomBeaconServiceImplV1);

  const initialize = keepRandomBeaconServiceImplV1.contract.methods
      .initialize(
          dkgContributionMargin,
          withdrawalDelay,
          Registry.address
      ).encodeABI();

  await deployer.deploy(
      KeepRandomBeaconService,
      KeepRandomBeaconServiceImplV1.address,
      initialize
  );

  await deployer.deploy(KeepRandomBeaconOperator, KeepRandomBeaconService.address, TokenStaking.address);
  await deployer.deploy(KeepRandomBeaconOperatorStatistics, KeepRandomBeaconOperator.address);
};
