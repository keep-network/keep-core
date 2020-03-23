const crypto = require("crypto")
const KeepRandomBeaconServiceImplV1 = artifacts.require("KeepRandomBeaconServiceImplV1.sol");
const KeepRandomBeaconService = artifacts.require('KeepRandomBeaconService.sol');

module.exports = async function() {

  const keepRandomBeaconService = await KeepRandomBeaconService.deployed()
  const contractInstance = await KeepRandomBeaconServiceImplV1.at(keepRandomBeaconService.address)

  try {
    let entryFeeEstimate = await contractInstance.entryFeeEstimate(0);
    let tx = await contractInstance.methods['requestRelayEntry()']({value: entryFeeEstimate});
    console.log('Successfully requested relay entry with RequestId =', tx.logs[0].args.requestId.toString())
    console.log(
      '\n---Transaction Summary---' + '\n' +
      'From:' + tx.receipt.from + '\n' +
      'To:' + tx.receipt.to + '\n' +
      'BlockNumber:' + tx.receipt.blockNumber + '\n' +
      'TotalGas:' + tx.receipt.cumulativeGasUsed + '\n' +
      'TransactionHash:' + tx.receipt.transactionHash + '\n' +
      '--------------------------'
    )
  } catch(error) {
    console.error('Request failed with', error)
  }

  process.exit()
}
