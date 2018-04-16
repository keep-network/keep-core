const KeepToken = artifacts.require("./KeepToken.sol");
const StakingProxy = artifacts.require("./StakingProxy.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const TokenGrant = artifacts.require("./TokenGrant.sol");
const KeepRandomBeaconImpl = artifacts.require("./KeepRandomBeaconImpl.sol");

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
      return deployer.deploy(KeepRandomBeaconImpl, StakingProxy.address, minPayment, minStake);
    });
};
