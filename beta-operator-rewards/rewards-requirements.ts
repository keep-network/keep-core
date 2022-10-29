import { BigNumber } from "@ethersproject/bignumber";
import { Contract } from "ethers";
import { program } from "commander";
import * as fs from "fs";
import { ethers } from "ethers";
import {
  abi as RandomBeaconABI,
  address as RandomBeaconAddress,
} from "@keep-network/random-beacon/artifacts/RandomBeacon.json";
import {
  abi as WalletRegistryABI,
  address as WalletRegistryAddress,
} from "@keep-network/ecdsa/artifacts/WalletRegistry.json";
import {
  abi as TokenStakingABI,
  address as TokenStakingAddress,
} from "@threshold-network/solidity-contracts/artifacts/TokenStaking.json";
import axios from "axios";
import {
  DEFAULT_NETWORK,
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
  REQUIRED_UPTIME_PERCENT,
  REQUIRED_MIN_PRE_PARAMS,
  ALLOWED_UPGRADE_DELAY,
  PRECISION,
  QUERY_STEP,
  QUERY_RESOLUTION,
  HUNDRED,
} from "./rewards-constants";

const provider = new ethers.providers.EtherscanProvider(
  DEFAULT_NETWORK,
  process.env.ETHERSCAN_TOKEN
);

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
  .requiredOption("-i, --interval <timestamp>", "scrape interval") // IMPORTANT! Must match Prometheus config
  .requiredOption("-a, --api <prometheus api>", "prometheus API")
  .requiredOption("-j, --job <prometheus job>", "prometheus job")
  .requiredOption(
    "-r, --releases <client releases in a rewards interval>",
    "client releases in a rewards interval"
  )
  .requiredOption("-o, --output <file>", "output JSON file")
  .parse(process.argv);

// Parse the program options
const options = program.opts();
const prometheusJob = options.job;
const prometheusAPI = options.api;
const clientReleases = options.releases.split("|"); // sorted from latest to oldest
const startRewardsTimestamp = parseInt(options.startTimestamp);
const endRewardsTimestamp = parseInt(options.endTimestamp);
const startRewardsBlock = parseInt(options.startBlock);
const endRewardsBlock = parseInt(options.endBlock);
const scrapeInterval = parseInt(options.interval); // TODO: might not be needed.
const peersDataFile = options.output;

const prometheusAPIQuery = `${prometheusAPI}/query`;
// Go back in time relevant to the current date to get data for the exact
// rewards interval dates.
const offset = Math.floor(Date.now() / 1000) - endRewardsTimestamp;

