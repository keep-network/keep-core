import { BigNumber } from "@ethersproject/bignumber"
import { Contract } from "ethers"
import { program } from "commander"
import * as fs from "fs"
import { ethers } from "ethers"
import {
  abi as RandomBeaconABI,
  address as RandomBeaconAddress,
} from "@keep-network/random-beacon/artifacts/RandomBeacon.json"
import {
  abi as WalletRegistryABI,
  address as WalletRegistryAddress,
} from "@keep-network/ecdsa/artifacts/WalletRegistry.json"
import {
  abi as TokenStakingABI,
  address as TokenStakingAddress,
} from "@threshold-network/solidity-contracts/artifacts/TokenStaking.json"
import axios from "axios"
import {
  BEACON_AUTHORIZATION,
  TBTC_AUTHORIZATION,
  UP_TIME_PERCENT,
  AVG_PRE_PARAMS,
  VERSION,
  IS_BEACON_AUTHORIZED,
  IS_TBTC_AUTHORIZED,
  IS_UP_TIME_SATISFIED,
  IS_PRE_PARAMS_SATISFIED,
  IS_VERSION_SATISFIED,
  ALLOWED_UPGRADE_DELAY,
  PRECISION,
  OPERATORS_SEARCH_QUERY_STEP,
  QUERY_RESOLUTION,
  HUNDRED,
  APR,
  SECONDS_IN_YEAR,
} from "./rewards-constants"

program
  .version("0.0.1")
  .requiredOption(
    "-s, --start-timestamp <timestamp>",
    "starting time for rewards calculation"
  )
  .requiredOption(
    "-e, --end-timestamp <timestamp>",
    "ending time for rewards calculation"
  )
  .requiredOption(
    "-b, --start-block <timestamp>",
    "start block for rewards calculation"
  )
  .requiredOption(
    "-z, --end-block <timestamp>",
    "end block for rewards calculation"
  )
  .requiredOption("-a, --api <prometheus api>", "prometheus API")
  .requiredOption("-j, --job <prometheus job>", "prometheus job")
  .requiredOption(
    "-c, --october17-block <october 17 block>",
    "october 17 block"
  )
  .requiredOption(
    "-t, --october17-timestamp <october 17 timestamp>",
    "october 17 timestamp"
  )
  .requiredOption(
    "-r, --releases <client releases in a rewards interval>",
    "client releases in a rewards interval"
  )
  .requiredOption("-n, --network <name>", "network name")
  .requiredOption("-o, --output <file>", "output JSON file")
  .requiredOption(
    "-d, --output-details-path <path>",
    "output JSON details path"
  )
  .requiredOption("-q, --required-pre-params <number>", "required pre params")
  .requiredOption("-m, --required-uptime <percent>", "required uptime")
  .parse(process.argv)

// Parse the program options
const options = program.opts()
const prometheusJob = options.job
const prometheusAPI = options.api
const clientReleases = options.releases.split("|") // sorted from latest to oldest
const startRewardsTimestamp = parseInt(options.startTimestamp)
const endRewardsTimestamp = parseInt(options.endTimestamp)
const startRewardsBlock = parseInt(options.startBlock)
const endRewardsBlock = parseInt(options.endBlock)
const october17Block = parseInt(options.october17Block)
const october17Timestamp = parseInt(options.october17Timestamp)
const rewardsDataOutput = options.output
const rewardsDetailsPath = options.outputDetailsPath
const network = options.network
const requiredPreParams = options.requiredPreParams
const requiredUptime = options.requiredUptime // percent

const prometheusAPIQuery = `${prometheusAPI}/query`
// Go back in time relevant to the current date to get data for the exact
// rewards interval dates.
const offset = Math.floor(Date.now() / 1000) - endRewardsTimestamp

