const KeepRandomBeaconProxy = artifacts.require('KeepRandomBeacon.sol');
const KeepGroup = artifacts.require("KeepGroupImplV1");
const KeepRandomBeacon = artifacts.require("KeepRandomBeaconImplV1");

module.exports = async function () {

  const keepRandomBeaconProxy = await KeepRandomBeaconProxy.deployed();
  const keepGroup = await KeepGroup.deployed();

  async function printLastRelayEntry() {
    let contractRef = await KeepRandomBeacon.at(keepRandomBeaconProxy.address);
    let lastEntry = await contractRef.previousEntry();

    console.log('Last relay entry: ' + lastEntry.toString());
  }

  async function printNumberOfGroups() {
    let groupsCount = await keepGroup.numberOfGroups();

    console.log('Number of active groups: ' + groupsCount.toString());
  }

  printLastRelayEntry();
  printNumberOfGroups();
}
