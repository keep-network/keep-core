#!/bin/bash
set -eou pipefail

LOG_START='\n\e[1;36m'           # new line + bold + color
LOG_END='\n\e[0m'                # new line + reset color
DONE_START='\n\e[1;32m'          # new line + bold + green
DONE_END='\n\n\e[0m'             # new line + reset
LOG_WARNING_START='\n\e\033[33m' # new line + bold + warning color
LOG_WARNING_END='\n\e\033[0m'    # new line + reset

PROMETHEUS_API_DEFAULT="https://monitoring.threshold.network/prometheus/api/v1"
PROMETHEUS_JOB_DEFAULT="keep-discovered-nodes"
REWARDS_JSON_DEFAULT="./rewards.json"
ETHERSCAN_API_DEFAULT="https://api.etherscan.io"
NETWORK_DEFAULT="mainnet"
KEEP_CORE_REPO="https://github.com/keep-network/keep-core"
# Special case when calculating rewards
OCTOBER_17="2022-10-17"

help() {
  echo -e "\nUsage: $0" \
    "--rewards-start-date <rewards-start-date Y-m-d>" \
    "--rewards-end-date <rewards-end-date Y-m-d>" \
    "--etherscan-token <etherscan-token>" \
    "--prometheus-api <prometheus-api-address>" \
    "--prometheus-job <prometheus-job-name>" \
    "--etherscan-api <etherscan-api-url>" \
    "--network <network-name>" \
    "--rewards-json <rewards-json-output-path>"
  echo -e "\nRequired command line arguments:\n"
  echo -e "\t--rewards-start-date: Rewards interval start date formatted as Y-m-d"
  echo -e "\t--rewards-end-date: Rewards interval end date formatted as Y-m-d"
  echo -e "\t--etherscan-token: Etherscan API key token"
  echo -e "\nOptional command line arguments:\n"
  echo -e "\t--prometheus-api: Prometheus API. Default: ${PROMETHEUS_API_DEFAULT}"
  echo -e "\t--prometheus-job: Prometheus service discovery job name. Default: ${PROMETHEUS_JOB_DEFAULT}"
  echo -e "\t--etherscan-api: Etherscan API url. Default: ${ETHERSCAN_API_DEFAULT}"
  echo -e "\t--network: Network name. Default: ${NETWORK_DEFAULT}"
  echo -e "\t--rewards-json: Rewards JSON output path. Default: ${REWARDS_JSON_DEFAULT}"
  echo -e ""
  exit 1 # Exit script after printing help
}

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
  "--rewards-start-date") set -- "$@" "-k" ;;
  "--rewards-end-date") set -- "$@" "-e" ;;
  "--etherscan-token") set -- "$@" "-t" ;;
  "--etherscan-api") set -- "$@" "-r" ;;
  "--prometheus-api") set -- "$@" "-a" ;;
  "--prometheus-job") set -- "$@" "-p" ;;
  "--network") set -- "$@" "-n" ;;
  "--rewards-json") set -- "$@" "-o" ;;
  "--help") set -- "$@" "-h" ;;
  *) set -- "$@" "$arg" ;;
  esac
done

# Parse short options
OPTIND=1
while getopts "k:e:t:r:a:p:n:o:h" opt; do
  case "$opt" in
  k) rewards_start_date="$OPTARG" ;;
  e) rewards_end_date="$OPTARG" ;;
  t) etherscan_token="$OPTARG" ;;
  r) etherscan_api="$OPTARG" ;;
  a) prometheus_api="$OPTARG" ;;
  p) prometheus_job="$OPTARG" ;;
  n) network="$OPTARG" ;;
  o) rewards_json="$OPTARG" ;;
  h) help ;;
  ?) help ;; # Print help in case parameter is non-existent
  esac
done
shift $(expr $OPTIND - 1) # remove options from positional parameters

ETHERSCAN_TOKEN=${etherscan_token:-""}
REWARDS_START_DATE=${rewards_start_date:-""}
REWARDS_END_DATE=${rewards_end_date:-""}
PROMETHEUS_API=${prometheus_api:-${PROMETHEUS_API_DEFAULT}}
PROMETHEUS_JOB=${prometheus_job:-${PROMETHEUS_JOB_DEFAULT}}
REWARDS_JSON=${rewards_json:-${REWARDS_JSON_DEFAULT}}
ETHERSCAN_API=${etherscan_api:-${ETHERSCAN_API_DEFAULT}}
NETWORK=${network:-${NETWORK_DEFAULT}}

