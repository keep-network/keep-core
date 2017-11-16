const KeepToken = artifacts.require("./KeepToken.sol");

module.exports = function(deployer) {
  deployer.deploy(KeepToken);
};
