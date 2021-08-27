const KeepRandomBeaconServiceImplV1 = artifacts.require(
  "KeepRandomBeaconServiceImplV1.sol"
)
const KeepRandomBeaconService = artifacts.require("KeepRandomBeaconService.sol")
const KeepRandomBeaconOperator = artifacts.require(
  "KeepRandomBeaconOperator.sol"
)

const watchRelayEntry = process.env.WATCH_RELAY_ENTRY

module.exports = async function () {
  const keepRandomBeaconService = await KeepRandomBeaconService.deployed()
  const keepRandomBeaconServiceImpl = await KeepRandomBeaconServiceImplV1.at(
    keepRandomBeaconService.address
  )
  const keepRandomBeaconOperator = await KeepRandomBeaconOperator.deployed()

  console.log(
    `Address of the KeepRandomBeaconService contract is ${keepRandomBeaconService.address}`
  )

  try {
    const entryFeeEstimate = await keepRandomBeaconServiceImpl.entryFeeEstimate(
      0
    )
    const tx = await keepRandomBeaconServiceImpl.methods["requestRelayEntry()"](
      {
        value: entryFeeEstimate,
      }
    )
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
    try {
      console.log(`Watch new relay entry generation...`)

      const iterationDelay = 30000 // 30s
      let entryGenerated = false

      // Wait 10 minutes for a relay entry to be generated.
      for (let i = 0; i < 20; i++) {
        await wait(iterationDelay)

        const block = await keepRandomBeaconOperator.currentRequestStartBlock()
        if (web3.utils.toBN(block).isZero()) {
          entryGenerated = true
          break
        }
      }

      if (!entryGenerated) {
        throw new Error(`New relay has not been generated in observed time`)
      }

      console.log(`New relay entry has been generated`)
    } catch (error) {
      console.error("New relay entry watch failed with", error)
      process.exit(1)
    }
  }

  process.exit(0)
}

function wait(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}
