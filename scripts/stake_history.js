// Script that retrieves the information of a particular staker, including the staking history
//Use: node scripts/stake_history.js <staking provider address>

const stakingRewards = require("../src/stakingrewards/stakingrewards.js")

const graphqlApi =
  "https://api.studio.thegraph.com/query/24143/main-threshold-subgraph/0.0.7"

async function main() {
  const args = process.argv.slice(2)
  const stakingProvAddress = args[0]

  const stakeHistory = await stakingRewards.getStakingHistory(
    graphqlApi,
    stakingProvAddress
  )

  console.log(JSON.stringify(stakeHistory, null, 2))
}

;(async () => {
  await main()
})()
