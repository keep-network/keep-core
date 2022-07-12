import "isomorphic-unfetch"
import { createClient, gql } from "@urql/core"
import { ethers } from "ethers"
import BigNumber from "bignumber.js"

const SECONDS_IN_YEAR = 31536000

async function getEpochById(gqlClient, epochId) {
  let data

  const FIRST_EPOCH_QUERY = gql`
    query FirstEpoch($id: String) {
      epoch(id: $id) {
        id
        timestamp
        duration
        totalAmount
        stakes(first: 1000) {
          stakingProvider
          owner
          amount
        }
      }
    }
  `

  await gqlClient
    .query(FIRST_EPOCH_QUERY, { id: epochId.toString() })
    .toPromise()
    .then((result) => {
      data = result.data.epoch
    })

  return data
}

async function getEpochByIdAndOwner(gqlClient, epochId, address) {
  let data

  // TODO: Max amount of items you can get in a query is 100.
  // adding 'first: 1000' is a WA to get more than 100 stakes,
  // but the most correct option is to use GraphQL pagination.
  const FIRST_EPOCH_QUERY = gql`
    query FirstEpoch($id: String, $address: String) {
      epoch(id: $id) {
        id
        timestamp
        duration
        totalAmount
        stakes(first: 1000, where: { owner: $address }) {
          amount
          stakingProvider
        }
      }
    }
  `

  await gqlClient
    .query(FIRST_EPOCH_QUERY, { id: epochId.toString(), address: address })
    .toPromise()
    .then((result) => {
      data = result.data.epoch
    })

  return data
}

async function getEpochsByStartTime(gqlClient, timestamp) {
  let data

  // TODO: Max amount of items you can get in a query is 100.
  // adding 'first: 1000' is a WA to get more than 100 stakes,
  // but the most correct option is to use GraphQL pagination.
  const ONGOING_STAKES_QUERY = gql`
    query OngoingStakesJuneFirst($timestamp: String) {
      epoches(
        first: 1000
        orderBy: timestamp
        where: { timestamp_gte: $timestamp }
      ) {
        id
        timestamp
        duration
        totalAmount
      }
    }
  `

  await gqlClient
    .query(ONGOING_STAKES_QUERY, { timestamp: timestamp.toString() })
    .toPromise()
    .then((result) => {
      data = result.data.epoches
    })

  return data
}

async function getEpochsBetweenDates(gqlClient, startTimestamp, endTimestamp) {
  let data

  // TODO: Max amount of items you can get in a query is 100.
  // adding 'first: 1000' is a WA to get more than 100 stakes,
  // but the most correct option is to use GraphQL pagination.
  const ONGOING_STAKES_QUERY = gql`
    query OngoingStakes($startTimestamp: String, $endTimestamp: String) {
      epoches(
        first: 1000
        orderBy: timestamp
        where: { timestamp_gte: $startTimestamp, timestamp_lte: $endTimestamp }
      ) {
        id
        timestamp
        duration
        totalAmount
        stakes(first: 1000) {
          stakingProvider
          owner
          amount
        }
      }
    }
  `

  await gqlClient
    .query(ONGOING_STAKES_QUERY, {
      startTimestamp: startTimestamp.toString(),
      endTimestamp: endTimestamp.toString(),
    })
    .toPromise()
    .then((result) => {
      data = result.data.epoches
    })

  return data
}

async function getEpochsBetweenDatesByOwner(
  gqlClient,
  address,
  startTimestamp,
  endTimestamp
) {
  let data

  // TODO: Max amount of items you can get in a query is 100.
  // adding 'first: 1000' is a WA to get more than 100 stakes,
  // but the most correct option is to use GraphQL pagination.
  const ONGOING_STAKES_QUERY = gql`
    query OngoingStakes(
      $address: String
      $startTimestamp: String
      $endTimestamp: String
    ) {
      epoches(
        first: 1000
        orderBy: timestamp
        where: { timestamp_gte: $startTimestamp, timestamp_lte: $endTimestamp }
      ) {
        id
        timestamp
        duration
        totalAmount
        stakes(first: 1000, where: { owner: $address }) {
          amount
          stakingProvider
        }
      }
    }
  `

  await gqlClient
    .query(ONGOING_STAKES_QUERY, {
      address: address,
      startTimestamp: startTimestamp.toString(),
      endTimestamp: endTimestamp.toString(),
    })
    .toPromise()
    .then((result) => {
      data = result.data.epoches
    })

  return data
}

async function getOpConfTimestamp(gqlClient, stakingProviderAddress) {
  let data

  const OP_CONF_TIMESTAMP = gql`
    query SimplePREApplication($address: String) {
      simplePREApplication(id: $address) {
        confirmedTimestamp
      }
    }
  `

  await gqlClient
    .query(OP_CONF_TIMESTAMP, { address: stakingProviderAddress })
    .toPromise()
    .then((result) => {
      data = result?.data?.simplePREApplication?.confirmedTimestamp
    })

  return data ? parseInt(data) : null
}

