// Script that generates a new Merkle Distribution and outputs the data to JSON files

const fs = require("fs")
const subgraph = require("../src/stakingrewards/subgraph.js")
const rewards = require("../src/stakingrewards/rewards.js")

const graphqlApi =
  "https://api.studio.thegraph.com/query/24143/main-threshold-subgraph/0.0.7"
const bonusWeight = 1.0 // Stakes receive the full bonus
const ongoingWeight = 1.0 // Stakes receive the full ongoing rewards
const startTime = 1654041600 // Jun 1st 2022 00:00:00 GMT
const endTime = 1664496000 // Sep 30th 2022 00:00:00 GMT
const endTimeDate = new Date(endTime * 1000).toISOString().slice(0, 10)
const distribution_path = "distributions/" + endTimeDate

async function main() {
  try {
    fs.mkdirSync(distribution_path)
  } catch (err) {
    console.error(err)
    return
  }

  const ongoingStakes = await subgraph.getOngoingStakes(
    graphqlApi,
    startTime,
    endTime
  )
  const ongoingRewards = await rewards.calculateOngoingRewards(ongoingStakes, ongoingWeight)
  const bonusStakes = await subgraph.getBonusStakes(graphqlApi)
  const bonusRewards = rewards.calculateBonusRewards(bonusStakes, bonusWeight)
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
