const KeepRandomBeaconFrontend = artifacts.require('KeepRandomBeaconFrontend.sol');
const KeepRandomBeaconBackend = artifacts.require("KeepRandomBeaconBackend");
const KeepRandomBeaconFrontendImplV1 = artifacts.require("KeepRandomBeaconFrontendImplV1");

module.exports = async function () {

  const keepRandomBeaconFrontendProxy = await KeepRandomBeaconFrontend.deployed();
  const keepRandomBeaconBackend = await KeepRandomBeaconBackend.deployed();

  async function printLastRelayEntry() {
    let contractRef = await KeepRandomBeaconFrontendImplV1.at(keepRandomBeaconFrontendProxy.address);
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