export async function calculateRewards() {
  if (Date.now() / 1000 < endRewardsTimestamp) {
    console.log("End time interval must be in the past")
    return "End time interval must be in the past"
  }

  const provider = new ethers.providers.EtherscanProvider(
    network,
    process.env.ETHERSCAN_TOKEN
  )

  const rewardsInterval = endRewardsTimestamp - startRewardsTimestamp
  // periodic rate rounded and adjusted because BigNumber can't operate on floating numbers.
  const periodicRate = Math.round(
    APR * (rewardsInterval / SECONDS_IN_YEAR) * PRECISION
  )
  const currentBlockNumber = await provider.getBlockNumber()

  // Query for bootstrap data that has peer instances grouped by operators
  const queryBootstrapData = `${prometheusAPI}/query_range`
  const paramsBootstrapData = {
    query: `sum by(chain_address)({job='${prometheusJob}'})`,
    start: startRewardsTimestamp,
    end: endRewardsTimestamp,
    step: OPERATORS_SEARCH_QUERY_STEP,
  }

  const bootstrapData = (
    await queryPrometheus(queryBootstrapData, paramsBootstrapData)
  ).data.result

  const operatorsData = new Array()
  const rewardsData: any = {}

  const randomBeacon = new Contract(
    RandomBeaconAddress,
    JSON.stringify(RandomBeaconABI),
    provider
  )

  const tokenStaking = new Contract(
    TokenStakingAddress,
    JSON.stringify(TokenStakingABI),
    provider
  )

  const walletRegistry = new Contract(
    WalletRegistryAddress,
    JSON.stringify(WalletRegistryABI),
    provider
  )

  console.log("Fetching AuthorizationIncreased events in rewards interval...")
  const intervalAuthorizationIncreasedEvents = await tokenStaking.queryFilter(
    "AuthorizationIncreased",
    startRewardsBlock,
    endRewardsBlock
  )

  console.log("Fetching AuthorizationDecreased events in rewards interval...")
  const intervalAuthorizationDecreasedEvents = await tokenStaking.queryFilter(
    "AuthorizationDecreaseApproved",
    startRewardsBlock,
    endRewardsBlock
  )

  console.log(
    "Fetching AuthorizationIncreased events after rewards interval..."
  )
  const postIntervalAuthorizationIncreasedEvents =
    await tokenStaking.queryFilter(
      "AuthorizationIncreased",
      endRewardsBlock,
      currentBlockNumber
    )

  console.log(
    "Fetching AuthorizationDecreased events after rewards interval..."
  )
  const postIntervalAuthorizationDecreasedEvents =
    await tokenStaking.queryFilter(
      "AuthorizationDecreaseApproved",
      endRewardsBlock,
      currentBlockNumber
    )

  for (let i = 0; i < bootstrapData.length; i++) {
    const operatorAddress = bootstrapData[i].metric.chain_address
    let authorizations = new Map<string, BigNumber>() // application: value
    let requirements = new Map<string, boolean>() // factor: true | false
    let instancesData = new Map<string, Map<string, string | number>>()
    let operatorData: any = {}

    // Staking provider should be the same for Beacon and TBTC apps
    const stakingProvider = await randomBeacon.operatorToStakingProvider(
      operatorAddress
    )
    const stakingProviderAddressForTbtc =
      await walletRegistry.operatorToStakingProvider(operatorAddress)

    if (stakingProvider !== stakingProviderAddressForTbtc) {
      console.log(
        `Staking providers for Beacon ${stakingProvider} and TBTC ${stakingProviderAddressForTbtc} must match. ` +
          `No Rewards were calculated for operator ${operatorAddress}`
      )
      continue
    }
    const { beneficiary } = await tokenStaking.rolesOf(stakingProvider)

    if (stakingProvider === ethers.constants.AddressZero) {
      console.log(
        "Staking provider cannot be zero address. " +
          `No Rewards were calculated for operator ${operatorAddress}`
      )
      continue
    }

    // Events that were emitted between the [start:end] rewards dates for a given
    // stakingProvider.
    let intervalEvents = intervalAuthorizationIncreasedEvents.concat(
      intervalAuthorizationDecreasedEvents
    )
    if (intervalEvents.length > 0) {
      intervalEvents = intervalEvents.filter(
        (event) => event.args!.stakingProvider === stakingProvider
      )
    }

    // Events that were emitted between the [end:now] dates for a given
    // stakingProvider.
    let postIntervalEvents = postIntervalAuthorizationIncreasedEvents.concat(
      postIntervalAuthorizationDecreasedEvents
    )
    if (postIntervalEvents.length > 0) {
      postIntervalEvents = postIntervalEvents.filter(
        (event) => event.args!.stakingProvider === stakingProvider
      )
    }

    /// Random Beacon application authorization requirement
    let beaconIntervalEvents = new Array()
    if (intervalEvents.length > 0) {
      beaconIntervalEvents = intervalEvents.filter(
        (obj) => obj.args!.application == randomBeacon.address
      )
    }

    let beaconPostIntervalEvents = new Array()
    if (postIntervalEvents.length > 0) {
      beaconPostIntervalEvents = postIntervalEvents.filter(
        (obj) => obj.args!.application == randomBeacon.address
      )
    }

    const beaconAuthorization = await getAuthorization(
      randomBeacon,
      beaconIntervalEvents,
      beaconPostIntervalEvents,
      stakingProvider,
      startRewardsBlock,
      endRewardsBlock,
      october17Block,
      currentBlockNumber
    )
    authorizations.set(BEACON_AUTHORIZATION, beaconAuthorization.toString())
    requirements.set(IS_BEACON_AUTHORIZED, !beaconAuthorization.isZero())

    /// tBTC application authorized requirement
    let tbtcIntervalEvents = new Array()
    if (intervalEvents.length > 0) {
      tbtcIntervalEvents = intervalEvents.filter(
        (obj) => obj.args!.application == walletRegistry.address
      )
    }

    let tbtcPostIntervalEvents = new Array()
    if (postIntervalEvents.length > 0) {
      tbtcPostIntervalEvents = postIntervalEvents.filter(
        (obj) => obj.args!.application == walletRegistry.address
      )
    }

    const tbtcAuthorization = await getAuthorization(
      walletRegistry,
      tbtcIntervalEvents,
      tbtcPostIntervalEvents,
      stakingProvider,
      startRewardsBlock,
      endRewardsBlock,
      october17Block,
      currentBlockNumber
    )

    authorizations.set(TBTC_AUTHORIZATION, tbtcAuthorization.toString())
    requirements.set(IS_TBTC_AUTHORIZED, !tbtcAuthorization.isZero())

    /// Off-chain client reqs

    // Populate instances for a given operator.
    await instancesForOperator(operatorAddress, rewardsInterval, instancesData)

    /// Uptime requirement
    let { uptimeCoefficient, isUptimeSatisfied } = await checkUptime(
      operatorAddress,
      rewardsInterval,
      instancesData
    )
    // BigNumbers cannot operate on floats. Coefficient needs to be multiplied
    // by PRECISION
    uptimeCoefficient = Math.floor(uptimeCoefficient * PRECISION)
    requirements.set(IS_UP_TIME_SATISFIED, isUptimeSatisfied)

    /// Pre-params requirement
    const isPrePramsSatisfied = await checkPreParams(
      operatorAddress,
      rewardsInterval,
      instancesData
    )

    requirements.set(IS_PRE_PARAMS_SATISFIED, isPrePramsSatisfied)

    // keep-core client already has at least 2 released versions
    const latestClient = clientReleases[0].split("_")
    const latestClientTag = latestClient[0]
    const latestClientTagTimestamp = Number(latestClient[1])
    const secondToLatestClient = clientReleases[1].split("_")
    const secondToLatestClientTag = secondToLatestClient[0]

    const instances = await processInstances(
      operatorAddress,
      rewardsInterval,
      instancesData
    )

    const upgradeCutoffDate = latestClientTagTimestamp + ALLOWED_UPGRADE_DELAY
    requirements.set(IS_VERSION_SATISFIED, true)
    if (upgradeCutoffDate < startRewardsTimestamp) {
      // v1-|-------v1 or v2------|------------------v2 only--------------|
      // ---|---------------------|---------|-----------------------------|--->
      //  v2tag                 cutoff     Feb1                          Feb28
      // All the instances must run on the latest version during the rewards
      // interval in Feb.
      for (let i = 0; i < instances.length; i++) {
        if (instances[i].buildVersion != latestClientTag) {
          requirements.set(IS_VERSION_SATISFIED, false)
        }
      }
    } else if (upgradeCutoffDate < endRewardsTimestamp) {
      // -v1-|-------v1 or v2---------|--------v2 only--------|
      // ----|---------|--------------|-----------------------|--->
      //   v2tag     Feb1          cutoff                  Feb28
      // All the instances between (upgradeCutoffDate : endRewardsTimestamp]
      // must run on the latest version
      for (let i = instances.length - 1; i >= 0; i--) {
        if (
          instances[i].lastRegisteredTimestamp > upgradeCutoffDate &&
          !instances[i].buildVersion.includes(latestClientTag)
        ) {
          // After the cutoff day a node operator still run an instance with an
          // older version. No rewards.
          requirements.set(IS_VERSION_SATISFIED, false)
          // No need to check further since at least one instance run on the
          // older version after the cutoff day.
          break
        } else {
          // It might happen that a node operator stopped an instance before the
          // upgrade cutoff date that happens to be right before the interval
          // end date. However, it might still be eligible for rewards because
          // of the uptime requirement.
          if (
            !(
              instances[i].buildVersion.includes(latestClientTag) ||
              instances[i].buildVersion.includes(secondToLatestClientTag)
            )
          ) {
            // Instance run on the older version than 2 latest.
            requirements.set(IS_VERSION_SATISFIED, false)
          }
          // No need to check other instances.
          break
        }
      }
    } else {
      // ------------v1------------|-----------v1 or v2---------|---v2 only--->
      // --|-----------------------|---------------|------------|-->
      //  Feb1                   v2tag           Feb28        cutoff
      // All the instances between [latestClientTagTimestamp : endRewardsTimestamp]
      // must run either on secondToLatest or the latest version
      for (let i = instances.length - 1; i >= 0; i--) {
        if (instances[i].lastRegisteredTimestamp >= latestClientTagTimestamp) {
          if (
            !(
              instances[i].buildVersion.includes(latestClientTag) ||
              instances[i].buildVersion.includes(secondToLatestClientTag)
            )
          ) {
            // A client run a version older than 2 latest allowed. No rewards.
            requirements.set(IS_VERSION_SATISFIED, false)
            break
          }
        } else {
          if (!instances[i].buildVersion.includes(secondToLatestClientTag)) {
            requirements.set(IS_VERSION_SATISFIED, false)
          }
          // No need to check other instances before the latestClientTagTimestamp.
          break
        }
      }
    }

    /// Start assembling peer data and weighted authorizations
    operatorData[stakingProvider] = {
      applications: Object.fromEntries(authorizations),
      instances: convertToObject(instancesData),
      requirements: Object.fromEntries(requirements),
    }

    if (
      requirements.get(IS_BEACON_AUTHORIZED) &&
      requirements.get(IS_TBTC_AUTHORIZED) &&
      requirements.get(IS_UP_TIME_SATISFIED) &&
      requirements.get(IS_PRE_PARAMS_SATISFIED) &&
      requirements.get(IS_VERSION_SATISFIED)
    ) {
      const beacon = BigNumber.from(authorizations.get(BEACON_AUTHORIZATION))
      const tbct = BigNumber.from(authorizations.get(TBTC_AUTHORIZATION))
      let minApplicationAuthorization = beacon
      if (beacon.gt(tbct)) {
        minApplicationAuthorization = tbct
      }

      rewardsData[stakingProvider] = {
        beneficiary: beneficiary,
        // amount = min(beaconWeightedAuthorization, tbtcWeightedAuthorization) * clientUptimeCoefficient * periodicRate
        amount: minApplicationAuthorization
          .mul(uptimeCoefficient)
          .mul(periodicRate)
          .div(PRECISION) // coefficient was multiplied by PRECISION earlier
          .div(PRECISION) // APR was multiplied by PRECISION earlier
          .div(HUNDRED) // APR is in %
          .toString(),
      }
    }

    operatorsData.push(operatorData)
  }

  fs.writeFileSync(rewardsDataOutput, JSON.stringify(rewardsData, null, 4))
  const detailsFileName = `${startRewardsTimestamp}-${endRewardsTimestamp}`
  fs.writeFileSync(
    rewardsDetailsPath + "/" + detailsFileName + ".json",
    JSON.stringify(operatorsData, null, 4)
  )
}

