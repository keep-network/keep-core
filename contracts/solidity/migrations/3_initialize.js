const KeepRandomBeaconService = artifacts.require("./KeepRandomBeaconService.sol");
const KeepRandomBeaconServiceImplV1 = artifacts.require("./KeepRandomBeaconServiceImplV1.sol");
const KeepRandomBeaconOperator = artifacts.require("./KeepRandomBeaconOperator.sol");
const Registry = artifacts.require("./Registry.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const TokenGrant = artifacts.require("./TokenGrant.sol");

const withdrawalDelay = 86400; // 1 day
const dkgContributionMargin = 1; // 1%

module.exports = async function(deployer, network) {
    const keepRandomBeaconService = await KeepRandomBeaconServiceImplV1.at(KeepRandomBeaconService.address);
    const keepRandomBeaconOperator = await KeepRandomBeaconOperator.deployed();
    const registry = await Registry.deployed();
    const tokenStaking = await TokenStaking.deployed();
    const tokenGrant = await TokenGrant.deployed();

    await tokenGrant.authorizeStakingContract(tokenStaking.address);

    keepRandomBeaconService.initialize(
        dkgContributionMargin,
        withdrawalDelay,
        registry.address
    );

    await registry.approveOperatorContract(keepRandomBeaconOperator.address);

    // Set service contract owner as operator contract upgrader by default
    const operatorContractUpgrader = await keepRandomBeaconService.owner();
    await registry.setOperatorContractUpgrader(keepRandomBeaconService.address, operatorContractUpgrader);
    keepRandomBeaconService.addOperatorContract(
        keepRandomBeaconOperator.address,
        {from: operatorContractUpgrader}
    );
};
