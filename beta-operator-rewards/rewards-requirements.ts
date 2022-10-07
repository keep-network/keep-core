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
import axios from "axios";
import {
  allowedUpgradeDelay,
  clientTimestampIndex,
  clientVersionIndex,
  defaultProvider,
  minPreParams,
  preParamsAvgInterval,
  preParamsResolution,
  queryStep,
  requiredUptime,
  upTimeRewardsCoefficient,
} from "./rewards-constants"

export async function calculateRewardsFactors() {
  program
    .version("0.0.1")
    .requiredOption(
      "-s, --start <timestamp>",
      "starting time for rewards calculation"
    )
    .requiredOption(
      "-e, --end <timestamp>",
      "ending time for rewards calculation"
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
  const prometheus_job = options.job;
  const prometheusAPI = options.api;
  const clientVersions = options.versions.split("|"); // sorted from latest to oldest
  const startRewardsTimestamp = parseInt(options.start);
  const endRewardsTimestamp = parseInt(options.end);
  const scrapeInterval = parseInt(options.interval);
  const peersDataFile = options.output;
  // End program option parsing

  const rewardsInterval = endRewardsTimestamp - startRewardsTimestamp;

  const factors = {
    isBeaconAuthorized: "isBeaconAuthorized",
    isTbtcAuthorized: "isTbtcAuthorized",
    upTime: "upTime",
    version: "version",
    preParams: "preParams",
  };

  const prometheusAPIQuery = `${prometheusAPI}/query`;
  const queryBootstrapData = `${prometheusAPI}/query_range`;

  // Query for bootstrap data that has peer instances
  const paramsBootstrapData = {
    query: `up{job='${prometheus_job}'}`,
    start: startRewardsTimestamp,
    end: endRewardsTimestamp,
    step: queryStep,
  };

  const bootstrapData = await queryPrometheus(
    queryBootstrapData,
    paramsBootstrapData
  );

  let peersData = new Map<string, Map<string, number>>(); // peer address -> {component name: factor}

  const provider = ethers.getDefaultProvider(defaultProvider);
  const randomBeacon = new Contract(
    RandomBeaconAddress,
    JSON.stringify(RandomBeaconABI),
    provider
  );

  const walletRegistry = new Contract(
    WalletRegistryAddress,
    JSON.stringify(WalletRegistryABI),
    provider
  );

  if (Date.now() / 1000 < endRewardsTimestamp) {
    console.log("End time interval must be in the past");
    return "End time interval must be in the past";
  }

  for (let i = 0; i < (await bootstrapData.data.result.length); i++) {
    const peer = bootstrapData.data.result[i];
    let peerData = new Map<string, any>(); // Map<factor_name: value> value is in the range from 0 to 1
    peerData.set("address", peer.metric.chain_address);

    /// Random Beacon application authorization requirement

    const stakingProviderAddressForBeacon =
      await randomBeacon.operatorToStakingProvider(peer.metric.chain_address);
    const eligibleStakeForBeacon = await walletRegistry.eligibleStake(
      stakingProviderAddressForBeacon
    );
    if (eligibleStakeForBeacon.isZero()) {
      peerData.set(factors.isBeaconAuthorized, 0);
    } else {
      peerData.set(factors.isBeaconAuthorized, 1);
    }

    /// tBTC application authorized requirement

    const stakingProviderAddressForTbtc =
      await walletRegistry.operatorToStakingProvider(peer.metric.chain_address);
    const eligibleStakeForTbtc = await walletRegistry.eligibleStake(
      stakingProviderAddressForTbtc
    );
    if (eligibleStakeForTbtc.isZero()) {
      peerData.set(factors.isTbtcAuthorized, 0);
    } else {
      peerData.set(factors.isTbtcAuthorized, 1);
    }

    /// Up time requirement

    // First registered 'up' metric in a given timeframe <start:end>. We start
    // evaluating uptime from this point.
    const firstRegisteredUptime = peer.values[0][0];
    const uptimeSearchRange = endRewardsTimestamp - firstRegisteredUptime;
    // Offset is set in case the end time interval is not aligned with execution
    // of this script. It "goes" back in time relevant to the current time.
    const offset = Math.floor(Date.now() / 1000) - endRewardsTimestamp;
    // Sum of all uptimes since the endpoint became available in a given
    // timeframe. "up" metric won't take into account when a node wasn't available,
    // hence we need to multiply the result by the scrape interval
    // (set in the config file) and divide by the uptime search range.
    const paramsUptime = {
      query: `sum_over_time(up{instance='${peer.metric.instance}', job='${prometheus_job}'}
              [${uptimeSearchRange}s] offset ${offset}s) * ${scrapeInterval} / ${uptimeSearchRange}`,
    };
    const resultUptime = await queryPrometheus(
      prometheusAPIQuery,
      paramsUptime
    );
    const resultUptimePercent = resultUptime.data.result[0].value[1] * 100;
    const upFactor = resultUptimePercent < requiredUptime ? 0 : 1;
    peerData.set(factors.upTime, upFactor);
    const upFactorCoefficient = upFactor
      ? uptimeSearchRange / rewardsInterval
      : 0;
    // Rewards should be adjusted by the upFactorCoefficient for a given peer if
    // that peer joins the network later relative to the rewards interval start.
    // Ex. if a peer joins mid month and all other factors are satisfied, then
    // the rewards are devided by half.
    peerData.set(upTimeRewardsCoefficient, upFactorCoefficient);

    /// Pre-param requirement

    // <func>(<metric>{<labels>}[<local_range>] offset <time>)[<global_range>:<resolution>]
    const paramsPreParams = {
      query: `avg_over_time(tbtc_pre_params_count{instance='${peer.metric.instance}', job='${prometheus_job}'}
              [${preParamsAvgInterval}] offset ${offset}s)[${rewardsInterval}s:${preParamsResolution}]`,
    };
    const resultPreParams = await queryPrometheus(
      prometheusAPIQuery,
      paramsPreParams
    );
    peerData.set(factors.preParams, 1);
    if (resultPreParams.data.result.length == 0) {
      peerData.set(factors.preParams, 0);
    } else {
      resultPreParams.data.result[0].values.forEach(function (
        peerPreParams: any
      ) {
        if (Number(peerPreParams[1]) < minPreParams) {
          peerData.set(factors.preParams, 0);
        }
      });
    }

    /// Version requirement (One-week delay in updates to the most recent version)

    // Check a peer's version at the end of the rewards interval
    const buildVersionParams = {
      query: `client_info{instance='${peer.metric.instance}', job='${prometheus_job}'} @ ${endRewardsTimestamp}`,
    };
    const resultBuildVersion = await queryPrometheus(
      prometheusAPIQuery,
      buildVersionParams
    );

    if (resultBuildVersion.data !== undefined && resultBuildVersion.data.result.length > 0) {
      const peerVersion =
        resultBuildVersion.data.result[0].metric.version.split("-")[0];
      const latestClientVersionInfo = clientVersions[0].split("_");
      if (clientVersions.length > 1) {
        const oneBeforeLatestClientVersionInfo = clientVersions[0].split("_");
        if (
          latestClientVersionInfo[clientTimestampIndex] <
          endRewardsTimestamp - allowedUpgradeDelay
        ) {
          // Latest version was released prior to a delay threshold.
          // Peer's version must be the latest client's version.
          if (peerVersion === latestClientVersionInfo[clientVersionIndex]) {
            peerData.set(factors.version, 1);
          } else {
            peerData.set(factors.version, 0);
          }
        } else {
          // Latest version was released in the allowed delay window.
          // Peer's version should match the latest or one before the latest client's
          // version.
          if (
            peerVersion === latestClientVersionInfo[clientVersionIndex] ||
            peerVersion === oneBeforeLatestClientVersionInfo[clientVersionIndex]
          ) {
            peerData.set(factors.version, 1);
          } else {
            peerData.set(factors.version, 0);
          }
        }
      } else {
        // Latest release was done prior to the start interval
        // Peer's version must be the latest
        if (peerVersion === latestClientVersionInfo[clientVersionIndex]) {
          peerData.set(factors.version, 1);
        } else {
          peerData.set(factors.version, 0);
        }
      }
    } else {
      // A peer doesn't metric any build versions
      peerData.set(factors.version, 0);
    }

    peersData.set(peer.metric.instance, peerData);
    console.log("peersDataFactors", peersData);
  }

  // TODO: calculate rewards for a given address
  //
  // for a given peer:
  // - check if all the requirements were satisfied (factors.* == 1)
  // - if all the reqs ^ are satisfied calculate the rewards:
  // -- peerRewards = (peer's authorized stake / total authorized stake) * rewardsForAGivenMonth * upTimeRewardsCoefficient

  const jsonObject = await convertToJSON(peersData);
  // Save to file
  fs.writeFileSync(peersDataFile, JSON.stringify(jsonObject, null, 2));
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
