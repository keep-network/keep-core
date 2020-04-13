const KeepRandomBeaconService = artifacts.require("./KeepRandomBeaconService.sol");
const KeepRandomBeaconServiceImplV1 = artifacts.require("./KeepRandomBeaconServiceImplV1.sol");
const KeepRandomBeaconOperator = artifacts.require("./KeepRandomBeaconOperator.sol");
const Registry = artifacts.require("./Registry.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const TokenGrant = artifacts.require("./TokenGrant.sol");
let governance, registryKeeper, panicButton, operatorContractUpgrader;

module.exports = async function(deployer, network, accounts) {
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

    const keepRandomBeaconService = await KeepRandomBeaconService.deployed();
    const keepRandomBeaconServiceImplV1 = await KeepRandomBeaconServiceImplV1.at(keepRandomBeaconService.address);
    const keepRandomBeaconOperator = await KeepRandomBeaconOperator.deployed();
    const registry = await Registry.deployed();
    const tokenStaking = await TokenStaking.deployed();
    const tokenGrant = await TokenGrant.deployed();

    if (!(await keepRandomBeaconServiceImplV1.initialized())) {
        throw Error("keep random beacon service not initialized")
    }

    // Set up roles
    await registry.setGovernance(governance)
    await registry.setRegistryKeeper(registryKeeper)
    await registry.setPanicButton(panicButton)

    // Governance as token grand manager has to authorize
    // token staking contract before issuing token grants
    await tokenGrant.authorizeStakingContract(
        tokenStaking.address,
        {from: governance}
    );

    // Only registry keeper can approve operator contracts
    await registry.approveOperatorContract(
        keepRandomBeaconOperator.address,
        {from: registryKeeper}
    );

    // Only governance can set operator contract upgrader for a service contract
    await registry.setOperatorContractUpgrader(
        keepRandomBeaconServiceImplV1.address, operatorContractUpgrader,
        {from: governance}
    );

    // Operator contract upgrader can add operator contract to the service contract
    keepRandomBeaconServiceImplV1.addOperatorContract(
        keepRandomBeaconOperator.address,
        {from: operatorContractUpgrader}
    );

    console.log(
        `\n Successfully deployed contracts and set roles\n\n`,
        `Roles set in Registry.sol\n`,
        `Governance: ${governance} \n`,
        `Registry Keeper: ${registryKeeper} \n`,
        `Panic Button: ${panicButton} \n`,
        `Operator Contract Upgrader: ${operatorContractUpgrader} \n\n`,
        `Roles set in the contracts\n`,
        `KeepToken.sol => Initial token supply is allocated to Governance: ${governance} \n`,
        `TokenGrant.sol => Token grant manager is set to Governance: ${governance} \n`,
        `KeepRandomBeaconService.sol => Implementation upgrader: ${operatorContractUpgrader} \n`,
        `KeepRandomBeaconService.sol => Contract admin (can add/remove operator contracts): ${operatorContractUpgrader} \n`,
        `KeepRandomBeaconService.sol => Withdrawal admin (can withdraw eth from the contract): ${operatorContractUpgrader} \n`,
        `KeepRandomBeaconOperator.sol => Contract admin (can add/remove service contracts from it): ${operatorContractUpgrader} \n\n`,
    )
};
