const KeepToken = artifacts.require("./KeepToken.sol");
const StakingProxy = artifacts.require("./StakingProxy.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const TokenGrant = artifacts.require("./TokenGrant.sol");

const withdrawalDelay = 86400; // 1 day
module.exports = function(deployer) {
  deployer.deploy(KeepToken)
    .then(function() {
      return deployer.deploy(StakingProxy);
    }).then(function() {
      return deployer.deploy(TokenStaking, KeepToken.address, StakingProxy.address, withdrawalDelay);
    }).then(function() {
      return deployer.deploy(TokenGrant, KeepToken.address, StakingProxy.address, withdrawalDelay);
    });
};