export async function runRewardsRequirements() {
  if (Date.now() / 1000 < endRewardsTimestamp) {
    console.log("End time interval must be in the past");
    return "End time interval must be in the past";
  }

  const currentBlockNumber = await provider.getBlockNumber();
  const rewardsInterval = endRewardsTimestamp - startRewardsTimestamp;

  // Query for bootstrap data that has peer instances grouped by operators
  const queryBootstrapData = `${prometheusAPI}/query_range`;
  const paramsBootstrapData = {
    query: `sum by(chain_address)({job='${prometheusJob}'})`,
    start: startRewardsTimestamp,
    end: endRewardsTimestamp,
    step: QUERY_STEP,
  };

  const bootstrapData = (
    await queryPrometheus(queryBootstrapData, paramsBootstrapData)
  ).data.result;

  let peersData = new Array();
  let weightedAuthorizations = new Array();

  const randomBeacon = new Contract(
    RandomBeaconAddress,
    JSON.stringify(RandomBeaconABI),
    provider
  );

  const tokenStaking = new Contract(
    TokenStakingAddress,
    JSON.stringify(TokenStakingABI),
    provider
  );

  const walletRegistry = new Contract(
    WalletRegistryAddress,
    JSON.stringify(WalletRegistryABI),
    provider
  );

  console.log("Fetching AuthorizationIncreased events in rewards interval...");
  const allIntevalAuthorizationIncreasedEvents = await tokenStaking.queryFilter(
    "AuthorizationIncreased",
    startRewardsBlock,
    endRewardsBlock
  );
  const intevalAuthorizationIncreasedEvents = filterEventsByApplications(
    allIntevalAuthorizationIncreasedEvents
  );

  console.log("Fetching AuthorizationDecreased events in rewards interval...");
  const allIntervalAuthorizationDecreasedEvents =
    await tokenStaking.queryFilter(
      "AuthorizationDecreaseApproved",
      startRewardsBlock,
      endRewardsBlock
    );
  const intervalAuthorizationDecreasedEvents = filterEventsByApplications(
    allIntervalAuthorizationDecreasedEvents
  );

  console.log(
    "Fetching AuthorizationIncreased events after rewards interval..."
  );
  const allPostIntervalAuthorizationIncreasedEvents =
    await tokenStaking.queryFilter(
      "AuthorizationIncreased",
      endRewardsBlock,
      currentBlockNumber
    );
  const postIntervalIncreasedEvents = filterEventsByApplications(
    allPostIntervalAuthorizationIncreasedEvents
  );

  console.log(
    "Fetching AuthorizationDecreased events after rewards interval..."
  );
  const allPostIntervalAuthorizationDecreasedEvents =
    await tokenStaking.queryFilter(
      "AuthorizationDecreaseApproved",
      endRewardsBlock,
      currentBlockNumber
    );
  const postIntervalDecreasedEvents = filterEventsByApplications(
    allPostIntervalAuthorizationDecreasedEvents
  );

  for (let i = 0; i < bootstrapData.length; i++) {
    const operatorAddress = bootstrapData[i].metric.chain_address;
    let authorizations = new Map<string, BigNumber>(); // application: value
    let requirements = new Map<string, boolean>(); // factor: true | false
    let instancesData = new Map<string, Map<string, string | number>>();
    let peerData: any = {};
    let weightedAuthorization: any = {};

    // Staking provider should be the same for Beacon and TBTC apps
    const stakingProvider = await randomBeacon.operatorToStakingProvider(
      operatorAddress
    );
    const stakingProviderAddressForTbtc =
      await walletRegistry.operatorToStakingProvider(operatorAddress);

    if (stakingProvider !== stakingProviderAddressForTbtc) {
      console.log(
        `Staking providers for Beacon ${stakingProvider} and TBTC ${stakingProviderAddressForTbtc} must match. ` +
          `No Rewards were calculated for operator ${operatorAddress}`
      );
      continue;
    }

    if (stakingProvider === ethers.constants.AddressZero) {
      console.log(
        `Staking provider cannot be zero address. ` +
          `No Rewards were calculated for operator ${operatorAddress}`
      );
      continue;
    }

    // Populate instances for a given operator.
    await instancesForOperator(operatorAddress, rewardsInterval, instancesData);

    // Events that were emitted between the [start:end] rewards dates for a given
    // stakingProvider.
    const intervalEvents = intevalAuthorizationIncreasedEvents
      .concat(intervalAuthorizationDecreasedEvents)
      .filter((event) => event.args.stakingProvider === stakingProvider);

    /// Random Beacon application authorization requirement
    let beaconAuthorization = await getAuthorization(
      randomBeacon,
      intervalEvents,
      postIntervalIncreasedEvents.concat(postIntervalDecreasedEvents),
      stakingProvider,
      startRewardsBlock,
      endRewardsBlock
    );

    authorizations.set(BEACON_AUTHORIZATION, beaconAuthorization.toString());
    requirements.set(IS_BEACON_AUTHORIZED, !beaconAuthorization.isZero());

    /// tBTC application authorized requirement
    const tbtcAuthorization = await getAuthorization(
      walletRegistry,
      intervalEvents,
      postIntervalIncreasedEvents.concat(postIntervalDecreasedEvents),
      stakingProvider,
      startRewardsBlock,
      endRewardsBlock
    );

    authorizations.set(TBTC_AUTHORIZATION, tbtcAuthorization.toString());
    requirements.set(IS_TBTC_AUTHORIZED, !tbtcAuthorization.isZero());

    /// Uptime requirement
    let { uptimeCoefficient, isUptimeSatisfied } = await checkUptime(
      operatorAddress,
      rewardsInterval,
      instancesData
    );
    // BigNumbers cannot operate on floats. Coefficient needs to be multiplied by 100
    uptimeCoefficient = Math.floor(uptimeCoefficient * HUNDRED);
    requirements.set(IS_UP_TIME_SATISFIED, isUptimeSatisfied);

    /// Pre-params requiremnt
    const isPrePramsSatisfied = await checkPreParams(
      operatorAddress,
      rewardsInterval,
      instancesData
    );

    requirements.set(IS_PRE_PARAMS_SATISFIED, isPrePramsSatisfied);

    /// Version requirement
    let {
      instanceWithLatestBuildVersion,
      latestRegisteredBuildVersionTimestmap,
    } = await checkVersion(operatorAddress, rewardsInterval, instancesData);

    if (instanceWithLatestBuildVersion === undefined) {
      console.log(
        `Cannot determine a client version. No Rewards were calculated for operator ${operatorAddress}`
      );
      continue;
    }

    // This is an example to illustrate a client's build version requirement.
    // A client must be either on the second to latest version or latest.
    //                       v1 or v2                v2
    //           |-------------------------------\|-------|
    // Timeline -|-------------*------------------*-------|->
    //         Sep1            v2             v2+delay   Sep30
    // Where:
    // v2 was released Sep10
    // delay = 14days
    // v2 + delay = Sep10 + 14days = Sep24
    // Between Sep1 - Sep24 a client is allowed to run v1 or v2
    // Between Sep24 - Sep30 a client is allowed to run only v2

    const buildVersion = instanceWithLatestBuildVersion.metric.version;
    const latestClientRelease = clientReleases[0].split("_");
    const latestClientTag = latestClientRelease[0];
    const latestClientReleaseTimestamp = latestClientRelease[1];
    if (clientReleases.length > 1) {
      // A client is allowed to be on either of the two latest releases.
      const secondToLatestClientRelease = clientReleases[1].split("_");
      const secondToLatestClientTag = secondToLatestClientRelease[0];

      let allowedDelayEndTimestamp =
        latestClientReleaseTimestamp + ALLOWED_UPGRADE_DELAY;
      if (allowedDelayEndTimestamp > endRewardsTimestamp) {
        allowedDelayEndTimestamp = endRewardsTimestamp;
      }

      if (latestRegisteredBuildVersionTimestmap <= allowedDelayEndTimestamp) {
        // A client's version can be on either 2 latest versions
        requirements.set(
          IS_VERSION_SATISFIED,
          buildVersion.includes(latestClientTag) ||
            buildVersion.includes(secondToLatestClientTag)
        );
      } else {
        // The allowed delay for an upgrade is over. A client should be on the
        // latest build version.
        requirements.set(
          IS_VERSION_SATISFIED,
          buildVersion.includes(latestClientTag)
        );
      }
    } else {
      requirements.set(
        IS_VERSION_SATISFIED,
        buildVersion.includes(latestClientTag)
      );
    }

    /// Start assembling peer data and weighted authorizations
    peerData[stakingProvider] = {
      applications: Object.fromEntries(authorizations),
      instances: convertToObject(instancesData),
      requirements: Object.fromEntries(requirements),
    };

    if (
      requirements.get(IS_BEACON_AUTHORIZED) &&
      requirements.get(IS_TBTC_AUTHORIZED) &&
      requirements.get(IS_UP_TIME_SATISFIED) &&
      requirements.get(IS_PRE_PARAMS_SATISFIED) &&
      requirements.get(IS_VERSION_SATISFIED)
    ) {
      const beacon = BigNumber.from(authorizations.get(BEACON_AUTHORIZATION));
      const tbct = BigNumber.from(authorizations.get(TBTC_AUTHORIZATION));
      let minApplicationAuthorization = beacon;
      if (beacon.gt(tbct)) {
        minApplicationAuthorization = tbct;
      }
      // TODO: - adjust by APR 15%, ie *1.25%
      //       - make APR an input var
      weightedAuthorization[stakingProvider] = {
        // beaneficiary: <address> TODO: implement
        weightedAuthorization: minApplicationAuthorization
          .mul(uptimeCoefficient)
          .div(HUNDRED)
          .toString(),
      };
      weightedAuthorizations.push(weightedAuthorization);
    }

    peersData.push(peerData);
  }

  console.log("peersData: ", JSON.stringify(peersData, null, 2));
  console.log(
    "weightedAuthorizations: ",
    JSON.stringify(weightedAuthorizations, null, 2)
  );
}

