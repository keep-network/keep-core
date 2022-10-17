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
  ALLOWED_UPGRADE_DELAY,
  CLIENT_TIMESTAMP_INDEX,
  CLIENT_VERSION_INDEX,
  DEFAULT_PROVIDER,
  IS_BEACON_AUTHORIZED_FACTOR,
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
  const prometheusJob = options.job;
  const prometheusAPI = options.api;
  const clientVersions = options.versions.split("|"); // sorted from latest to oldest
  const startRewardsTimestamp = parseInt(options.start);
  const endRewardsTimestamp = parseInt(options.end);
  const scrapeInterval = parseInt(options.interval);
  const peersDataFile = options.output;
  // End program option parsing

  if (Date.now() / 1000 < endRewardsTimestamp) {
    console.log("End time interval must be in the past");
    return "End time interval must be in the past";
  }

  const rewardsInterval = endRewardsTimestamp - startRewardsTimestamp;
  const offset = Math.floor(Date.now() / 1000) - endRewardsTimestamp;

  const prometheusAPIQuery = `${prometheusAPI}/query`;
  const queryBootstrapData = `${prometheusAPI}/query_range`;

  // Query for bootstrap data that has peer instances
  const paramsBootstrapData = {
    query: `up{job='${prometheusJob}'}`,
    start: startRewardsTimestamp,
    end: endRewardsTimestamp,
    step: QUERY_STEP,
  };

  const bootstrapData = (await queryPrometheus(
    queryBootstrapData,
    paramsBootstrapData
  )).data.result;

  const provider = ethers.getDefaultProvider(DEFAULT_PROVIDER);
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

  let operatorIds = new Set<string>();
  for (let i = 0; i < bootstrapData.length; i++) {
    operatorIds.add(bootstrapData[i].metric.chain_address)
  }

  let rewardData = new Map<string, Map<string, number>>()
  for (let operatorId of operatorIds) {
    let operatorData = new Map<string, number>()

    const paramsUptime = {
      query: `sum_over_time(up{chain_address='${operatorId}'}[${rewardsInterval}s:2m]`+
             `offset ${offset}s) / ${rewardsInterval} * 120`
    };

    const resultUptime = (await queryPrometheus(
      prometheusAPIQuery,
      paramsUptime
    )).data.result;

    let uptimeByInstance = new Map<string, number>
    let totalUptime: number = 0
    for (let uptimeData of resultUptime) {
      const uptime = Number(uptimeData.value[1])
      uptimeByInstance.set(uptimeData.metric.instance, uptime)
      totalUptime += uptime
    }

    let normalizedUptimeByInstance = new Map<string, number>
    for (let [instance, uptime] of uptimeByInstance) {
      normalizedUptimeByInstance.set(instance, uptime / totalUptime)
    }


    console.log(uptimeByInstance)
    console.log(normalizedUptimeByInstance)
  }
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
