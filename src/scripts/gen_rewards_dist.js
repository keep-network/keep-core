// Script that generates a new Merkle Distribution for rewards and outputs the
// data to JSON files
// Use: node src/scripts/gen_merkle_dist.js

require("dotenv").config()
const fs = require("fs")
const shell = require("shelljs")
const Subgraph = require("../pre-rewards/subgraph.js")
const Rewards = require("../pre-rewards/rewards.js")
const MerkleDist = require("../merkle_dist/merkle_dist.js")

// The following parameters must be modified for each distribution
const bonusWeight = 0.0
const preWeight = 0.875
const tbtcv2Weight = 0.125
const startTime = 1664496000 // Sep 30th 2022 00:00:00 GMT
const endTime = 1667260800 // Nov 1st 2022 00:00:00 GMT
const lastDistribution = "2022-09-30"

const tbtcv2ScriptPath = "src/tbtcv2-rewards/"
const graphqlApi =
  "https://api.studio.thegraph.com/query/24143/main-threshold-subgraph/0.0.7"

async function main() {
  let earnedBonusRewards = {}
  let earnedPreRewards = {}
  let earnedTbtcv2Rewards = {}
  let bonusRewards = {}
  let preRewards = {}
  let tbtcv2Rewards = {}
  const endDate = new Date(endTime * 1000).toISOString().slice(0, 10)
  const distPath = `distributions/${endDate}`
  const lastDistPath = `distributions/${lastDistribution}`
  const tbtcv2Script =
    "./rewards.sh " +
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

  if (preWeight > 0) {
    console.log("Calculating PRE rewards...")
    const preStakes = await Subgraph.getPreStakes(
      graphqlApi,
      startTime,
      endTime
    )
    earnedPreRewards = await Rewards.calculatePreRewards(preStakes, preWeight)
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
    bonusRewards = JSON.parse(
      fs.readFileSync(`${lastDistPath}/MerkleInputBonusRewards.json`)
    )
    bonusRewards = MerkleDist.combineMerkleInputs(
      bonusRewards,
      earnedBonusRewards
    )
    fs.writeFileSync(
      distPath + "/MerkleInputBonusRewards.json",
      JSON.stringify(bonusRewards, null, 4)
    )
    // TODO: change MerkleInputOngoingRewards.json to MerkleInputPreRewards.json
    // for December distribution (the distribution that is going to be released
    // on 2023/01/01)
    preRewards = JSON.parse(
      fs.readFileSync(`${lastDistPath}/MerkleInputOngoingRewards.json`)
    )
    preRewards = MerkleDist.combineMerkleInputs(preRewards, earnedPreRewards)
    fs.writeFileSync(
      distPath + "/MerkleInputPreRewards.json",
      JSON.stringify(preRewards, null, 4)
    )
    if (fs.existsSync(`${lastDistPath}/MerkleInputTbtcv2Rewards.json`)) {
      tbtcv2Rewards = JSON.parse(
        fs.readFileSync(`${lastDistPath}/MerkleInputTbtcv2Rewards.json`)
      )
    } else {
      tbtcv2Rewards = {}
    }
    tbtcv2Rewards = MerkleDist.combineMerkleInputs(
      tbtcv2Rewards,
      earnedTbtcv2Rewards
    )
    fs.writeFileSync(
      distPath + "/MerkleInputTbtcv2Rewards.json",
      JSON.stringify(earnedTbtcv2Rewards, null, 4)
    )
  } catch (err) {
    console.error(err)
    return
  }

  let merkleInput = MerkleDist.combineMerkleInputs(bonusRewards, preRewards)
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
