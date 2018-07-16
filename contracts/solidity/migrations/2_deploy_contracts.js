const KeepToken = artifacts.require("./KeepToken.sol");
const StakingProxy = artifacts.require("./StakingProxy.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const TokenGrant = artifacts.require("./TokenGrant.sol");
const KeepRandomBeaconImplV1 = artifacts.require("./KeepRandomBeaconImplV1.sol");
const KeepRandomBeacon = artifacts.require("./KeepRandomBeacon.sol");
const KeepGroupImplV1 = artifacts.require("./KeepGroupImplV1.sol");
const KeepGroup = artifacts.require("./KeepGroup.sol");

const withdrawalDelay = 86400; // 1 day
const minPayment = 1;
const minStake = 1;

module.exports = function(deployer) {
  deployer.deploy(KeepToken)
    .then(function() {
      return deployer.deploy(StakingProxy);
    }).then(function() {
      return deployer.deploy(TokenStaking, KeepToken.address, StakingProxy.address, withdrawalDelay);
    }).then(function() {
      return deployer.deploy(TokenGrant, KeepToken.address, StakingProxy.address, withdrawalDelay);
    }).then(function() {
      return deployer.deploy(KeepRandomBeaconImplV1, StakingProxy.address, minPayment, minStake);
    }).then(function() {
      return deployer.deploy(KeepRandomBeacon, 'v1', KeepRandomBeaconImplV1.address);
    }).then(function() {
      return deployer.deploy(KeepGroupImplV1);
    }).then(function() {
      return deployer.deploy(KeepGroup, 'v1', KeepGroupImplV1.address);
    });
};
