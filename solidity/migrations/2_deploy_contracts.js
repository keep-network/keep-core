const KeepToken = artifacts.require("./KeepToken.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const TokenVesting = artifacts.require("./TokenVesting.sol");

module.exports = function(deployer) {
  deployer.deploy(KeepToken).then(function() {
    return deployer.deploy(TokenVesting, KeepToken.address);
  });
};