async function getOperatorsConfirmedBeforeDate(gqlClient, timestamp) {
  let data

  // TODO: Max amount of items you can get in a query is 100.
  // adding 'first: 1000' is a WA to get more than 100 stakes,
  // but the most correct option is to use GraphQL pagination.
  const OPS_CONF_BETWEEN_DATES = gql`
    query SimplePREApplications($timestamp: String) {
      simplePREApplications(
        first: 1000
        where: { confirmedTimestamp_lte: $timestamp }
      ) {
        bondedTimestamp
        confirmedTimestamp
        id
        operator
      }
    }
  `

  await gqlClient
    .query(OPS_CONF_BETWEEN_DATES, {
      timestamp: timestamp.toString(),
    })
    .toPromise()
    .then((result) => {
      data = result.data?.simplePREApplications
    })

  return data
}

async function getStakeDatasInfo(gqlClient) {
  let data

  // TODO: Max amount of items you can get in a query is 100.
  // adding 'first: 1000' is a WA to get more than 100 stakes,
  // but the most correct option is to use GraphQL pagination.
  const STAKES_DATA_INFO = gql`
    query stakeDatasInfo {
      stakeDatas(first: 1000) {
        beneficiary
        id
        authorizer
        owner {
          id
        }
      }
    }
  `

  await gqlClient
    .query(STAKES_DATA_INFO)
    .toPromise()
    .then((result) => {
      data = result.data?.stakeDatas
    })

  return data
}

/**
 * Get the total ongoing rewards generated by stakes since a specific date
 * @param {string}  gqlURL          Subgraph GraphQL API URL
 * @param {Number} startTimestamp   Start date UNIX timestamp
 * @return {BigNumber}              The amount of generated rewards
 */
async function getTotalOngoingRewards(gqlUrl, timestamp) {
  const currentTime = parseInt(Date.now() / 1000)
  const gqlClient = createClient({ url: gqlUrl })

  let epochs = await getEpochsByStartTime(gqlClient, timestamp)
  const firstEpochId = parseInt(epochs[0].id) > 0 ? epochs[0].id - 1 : 0
  let firstEpoch = await getEpochById(gqlClient, firstEpochId)
  epochs = [firstEpoch, ...epochs]

  epochs[0].timestamp = timestamp.toString()
  epochs[0].duration = (epochs[1].timestamp - timestamp).toString()

  const reward = epochs.reduce((total, epoch) => {
    const amount = new BigNumber(epoch.totalAmount)
    const duration = epoch.duration
      ? parseInt(epoch.duration)
      : currentTime - epoch.timestamp
    const epochReward = amount
      .times(15)
      .times(duration)
      .div(SECONDS_IN_YEAR * 100)
    return total.plus(epochReward)
  }, new BigNumber(0))

  return reward
}

/**
 * Get the ongoing rewards of an stake generated between two dates
 * @param {string}  gqlURL          Subgraph GraphQL API URL
 * @param {string}  ownerAddress    Stake owner address
 * @param {Number}  startTimestamp  Start date UNIX timestamp
 * @param {Number}  endTimestamp    End date UNIX timestamp
 * @return {BigNumber}              The amount of generated rewards
 */
