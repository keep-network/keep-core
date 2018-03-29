var KeepRelayBeacon = artifacts.require("./KeepRelayBeacon.sol");

module.exports = function(deployer) {
    deployer.deploy(KeepRelayBeacon);
};
