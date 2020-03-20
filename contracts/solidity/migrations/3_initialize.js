const KeepRandomBeaconService = artifacts.require("./KeepRandomBeaconService.sol");
const KeepRandomBeaconServiceImplV1 = artifacts.require("./KeepRandomBeaconServiceImplV1.sol");
const KeepRandomBeaconOperator = artifacts.require("./KeepRandomBeaconOperator.sol");
const Registry = artifacts.require("./Registry.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const TokenGrant = artifacts.require("./TokenGrant.sol");

module.exports = async function(deployer, network) {
    const keepRandomBeaconService = await KeepRandomBeaconService.deployed();
    const keepRandomBeaconServiceImplV1 = await KeepRandomBeaconServiceImplV1.at(keepRandomBeaconService.address);
    const keepRandomBeaconOperator = await KeepRandomBeaconOperator.deployed();
    const registry = await Registry.deployed();
    const tokenStaking = await TokenStaking.deployed();
    const tokenGrant = await TokenGrant.deployed();

    if (!(await keepRandomBeaconServiceImplV1.initialized())) {
        throw Error("keep random beacon service not initialized")
    }

    await tokenGrant.authorizeStakingContract(tokenStaking.address);

    await registry.approveOperatorContract(keepRandomBeaconOperator.address);

    // Set service contract owner as operator contract upgrader by default
    const operatorContractUpgrader = await keepRandomBeaconService.admin();
    await registry.setOperatorContractUpgrader(keepRandomBeaconServiceImplV1.address, operatorContractUpgrader);
    keepRandomBeaconServiceImplV1.addOperatorContract(
        keepRandomBeaconOperator.address,
        {from: operatorContractUpgrader}
    );
};
