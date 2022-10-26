const BigNumber = require("bignumber.js")

/**
 * Calculate the bonus rewards earned by each stake
 * reward = 0.03 * initial_amount
 * @param {Object} stakes     Stakes with staked T amount
 * @param {Number} weight     The weight of this type of reward
 * @return {Object}           The stakes including reward amount
 */
exports.calculateBonusRewards = function (stakes, weight) {
  Object.keys(stakes).map((stakingProvider) => {
    const amount = BigNumber(stakes[stakingProvider].amount)
    const reward = amount.times(0.03)
    const weightedReward = reward.times(weight)
    stakes[stakingProvider].amount = weightedReward.toFixed(0)
  })
  return stakes
}