async function getAuthorization(
  application: Contract,
  intervalEvents: any[],
  postEvents: any[],
  stakingProvider: string,
  startRewardsBlock: number,
  endRewardsBlock: number
) {
  // When there were no events during the rewards interval, then we fetch the
  // authorization after the interval which was the same as for the actual
  // rewards interval.
  if (intervalEvents.length == 0) {
    // Events that were emitted between the [end:firstEventDate|currentDate] dates.
    // This is used to fetch the authorization that was allocated during the rewards
    // interval.
    const postIntervalEvents = postEvents.filter(
      (event) => event.args.stakingProvider === stakingProvider
    );

    return await authorizationPostRewardsInterval(
      postIntervalEvents,
      application,
      stakingProvider
    );
  }

  // There is at least one event emitted during the rewards interval
  const applicationEvents = intervalEvents.filter((obj) => {
    return obj.args.application == application.address;
  });

  return authorizationForRewardsInterval(
    applicationEvents,
    startRewardsBlock,
    endRewardsBlock
  );
}

function filterEventsByApplications(events: any[]) {
  return events.filter((event) => {
    return (
      event.args.application === RandomBeaconAddress ||
      event.args.application === WalletRegistryAddress
    );
  });
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
// authorization = (3-0)/30*100 + (8-3)/30*110 + (14-8)/30*135
//               + (18-14)/30*120 + (30-18)/30*100
function authorizationForRewardsInterval(
  events: any[],
  startRewardsBlock: number,
  endRewardsBlock: number
) {
  if (events.length == 0) {
    return BigNumber.from("0");
  }

  let authorization = BigNumber.from("0");
  const deltaRewardsBlock = endRewardsBlock - startRewardsBlock;
  // ascending order
  events.sort((a, b) => a.blockNumber - b.blockNumber);

  let tmpBlock = startRewardsBlock; // prev tmp block
  for (let i = 0; i < events.length; i++) {
    const event = events[i];
    const coefficient = Math.floor(
      ((event.blockNumber - tmpBlock) / deltaRewardsBlock) * PRECISION
    );
    authorization = authorization.add(event.args.fromAmount.mul(coefficient));
    tmpBlock = event.blockNumber;
  }
  authorization = authorization.div(PRECISION);

  // calculating authorization for the last sub-interval
  const coefficient = Math.floor(
    ((endRewardsBlock - tmpBlock) / deltaRewardsBlock) * PRECISION
  );
  authorization = authorization.add(
    events[events.length - 1].args.toAmount.mul(coefficient)
  );

  return authorization.div(PRECISION);
}

// Get the authorization from the first event that occured after the rewards
// interval. If no events were emitted, then get the authorization from the current
// block.
async function authorizationPostRewardsInterval(
  events: any[],
  application: Contract,
  stakingProvider: string
) {
  const currentBlockNumber = await provider.getBlockNumber();
  const randomBeaconEvents = events.filter((obj) => {
    return obj.application == RandomBeaconABI;
  });

  // Sort events in ascending order
  randomBeaconEvents.sort((a, b) => a.blockNumber - b.blockNumber);

  if (events.length > 0 && events[0].blockNumber < currentBlockNumber) {
    // There are events present after the rewards interval and before the
    // current block.
    return events[0].args.fromAmount;
  }

  // There were no authorization events after the rewards interval and before
  // the current block.
  // Current authorization is the same as the authorization at the end of the
  // rewards interval.
  return await application.eligibleStake(stakingProvider);
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
  };
  const instancesDataByOperator = (
    await queryPrometheus(prometheusAPIQuery, instancesDataByOperatorParams)
  ).data.result;

  instancesDataByOperator.forEach(
    (element: { metric: { instance: string } }) => {
      const instanceData = new Map<string, string | number>();
      instancesData.set(element.metric.instance, instanceData);
    }
  );
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
  };

  const instances = (
    await queryPrometheus(prometheusAPIQuery, paramsOperatorUptime)
  ).data.result;

  // First registered 'up' metric in a given interval <start:end> for a given
  // operator. Start evaluating uptime from this point.
  const firstRegisteredUptime = instances.reduce(
    (currentMin: number, instance: any) =>
      Math.min(instance.values[0][0], currentMin),
    Number.MAX_VALUE
  );

  const uptimeSearchRange = endRewardsTimestamp - firstRegisteredUptime;

  const paramsSumUptimes = {
    query: `sum_over_time(up{chain_address="${operatorAddress}", job="${prometheusJob}"}
            [${uptimeSearchRange}s:${QUERY_RESOLUTION}s] offset ${offset}s) 
            * ${QUERY_RESOLUTION} / ${uptimeSearchRange}`,
  };

  const uptimesByInstance = (
    await queryPrometheus(prometheusAPIQuery, paramsSumUptimes)
  ).data.result;

  let sumUptime = 0;
  for (let i = 0; i < uptimesByInstance.length; i++) {
    const instance = uptimesByInstance[i];
    const uptime = instance.value[1] * HUNDRED;
    const dataInstance = instancesData.get(instance.metric.instance);
    if (dataInstance !== undefined) {
      dataInstance.set(UP_TIME_PERCENT, uptime);
    } else {
      // Should not happen
      console.error("Instance must be present for a given rewards interval.");
    }

    sumUptime += uptime;
  }

  const isUptimeSatisfied = sumUptime >= REQUIRED_UPTIME_PERCENT;

  const uptimeCoefficient = isUptimeSatisfied
    ? uptimeSearchRange / rewardsInterval
    : 0;
  return { uptimeCoefficient, isUptimeSatisfied };
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
  };

  const preParamsAvgByInstance = (
    await queryPrometheus(prometheusAPIQuery, paramsPreParams)
  ).data.result;

  let sumPreParams = 0;
  for (let i = 0; i < preParamsAvgByInstance.length; i++) {
    const instance = preParamsAvgByInstance[i];
    const preParams = parseInt(instance.value[1]); // [timestamp, value]
    const dataInstance = dataInstances.get(instance.metric.instance);
    if (dataInstance !== undefined) {
      dataInstance.set(AVG_PRE_PARAMS, preParams);
    } else {
      // Should not happen
      console.error("Instance must be present for a given rewards interval.");
    }

    sumPreParams += preParams;
  }

  const preParamsAvg = sumPreParams / preParamsAvgByInstance.length;
  return preParamsAvg >= REQUIRED_MIN_PRE_PARAMS;
}

