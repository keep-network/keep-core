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
  ALLOWED_UPGRADE_DELAY,
  CLIENT_TIMESTAMP_INDEX,
  CLIENT_VERSION_INDEX,
  DEFAULT_NETWORK,
  BEACON_AUTHORIZATION,
  IS_BEACON_AUTHORIZED_FACTOR,
  TBTC_AUTHORIZATION,
  IS_TBTC_AUTHORIZED_FACTOR,
  MIN_PRE_PARAMS,
  PRE_PARAMS_AVG_INTERVAL,
  PRE_PARAMS_FACTOR,
  PRE_PARAMS_RESOLUTION,
  QUERY_STEP,
  REQUIRED_UPTIME,
  UPTIME_REWARDS_COEFFICIENT,
  UP_TIME_FACTOR,
  VERSION_FACTOR,
  PRECISION,
} from "./rewards-constants";

const provider = new ethers.providers.EtherscanProvider(
  DEFAULT_NETWORK,
  process.env.ETHERSCAN_TOKEN
);

export async function calculateRewardsFactors() {
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
      "-v, --versions <client version(s) in a rewards interval>",
      "client version(s) in a rewards interval"
    )
    .requiredOption("-o, --output <file>", "output JSON file")
    .parse(process.argv);

  // Parse the program options
  const options = program.opts();
  const prometheusJob = options.job;
  const prometheusAPI = options.api;
  const clientVersions = options.versions.split("|"); // sorted from latest to oldest
  const startRewardsTimestamp = parseInt(options.startTimestamp);
  const endRewardsTimestamp = parseInt(options.endTimestamp);
  const startRewardsBlock = parseInt(options.startBlock);
  const endRewardsBlock = parseInt(options.endBlock);
  const scrapeInterval = parseInt(options.interval);
  const peersDataFile = options.output;
  const rewardsIntervalBlocksDelta = endRewardsBlock - startRewardsBlock;
  // End program option parsing

  if (Date.now() / 1000 < endRewardsTimestamp) {
    console.log("End time interval must be in the past");
    return "End time interval must be in the past";
  }

  const currentBlockNumber = await provider.getBlockNumber();
  const rewardsInterval = endRewardsTimestamp - startRewardsTimestamp;

  const prometheusAPIQuery = `${prometheusAPI}/query`;
  const queryBootstrapData = `${prometheusAPI}/query_range`;

  // Query for bootstrap data that has peer instances
  const paramsBootstrapData = {
    query: `sum by(chain_address)({job='${prometheusJob}'})`,
    start: startRewardsTimestamp,
    end: endRewardsTimestamp,
    step: QUERY_STEP,
  };

  const bootstrapDataByOperator = (
    await queryPrometheus(queryBootstrapData, paramsBootstrapData)
  ).data.result;

  // TODO: filter by 1 address and fetch authorizations by 1 address
  // you have multiple addresses because of the instances
  // console.log("bootstrapData", bootstrapData)

  let peersData = new Array();

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

  console.log("Fetching AuthorizationIncreased events in rewards interval..");
  const allIntevalAuthorizationIncreasedEvents = await tokenStaking.queryFilter(
    "AuthorizationIncreased",
    startRewardsBlock,
    endRewardsBlock
  );
  const intevalAuthorizationIncreasedEvents = filterEventsByApplications(
    allIntevalAuthorizationIncreasedEvents
  );

  console.log("Fetching AuthorizationDecreased events in rewards interval..");
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
    "Fetching AuthorizationIncreased events after rewards interval.."
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
    "Fetching AuthorizationDecreased events after rewards interval.."
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

  for (let i = 0; i < bootstrapDataByOperator.length; i++) {
    const peer = bootstrapDataByOperator[i];
    let authorizations = new Map<string, BigNumber>(); // application: value
    let requirements = new Map<string, boolean>(); // factor: true | false
    let peerData: any = {};

    // Staking provider should be the same for Beacon and TBTC apps
    const stakingProvider = await randomBeacon.operatorToStakingProvider(
      peer.metric.chain_address
    );
    const stakingProviderAddressForTbtc =
      await walletRegistry.operatorToStakingProvider(peer.metric.chain_address);

    if (stakingProvider !== stakingProviderAddressForTbtc) {
      console.log(
        `Staking providers for Beacon ${stakingProvider} and TBTC ${stakingProviderAddressForTbtc} must match. ` +
          `No Rewards were calculated for operator ${peer.metric.chain_address}`
      );
      continue;
    }

    if (stakingProvider === ethers.constants.AddressZero) {
      console.log(
        `Staking provider cannot be zero address. ` +
          `No Rewards were calculated for operator ${peer.metric.chain_address}`
      );
      continue;
    }

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
    requirements.set(
      IS_BEACON_AUTHORIZED_FACTOR,
      !beaconAuthorization.isZero()
    );

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
    requirements.set(IS_TBTC_AUTHORIZED_FACTOR, !tbtcAuthorization.isZero());

    /// Start assembling peer data and weighted authorizations
    peerData[stakingProvider] = {
      applications: Object.fromEntries(authorizations),
      requirements: Object.fromEntries(requirements),
    };

    // TODO: iterate over instances for a given operator address
    // Assemble
    //   "instances":[
    //     {
    //        "uptimePercent":4.251766217084136,
    //        "preParams":132,
    //        "version":"2.0.0-1m",
    //        "ip":"10.102.0.30:9701"
    //     },
    //    (...)
    //    ],

    peersData.push(peerData);
  }

  console.log("peersData: ", JSON.stringify(peersData, null, 2));

  // TODO: calculate weighted authorization for a given staking provider
  // <staking provider address>: {
  //   beneficiary: <address>,
  //   weightedAuthorization: <amount> // uptimeCoeficient * min(beaconAuthorization, tbtcAuthorization)
  // }

  // TODO: Save to file
  // - all requirements
  // - weighted authorization for NU team
  // const jsonObject = await convertToJSON(peersData);
  // fs.writeFileSync(peersDataFile, JSON.stringify(jsonObject, null, 2));
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
  randomBeaconEvents.sort(function (a, b) {
    return a.blockNumber < b.blockNumber
      ? -1
      : a.blockNumber == b.blockNumber
      ? 0
      : 1;
  });

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

async function convertToJSON(map: Map<string, Map<string, any>>) {
  let json: { [k: string]: any } = {};
  map.forEach((value: Map<string, any>, key: string) => {
    const result = Object.fromEntries(value);
    json[key] = result;
  });

  return json;
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

calculateRewardsFactors();
