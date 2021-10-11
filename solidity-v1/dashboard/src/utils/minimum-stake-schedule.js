import moment from "moment"

const dateFormat = "MM/DD/YYYY"

// Minimum stake diminishing schedule can be found here:
// https://staking.keep.network/about-staking/staking-minimums
const minimumStakeDiminishingSchedule = [
  {
    date: moment.utc("04/28/2020", dateFormat),
    value: 100000,
  },
  {
    date: moment.utc("04/28/2020", dateFormat),
    value: 90000,
  },
  {
    date: moment.utc("09/21/2020", dateFormat),
    value: 80000,
  },
  {
    date: moment.utc("03/12/2020", dateFormat),
    value: 70000,
  },
  {
    date: moment.utc("02/14/2021", dateFormat),
    value: 60000,
  },
  {
    date: moment.utc("04/28/2021", dateFormat),
    value: 50000,
  },
  {
    date: moment.utc("07/10/2021", dateFormat),
    value: 40000,
  },
  {
    date: moment.utc("21/09/2021", dateFormat),
    value: 30000,
  },
  {
    date: moment.utc("12/03/2021", dateFormat),
    value: 20000,
  },
  {
    date: moment.utc("02/14/2022", dateFormat),
    value: 10000,
  },
]

export const getNextMinStake = () => {
  const currentDate = moment().utc()

  for (const minimumStakeInfo of minimumStakeDiminishingSchedule) {
    if (!currentDate.isSameOrAfter(minimumStakeInfo.date)) {
      return {
        date: minimumStakeInfo.date.format(dateFormat),
        value: minimumStakeInfo.value,
      }
    }
  }

  return minimumStakeDiminishingSchedule[
    minimumStakeDiminishingSchedule.length - 1
  ]
}