if [ "$REWARDS_START_DATE" == "" ]; then
  printf "${LOG_WARNING_START}Rewards start date must be provided.${LOG_WARNING_END}"
  help
fi

if [ "$REWARDS_END_DATE" == "" ]; then
  printf "${LOG_WARNING_START}Rewards end date must be provided.${LOG_WARNING_END}"
  help
fi

if [ "$ETHERSCAN_TOKEN" == "" ]; then
  printf "${LOG_WARNING_START}Etherscan API key token must be provided.${LOG_WARNING_END}"
  help
fi

rewardsStartDate=$(TZ=UTC date -j -f "%Y-%m-%d %H:%M:%S" "${REWARDS_START_DATE} 00:00:00" "+%s")
rewardsEndDate=$(TZ=UTC date -j -f "%Y-%m-%d %H:%M:%S" "${REWARDS_END_DATE} 23:59:59" "+%s")
october17=$(TZ=UTC date -j -f "%Y-%m-%d %H:%M:%S" "${OCTOBER_17} 23:59:59" "+%s")

startBlockApiCall="${ETHERSCAN_API}/api?\
module=block&\
action=getblocknobytime&\
timestamp=$rewardsStartDate&\
closest=after&\
apikey=${ETHERSCAN_TOKEN}"

endBlockApiCall="${ETHERSCAN_API}/api?\
module=block&\
action=getblocknobytime&\
timestamp=$rewardsEndDate&\
closest=after&\
apikey=${ETHERSCAN_TOKEN}"

october17ApiCall="${ETHERSCAN_API}/api?\
module=block&\
action=getblocknobytime&\
timestamp=$october17&\
closest=after&\
apikey=${ETHERSCAN_TOKEN}"

startRewardsBlock=$(curl -s $startBlockApiCall | jq '.result|tonumber')
endRewardsBlock=$(curl -s $endBlockApiCall | jq '.result|tonumber')
october17Block=$(curl -s $october17ApiCall | jq '.result|tonumber')

printf "${LOG_START}Installing yarn dependencies...${LOG_END}"
yarn install

printf "${LOG_START}Retrieving client release tags...${LOG_END}"
merkleDistributionRepo=$(git remote get-url origin)
git remote set-url origin ${KEEP_CORE_REPO}
git fetch --all --tags --prune
allTags=($(git tag --sort=-version:refname --list 'v[0-9]*.*-m[0-9]*'))
latestTag=${allTags[0]}
latestTimestamp=$(git show -s --format=%ct ${latestTag}^{commit})
latestTagTimestamp="${latestTag}_$latestTimestamp"

tagsInRewardInterval=()

if [ ${#allTags[@]} -gt 1 ]; then
  secondToLatestTag=${allTags[1]}
  secondToLatestTagTimestamp="${secondToLatestTag}_$(git show -s --format=%ct ${secondToLatestTag}^{commit})"
  if [ $latestTimestamp -gt $rewardsStartDate ] && [ $latestTimestamp -lt $rewardsEndDate ]; then
    # The latest tag was created within the rewards interval dates.
    tagsInRewardInterval+=($latestTagTimestamp)
    tagsInRewardInterval+=($secondToLatestTagTimestamp)
  elif [ $latestTimestamp -gt $rewardsEndDate ]; then
    # The latest tag was created after the given rewards interval. Take a
    # second to latest tag.
    tagsInRewardInterval+=($secondToLatestTagTimestamp)
  fi
fi

if [ ${#tagsInRewardInterval[@]} -eq 0 ]; then
  # There is only one tag or the latest tag was created before the start rewards
  # interval and it is continue being the latest tag. (No new releases)
  tagsInRewardInterval+=($latestTagTimestamp)
fi

# Converting array to string so we can pass to the rewards-requirements.ts
printf -v tags '%s|' "${tagsInRewardInterval[@]}"
tagsTrimmed="${tags%?}" # remove "|" at the end

# Setting remote back to merkle distribution repo
git remote set-url origin $merkleDistributionRepo

# Run script
printf "${LOG_START}Fetching peers data...${LOG_END}"

ETHERSCAN_TOKEN=${ETHERSCAN_TOKEN} yarn rewards \
  --api ${PROMETHEUS_API} \
  --job ${PROMETHEUS_JOB} \
  --start-timestamp $rewardsStartDate \
  --end-timestamp $rewardsEndDate \
  --start-block $startRewardsBlock \
  --end-block $endRewardsBlock \
  --october17 $october17Block \
  --releases $tagsTrimmed \
  --network ${NETWORK} \
  --output ${REWARDS_JSON}

printf "${DONE_START}Complete!${DONE_END}"