async function getAuthorization(
  application: Contract,
  intervalEvents: any[],
  postIntervalEvents: any[],
  stakingProvider: string,
  startRewardsBlock: number,
  endRewardsBlock: number,
  october17Block: number,
  currentBlockNumber: number
) {
  if (intervalEvents.length > 0) {
    return authorizationForRewardsInterval(
      intervalEvents,
      startRewardsBlock,
      endRewardsBlock,
      october17Block
    )
  }

  // Events that were emitted between the [end:firstEventDate|currentDate] dates.
  // This is used to fetch the authorization that was set during the rewards
  // interval.
  return await authorizationPostRewardsInterval(
    postIntervalEvents,
    application,
    stakingProvider,
    currentBlockNumber
  )
}

// Calculates the weighted authorization for rewards interval based on events.
// The general idea of weighted rewards calculation describes the following example.
// Please note that this example operates on dates for simplicity purposes,
// however the actual calculation is based on block numbers.
// Ex.
// events:         ^     ^      *    *   ^
// timeline:  -|---|-----|------|----|---|--------|-->
//         Sep 0   3     8      14  18   22       30
// where: '^' denotes increase in authorization
//        '*' denotes decrease in authorization
//         0 -> Sep 1 at 00:00:00
// event authorizations:
//  Sep 0 - 3 from 100k to 110k
//  Sep 3 - 8 from 110k to 135k
//  Sep 8 - 14 from 135k to 120k
//  Sep 14 - 18 from 120k to 100k
//  Sep 18 - 30 constant 100k (last sub-interval)
// Weighted authorization = (3-0)/30*100 + (8-3)/30*110 + (14-8)/30*135
//                        + (18-14)/30*120 + (30-18)/30*100
// October 2022 is a special month for rewards calculation.  If a node was set
// after Oct 1 but prior to Oct 17, then we calculate the rewards for the entire
// month.
// See https://blog.threshold.network/tbtc-v2-hits-its-first-launch-milestone/
// E.g. 1
// First and only event was on Oct 10. The authorization is calculated for the
// entire month.
// E.g. 2
// First increase event was on Oct 10 from 0 to 100k
// Second increase was on Oct 15 from 100k to 150k
// Authorization of 100k is interpolated for the dates between Oct 1 - Oct 10
// Authorization for Oct 1 - Oct 15 is now 100k; coefficient 15/30
// Authorization between Oct 15 - Oct 30 is 150k; coefficent 15/30
// Weighted authorization: 15/30 * 100k + 15/30 * 150k
function authorizationForRewardsInterval(
  intervalEvents: any[],
  startRewardsBlock: number,
  endRewardsBlock: number,
  october17Block: number
) {
  let authorization = BigNumber.from("0")
  const deltaRewardsBlock = endRewardsBlock - startRewardsBlock
  // ascending order
  intervalEvents.sort((a, b) => a.blockNumber - b.blockNumber)

  let tmpBlock = startRewardsBlock // prev tmp block
  let firstEventBlock = intervalEvents[0].blockNumber

  let index = 0
  if (firstEventBlock < october17Block) {
    index = 1
  }

  for (let i = index; i < intervalEvents.length; i++) {
    const eventBlock = intervalEvents[i].blockNumber
    const coefficient = Math.floor(
      ((eventBlock - tmpBlock) / deltaRewardsBlock) * PRECISION
    )
    authorization = authorization.add(
      intervalEvents[i].args.fromAmount.mul(coefficient)
    )
    tmpBlock = eventBlock
  }

  // calculating authorization for the last sub-interval
  const coefficient = Math.floor(
    ((endRewardsBlock - tmpBlock) / deltaRewardsBlock) * PRECISION
  )
  authorization = authorization.add(
    intervalEvents[intervalEvents.length - 1].args.toAmount.mul(coefficient)
  )

  return authorization.div(PRECISION)
}

