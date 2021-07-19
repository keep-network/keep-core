import BigNumber from "bignumber.js"

const WEEKS_IN_YEAR = 52

class APYCalculator {
  /**
   * Calculates pool reward rate per number of compouning periods.
   * @param {number | string | BigNumber} rewardTokenPriceInUSD The price of the reward token in
   * USD.
   * @param {number | string | BigNumber} rewardPoolInRewardTokenPerPeriod The
   * reward pool in the reward token unit per a given period.
   * @param {number | string | BigNumber} totalStakedInUSD Total value locked in usd.
   * @return {BigNumber} The pool reward rate per number of compouning periods.
   */
  static calculatePoolRewardRate = (
    rewardTokenPriceInUSD,
    rewardPoolInRewardTokenPerPeriod,
    totalStakedInUSD
  ) => {
    return new BigNumber(rewardTokenPriceInUSD)
      .multipliedBy(rewardPoolInRewardTokenPerPeriod)
      .div(totalStakedInUSD)
  }

  /**
   * Calculates the APY- formula: (1 + r/n)^n -1
   * @param {number | string | BigNumber} poolRewardRate The pool reward rate - it equals
   * to r/n. It can be calculated by {@link APYCalculator.calculatePoolRewardRate}.
   * @param {number | string | BigNumber} n Number of compounding periods.
   * @return {BigNumber} The APY value.
   */
  static calculateAPY = (poolRewardRate, n = WEEKS_IN_YEAR) => {
    return new BigNumber(poolRewardRate).plus(1).pow(n).minus(1)
  }
}

export default APYCalculator
