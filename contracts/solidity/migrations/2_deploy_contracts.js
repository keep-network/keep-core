const KeepToken = artifacts.require("./KeepToken.sol");
const StakingProxy = artifacts.require("./StakingProxy.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const TokenGrant = artifacts.require("./TokenGrant.sol");
const KeepRandomBeaconImplV1 = artifacts.require("./KeepRandomBeaconImplV1.sol");
const KeepGroupImplV1 = artifacts.require("./KeepGroupImplV1.sol");
const EternalStorage = artifacts.require("./EternalStorage.sol");
const KeepGroup = artifacts.require("./KeepGroup.sol");
const KeepRandomBeacon = artifacts.require("./KeepRandomBeacon.sol");

const withdrawalDelay = 86400; // 1 day
const minPayment = 1;
const minStake = 1;

module.exports = function(deployer) {
  deployer.deploy(KeepToken)
    .then(function() {
      return deployer.deploy(StakingProxy);
    }).then(function() {
      return deployer.deploy(EternalStorage);
    }).then(function() {
      return deployer.deploy(TokenStaking, KeepToken.address, StakingProxy.address, withdrawalDelay);
    }).then(function() {
      return deployer.deploy(TokenGrant, KeepToken.address, StakingProxy.address, withdrawalDelay);
    }).then(function() {
      return deployer.deploy(KeepRandomBeaconImplV1, StakingProxy.address, minPayment, minStake);
    }).then(function() {
      return deployer.deploy(KeepGroupImplV1);
    }).then(function() {
      return deployer.deploy(KeepGroup, "v1.0.0", KeepGroupImplV1.address); // TODO - constants
    }).then(function() {
      return deployer.deploy(KeepRandomBeacon, "v1.0.0", KeepRandomBeaconImplV1.address ); // TODO - constants
    }).then(function() {
	  return KeepGroupImplV1.new(KeepGroupImplV1.address);
	}).then(function(instance) {
	  return instance.initialize(10, 4, KeepRandomBeaconImplV1.address); // TODO - really should have constants
    }).then(function() {
	  return KeepRandomBeaconImplV1.new(KeepRandomBeaconImplV1.address);
	}).then(function(instance) {
	  return instance.initialize(StakingProxy.address, 0, 0, 0); // TODO - really should have constants
	});
};