// Get the authorization from the first event that occurred after the rewards
// interval. If no events were emitted, then get the authorization from the current
// block.
async function authorizationPostRewardsInterval(
  postIntervalEvents: any[],
  application: Contract,
  stakingProvider: string,
  currentBlockNumber: number
) {
  // Sort events in ascending order
  postIntervalEvents.sort((a, b) => a.blockNumber - b.blockNumber)

  if (
    postIntervalEvents.length > 0 &&
    postIntervalEvents[0].blockNumber < currentBlockNumber
  ) {
    // There are events (increase or decrease) present after the rewards interval
    // and before the current block. Take the "fromAmount", because it was the
    // same as for the rewards interval dates.
    return postIntervalEvents[0].args.fromAmount
  }

  // There were no authorization events after the rewards interval and before
  // the current block.
  // Current authorization is the same as the authorization at the end of the
  // rewards interval.
  const authorization = await application.eligibleStake(stakingProvider)
  return authorization
}

async function instancesForOperator(
  operatorAddress: any,
  rewardsInterval: number,
  instancesData: Map<string, Map<string, string | number>>
) {
  // Resolution is defaulted to Prometheus settings.
  const instancesDataByOperatorParams = {
    query: `present_over_time(up{chain_address="${operatorAddress}", job="${prometheusJob}"}
                [${rewardsInterval}s] offset ${offset}s)`,
  }
  const instancesDataByOperator = (
    await queryPrometheus(prometheusAPIQuery, instancesDataByOperatorParams)
  ).data.result

  instancesDataByOperator.forEach(
    (element: { metric: { instance: string } }) => {
      const instanceData = new Map<string, string | number>()
      instancesData.set(element.metric.instance, instanceData)
    }
  )
}

