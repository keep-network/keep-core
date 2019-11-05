const KeepToken = artifacts.require("./KeepToken.sol");
const ModUtils = artifacts.require("./utils/ModUtils.sol");
const AltBn128 = artifacts.require("./cryptography/AltBn128.sol");
const BLS = artifacts.require("./cryptography/BLS.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const TokenGrant = artifacts.require("./TokenGrant.sol");
const KeepRandomBeaconService = artifacts.require("./KeepRandomBeaconService.sol");
const KeepRandomBeaconServiceImplV1 = artifacts.require("./KeepRandomBeaconServiceImplV1.sol");
const KeepRandomBeaconOperator = artifacts.require("./KeepRandomBeaconOperator.sol");
const KeepRandomBeaconOperatorGroups = artifacts.require("./KeepRandomBeaconOperatorGroups.sol");
const GroupSelection = artifacts.require("./libraries/GroupSelection.sol");

const withdrawalDelay = 86400; // 1 day
const priceFeedEstimate = web3.utils.toBN(20).mul(web3.utils.toBN(10**9)); // (20 Gwei = 20 * 10^9 wei)
const fluctuationMargin = 50; // 50%
const dkgContributionMargin = 10; // 10%

module.exports = async function(deployer) {
  await deployer.deploy(ModUtils);
  await deployer.link(ModUtils, AltBn128);
  await deployer.deploy(AltBn128);
  await deployer.link(AltBn128, BLS);
  await deployer.deploy(BLS);
  await deployer.deploy(KeepToken);
  await deployer.deploy(TokenStaking, KeepToken.address, withdrawalDelay);
  await deployer.deploy(TokenGrant, KeepToken.address, TokenStaking.address);
  await deployer.deploy(GroupSelection);
  await deployer.link(GroupSelection, KeepRandomBeaconOperator);
  await deployer.link(BLS, KeepRandomBeaconOperator);
  await deployer.deploy(KeepRandomBeaconServiceImplV1);
  await deployer.deploy(KeepRandomBeaconService, KeepRandomBeaconServiceImplV1.address);
  await deployer.deploy(KeepRandomBeaconOperatorGroups);

  // TODO: replace with a secure authorization protocol (addressed in RFC 11).
  await deployer.deploy(KeepRandomBeaconOperator, KeepRandomBeaconService.address, TokenStaking.address, KeepRandomBeaconOperatorGroups.address);

  const keepRandomBeaconService = await KeepRandomBeaconServiceImplV1.at(KeepRandomBeaconService.address);
  const keepRandomBeaconOperator = await KeepRandomBeaconOperator.deployed();

  const keepRandomBeaconOperatorGroups = await KeepRandomBeaconOperatorGroups.deployed();
  await keepRandomBeaconOperatorGroups.setOperatorContract(keepRandomBeaconOperator.address);

  keepRandomBeaconService.initialize(
    priceFeedEstimate,
    fluctuationMargin,
    dkgContributionMargin,
    withdrawalDelay,
    keepRandomBeaconOperator.address
  );
};
