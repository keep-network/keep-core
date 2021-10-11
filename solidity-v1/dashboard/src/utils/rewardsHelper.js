import moment from "moment"
import BigNumber from "bignumber.js"
import { formatDate } from "../utils/general.utils"

class RewardsHelper {
  firstIntervalStart = 0
  termLength = 0
  minimumKeepsPerInterval = 0

  constructor(_firstIntervalStart, _termLength, _minimumKeepsPerInterval) {
    this.firstIntervalStart = _firstIntervalStart
    this.termLength = _termLength
    this.minimumKeepsPerInterval = _minimumKeepsPerInterval
    this.currentInterval = Math.floor(
      (moment().unix() - this.firstIntervalStart) / this.termLength
    )
  }

  intervalStartOf = (interval) => {
    return moment
      .unix(this.firstIntervalStart)
      .add(interval * this.termLength, "seconds")
  }

  intervalEndOf = (interval) => {
    return this.intervalStartOf(interval + 1)
  }

  periodOf = (interval) => {
    const startDate = formatDate(this.intervalStartOf(interval))
    const endDate = formatDate(this.intervalEndOf(interval))

    return `${startDate} - ${endDate}`
  }
}

class BeaconRewards extends RewardsHelper {
  keepAllocationsInInterval = [
    /* eslint-disable*/
    792000, 1520640, 1748736, 1888635, 2077498, 1765874, 1500993, 1275844,
    1084467, 921797, 783528, 665998, 566099, 481184, 409006, 347655, 295507,
    251181, 213504, 181478, 154257, 131118, 111450, 94733,
    /* eslint-enable*/
  ]

  constructor() {
    // Beacon genesis date, 2020-09-24, is the first interval start.
    // https://etherscan.io/tx/0xe2e8ab5631473a3d7d8122ce4853c38f5cc7d3dcbfab3607f6b27a7ef3b86da2
    const beaconFirstIntervalStart = 1600905600

    // Each interval is 30 days long.
    const beaconTermLength = moment.duration(30, "days").asSeconds()

    // There has to be at least 2 groups per interval to meet the group quota
    // and distribute the full reward for the given interval.
    const minimumBeaconGroupsPerInterval = 2

    super(
      beaconFirstIntervalStart,
      beaconTermLength,
      minimumBeaconGroupsPerInterval
    )
  }
}

class ECDSARewards extends RewardsHelper {
  // The new rewards mechanism has 88 intervals started from 13.11.2020. Each
  // interval is 7 days long. There might be more intervals than 88 because
  // unallocated rewards (eg. when some operrator does not meet SLA) remain in
  // the pool and can be used for other intervals.
  intervals = 88
  ethScoreThreshold = new BigNumber(3000).multipliedBy(new BigNumber(1e18)) // 3000 ETH
  constructor() {
    // 13.11.2020 interval started
    const ecdsaFirstIntervalStart = 1605225600

    // Each interval is 7 days long.
    const termLength = moment.duration(7, "days").asSeconds()

    const minimumECDSAKeepsPerInterval = 1000

    super(ecdsaFirstIntervalStart, termLength, minimumECDSAKeepsPerInterval)
  }

  calculateETHScore(ethTotal) {
    if (ethTotal.isLessThan(this.ethScoreThreshold)) {
      return ethTotal
    }

    const sqrt = this.ethScoreThreshold.multipliedBy(ethTotal).squareRoot()

    return new BigNumber(2).multipliedBy(sqrt).minus(this.ethScoreThreshold)
  }

  calculateBoost(keepStaked, ethTotal, minimumStake) {
    const a = keepStaked.dividedBy(minimumStake)
    const b = ethTotal.isGreaterThan(new BigNumber(0))
      ? keepStaked
          .dividedBy(ethTotal.multipliedBy(new BigNumber(500)))
          .squareRoot()
      : new BigNumber(0)

    return new BigNumber(1).plus(BigNumber.minimum(a, b))
  }
}

const BeaconRewardsHelper = new BeaconRewards()
const ECDSARewardsHelper = new ECDSARewards()

export default BeaconRewardsHelper

export { ECDSARewardsHelper }