// Peer uptime requirement. The total uptime for all the instances for a given
// operator has to be greater than 96% to receive the rewards.
async function checkUptime(
  operatorAddress: string,
  rewardsInterval: number,
  instancesData: Map<string, Map<string, string | number>>
) {
  const paramsOperatorUptime = {
    query: `up{chain_address="${operatorAddress}", job="${prometheusJob}"}
            [${rewardsInterval}s:${QUERY_RESOLUTION}s] offset ${offset}s`,
  }

  const instances = (
    await queryPrometheus(prometheusAPIQuery, paramsOperatorUptime)
  ).data.result

  // First registered 'up' metric in a given interval <start:end> for a given
  // operator. Start evaluating uptime from this point.
  const firstRegisteredUptime = instances.reduce(
    (currentMin: number, instance: any) =>
      Math.min(instance.values[0][0], currentMin),
    Number.MAX_VALUE
  )

  let uptimeSearchRange = endRewardsTimestamp - firstRegisteredUptime

  const paramsSumUptimes = {
    query: `sum_over_time(up{chain_address="${operatorAddress}", job="${prometheusJob}"}
            [${uptimeSearchRange}s:${QUERY_RESOLUTION}s] offset ${offset}s)
            * ${QUERY_RESOLUTION} / ${uptimeSearchRange}`,
  }

  const uptimesByInstance = (
    await queryPrometheus(prometheusAPIQuery, paramsSumUptimes)
  ).data.result

  let sumUptime = 0
  for (let i = 0; i < uptimesByInstance.length; i++) {
    const instance = uptimesByInstance[i]
    const uptime = instance.value[1] * HUNDRED
    const dataInstance = instancesData.get(instance.metric.instance)
    if (dataInstance !== undefined) {
      dataInstance.set(UP_TIME_PERCENT, uptime)
    } else {
      // Should not happen
      console.error("Instance must be present for a given rewards interval.")
    }

    sumUptime += uptime
  }

  const isUptimeSatisfied = sumUptime >= requiredUptime
  // October is a special month for rewards calculation. If a node was set before
  // October 17th, then it is eligible for the entire month of rewards. Uptime of
  // a running node still need to meet the uptime requirement after it was set.
  if (firstRegisteredUptime < october17Timestamp) {
    uptimeSearchRange = rewardsInterval
  }

  const uptimeCoefficient = isUptimeSatisfied
    ? uptimeSearchRange / rewardsInterval
    : 0
  return { uptimeCoefficient, isUptimeSatisfied }
}

