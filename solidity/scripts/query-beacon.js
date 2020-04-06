const KeepRandomBeaconService = artifacts.require('KeepRandomBeaconService.sol');
const KeepRandomBeaconOperator = artifacts.require("KeepRandomBeaconOperator.sol");
const KeepRandomBeaconServiceImplV1 = artifacts.require("KeepRandomBeaconServiceImplV1.sol");

module.exports = async function () {

  const keepRandomBeaconService = await KeepRandomBeaconService.deployed();
  const keepRandomBeaconOperator = await KeepRandomBeaconOperator.deployed();

  async function printLastRelayEntry() {
    let contractRef = await KeepRandomBeaconServiceImplV1.at(keepRandomBeaconService.address);
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
