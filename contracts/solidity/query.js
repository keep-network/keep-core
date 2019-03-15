var keepGroupAddress = "0x6B9e133df21b3F7Ed43f333c3253EfB598D39B4d";
var keepRandomBeaconAddress = "0x1dff256549dBFeeE5116444D25c43392E0166f6A";

var KeepGroup = artifacts.require("KeepGroupImplV1");
var KeepRandomBeacon = artifacts.require("KeepRandomBeaconImplV1");

module.exports = function () {

  async function printLastRelayEntry() {
    let contractRef = await KeepRandomBeacon.at(keepRandomBeaconAddress);
    let lastEntry = await contractRef.lastEntryValue();

    console.log('Last relay entry: ' + lastEntry.toString());
  }

  async function printNumberOfGroups() {
    let contractRef = await KeepGroup.at(keepGroupAddress);
    let groupsCount = await contractRef.numberOfGroups();

    console.log('Number of active groups: ' + groupsCount.toString());
  }

  printLastRelayEntry()
  printNumberOfGroups();
}