async function checkPreParams(
  operatorAddress: string,
  rewardsInterval: number,
  dataInstances: Map<string, Map<string, string | number>>
) {
  // Avg of pre-params across all the instances for a given operator in the rewards
  // interval dates. Resolution is defaulted to Prometheus settings.
  const paramsPreParams = {
    query: `avg_over_time(tbtc_pre_params_count{chain_address="${operatorAddress}", job="${prometheusJob}"}
              [${rewardsInterval}s:${QUERY_RESOLUTION}s] offset ${offset}s)`,
  }

  const preParamsAvgByInstance = (
    await queryPrometheus(prometheusAPIQuery, paramsPreParams)
  ).data.result

  let sumPreParams = 0
  for (let i = 0; i < preParamsAvgByInstance.length; i++) {
    const instance = preParamsAvgByInstance[i]
    const preParams = parseInt(instance.value[1]) // [timestamp, value]
    const dataInstance = dataInstances.get(instance.metric.instance)
    if (dataInstance !== undefined) {
      dataInstance.set(AVG_PRE_PARAMS, preParams)
    } else {
      // Should not happen
      console.error("Instance must be present for a given rewards interval.")
    }

    sumPreParams += preParams
  }

  const preParamsAvg = sumPreParams / preParamsAvgByInstance.length
  return preParamsAvg >= requiredPreParams
}

