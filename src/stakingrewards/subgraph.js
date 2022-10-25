require("isomorphic-unfetch")
const { createClient, gql } = require("@urql/core")
const { ethers } = require("ethers")
const { MerkleTree } = require("merkletreejs")
const keccak256 = require("keccak256")
const BigNumber = require("bignumber.js")

// The Graph limits GraphQL queries to 1000 results max
const RESULTS_PER_QUERY = 1000
const SECONDS_IN_YEAR = 31536000

async function getEpochById(gqlClient, epochId) {
  let epoch
  let lastId = ""
  let epochStakes = []
  let data = []

  const FIRST_EPOCH_QUERY = gql`
    query FirstEpoch($id: String) {
      epoch(id: $id) {
        id
        timestamp
        duration
        totalAmount
      }
    }
  `

  const EPOCH_STAKES_QUERY = gql`
    query EpochStakes(
      $epochIds: [String!]
      $resultsPerQuery: Int
      $lastId: String
    ) {
      epochStakes(
        first: $resultsPerQuery
        where: { epoch_in: $epochIds, id_gt: $lastId }
      ) {
        id
        owner
        stakingProvider
        amount
        epoch {
          id
        }
      }
    }
  `

  await gqlClient
    .query(FIRST_EPOCH_QUERY, { id: epochId.toString() })
    .toPromise()
    .then((result) => {
      if (result.error) console.error(result.error)
      epoch = result.data.epoch
    })

  const epochIds = [epoch.id]

  do {
    await gqlClient
      .query(EPOCH_STAKES_QUERY, {
        epochIds: epochIds,
        resultsPerQuery: RESULTS_PER_QUERY,
        lastId: lastId,
      })
      .toPromise()
      .then((result) => {
        if (result.error) console.error(result.error)
        data = result.data?.epochStakes
        if (data.length > 0) {
          epochStakes = epochStakes.concat(data)
          lastId = data[data.length - 1].id
        }
      })
  } while (data.length > 0)

  epoch.stakes = epochStakes

  return epoch
}

async function getEpochsBetweenDates(gqlClient, startTimestamp, endTimestamp) {
  let lastTimestamp = startTimestamp - 1
  let lastId = ""
  let epochs = []
  let epochStakes = []
  let data = []

  const EPOCHS_QUERY = gql`
    query Epochs(
      $lastTimestamp: String
      $endTimestamp: String
      $resultsPerQuery: Int
    ) {
      epoches(
        first: $resultsPerQuery
        orderBy: timestamp
        where: { timestamp_gt: $lastTimestamp, timestamp_lte: $endTimestamp }
      ) {
        id
        timestamp
        duration
        totalAmount
      }
    }
  `

  const EPOCH_STAKES_QUERY = gql`
    query EpochStakes(
      $epochIds: [String!]
      $resultsPerQuery: Int
      $lastId: String
    ) {
      epochStakes(
        first: $resultsPerQuery
        where: { epoch_in: $epochIds, id_gt: $lastId }
      ) {
        id
        owner
        stakingProvider
        amount
        epoch {
          id
        }
      }
    }
  `

  do {
    await gqlClient
      .query(EPOCHS_QUERY, {
        lastTimestamp: lastTimestamp.toString(),
        endTimestamp: endTimestamp.toString(),
        resultsPerQuery: RESULTS_PER_QUERY,
      })
      .toPromise()
      .then((result) => {
        if (result.error) console.error(result.error)
        data = result.data?.epoches
        if (data.length > 0) {
          epochs = epochs.concat(data)
          lastTimestamp = data[data.length - 1].timestamp
        }
      })
  } while (data.length > 0)

  const epochIds = epochs.map((epoch) => epoch.id)

  do {
    await gqlClient
      .query(EPOCH_STAKES_QUERY, {
        epochIds: epochIds,
        resultsPerQuery: RESULTS_PER_QUERY,
        lastId: lastId,
      })
      .toPromise()
      .then((result) => {
        if (result.error) console.error(result.error)
        data = result.data?.epochStakes
        if (data.length > 0) {
          epochStakes = epochStakes.concat(data)
          lastId = data[data.length - 1].id
        }
      })
  } while (data.length > 0)

  epochs.forEach((epoch) => (epoch.stakes = []))
  epochStakes.forEach((epochStake) => {
    const i = epochIds.findIndex((epochId) => epochId === epochStake.epoch.id)
    epochs[i].stakes.push(epochStake)
  })

  return epochs
}

