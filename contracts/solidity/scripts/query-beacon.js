const KeepRandomBeaconFrontend = artifacts.require('KeepRandomBeaconFrontend.sol');
const KeepRandomBeaconOperator = artifacts.require("KeepRandomBeaconOperator");
const KeepRandomBeaconFrontendImplV1 = artifacts.require("KeepRandomBeaconFrontendImplV1");

module.exports = async function () {

  const keepRandomBeaconFrontendProxy = await KeepRandomBeaconFrontend.deployed();
  const keepRandomBeaconOperator = await KeepRandomBeaconOperator.deployed();

  async function printLastRelayEntry() {
    let contractRef = await KeepRandomBeaconFrontendImplV1.at(keepRandomBeaconFrontendProxy.address);
    let lastEntry = await contractRef.previousEntry();

    console.log('Last relay entry: ' + lastEntry.toString());
  }

  async function printNumberOfGroups() {
    let groupsCount = await keepRandomBeaconOperator.numberOfGroups();

    console.log('Number of active groups: ' + groupsCount.toString());
  }

  printLastRelayEntry();
  printNumberOfGroups();
}
