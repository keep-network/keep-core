#!/bin/bash
set -eou pipefail

LOG_START='\n\e[1;36m'  # new line + bold + color
LOG_END='\n\e[0m'       # new line + reset color
DONE_START='\n\e[1;32m' # new line + bold + green
DONE_END='\n\n\e[0m'    # new line + reset

PROMETHEUS_API_DEFAULT="http://prometheus.monitoring.svc.cluster.local:9090/api/v1"
PROMETHEOUS_JOB_DEFAULT="keep-discovered-nodes"
REWARDS_END_DATE_DEFAULT=$(date +"%Y-%m-%d")
OUTPUT_JSON_FILE="peersData.json"
BUCKET_NAME_DEFAULT="diagnostics_test"
CLIENT_UPGRADE_DELAY_ACCEPTANCE=604800 # 7days in sec

help() {
  echo -e "\nUsage: $0" \
    "--prometheus-api <prometheus-api-address>" \
    "--prometheus-job <prometheus-job>" \
    "--rewards-start-date <rewards-start-date Y-m-d>" \
    "--rewards-end-date <rewards-end-date Y-m-d>" \
    "--bucket-name <GCP-bucket-name>"
  echo -e "\nCommand line arguments:\n"
  echo -e "\t--prometheus-api: Prometheus API"
  echo -e "\t--prometheus-job: Prometheus service discovery job name"
  echo -e "\t--rewards-start-date: Rewards start date Y-m-d"
  echo -e "\t--rewards-end-date: Rewards end date Y-m-d"
  echo -e "\t--bucket-name: GCP destination bucket name where peer data are stored\n"
  exit 1 # Exit script after printing help
}

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
  "--prometheus-api") set -- "$@" "-i" ;;
  "--prometheus-job") set -- "$@" "-p" ;;
  "--rewards-start-date") set -- "$@" "-k" ;;
  "--rewards-end-date") set -- "$@" "-e" ;;
  "--bucket-name") set -- "$@" "-b" ;;
  "--help") set -- "$@" "-h" ;;
  *) set -- "$@" "$arg" ;;
  esac
done

# Parse short options
OPTIND=1
while getopts "a:d:p:k:e:b:h" opt; do
  case "$opt" in
  a) prometheus_api="$OPTARG" ;;
  p) prometheous_job="$OPTARG" ;;
  k) rewards_start_date="$OPTARG" ;;
  e) rewards_end_date="$OPTARG" ;;
  b) bucket_name="$OPTARG" ;;
  h) help ;;
  ?) help ;; # Print help in case parameter is non-existent
  esac
done
shift $(expr $OPTIND - 1) # remove options from positional parameters

REWARDS_START_DATE=${rewards_start_date:-""}
REWARDS_END_DATE=${rewards_end_date:-${REWARDS_END_DATE_DEFAULT}}
PROMETHEUS_API=${prometheus_api:-${PROMETHEUS_API_DEFAULT}}
PROMETHEOUS_JOB=${prometheous_job:-${PROMETHEOUS_JOB_DEFAULT}}
BUCKET_NAME=${bucket_name:-${BUCKET_NAME_DEFAULT}}

if [ "$REWARDS_START_DATE" == "" ]; then
  printf "${LOG_WARNING_START}Rewards start date must be provided.${LOG_WARNING_END}"
  exit 1
fi

rewardsStartDate=$(date -j -f "%Y-%m-%d" ${REWARDS_START_DATE} "+%s")
rewardsEndDate=$(date -j -f "%Y-%m-%d" ${REWARDS_END_DATE} "+%s")

printf "${LOG_START}Installing yarn dependencies...${LOG_END}"
yarn install

# latestTags=$(git tag --sort=version:refname --list 'v[0-9]*') # sorted oldest -> latest
latestTags=$(git tag --sort=-version:refname --list 'v[0-9]*') # sorted latest -> oldest
# latestTags=$(git tag --sort=version:refname --list 'v[0-9]*.*-m[0-9]*') # sorted oldest -> latest
tagsInRewardInterval=()
for tag in $latestTags; do
  timestamp=$(git show -s --format=%ct $tag^{commit})
  # Getting the client tags that were released 7 days prior to the rewards interval
  # start time and 7 days prior to the interval end. If a new client is released
  # ex. 3 days before the end of the interval and a peer didn't upgrade the client
  # it is not taken into account for the current rewards interval. It will be
  # measured only in the next reward interval. In short, the rewards interval is
  # offset by 7 days.
  rewardsStartOffset=$(($rewardsStartDate - ${CLIENT_UPGRADE_DELAY_ACCEPTANCE}))
  rewardsEndOffset=$(($rewardsEndDate - ${CLIENT_UPGRADE_DELAY_ACCEPTANCE}))
  if [ "$rewardsStartOffset" -lt "$timestamp" ] && [ "$timestamp" -le "$rewardsEndOffset" ]; then
    tagTimestamp="${tag}_$(git show -s --format=%ct $tag^{commit})"
    tagsInRewardInterval+=($tagTimestamp)
  elif [ "$timestamp" -le "$rewardsStartOffset" ]; then
    # current tag (in case no new tags in a given rewards interval)
    # or one before the latest to check what version was before the newest tag
    tagTimestamp="${tag}_$(git show -s --format=%ct $tag^{commit})"
    tagsInRewardInterval+=($tagTimestamp)
    break
  fi
done

# Converting array to string so we can pass to the rewards-requirements.ts
printf -v tags '|%s' "${tagsInRewardInterval[@]}"

echo $tags

# Run script
printf "${LOG_START}Fetching peers data...${LOG_END}"

yarn rewards-requirements \
  --api ${PROMETHEUS_API} \
  --job ${PROMETHEOUS_JOB} \
  --start $rewardsStartDate \
  --end $rewardsEndDate \
  --interval 5 \
  --versions $tags \
  --output ${OUTPUT_JSON_FILE}

# TODO: do we want to upload the output file to a GCP bucket?

printf "${DONE_START}Complete!${DONE_END}"
