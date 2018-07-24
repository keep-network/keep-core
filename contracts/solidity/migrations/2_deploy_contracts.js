const KeepToken = artifacts.require("./KeepToken.sol");
const StakingProxy = artifacts.require("./StakingProxy.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const TokenGrant = artifacts.require("./TokenGrant.sol");
const KeepRandomBeaconImplV1 = artifacts.require("./KeepRandomBeaconImplV1.sol");
const KeepGroupImplV1 = artifacts.require("./KeepGroupImplV1.sol");
const KeepGroup = artifacts.require("./KeepGroup.sol");
const KeepRandomBeacon = artifacts.require("./KeepRandomBeacon.sol");

const withdrawalDelay = 86400; // 1 day
const minPayment = 1;
const minStake = 1;

const groupThreshold = 2;
const groupSize = 5;

module.exports = (deployer) => {
  deployer.then(async () => {
    await deployer.deploy(KeepToken);
    await deployer.deploy(StakingProxy);
    await deployer.deploy(TokenStaking, KeepToken.address, StakingProxy.address, withdrawalDelay);
    await deployer.deploy(TokenGrant, KeepToken.address, StakingProxy.address, withdrawalDelay);
    await deployer.deploy(KeepRandomBeaconImplV1);
    await deployer.deploy(KeepRandomBeacon, "v1.0.0", KeepRandomBeaconImplV1.address);
    await deployer.deploy(KeepGroupImplV1);
    await deployer.deploy(KeepGroup, "v1.0.0", KeepGroupImplV1.address);
    await KeepRandomBeaconImplV1.at(KeepRandomBeacon.address).initialize(StakingProxy.address, minPayment, minStake, withdrawalDelay);
    await KeepGroupImplV1.at(KeepGroup.address).initialize(groupThreshold, groupSize, KeepRandomBeaconImplV1.address);
  });
};
