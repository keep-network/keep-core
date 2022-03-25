import moment from "moment"
import BigNumber from "bignumber.js"
import { formatDate } from "../utils/general.utils"

class RewardsInterface {
  intervalStartOf = (interval) => {
    throw Error("Implement first")
  }

  intervalEndOf = (interval) => {
    throw Error("Implement first")
  }

  get currentInterval() {
    throw Error("Implement first")
  }

  periodOf = (interval) => {
    const startDate = formatDate(moment.unix(this.intervalStartOf(interval)))
    const endDate = formatDate(moment.unix(this.intervalEndOf(interval)))

    return `${startDate} - ${endDate}`
  }
}

class RewardsHelper extends RewardsInterface {
  firstIntervalStart = 0
  termLength = 0
  minimumKeepsPerInterval = 0

  constructor(_firstIntervalStart, _termLength, _minimumKeepsPerInterval) {
    super()
    this.firstIntervalStart = _firstIntervalStart
    this.termLength = _termLength
    this.minimumKeepsPerInterval = _minimumKeepsPerInterval
  }

  intervalStartOf = (interval) => {
    return moment
      .unix(this.firstIntervalStart)
      .add(interval * this.termLength, "seconds")
      .unix()
  }

  intervalEndOf = (interval) => {
    return this.intervalStartOf(interval + 1)
  }

  get currentInterval() {
    return Math.floor(
      (moment().unix() - this.firstIntervalStart) / this.termLength
    )
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

class ExtendedECDSARewards extends RewardsInterface {
  // https://forum.keep.network/t/proposal-to-extend-stakedrop-rewards-for-an-additional-4-months/351
  intervalsDates = [
    // 2021-12-03 -> 2021-12-21
    {
      start: 1638489600,
      end: 1640044800,
      merkleRoot:
        "0x83f51475c210aff536867981fcc803ace41787a6b0b256fced9802ec37127dd7",
    },
    // 2021-12-21 -> 2022-01-22
    {
      start: 1640044800,
      end: 1642809600,
      merkleRoot:
        "0xdc8026bc52c1d200477e8aa8d374e934e57c79d1d0c9fa65d121a8f6607987b0",
    },
    // 2021-01-22 -> 2022-02-22
    {
      start: 1642809600,
      end: 1645488000,
      merkleRoot:
        "0xa48918393536de2ba2dfb10b66ea91abb7f66352156d34f2f3ffcbb0b976ba2b",
    },
    // 2021-02-22 -> 2022-03-22
    { start: 1645488000, end: 1647907200, merkleRoot: "" },
  ]

  get currentInterval() {
    const currentTimestamp = moment().unix()
    const index = this.intervalsDates.findIndex(
      (_) => _.start <= currentTimestamp && _.end >= currentTimestamp
    )
    if (index < 0) return this.intervals

    return index + 1
  }

  get intervals() {
    return this.intervalsDates.length
  }

  intervalStartOf = (interval) => {
    const _interval = interval - 1

    if (_interval < 0) return this.intervalsDates[0].start
    else if (_interval > this.intervalsDates.length)
      return this.intervalsDates[this.intervalsDates.length - 1].start

    return this.intervalsDates[_interval].start
  }

  intervalEndOf = (interval) => {
    const _interval = interval - 1

    if (_interval < 0) return this.intervalsDates[0].end
    if (_interval > this.intervalsDates.length)
      return this.intervalsDates[this.intervalsDates.length - 1].end

    return this.intervalsDates[_interval].end
  }
}

const BeaconRewardsHelper = new BeaconRewards()
const ECDSARewardsHelper = new ECDSARewards()
const ExtendedECDSARewardsHelper = new ExtendedECDSARewards()

export class ECDSAPeriodOfResolver {
  static resolve(interval, merkleRoot) {
    const index = ExtendedECDSARewardsHelper.intervalsDates.findIndex(
      (_) => _.merkleRoot === merkleRoot
    )
    if (index >= 0) {
      return ExtendedECDSARewardsHelper.periodOf(index + 1)
    }

    return ECDSARewardsHelper.periodOf(interval)
  }
}

export default BeaconRewardsHelper

export { ECDSARewardsHelper, ExtendedECDSARewardsHelper }
