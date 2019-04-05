const KeepToken = artifacts.require("./KeepToken.sol");
const ModUtils = artifacts.require("./utils/ModUtils.sol");
const AltBn128 = artifacts.require("./AltBn128.sol");
const BLS = artifacts.require("./BLS.sol");
const StakingProxy = artifacts.require("./StakingProxy.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const TokenGrant = artifacts.require("./TokenGrant.sol");
const KeepRandomBeaconImplV1 = artifacts.require("./KeepRandomBeaconImplV1.sol");
const KeepRandomBeaconUpgradeExample = artifacts.require("./KeepRandomBeaconUpgradeExample.sol");
const KeepGroupImplV1 = artifacts.require("./KeepGroupImplV1.sol");
const KeepGroup = artifacts.require("./KeepGroup.sol");
const KeepRandomBeacon = artifacts.require("./KeepRandomBeacon.sol");

const withdrawalDelay = 86400; // 1 day
const minPayment = 1;
const minStake = web3.utils.toBN(200000).mul(web3.utils.toBN(10**18));

const groupThreshold = 3;
const groupSize = 5;
const timeoutInitial = 4;
const timeoutSubmission = 4;
const timeoutChallenge = 4;
const dkgSubmissionTimeout = 4;

module.exports = async function(deployer) {
  await deployer.deploy(ModUtils);
  await deployer.link(ModUtils, AltBn128);
  await deployer.deploy(AltBn128);
  await deployer.link(AltBn128, BLS);
  await deployer.deploy(BLS);
  await deployer.deploy(KeepToken);
  await deployer.deploy(StakingProxy);
  await deployer.deploy(TokenStaking, KeepToken.address, StakingProxy.address, withdrawalDelay);
  await deployer.deploy(TokenGrant, KeepToken.address, StakingProxy.address, withdrawalDelay);
  await deployer.link(BLS, KeepRandomBeaconImplV1);
  await deployer.link(BLS, KeepRandomBeaconUpgradeExample);
  await deployer.deploy(KeepRandomBeaconImplV1);
  await deployer.deploy(KeepRandomBeacon, KeepRandomBeaconImplV1.address);
  await deployer.deploy(KeepGroupImplV1);
  await deployer.deploy(KeepGroup, KeepGroupImplV1.address);

  const keepRandomBeacon = await KeepRandomBeaconImplV1.at(KeepRandomBeacon.address);
  const keepGroup = await KeepGroupImplV1.at(KeepGroup.address);
  await keepGroup.initialize(
    StakingProxy.address, KeepRandomBeacon.address, minStake, groupThreshold, groupSize, timeoutInitial, timeoutSubmission, timeoutChallenge, dkgSubmissionTimeout
  );
  // Initialize contract genesis entry value and genesis group defined in Go client submitGenesisRelayEntry()
  await keepRandomBeacon.initialize(
    minPayment,
    withdrawalDelay,
    web3.utils.toBN('31415926535897932384626433832795028841971693993751058209749445923078164062862'),
    "0x1f1954b33144db2b5c90da089e8bde287ec7089d5d6433f3b6becaefdb678b1b2a9de38d14bef2cf9afc3c698a4211fa7ada7b4f036a2dfef0dc122b423259d0",
    KeepGroup.address
  );
};
