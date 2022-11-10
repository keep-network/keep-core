// Script that retrieves the information of a particular staker, including the staking history
// Use: node scripts/stake_history.js <staking provider address>

const Subgraph = require("../pre-rewards/subgraph.js")
const { ethers } = require("ethers")

const graphqlApi =
  "https://api.studio.thegraph.com/query/24143/main-threshold-subgraph/0.0.7"

async function main() {
  const args = process.argv.slice(2)
  const stakingProvAddress = args[0]

  if(!ethers.utils.isAddress(stakingProvAddress)) {
    console.error("Invalid address")
    return
  }

  const stakeHistory = await Subgraph.getStakingHistory(
    graphqlApi,
    stakingProvAddress
  )

  console.log(JSON.stringify(stakeHistory, null, 2))
}

;(async () => {
  await main()
})()
