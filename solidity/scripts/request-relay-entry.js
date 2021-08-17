const KeepRandomBeaconServiceImplV1 = artifacts.require(
  "KeepRandomBeaconServiceImplV1.sol"
)
const KeepRandomBeaconService = artifacts.require("KeepRandomBeaconService.sol")

const watchRelayEntry = process.env.WATCH_RELAY_ENTRY

module.exports = async function () {
  const keepRandomBeaconService = await KeepRandomBeaconService.deployed()
  const contractInstance = await KeepRandomBeaconServiceImplV1.at(
    keepRandomBeaconService.address
  )

  try {
    const entryFeeEstimate = await contractInstance.entryFeeEstimate(0)
    const tx = await contractInstance.methods["requestRelayEntry()"]({
      value: entryFeeEstimate,
    })
    console.log(
      "Successfully requested relay entry with RequestId =",
      tx.logs[0].args.requestId.toString()
    )
    console.log(
      "\n---Transaction Summary---" +
        "\n" +
        "From:" +
        tx.receipt.from +
        "\n" +
        "To:" +
        tx.receipt.to +
        "\n" +
        "BlockNumber:" +
        tx.receipt.blockNumber +
        "\n" +
        "TotalGas:" +
        tx.receipt.cumulativeGasUsed +
        "\n" +
        "TransactionHash:" +
        tx.receipt.transactionHash +
        "\n" +
        "--------------------------"
    )
  } catch (error) {
    console.error("Request failed with", error)
    process.exit(1)
  }

  if (watchRelayEntry === "true") {
    console.log(`Watch new relay entry generation...`)

    const event = await watchRelayEntryGenerated(contractInstance)
    const newRelayEntry = web3.utils.toBN(event.returnValues.entry)

    console.log(`New relay entry has been generated: ${newRelayEntry}`)
  }

  process.exit(0)
}

function watchRelayEntryGenerated(keepRandomBeaconService) {
  return new Promise(async (resolve) => {
    keepRandomBeaconService.RelayEntryGenerated().on("data", (event) => {
      resolve(event)
    })
  })
}
