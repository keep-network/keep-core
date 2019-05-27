const KeepRandomBeaconFrontendImplV1 = artifacts.require("KeepRandomBeaconFrontendImplV1");
const KeepRandomBeaconFrontend = artifacts.require('KeepRandomBeaconFrontend.sol');

module.exports = async function() {

  const keepRandomBeaconFrontendProxy = await KeepRandomBeaconFrontend.deployed();

  async function requestRelayEntry() {

    let crypto = require("crypto");
    let KeepRandomBeaconContractAddress = await KeepRandomBeaconFrontendImplV1.at(keepRandomBeaconFrontendProxy.address);

    // Generate 32 byte sort of random number
    try {
      relayEntrySeed = crypto.randomBytes(32);
    }
    catch(error) {
      console.error(error);
    }

    try {
      let requestEntry = await KeepRandomBeaconContractAddress.requestRelayEntry(relayEntrySeed, {value: 2});
      console.log(
        '---Transaction Summary---' + '\n' +
        'From:' + requestEntry.receipt.from + '\n' +
        'To:' + requestEntry.receipt.to + '\n' +
        'BlockNumber:' + requestEntry.receipt.blockNumber + '\n' +
        'TotalGas:' + requestEntry.receipt.cumulativeGasUsed + '\n' +
        'TransactionHash:' + requestEntry.receipt.transactionHash + '\n' +
        '--------------------------'
      );
    }
    catch(error) {
      console.log('Request Failed:');
      console.error(error);
    }
  }

 requestRelayEntry();
}