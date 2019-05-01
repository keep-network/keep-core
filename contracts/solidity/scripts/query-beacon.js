const KeepRandomBeaconFrontendProxy = artifacts.require('KeepRandomBeaconFrontendProxy.sol');
const KeepRandomBeaconBackend = artifacts.require("KeepRandomBeaconBackend");
const KeepRandomBeaconFrontend = artifacts.require("KeepRandomBeaconFrontendImplV1");

module.exports = async function () {

  const keepRandomBeaconFrontendProxy = await KeepRandomBeaconFrontendProxy.deployed();
  const keepRandomBeaconBackend = await KeepRandomBeaconBackend.deployed();

  async function printLastRelayEntry() {
    let contractRef = await KeepRandomBeaconFrontend.at(keepRandomBeaconFrontendProxy.address);
    let lastEntry = await contractRef.previousEntry();

    console.log('Last relay entry: ' + lastEntry.toString());
  }

  async function printNumberOfGroups() {
    let groupsCount = await keepRandomBeaconBackend.numberOfGroups();

    console.log('Number of active groups: ' + groupsCount.toString());
  }

  printLastRelayEntry();
  printNumberOfGroups();
}
