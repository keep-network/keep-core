const KeepRandomBeaconService = artifacts.require("./KeepRandomBeaconService.sol");
const KeepRandomBeaconServiceImplV1 = artifacts.require("./KeepRandomBeaconServiceImplV1.sol");
const KeepRandomBeaconOperator = artifacts.require("./KeepRandomBeaconOperator.sol");

const withdrawalDelay = 86400; // 1 day
let priceFeedEstimate = web3.utils.toBN(20).mul(web3.utils.toBN(10**9)); // (20 Gwei = 20 * 10^9 wei)
const fluctuationMargin = 50; // 50%
const dkgContributionMargin = 1; // 1%

module.exports = async function(deployer, network) {
    const keepRandomBeaconService = await KeepRandomBeaconServiceImplV1.at(KeepRandomBeaconService.address);
    const keepRandomBeaconOperator = await KeepRandomBeaconOperator.deployed();

    // Set the price feed estimate to 1 Gwei for Ropsten network
    if (network === 'ropsten') {
        priceFeedEstimate = web3.utils.toBN(1).mul(web3.utils.toBN(10 ** 9)); // (1 Gwei = 1 * 10^9 wei)
    }

    keepRandomBeaconService.initialize(
        priceFeedEstimate,
        fluctuationMargin,
        dkgContributionMargin,
        withdrawalDelay,
        keepRandomBeaconOperator.address
    );

    await keepRandomBeaconOperator.setPriceFeedEstimate(priceFeedEstimate);
};