async function getOperatorsConfirmedBeforeDate(gqlClient, timestamp) {
  let lastId = ""
  let operators = []
  let data = []

  const OPS_CONF_BETWEEN_DATES = gql`
    query SimplePREApplications(
      $timestamp: String
      $resultsPerQuery: Int
      $lastId: String
    ) {
      simplePREApplications(
        first: $resultsPerQuery
        where: { confirmedTimestamp_lte: $timestamp, id_gt: $lastId }
      ) {
        bondedTimestamp
        confirmedTimestamp
        id
        operator
      }
    }
  `

  do {
    await gqlClient
      .query(OPS_CONF_BETWEEN_DATES, {
        timestamp: timestamp.toString(),
        resultsPerQuery: RESULTS_PER_QUERY,
        lastId: lastId,
      })
      .toPromise()
      .then((result) => {
        if (result.error) console.error(result.error)
        data = result.data?.simplePREApplications
        if (data.length > 0) {
          operators = operators.concat(data)
          lastId = data[data.length - 1].id
        }
      })
  } while (data.length > 0)

  return operators
}

async function getStakeDatasInfo(gqlClient) {
  let lastId = ""
  let stakeDatas = []
  let data = []

  const STAKES_DATA_INFO = gql`
    query stakeDatasInfo($resultsPerQuery: Int, $lastId: String) {
      stakeDatas(first: $resultsPerQuery, where: { id_gt: $lastId }) {
        beneficiary
        id
        authorizer
        owner {
          id
        }
      }
    }
  `

  do {
    await gqlClient
      .query(STAKES_DATA_INFO, {
        resultsPerQuery: RESULTS_PER_QUERY,
        lastId: lastId,
      })
      .toPromise()
      .then((result) => {
        if (result.error) console.error(result.error)
        data = result.data?.stakeDatas
        if (data.length > 0) {
          stakeDatas = stakeDatas.concat(data)
          lastId = data[data.length - 1].id
        }
      })
  } while (data.length > 0)

  return stakeDatas
}

/**
 * Combine two Merkle distribution inputs, adding the amounts and taking the
 * beneficiary of the second input
 * @param {Object} baseMerkleInput  Merkle input used as base
 * @param {Object} addedMerkleInput Merkle input to be added
 * @return {Object}                 Combination of two Merkle inputs
 */
exports.combineMerkleInputs = function (baseMerkleInput, addedMerkleInput) {
  const combined = JSON.parse(JSON.stringify(baseMerkleInput))
  Object.keys(addedMerkleInput).map((stakingProvider) => {
    const combinedClaim = combined[stakingProvider]
    const addedClaim = addedMerkleInput[stakingProvider]
    if (combinedClaim) {
      combinedClaim.beneficiary = addedClaim.beneficiary
      combinedClaim.amount = BigNumber(combinedClaim.amount)
        .plus(BigNumber(addedClaim.amount))
        .toFixed()
    } else {
      combined[stakingProvider] = {
        beneficiary: addedClaim.beneficiary,
        amount: addedClaim.amount,
      }
    }
  })
  return combined
}

/**
 * Generate a Merkle distribution from Merkle distribution input
 * @param {Object} merkleInput      Merkle input generated from rewards
 * @return {Object}                 Merkle distribution
 */
