const KeepRandomBeaconServiceImplV1 = artifacts.require("KeepRandomBeaconServiceImplV1");
const KeepRandomBeaconService = artifacts.require('KeepRandomBeaconService.sol');

module.exports = async function() {

  const keepRandomBeaconServiceProxy = await KeepRandomBeaconService.deployed();

  async function requestRelayEntry() {

    let crypto = require("crypto");
    let KeepRandomBeaconContractAddress = await KeepRandomBeaconServiceImplV1.at(keepRandomBeaconServiceProxy.address);

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