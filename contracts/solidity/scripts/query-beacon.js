const KeepRandomBeaconProxy = artifacts.require('KeepRandomBeacon.sol');
const KeepRandomBeaconBackend = artifacts.require("KeepRandomBeaconBackend");
const KeepRandomBeacon = artifacts.require("KeepRandomBeaconImplV1");

module.exports = async function () {

  const keepRandomBeaconProxy = await KeepRandomBeaconProxy.deployed();
  const keepRandomBeaconBackend = await KeepRandomBeaconBackend.deployed();

  async function printLastRelayEntry() {
    let contractRef = await KeepRandomBeacon.at(keepRandomBeaconProxy.address);
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