async function getOngoingRewards(
  gqlUrl,
  ownerAddress,
  startTimestamp,
  endTimestamp
) {
  const currentTime = parseInt(Date.now() / 1000)
  const gqlClient = createClient({ url: gqlUrl })

  let epochs = await getEpochsBetweenDatesByOwner(
    gqlClient,
    ownerAddress,
    startTimestamp,
    endTimestamp
  )

  const firstEpochId = parseInt(epochs[0].id) > 0 ? epochs[0].id - 1 : 0
  const firstEpoch = await getEpochByIdAndOwner(
    gqlClient,
    firstEpochId,
    ownerAddress
  )

  epochs = [firstEpoch, ...epochs]
  epochs[0].timestamp = startTimestamp.toString()
  epochs[0].duration = (epochs[1].timestamp - startTimestamp).toString()
  const lastEpochIndex = epochs.length - 1 > 0 ? epochs.length - 1 : 0
  epochs[lastEpochIndex].duration =
    endTimestamp - epochs[lastEpochIndex].timestamp

  // Clean the empty epochs
  epochs = epochs.filter((epoch) => {
    return epoch.stakes.length > 0
  })

  // Sort the epoch's stakes by staking provider
  const stakeList = {}
  epochs.forEach((epoch) => {
    epoch.stakes.forEach((stake) => {
      const stakeData = {}
      stakeData.epochTotalAmount = epoch.totalAmount
      stakeData.epochDuration = epoch.duration
      stakeData.epochTimestamp = epoch.timestamp
      stakeData.amount = stake.amount
      stakeList[stake.stakingProvider] = stakeList[stake.stakingProvider]
        ? stakeList[stake.stakingProvider]
        : []
      stakeList[stake.stakingProvider].push(stakeData)
    })
  })

  const rewards = []
  // Calculate the rewards of each stake
  // Rewards formula: r = (s_1 * y_t) * t / 365; where y_t is 0.15
  await Promise.all(
    Object.keys(stakeList).map(async (stakingProvider) => {
      let reward = BigNumber(0)

      const stake = stakeList[stakingProvider]
      // Check if operator is confirmed and when
      const opConfTimestamp = await getOpConfTimestamp(
        gqlClient,
        stakingProvider
      )
      if (opConfTimestamp) {
        reward = stake.reduce((total, epochStake) => {
          let epochReward = BigNumber(0)
          const stakeAmount = BigNumber(epochStake.amount)
          const epochTimestamp = parseInt(epochStake.epochTimestamp)
          let epochDuration = epochStake.epochDuration
            ? parseInt(epochStake.epochDuration)
            : currentTime - epochStake.epochTimestamp

          // If the operator was confirmed in the middle of this epoch...
          if (
            opConfTimestamp > epochTimestamp &&
            opConfTimestamp <= epochTimestamp + epochDuration
          ) {
            epochDuration = epochTimestamp + epochDuration - opConfTimestamp
          }

          epochReward = stakeAmount
            .times(15)
            .times(epochDuration)
            .div(SECONDS_IN_YEAR * 100)

          return total.plus(epochReward)
        }, BigNumber(0))
      }

      const rewardsItem = {
        stakingProvider: stakingProvider,
        reward: reward,
      }
      rewards.push(rewardsItem)
    })
  )
  return rewards
}

/**
 * Get the ongoing rewards generated by each stake since a specific date and
 * export it to a JSON file
 * @param {string}  gqlURL          Subgraph GraphQL API URL
 * @param {Number} startTimestamp   Start date UNIX timestamp
 * @param {Number}  endTimestamp    End date UNIX timestamp
 * @return {Object}                 The ongoing rewards of each stake
 */
async function getOngoingRewardsMekleInput(
  gqlUrl,
  startTimestamp,
  endTimestamp
) {
  const currentTime = parseInt(Date.now() / 1000)
  const gqlClient = createClient({ url: gqlUrl })

  // Get the list of operators confirmed between dates
  const opsConfirmed = await getOperatorsConfirmedBeforeDate(
    gqlClient,
    startTimestamp
  )

  // Get the stakes information
  const stakeDatas = await getStakeDatasInfo(gqlClient)

  let epochs = await getEpochsBetweenDates(
    gqlClient,
    startTimestamp,
    endTimestamp
  )
  const firstEpochId = parseInt(epochs[0].id) > 0 ? epochs[0].id - 1 : 0
  let firstEpoch = await getEpochById(gqlClient, firstEpochId)

  epochs = [firstEpoch, ...epochs]
  epochs[0].timestamp = startTimestamp.toString()
  epochs[0].duration = (epochs[1].timestamp - startTimestamp).toString()
  const lastEpochIndex = epochs.length - 1 > 0 ? epochs.length - 1 : 0
  epochs[lastEpochIndex].duration =
    endTimestamp - epochs[lastEpochIndex].timestamp

  // Clean the empty epochs
  epochs = epochs.filter((epoch) => {
    return epoch.stakes.length > 0
  })

  // Sort the epoch's stakes by staking provider
  const stakeList = {}
  epochs.forEach((epoch) => {
    epoch.stakes.forEach((epochStake) => {
      const stakeData = {}
      stakeData.epochTotalAmount = epoch.totalAmount
      stakeData.epochDuration = epoch.duration
      stakeData.epochTimestamp = epoch.timestamp
      stakeData.amount = epochStake.amount
      stakeData.epochId = epoch.id
      stakeList[epochStake.stakingProvider] = stakeList[
        epochStake.stakingProvider
      ]
        ? stakeList[epochStake.stakingProvider]
        : []
      stakeList[epochStake.stakingProvider].push(stakeData)
    })
  })

  // Calculate the reward of each stake
  // Rewards formula: r = (s_1 * y_t) * t / 365; where y_t is 0.15
  const rewards = []
  Object.keys(stakeList).map((stakingProvider) => {
    let reward = BigNumber(0)

    const stake = stakeList[stakingProvider]

    // Check if operator is confirmed and when
    const opConf = opsConfirmed.find((op) => op.id === stakingProvider)
    const opConfTimestamp = opConf ? opConf.confirmedTimestamp : undefined
    if (opConfTimestamp) {
      reward = stake.reduce((total, epochStake) => {
        let epochReward = BigNumber(0)
        const stakeAmount = BigNumber(epochStake.amount)
        const epochTimestamp = parseInt(epochStake.epochTimestamp)
        let epochDuration = epochStake.epochDuration
          ? parseInt(epochStake.epochDuration)
          : currentTime - epochStake.epochTimestamp

        // If the operator was confirmed in the middle of this epoch...
        if (
          opConfTimestamp > epochTimestamp &&
          opConfTimestamp <= epochTimestamp + epochDuration
        ) {
          epochDuration = epochTimestamp + epochDuration - opConfTimestamp
        }

        epochReward = stakeAmount
          .times(15)
          .times(epochDuration)
          .div(SECONDS_IN_YEAR * 100)

        return total.plus(epochReward)
      }, BigNumber(0))
    }

    if (!reward.isZero()) {
      // Find the beneficiary of this reward
      const stakeDatasItem = stakeDatas.find(
        (stake) => stake.id === stakingProvider
      )
      const beneficiary = stakeDatasItem.beneficiary

      const rewardsItem = {
        stakingProvider: ethers.utils.getAddress(stakingProvider),
        reward: reward,
        beneficiary: ethers.utils.getAddress(beneficiary),
      }
      rewards.push(rewardsItem)
    }
  })

  return rewards
}

