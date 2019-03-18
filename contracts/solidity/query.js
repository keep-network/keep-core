const KeepGroupProxy = artifacts.require('KeepGroup.sol');
const KeepRandomBeaconProxy = artifacts.require('KeepRandomBeacon.sol');
const KeepGroup = artifacts.require("KeepGroupImplV1");
const KeepRandomBeacon = artifacts.require("KeepRandomBeaconImplV1");

module.exports = async function () {

  const keepRandomBeaconProxy = await KeepRandomBeaconProxy.deployed();
  const keepGroupProxy = await KeepGroupProxy.deployed();

  async function printLastRelayEntry() {
    let contractRef = await KeepRandomBeacon.at(keepRandomBeaconProxy.address);
    let lastEntry = await contractRef.lastEntryValue();

    console.log('Last relay entry: ' + lastEntry.toString());
  }

  async function printNumberOfGroups() {
    let contractRef = await KeepGroup.at(keepGroupProxy.address);
    let groupsCount = await contractRef.numberOfGroups();

    console.log('Number of active groups: ' + groupsCount.toString());
  }

  printLastRelayEntry()
  printNumberOfGroups();
}
