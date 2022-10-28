#!/bin/bash
set -eou pipefail

LOG_START='\n\e[1;36m'           # new line + bold + color
LOG_END='\n\e[0m'                # new line + reset color
DONE_START='\n\e[1;32m'          # new line + bold + green
DONE_END='\n\n\e[0m'             # new line + reset
LOG_WARNING_START='\n\e\033[33m' # new line + bold + warning color
LOG_WARNING_END='\n\e\033[0m'    # new line + reset

PROMETHEUS_API_DEFAULT="https://monitoring.test.threshold.network/prometheus/api/v1"
PROMETHEUS_JOB_DEFAULT="keep-discovered-nodes"
PROMETHEUS_SCRAPE_INTERVAL_DEFAULT=60 # in sec
OUTPUT_JSON_FILE="peersData.json"
ETHERSCAN_API="https://api-goerli.etherscan.io" # TODO: change to mainnet https://api.etherscan.io/

help() {
  echo -e "\nUsage: $0" \
    "--rewards-start-date <rewards-start-date Y-m-d>" \
    "--rewards-end-date <rewards-end-date Y-m-d>" \
    "--etherscan-token <etherscan-token>" \
    "--prometheus-api <prometheus-api-address>" \
    "--prometheus-job <prometheus-job-name>" \
    "--prometheus-scrape-interval <prometheus-scrape-interval-in-sec>"
  echo -e "\nRequired command line arguments:\n"
  echo -e "\t--rewards-start-date: Rewards interval start date formatted as Y-m-d"
  echo -e "\t--rewards-end-date: Rewards interval end date formatted as Y-m-d"
  echo -e "\t--etherscan-token: Etherscan API key token"
  echo -e "\nOptional command line arguments:\n"
  echo -e "\t--prometheus-api: Prometheus API. Default: ${PROMETHEUS_API_DEFAULT}"
  echo -e "\t--prometheus-job: Prometheus service discovery job name. Default: ${PROMETHEUS_JOB_DEFAULT}"
  echo -e "\t--prometheus-interval: Prometheus scrape interval. Default: ${PROMETHEUS_SCRAPE_INTERVAL_DEFAULT} sec."
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
  "--prometheus-api") set -- "$@" "-a" ;;
  "--prometheus-job") set -- "$@" "-p" ;;
  "--prometheus-interval") set -- "$@" "-i" ;;
  "--help") set -- "$@" "-h" ;;
  *) set -- "$@" "$arg" ;;
  esac
done

# Parse short options
OPTIND=1
while getopts "k:e:t:a:p:i:h" opt; do
  case "$opt" in
  k) rewards_start_date="$OPTARG" ;;
  e) rewards_end_date="$OPTARG" ;;
  t) etherscan_token="$OPTARG" ;;
  a) prometheus_api="$OPTARG" ;;
  p) prometheus_job="$OPTARG" ;;
  i) prometheus_scrape_interval="$OPTARG" ;;
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
PROMETHEUS_SCRAPE_INTERVAL=${prometheus_scrape_interval:-${PROMETHEUS_SCRAPE_INTERVAL_DEFAULT}}

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

startRewardsBlock=$(curl -s $startBlockApiCall | jq '.result|tonumber')
endRewardsBlock=$(curl -s $endBlockApiCall | jq '.result|tonumber')

printf "${LOG_START}Installing yarn dependencies...${LOG_END}"
yarn install

# allTags=($(git tag --sort=-version:refname --list 'v[0-9]*.*-m[0-9]*')) # TODO: REPLACE, this is correct regex to match latest versions
allTags=($(git tag --sort=-version:refname --list 'v[0-9]*')) # sorted latest -> oldest # TODO: For testing
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

# Run script
printf "${LOG_START}Fetching peers data...${LOG_END}"

ETHERSCAN_TOKEN=${ETHERSCAN_TOKEN} yarn rewards-requirements \
  --api ${PROMETHEUS_API} \
  --job ${PROMETHEUS_JOB} \
  --start-timestamp $rewardsStartDate \
  --end-timestamp $rewardsEndDate \
  --start-block $startRewardsBlock \
  --end-block $endRewardsBlock \
  --interval ${PROMETHEUS_SCRAPE_INTERVAL} \
  --releases $tagsTrimmed \
  --output ${OUTPUT_JSON_FILE}

printf "${DONE_START}Complete!${DONE_END}"