// Query Prometheus and fetch instances that run on either of two latest client
// versions and mark their first and last registered timestamp.
async function processInstances(
  operatorAddress: string,
  rewardsInterval: number,
  instancesData: Map<string, Map<string, string | number>>
) {
  const buildVersionInstancesParams = {
    query: `client_info{chain_address="${operatorAddress}", job="${prometheusJob}"}[${rewardsInterval}s:${QUERY_RESOLUTION}s] offset ${offset}s`,
  }
  // Get instances data for a given rewards interval
  const queryBuildVersionInstances = (
    await queryPrometheus(prometheusAPIQuery, buildVersionInstancesParams)
  ).data.result

  let instances = []

  // Determine client's build version for it all it's instances
  for (let i = 0; i < queryBuildVersionInstances.length; i++) {
    const instance = queryBuildVersionInstances[i]

    const instanceTimestampsVersionInfo = {
      // First timestamp registered by Prometheus for a given instance
      firstRegisteredTimestamp: instance.values[0][0],
      // Last timestamp registered by Prometheus for a given instance
      lastRegisteredTimestamp: instance.values[instance.values.length - 1][0],
      buildVersion: instance.metric.version,
    }

    instances.push(instanceTimestampsVersionInfo)

    const dataInstance = instancesData.get(instance.metric.instance)
    if (dataInstance !== undefined) {
      dataInstance.set(VERSION, instance.metric.version)
    } else {
      // Should not happen
      console.error("Instance must be present for a given rewards interval.")
    }
  }

  // Sort instances in ascending order by first registration timestamp
  instances.sort((a, b) =>
    a.firstRegisteredTimestamp > b.firstRegisteredTimestamp ? 1 : -1
  )

  return instances
}

function convertToObject(map: Map<string, Map<string, any>>) {
  let obj: { [k: string]: any } = {}
  map.forEach((value: Map<string, any>, key: string) => {
    const result = Object.fromEntries(value)
    obj[key] = result
  })

  return obj
}

async function queryPrometheus(url: string, params: any): Promise<any> {
  try {
    const { data } = await axios.get(url, { params: params })

    return data
  } catch (error) {
    if (axios.isAxiosError(error)) {
      console.log("error message: ", error.message)
      return error.message
    } else {
      console.log("unexpected error: ", error)
      return "An unexpected error occurred"
    }
  }
}

calculateRewards()
