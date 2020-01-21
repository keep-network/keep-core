const KeepRandomBeaconService = artifacts.require("./KeepRandomBeaconService.sol");
const KeepRandomBeaconServiceImplV1 = artifacts.require("./KeepRandomBeaconServiceImplV1.sol");
const KeepRandomBeaconOperator = artifacts.require("./KeepRandomBeaconOperator.sol");
const RegistryKeeper = artifacts.require("./RegistryKeeper.sol");

const withdrawalDelay = 86400; // 1 day
let priceFeedEstimate = web3.utils.toWei(web3.utils.toBN(20), 'Gwei');
const fluctuationMargin = 50; // 50%
const dkgContributionMargin = 1; // 1%

module.exports = async function(deployer, network) {
    const keepRandomBeaconService = await KeepRandomBeaconServiceImplV1.at(KeepRandomBeaconService.address);
    const keepRandomBeaconOperator = await KeepRandomBeaconOperator.deployed();
    const registryKeeper = await RegistryKeeper.deployed();

    // Set the price feed estimate to 1 Gwei for Ropsten network
    if (network === 'ropsten') {
        priceFeedEstimate = web3.utils.toWei( web3.utils.toBN(1), 'Gwei');
    }

    keepRandomBeaconService.initialize(
        priceFeedEstimate,
        fluctuationMargin,
        dkgContributionMargin,
        withdrawalDelay,
        registryKeeper.address
    );

    await keepRandomBeaconOperator.setPriceFeedEstimate(priceFeedEstimate);
    await registryKeeper.approveOperatorContract(keepRandomBeaconOperator.address);
    keepRandomBeaconService.addOperatorContract(
        keepRandomBeaconOperator.address
    );
};
