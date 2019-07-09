const crypto = require("crypto")
const KeepRandomBeaconService = artifacts.require("KeepRandomBeaconService.sol");
const KeepRandomBeaconServiceProxy = artifacts.require('KeepRandomBeaconServiceProxy.sol');

module.exports = async function() {

  const keepRandomBeaconServiceProxy = await KeepRandomBeaconServiceProxy.deployed()
  const contractInstance = await KeepRandomBeaconService.at(keepRandomBeaconServiceProxy.address)

  try {
    let tx = await contractInstance.requestRelayEntry(crypto.randomBytes(32), {value: 2})
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
