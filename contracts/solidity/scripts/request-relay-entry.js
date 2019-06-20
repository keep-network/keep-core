const crypto = require("crypto")
const KeepRandomBeacon = artifacts.require("KeepRandomBeaconImplV1")
const KeepRandomBeaconProxy = artifacts.require('KeepRandomBeacon.sol')

module.exports = async function() {

  const keepRandomBeaconProxy = await KeepRandomBeaconProxy.deployed()
  const contractInstance = await KeepRandomBeacon.at(keepRandomBeaconProxy.address)

  try {
    let tx = await contractInstance.requestRelayEntry(crypto.randomBytes(32), {value: 2})
    console.log('Successfully requested relay entry with RequestId =', tx.logs[0].args.requestID.toString())
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