exports.genMerkleDist = function (merkleInput) {
  const stakingProviders = Object.keys(merkleInput)
  const data = Object.values(merkleInput)

  const elements = stakingProviders.map(
    (stakingProvider, i) =>
      stakingProvider +
      data[i].beneficiary.substr(2) +
      BigNumber(data[i].amount).toString(16).padStart(64, "0")
  )

  const tree = new MerkleTree(elements, keccak256, {
    hashLeaves: true,
    sort: true,
  })

  const root = tree.getHexRoot()
  const leaves = tree.getHexLeaves()
  const proofs = leaves.map(tree.getHexProof, tree)

  const totalAmount = data
    .map((claim) => BigNumber(claim.amount))
    .reduce((a, b) => a.plus(b))
    .toFixed()

  const claims = Object.entries(merkleInput).map(([stakingProvider, data]) => {
    const leaf = MerkleTree.bufferToHex(
      keccak256(
        stakingProvider +
          data.beneficiary.substr(2) +
          BigNumber(data.amount).toString(16).padStart(64, "0")
      )
    )
    return {
      stakingProvider: stakingProvider,
      beneficiary: data.beneficiary,
      amount: data.amount,
      proof: proofs[leaves.indexOf(leaf)],
    }
  })

  const dist = {
    totalAmount: totalAmount,
    merkleRoot: root,
    claims: claims.reduce(
      (a, { stakingProvider, beneficiary, amount, proof }) => {
        a[stakingProvider] = { beneficiary, amount, proof }
        return a
      },
      {}
    ),
  }

  return dist
}

/**
 * Generate the ongoing rewards earned by stakes since a specific date and
 * return it in Merkle distribution input format
 * @param {string}  gqlURL          Subgraph GraphQL API URL
 * @param {Number} startTimestamp   Start date UNIX timestamp
 * @param {Number}  endTimestamp    End date UNIX timestamp
 * @return {Object}                 The ongoing rewards of each stake
 */
exports.getOngoingMekleInput = async function (
  gqlUrl,
  startTimestamp,
  endTimestamp
) {
  const currentTime = parseInt(Date.now() / 1000)
  const gqlClient = createClient({ url: gqlUrl })

  // Get the list of operators confirmed between dates
  const opsConfirmed = await getOperatorsConfirmedBeforeDate(
    gqlClient,
    endTimestamp
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
  const rewards = {}
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

        if (
          // If the operator was confirmed in the middle of this epoch...
          opConfTimestamp > epochTimestamp &&
          opConfTimestamp <= epochTimestamp + epochDuration
        ) {
          epochDuration = epochTimestamp + epochDuration - opConfTimestamp
        } else if (opConfTimestamp >= epochTimestamp + epochDuration) {
          // No rewards if the operator was not yet confirmed
          epochDuration = 0
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

      const stProvCheckSum = ethers.utils.getAddress(stakingProvider)
      rewards[stProvCheckSum] = {
        beneficiary: ethers.utils.getAddress(beneficiary),
        amount: reward.toFixed(0),
      }
    }
  })

  return rewards
}

/**
 * Generate the bonus rewards earned by stakes between June 1st and July 15th
 * and return it in Merkle distribution input format
 * @param {string}  gqlURL          Subgraph GraphQL API URL
 * @return {BigNumber}              The amount of generated rewards
 */
exports.getBonusMerkleInput = async function (gqlUrl) {
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

  const rewards = {}

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

      // Calculate the earning of this stake r = 0.03 * initial_amount
      const reward = epochStakes[0].amount.times(0.03)

      const stProvCheckSum = ethers.utils.getAddress(stakingProvider)
      rewards[stProvCheckSum] = {
        beneficiary: ethers.utils.getAddress(beneficiary),
        amount: reward.toFixed(0),
      }
    }
  })

  return rewards
}

/**
 * Retrieve the information of a particular staker, including the staking history.
 * @param {string}  gqlURL            Subgraph's GraphQL API URL
 * @param {string}  stakingProvider   Staking provider address
 * @return {Object}                   The stake's data
 */
