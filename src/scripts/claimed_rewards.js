// Script that returns the amount of already claimed Threshold rewards
// Use: node src/scripts/gen_merkle_dist.js

require("dotenv").config()
const axios = require("axios")
const keccak256 = require("keccak256")

const merkleDistContractAdd = "0xeA7CA290c7811d1cC2e79f8d706bD05d8280BD37"
const eventTopic =
  "0x" + keccak256("Claimed(address,uint256,address,bytes32)").toString("hex")

axios
  .get("https://api.etherscan.io/api", {
    params: {
      apikey: process.env.ETHERSCAN_TOKEN,
      module: "logs",
      action: "getLogs",
      address: merkleDistContractAdd,
      fromBlock: 15146501,
      page: 1,
      offset: 1000,
      topic0: eventTopic,
    },
  })
  .then(function (response) {
    if (response.data.status === "1") {
      const events = response.data.result
      const claimedAmount = events.reduce((cum, event) => {
        return cum + BigInt(event.data.slice(0, 66))
      }, BigInt(0))
      console.log("Threshold rewards already claimed:")
      console.log(claimedAmount.toString())
    } else {
      console.error("Error: " + response.data.message)
      console.error(response.data.result)
    }
  })
  .catch(function (error) {
    console.error(error.toJSON())
  })
