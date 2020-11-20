import moment from "moment"

class BeaconRewardsHelper {
  static keepAllocationsInInterval = [
    /* eslint-disable*/
          792000,     1520640,    1748736,    1888635,    2077498,    1765874,
          1500993,    1275844,    1084467,    921797,     783528,     665998,
          566099,     481184,     409006,     347655,     295507,     251181,
          213504,     181478,     154257,     131118,     111450,     94733
        /* eslint-enable*/
  ]

  // Beacon genesis date, 2020-09-24, is the first interval start.
  // https://etherscan.io/tx/0xe2e8ab5631473a3d7d8122ce4853c38f5cc7d3dcbfab3607f6b27a7ef3b86da2
  static beaconFirstIntervalStart = 1600905600

  // Each interval is 30 days long.
  static beaconTermLength = moment.duration(30, "days").asSeconds()

  static currentInterval = Math.floor(
    (moment().unix() - this.beaconFirstIntervalStart) / this.beaconTermLength
  )

  // There has to be at least 2 groups per interval to meet the group quota
  // and distribute the full reward for the given interval.
  static minimumBeaconGroupsPerInterval = 2

  static intervalStartOf = (interval) => {
    return moment
      .unix(this.beaconFirstIntervalStart)
      .add(interval * this.beaconTermLength, "seconds")
  }
}

export default BeaconRewardsHelper
