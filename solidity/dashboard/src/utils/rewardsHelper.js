import moment from "moment"

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
}

class BeaconRewards extends RewardsHelper {
  keepAllocationsInInterval = [
    /* eslint-disable*/
          792000,     1520640,    1748736,    1888635,    2077498,    1765874,
          1500993,    1275844,    1084467,    921797,     783528,     665998,
          566099,     481184,     409006,     347655,     295507,     251181,
          213504,     181478,     154257,     131118,     111450,     94733
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
  intervals = 24
  constructor() {
    // BondedECDSAKeepFactory deployment date, Sep-14-2020 interval started.
    // https://etherscan.io/address/0xA7d9E842EFB252389d613dA88EDa3731512e40bD
    const ecdsaFirstIntervalStart = 1600041600

    // Each interval is 30 days long.
    const termLength = moment.duration(30, "days").asSeconds()

    const minimumECDSAKeepsPerInterval = 1000

    super(ecdsaFirstIntervalStart, termLength, minimumECDSAKeepsPerInterval)
  }
}

const BeaconRewardsHelper = new BeaconRewards()
const ECDSARewardsHelper = new ECDSARewards()

export default BeaconRewardsHelper

export { ECDSARewardsHelper }
