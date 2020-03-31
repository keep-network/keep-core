const KeepRandomBeaconServiceImplV1 = artifacts.require("KeepRandomBeaconServiceImplV1.sol");
const KeepRandomBeaconService = artifacts.require('KeepRandomBeaconService.sol');

// Example usage:
// truffle exec ./scripts/request-relay-entry-with-callback.js yourContractAddress callbackGas
// truffle exec ./scripts/request-relay-entry-with-callback.js 0x9F57C01059057d821c6b4B04A4598322661C934F 20000

module.exports = async function() {

  const keepRandomBeaconService = await KeepRandomBeaconService.deployed()
  const contractInstance = await KeepRandomBeaconServiceImplV1.at(keepRandomBeaconService.address)

  try {
    let entryFeeEstimate = await contractInstance.entryFeeEstimate(process.argv[5]);
    let tx = await contractInstance.methods['requestRelayEntry(address,uint256)'](
      process.argv[4],
      process.argv[5],
      {value: entryFeeEstimate}
    )
    console.log('Successfully requested relay entry with a callback. RequestId =', tx.logs[0].args.requestId.toString())
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
