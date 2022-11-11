const BigNumber = require("bignumber.js")

const SECONDS_IN_YEAR = 31536000

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

/**
 * Calculate the PRE rewards earned by each stake
 * reward = (s_1 * y_t) * t / 365;
 * where y_t (target APY) is 0.15 and s_1 is the T amount staked
 * @param {Object} stakes     Stakes with staked T amount
 * @param {Number} weight     The weight of this type of reward
 * @return {Object}           The stakes including reward amount
 */
exports.calculatePreRewards = function (stakes, weight) {
  const preRewards = {}
  Object.keys(stakes).map((stakingProvider) => {
    const epochStakes = stakes[stakingProvider].epochStakes
    const reward = epochStakes.reduce((total, epochStake) => {
      const epochStakeAmount = BigNumber(epochStake.amount)
      const epochDuration = epochStake.epochDuration
      const epochReward = epochStakeAmount
        .times(15)
        .times(epochDuration)
        .div(SECONDS_IN_YEAR * 100)
      return total.plus(epochReward)
    }, BigNumber(0))

    if (!reward.isZero()) {
      const weightedReward = reward.times(weight)
      preRewards[stakingProvider] = {
        beneficiary: stakes[stakingProvider].beneficiary,
        amount: weightedReward.toFixed(0),
      }
    }
  })
  return preRewards
}

/**
 * Calculate the tBTCv2 weighted rewards earned by each stake
 * @param {Object} stakes     Stakes with staked T amount
 * @param {Number} weight     The weight of this type of reward
 * @return {Object}           The stakes including reward amount
 */
exports.calculateTbtcv2Rewards = function (stakes, weight) {
  Object.keys(stakes).map((stakingProvider) => {
    const amount = BigNumber(stakes[stakingProvider].amount)
    const weightedReward = amount.times(weight)
    stakes[stakingProvider].amount = weightedReward.toFixed(0)
  })
  return stakes
}