exports.getStakingHistory = async function (gqlUrl, stakingProvider) {
  let lastId = ""
  let data = []
  let epochStakes = []
  let amount = 0
  let stakeData = {
    data: {},
    stake: {},
    operator: {},
    stakingHistory: [],
  }

  const gqlClient = createClient({ url: gqlUrl })

  const STAKE_DATA_QUERY = gql`
    query StakeData($stakingProvider: String) {
      stakeData(id: $stakingProvider) {
        id
        totalStaked
        authorizer
        beneficiary
        keepInTStake
        nuInTStake
        tStake
        owner {
          id
        }
      }
    }
  `

  const OPERATOR_QUERY = gql`
    query Operator($stakingProvider: String) {
      simplePREApplication(id: $stakingProvider) {
        operator
        bondedTimestamp
        confirmedTimestamp
      }
    }
  `

  const EPOCH_STAKES_QUERY = gql`
    query EpochStakes(
      $stakingProvider: String
      $resultsPerQuery: Int
      $lastId: String
    ) {
      epochStakes(
        first: $resultsPerQuery
        where: { stakingProvider: $stakingProvider, id_gt: $lastId }
      ) {
        id
        amount
        epoch {
          id
          timestamp
        }
      }
    }
  `

  await gqlClient
    .query(STAKE_DATA_QUERY, {
      stakingProvider: stakingProvider.toLowerCase(),
    })
    .toPromise()
    .then((result) => {
      if (result.error) console.error(result.error)
      const data = result.data.stakeData
      stakeData.data.stakingProvider = data.id
      stakeData.data.owner = data.owner.id
      stakeData.data.beneficiary = data.beneficiary
      stakeData.data.authorizer = data.authorizer
      stakeData.stake.totalStaked = parseInt(data.totalStaked / 10 ** 18)
      stakeData.stake.tStake = parseInt(data.tStake / 10 ** 18)
      stakeData.stake.nuInTStake = parseInt(data.nuInTStake / 10 ** 18)
      stakeData.stake.keepInTStake = parseInt(data.keepInTStake / 10 ** 18)
    })

  await gqlClient
    .query(OPERATOR_QUERY, {
      stakingProvider: stakingProvider.toLowerCase(),
    })
    .toPromise()
    .then((result) => {
      if (result.error) console.error(result.error)
      const data = result.data.simplePREApplication
      stakeData.operator.operator = data.operator
      stakeData.operator.bondedDate = new Date(
        data.bondedTimestamp * 1000
      ).toISOString()
      stakeData.operator.confirmedDate = new Date(
        data.confirmedTimestamp * 1000
      ).toISOString()
    })

  do {
    await gqlClient
      .query(EPOCH_STAKES_QUERY, {
        stakingProvider: stakingProvider.toLowerCase(),
        resultsPerQuery: RESULTS_PER_QUERY,
        lastId: lastId,
      })
      .toPromise()
      .then((result) => {
        if (result.error) console.error(result.error)
        data = result.data?.epochStakes
        if (data.length > 0) {
          epochStakes = epochStakes.concat(data)
          lastId = data[data.length - 1].id
        }
      })
  } while (data.length > 0)

  epochStakes = epochStakes.sort((epochA, epochB) => {
    return parseInt(epochA.epoch.id) - parseInt(epochB.epoch.id)
  })

  epochStakes.forEach((epoch, index) => {
    const epochAmount = parseInt(epoch.amount)
    if (epochAmount !== amount) {
      const histElem = { epoch: epoch.epoch.id }
      if (stakeData.stakingHistory.length == 0) {
        histElem.event = "staked"
      } else {
        histElem.event = epochAmount > amount ? "topped-up" : "unstaked"
      }
      histElem.staked = (epochAmount / 10 ** 18).toFixed()
      histElem.timestamp = new Date(
        epoch.epoch.timestamp * 1000
      ).toISOString()
      stakeData.stakingHistory.push(histElem)
      amount = epochAmount
    }
  })

  return stakeData
}
