// Script that generates a new Merkle Distribution and outputs the data to JSON files
// Use: node scripts/gen_merkle_dist.js

const fs = require("fs")
const shell = require("shelljs")
const dotenv = require("dotenv").config()
const Subgraph = require("../src/stakingrewards/subgraph.js")
const Rewards = require("../src/stakingrewards/rewards.js")
const MerkleDist = require("../src/stakingrewards/merkle_dist.js")

// The following parameters must be modified for each distribution
const bonusWeight = 0.0
const ongoingWeight = 0.875
const tbtcv2Weight = 0.125
const startTime = 1664496000 // Sep 30th 2022 00:00:00 GMT
const endTime = 1667260800 // Nov 1st 2022 00:00:00 GMT
const lastDistribution = "2022-09-30"

const tbtcv2ScriptPath = "src/tbtcv2-rewards/"
const graphqlApi =
  "https://api.studio.thegraph.com/query/24143/main-threshold-subgraph/0.0.7"

async function main() {
  let earnedBonusRewards = {}
  let earnedOngoingRewards = {}
  let earnedTbtcv2Rewards = {}
  let bonusRewards = {}
  let ongoingRewards = {}
  let tbtcv2Rewards = {}
  const endDate = new Date(endTime * 1000).toISOString().slice(0, 10)
  const distPath = `distributions/${endDate}`
  const lastDistPath = `distributions/${lastDistribution}`
  const tbtcv2Script =
  `./rewards.sh ` +
  `--rewards-start-date ${startTime} ` +
  `--rewards-end-date ${endTime} ` +
  `--etherscan-token ${process.env.ETHERSCAN_TOKEN}`

  try {
    fs.mkdirSync(distPath)
  } catch (err) {
    console.error(err)
    return
  }

  if (bonusWeight > 0) {
    console.log("Calculating bonus rewards...")
    const bonusStakes = await Subgraph.getBonusStakes(graphqlApi)
    earnedBonusRewards = Rewards.calculateBonusRewards(bonusStakes, bonusWeight)
  }

  if (ongoingWeight > 0) {
    console.log("Calculating ongoing rewards...")
    const ongoingStakes = await Subgraph.getOngoingStakes(
      graphqlApi,
      startTime,
      endTime
    )
    earnedOngoingRewards = await Rewards.calculateOngoingRewards(
      ongoingStakes,
      ongoingWeight
    )
  }

  if (tbtcv2Weight > 0) {
    console.log("Calculating tBTCv2 rewards...")
    shell.exec(`cd ${tbtcv2ScriptPath} && ${tbtcv2Script}`)
    const tbtcv2RewardsRaw = JSON.parse(
      fs.readFileSync("./src/tbtcv2-rewards/rewards.json")
    )
    earnedTbtcv2Rewards = Rewards.calculateTbtcv2Rewards(
      tbtcv2RewardsRaw,
      tbtcv2Weight
    )
  }

  try {
    bonusRewards = JSON.parse(fs.readFileSync(`${lastDistPath}/MerkleInputBonusRewards.json`))
    bonusRewards = MerkleDist.combineMerkleInputs(bonusRewards, earnedBonusRewards)
    fs.writeFileSync(
      distPath + "/MerkleInputBonusRewards.json",
      JSON.stringify(bonusRewards, null, 4)
    )
    ongoingRewards = JSON.parse(fs.readFileSync(`${lastDistPath}/MerkleInputOngoingRewards.json`))
    ongoingRewards = MerkleDist.combineMerkleInputs(ongoingRewards, earnedOngoingRewards)
    fs.writeFileSync(
      distPath + "/MerkleInputOngoingRewards.json",
      JSON.stringify(ongoingRewards, null, 4)
    )
    if(fs.existsSync(`${lastDistPath}/MerkleInputTbtcv2Rewards.json`)) {
      tbtcv2Rewards = JSON.parse(fs.readFileSync(`${lastDistPath}/MerkleInputTbtcv2Rewards.json`))
    } else {
      tbtcv2Rewards = {}
    }
    tbtcv2Rewards = MerkleDist.combineMerkleInputs(tbtcv2Rewards, earnedTbtcv2Rewards)
    fs.writeFileSync(
      distPath + "/MerkleInputTbtcv2Rewards.json",
      JSON.stringify(earnedTbtcv2Rewards, null, 4)
    )
  } catch (err) {
    console.error(err)
    return
  }

  let merkleInput = MerkleDist.combineMerkleInputs(bonusRewards, ongoingRewards)
  merkleInput = MerkleDist.combineMerkleInputs(merkleInput, tbtcv2Rewards)

  const merkleDist = MerkleDist.genMerkleDist(merkleInput)

  try {
    fs.writeFileSync(
      distPath + "/MerkleInputTotalRewards.json",
      JSON.stringify(merkleInput, null, 4)
    )
    fs.writeFileSync(
      distPath + "/MerkleDist.json",
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