/**
 * Get the bonus rewards generated by each stake between June 1st and July 15th
 * and export it to a JSON file
 * @param {string}  gqlURL          Subgraph GraphQL API URL
 * @return {BigNumber}              The amount of generated rewards
 */
async function getBonusRewardsMerkleInput(gqlUrl) {
  const startTimestamp = 1654041600 // Jun 1st 2022 00:00:00 GMT
  const endTimestamp = 1657843200 // Jul 15th 2022 00:00:00 GMT
  const gqlClient = createClient({ url: gqlUrl })

  // Get the list of operators confirmed between dates
  const opsConfirmed = await getOperatorsConfirmedBeforeDate(
    gqlClient,
    startTimestamp
  )

  // Get the stakes information
  const stakeDatas = await getStakeDatasInfo(gqlClient)

  let epochs = await getEpochsBetweenDates(
    gqlClient,
    startTimestamp,
    endTimestamp
  )

  const firstEpochId = parseInt(epochs[0].id) > 0 ? epochs[0].id - 1 : 0
  let firstEpoch = await getEpochById(gqlClient, firstEpochId)
  epochs = [firstEpoch, ...epochs]

  // Sort the epoch's stakes by staking provider
  const stakeList = {}
  epochs.forEach((epoch) => {
    epoch.stakes.forEach((epochStake) => {
      const stakeData = {}
      stakeData.amount = BigNumber(epochStake.amount)
      stakeData.epochId = Number(epoch.id)
      stakeData.epochTimestamp = Number(epoch.timestamp)
      stakeList[epochStake.stakingProvider] = stakeList[
        epochStake.stakingProvider
      ]
        ? stakeList[epochStake.stakingProvider]
        : []
      stakeList[epochStake.stakingProvider].push(stakeData)
    })
  })

  const rewards = []

  // Filter the stakes that are not elegible for bonus
  Object.keys(stakeList).map((stakingProvider) => {
    const epochStakes = stakeList[stakingProvider]
    epochStakes.sort((a, b) => a.epochId - b.epochId)

    // stake must have started before the start date
    let elegible = epochStakes[0].epochTimestamp <= startTimestamp
    // stake must have confirmed operator before startDate
    const opConf = opsConfirmed.find((op) => op.id === stakingProvider)
    const opConfTimestamp = opConf
      ? Number(opConf.confirmedTimestamp)
      : undefined
    elegible = elegible && opConfTimestamp <= startTimestamp
    // stake must not have unstaked tokens
    epochStakes.reduce((acc, cur) => {
      if (elegible && cur.amount.gte(acc)) {
        return cur.amount
      } else {
        elegible = false
      }
    }, BigNumber(0))

    if (elegible) {
      // Find the beneficiary of this reward
      const stakeDatasItem = stakeDatas.find(
        (stake) => stake.id === stakingProvider
      )
      const beneficiary = stakeDatasItem.beneficiary
      // Calculate the earning of this stake
      // r = 0.03 * initial_amount
      const reward = epochStakes[0].amount.times(0.03)
      const rewardsItem = {
        stakingProvider: ethers.utils.getAddress(stakingProvider),
        reward: reward,
        beneficiary: ethers.utils.getAddress(beneficiary),
      }
      rewards.push(rewardsItem)
    }
  })

  return rewards
}

export {
  getTotalOngoingRewards,
  getOngoingRewards,
  getOngoingRewardsMekleInput,
  getBonusRewardsMerkleInput,
}
