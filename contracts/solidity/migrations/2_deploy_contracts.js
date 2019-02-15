const KeepToken = artifacts.require("./KeepToken.sol");
const ModUtils = artifacts.require("./utils/ModUtils.sol");
const AltBn128 = artifacts.require("./AltBn128.sol");
const BLS = artifacts.require("./BLS.sol");
const StakingProxy = artifacts.require("./StakingProxy.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const TokenGrant = artifacts.require("./TokenGrant.sol");
const KeepRandomBeaconImplV1 = artifacts.require("./KeepRandomBeaconImplV1.sol");
const KeepGroupImplV1 = artifacts.require("./KeepGroupImplV1.sol");
const KeepGroup = artifacts.require("./KeepGroup.sol");
const KeepRandomBeacon = artifacts.require("./KeepRandomBeacon.sol");

const withdrawalDelay = 86400; // 1 day
const minPayment = 1;
const minStake = 200000 * (10**18);

const groupThreshold = 3;
const groupSize = 5;
const timeoutInitial = 4;
const timeoutSubmission = 4;
const timeoutChallenge = 4;

module.exports = (deployer) => {
  deployer.then(async () => {
    await deployer.deploy(ModUtils);
    await deployer.link(ModUtils, AltBn128);
    await deployer.deploy(AltBn128);
    await deployer.link(AltBn128, BLS);
    await deployer.deploy(BLS);
    await deployer.deploy(KeepToken);
    await deployer.deploy(StakingProxy);
    await deployer.deploy(TokenStaking, KeepToken.address, StakingProxy.address, withdrawalDelay);
    await deployer.deploy(TokenGrant, KeepToken.address, StakingProxy.address, withdrawalDelay);
    await deployer.deploy(KeepRandomBeaconImplV1);
    await deployer.deploy(KeepRandomBeacon, KeepRandomBeaconImplV1.address);
    await deployer.deploy(KeepGroupImplV1);
    await deployer.deploy(KeepGroup, KeepGroupImplV1.address);
    await KeepRandomBeaconImplV1.at(KeepRandomBeacon.address).initialize(minPayment, withdrawalDelay);
    await KeepGroupImplV1.at(KeepGroup.address).initialize(
      StakingProxy.address, KeepRandomBeacon.address, minStake, groupThreshold, groupSize, timeoutInitial, timeoutSubmission, timeoutChallenge
    );
    await KeepRandomBeaconImplV1.at(KeepRandomBeacon.address).setGroupContract(KeepGroup.address);
  });
};