async function checkVersion(
  operatorAddress: string,
  rewardsInterval: number,
  instancesData: Map<string, Map<string, string | number>>
) {
  const buildVersionInstancesParams = {
    query: `client_info{chain_address="${operatorAddress}", job="${prometheusJob}"}[${rewardsInterval}s:${QUERY_RESOLUTION}s] offset ${offset}s`,
  };
  // Get build versions of instances in rewards interval
  const queryBuildVersionInstances = (
    await queryPrometheus(prometheusAPIQuery, buildVersionInstancesParams)
  ).data.result;

  let instanceWithLatestBuildVersion;
  let latestRegisteredBuildVersionTimestmap = 0; // min number

  // Determine client's build version for it all it's instances
  for (let i = 0; i < queryBuildVersionInstances.length; i++) {
    const instance = queryBuildVersionInstances[i];
    // Find latest registered timestamp in a given instance
    if (
      instance.values[instance.values.length - 1][0] >
      latestRegisteredBuildVersionTimestmap
    ) {
      latestRegisteredBuildVersionTimestmap =
        instance.values[instance.values.length - 1][0];
      instanceWithLatestBuildVersion = instance;
    }
    const dataInstance = instancesData.get(instance.metric.instance);
    if (dataInstance !== undefined) {
      dataInstance.set(VERSION, instance.metric.version);
    } else {
      // Should not happen
      console.error("Instance must be present for a given rewards interval.");
    }
  }
  return {
    instanceWithLatestBuildVersion,
    latestRegisteredBuildVersionTimestmap,
  };
}

function convertToObject(map: Map<string, Map<string, any>>) {
  let obj: { [k: string]: any } = {};
  map.forEach((value: Map<string, any>, key: string) => {
    const result = Object.fromEntries(value);
    obj[key] = result;
  });

  return obj;
}

async function queryPrometheus(url: string, params: any): Promise<any> {
  try {
    const { data } = await axios.get(url, { params: params });

    return data;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      console.log("error message: ", error.message);
      return error.message;
    } else {
      console.log("unexpected error: ", error);
      return "An unexpected error occurred";
    }
  }
}

runRewardsRequirements();
