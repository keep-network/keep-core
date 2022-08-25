// Script that generates a new Merkle Distribution and outputs the data to JSON files

const fs = require("fs")
const stakingRewards = require("../src/stakingrewards/stakingrewards.js")

const graphqlApi =
  "https://api.studio.thegraph.com/query/24143/main-threshold-subgraph/0.0.7"
const startTime = 1654041600  // Jun 1st 2022 00:00:00 GMT
const endTime = 1664496000    // Sep 30th 2022 00:00:00 GMT
const endTimeDate = new Date(endTime * 1000).toISOString().slice(0, 10)
const distribution_path = "distributions/" + endTimeDate

async function main() {
  try {
    fs.mkdirSync(distribution_path)
  } catch (err) {
    console.error(err)
    return
  }

  const ongoingRewards = await stakingRewards.getOngoingMekleInput(
    graphqlApi,
    startTime,
    endTime
  )
  const bonusRewards = await stakingRewards.getBonusMerkleInput(graphqlApi)
  const merkleInput = stakingRewards.combineMerkleInputs(
    ongoingRewards,
    bonusRewards
  )
  const merkleDist = stakingRewards.genMerkleDist(merkleInput)

  try{
    fs.writeFileSync(
      distribution_path + "/MerkleInputOngoingRewards.json",
      JSON.stringify(ongoingRewards, null, 4)
    )
    fs.writeFileSync(
      distribution_path + "/MerkleInputBonusRewards.json",
      JSON.stringify(bonusRewards, null, 4)
    )
    fs.writeFileSync(
      distribution_path + "/MerkleInputTotalRewards.json",
      JSON.stringify(merkleInput, null, 4)
    )
    fs.writeFileSync(
      distribution_path + "/MerkleDist.json",
      JSON.stringify(merkleDist, null, 4)
    )
  } catch (err) {
    console.error(err)
    return
  }

  console.log("Total amount of rewards: ", merkleDist.totalAmount)
}

;(async () => {
  await main()
})()
